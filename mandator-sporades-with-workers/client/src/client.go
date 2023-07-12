package src

import (
	"async-consensus/common"
	"async-consensus/configuration"
	"async-consensus/proto"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

/*
	This file defines the client struct and the new method that is invoked when creating a new client by the main
*/

type Client struct {
	clientName int32 // unique client identifier as defined in the configuration.yml
	numWorkers int   // number of worker nodes

	workerAddrList             map[int32]string // map with the IP:port address of every worker node
	workerArrayIndex           map[int32]int    // map each worker name to a unique index, this is required because the identifiers can be arbitary
	incomingWorkerReaders      []*bufio.Reader  // socket readers for each worker
	outgoingWorkerWriters      []*bufio.Writer  // socket writer for each worker
	outgoingWorkerWriterMutexs []sync.Mutex     // for mutual exclusion for each buffio.writer outgoingWorkerWriters

	rpcTable     map[uint8]*common.RPCPair // map each RPC type (message type) to its unique number
	incomingChan chan *common.RPCPair      // used to collect ClientBatch messages for responses and Status messages for responses (basically all the incoming messages)

	messageCodes proto.MessageCode
	logFilePath  string // the path to write the requests and responses, used for sanity checks

	clientBatchSize int // maximum client side batch size
	clientBatchTime int // maximum client side batch time in micro seconds

	outgoingMessageChan chan *common.OutgoingRPC // buffer for messages that are written to the wire

	defaultWorkers       []int32        // ids of the default workers to which the client sends the request batches
	defaultWorkerMutexes []sync.RWMutex // Since the default worker array will be accessed by sending thread and the failure detector, we need mutual exclusion

	workerTimeout       int          // in seconds: each client has a unique set of workers assigned, upon a timeout, the client thread changes its default workers
	lastSeenTimes       []time.Time  // the last time each worker sent a response (includes all workers not just default writers)
	lastSeenTimeMutexes []sync.Mutex // mutexes for lastSeenTimes

	debugOn    bool // if turned on, the debug messages will be printed on the console
	debugLevel int  // current debug level

	requestSize  int // size of the request payload in bytes (applicable only for the no-op application / testing purposes)
	testDuration int // test duration in seconds
	arrivalRate  int // poisson rate of the arrivals (requests per second)

	arrivalTimeChan     chan int64              // channel to which the poisson process adds new request arrival times in nanoseconds w.r.t test start time
	arrivalChan         chan bool               // channel to which the main scheduler adds new request indications, to be consumed by the request generation threads
	RequestType         string                  // [request] for sending a stream of client requests, [status] for sending a status request
	OperationType       int                     // status operation type 1 (bootstrap server), 2: print log, 3: start consensus
	sentRequests        [][]requestBatch        // generator i updates sentRequests[i] :this is to avoid concurrent access to the same array
	receivedResponses   map[string]requestBatch // set of received client response batches from workers: a map is used for fast lookup
	startTime           time.Time               // test start time
	clientListenAddress string                  // TCP address to which the client listens to new incoming TCP connections
	keyLen              int                     // length of key
	valueLen            int                     // length of value
}

/*
	requestBatch contains a batch that was written to wire, and the time it was written
*/

type requestBatch struct {
	batch proto.ClientBatch
	time  time.Time
}

const statusTimeout = 20              // time to wait for a status request in seconds
const numOutgoingThreads = 100        // number of wire writers: since the I/O writing is expensive we delegate that task to a thread pool and separate from the critical path
const numRequestGenerationThreads = 4 // number of  threads that generate client requests upon receiving an arrival indication todo try different values for this: lower values result in big batches
const incomingBufferSize = 1000000    // the size of the buffer which receives all the incoming messages (client response batch messages and client status response message)
const outgoingBufferSize = 1000000    // size of the buffer that collects messages to be written to the wire
const arrivalBufferSize = 1000000     // size of the buffer that collects new request arrivals

/*
	Instantiate a new Client instance, allocate the buffers
*/

func New(name int32, cfg *configuration.InstanceConfig, logFilePath string, clientBatchSize int, clientBatchTime int, defaultWorkers []int32, workerTimeout int, requestSize int, testDuration int, arrivalRate int, requestType string, operationType int, debugOn bool, debugLevel int, keyLen int, valLen int) *Client {
	cl := Client{
		clientName:                 name,
		numWorkers:                 len(cfg.Workers),
		workerAddrList:             make(map[int32]string),
		workerArrayIndex:           make(map[int32]int),
		incomingWorkerReaders:      make([]*bufio.Reader, len(cfg.Workers)),
		outgoingWorkerWriters:      make([]*bufio.Writer, len(cfg.Workers)),
		outgoingWorkerWriterMutexs: make([]sync.Mutex, len(cfg.Workers)),
		rpcTable:                   make(map[uint8]*common.RPCPair),
		incomingChan:               make(chan *common.RPCPair, incomingBufferSize),
		messageCodes:               proto.GetRPCCodes(),
		logFilePath:                logFilePath,
		clientBatchSize:            clientBatchSize,
		clientBatchTime:            clientBatchTime,
		outgoingMessageChan:        make(chan *common.OutgoingRPC, outgoingBufferSize),

		defaultWorkers:       defaultWorkers,
		defaultWorkerMutexes: make([]sync.RWMutex, len(defaultWorkers)),
		workerTimeout:        workerTimeout,
		lastSeenTimes:        make([]time.Time, len(cfg.Workers)),
		lastSeenTimeMutexes:  make([]sync.Mutex, len(cfg.Workers)),

		debugOn:    debugOn,
		debugLevel: debugLevel,

		requestSize:         requestSize,
		testDuration:        testDuration,
		arrivalRate:         arrivalRate,
		arrivalTimeChan:     make(chan int64, arrivalBufferSize),
		arrivalChan:         make(chan bool, arrivalBufferSize),
		RequestType:         requestType,
		OperationType:       operationType,
		sentRequests:        make([][]requestBatch, numRequestGenerationThreads),
		receivedResponses:   make(map[string]requestBatch),
		startTime:           time.Now(),
		clientListenAddress: common.GetAddress(cfg.Clients, name),
		keyLen:              keyLen,
		valueLen:            valLen,
	}

	common.Debug("Created a new client instance", 0, cl.debugLevel, cl.debugOn)

	// initialize workerAddrList
	for i := 0; i < len(cfg.Workers); i++ {
		int32Name, _ := strconv.ParseInt(cfg.Workers[i].Name, 10, 32)
		cl.workerAddrList[int32(int32Name)] = cfg.Workers[i].Address
		cl.workerArrayIndex[int32(int32Name)] = i
		cl.outgoingWorkerWriterMutexs[i] = sync.Mutex{}
	}

	/*
		Register the rpcs
	*/
	cl.RegisterRPC(new(proto.ClientBatch), cl.messageCodes.ClientBatchRpc)
	cl.RegisterRPC(new(proto.Status), cl.messageCodes.StatusRPC)
	cl.RegisterRPC(new(proto.MemPoolMini), cl.messageCodes.MemPoolMiniRPC)
	cl.RegisterRPC(new(proto.MemPool), cl.messageCodes.MemPoolRPC)

	common.Debug("Registered RPCs in the table", 0, cl.debugLevel, cl.debugOn)

	// Set random seed
	rand.Seed(time.Now().UTC().UnixNano())

	// set mutexes for default workers
	for i := 0; i < len(defaultWorkers); i++ {
		cl.defaultWorkerMutexes[i] = sync.RWMutex{}
	}

	pid := os.Getpid()
	fmt.Printf("initialized client %v with process id: %v \n", name, pid)

	return &cl
}

/*
	Fill the RPC table by assigning a unique id to each message type
*/

func (cl *Client) RegisterRPC(msgObj proto.Serializable, code uint8) {
	cl.rpcTable[code] = &common.RPCPair{Code: code, Obj: msgObj}
}
