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
	This file defines the worker struct and the new method that is invoked when creating a new worker
*/

type Worker struct {
	name          int32  // unique worker identifier as defined in the configuration.yml
	listenAddress string // TCP address to which the worker listens to new incoming TCP connections

	workerAssignment map[int32][]int32 // maps replica to its workers

	workerAddrList             map[int32]string // map with the IP:port address of every worker node
	workerArrayIndex           map[int32]int    // map each worker name to a unique index, this is required because the identifiers can be arbitary
	incomingWorkerReaders      []*bufio.Reader  // socket readers for each worker
	outgoingWorkerWriters      []*bufio.Writer  // socket writer for each worker
	outgoingWorkerWriterMutexs []sync.Mutex     // for mutual exclusion for each buffio.writer outgoingWorkerWriters

	clientAddrList             map[int32]string // map with the IP:port address of every client node
	clientArrayIndex           map[int32]int    // map each client name to a unique index, this is required because the identifiers can be arbitary
	incomingClientReaders      []*bufio.Reader  // socket readers for each client
	outgoingClientWriters      []*bufio.Writer  // socket writer for each client
	outgoingClientWriterMutexs []sync.Mutex     // for mutual exclusion for each buffio.writer outgoingClientWriters

	defaultReplicaName                int32         // name of default replica
	defaultReplicaAddr                string        // address of the default replica
	incomingDefaultReplicaReader      *bufio.Reader // socket readers for default replica
	outgoingDefaultReplicaWriter      *bufio.Writer // socket writer for default replica
	outgoingDefaultReplicaWriterMutex sync.Mutex    // for mutual exclusion for each buffio.writer  default replica writer

	rpcTable     map[uint8]*common.RPCPair // map each RPC type (message type) to its unique number
	messageCodes proto.MessageCode

	incomingChan chan *common.RPCPair // used to collect all the incoming messages

	logFilePath string // the path to write the requests and responses, used for sanity checks

	workerBatchSize int // maximum worker side batch size (mini pool batch size)
	workerBatchTime int // maximum worker side batch time in micro seconds

	outgoingMessageChan chan *common.OutgoingRPC // buffer for messages that are written to the wire

	debugOn    bool // if turned on, the debug messages will be printed on the console
	debugLevel int  // current debug level

	serverStarted bool         // to bootstrap
	miniPool      *MiniMemPool // all the data structs related to worker
}

const numOutgoingThreads = 100      // number of wire writers: since the I/O writing is expensive we delegate that task to a thread pool and separate from the critical path
const incomingBufferSize = 10000000 // the size of the buffer which receives all the incoming messages
const outgoingBufferSize = 10000000 // size of the buffer that collects messages to be written to the wire

/*
	Instantiate a new worker instance, allocate the buffers
*/

func New(name int32, cfg *configuration.InstanceConfig, workerAssignment map[int32][]int32, logFilePath string, workerBatchSize int, workerBatchTime int, debugOn bool, mode int, debuglevel int, window int) *Worker {
	wr := Worker{
		name:          name,
		listenAddress: common.GetAddress(cfg.Workers, name),

		workerAssignment: workerAssignment,

		workerAddrList:             make(map[int32]string),
		workerArrayIndex:           make(map[int32]int),
		incomingWorkerReaders:      make([]*bufio.Reader, len(cfg.Workers)),
		outgoingWorkerWriters:      make([]*bufio.Writer, len(cfg.Workers)),
		outgoingWorkerWriterMutexs: make([]sync.Mutex, len(cfg.Workers)),

		clientAddrList:             make(map[int32]string),
		clientArrayIndex:           make(map[int32]int),
		incomingClientReaders:      make([]*bufio.Reader, len(cfg.Clients)),
		outgoingClientWriters:      make([]*bufio.Writer, len(cfg.Clients)),
		outgoingClientWriterMutexs: make([]sync.Mutex, len(cfg.Clients)),

		defaultReplicaName:                getDefaultReplica(workerAssignment, name),
		defaultReplicaAddr:                common.GetAddress(cfg.Peers, getDefaultReplica(workerAssignment, name)),
		incomingDefaultReplicaReader:      nil,
		outgoingDefaultReplicaWriter:      nil,
		outgoingDefaultReplicaWriterMutex: sync.Mutex{},

		rpcTable:     make(map[uint8]*common.RPCPair),
		messageCodes: proto.GetRPCCodes(),

		incomingChan: make(chan *common.RPCPair, incomingBufferSize),

		logFilePath: logFilePath,

		workerBatchSize: workerBatchSize,
		workerBatchTime: workerBatchTime,

		outgoingMessageChan: make(chan *common.OutgoingRPC, outgoingBufferSize),
		debugOn:             debugOn,
		debugLevel:          debuglevel,
		serverStarted:       false,
	}

	wr.miniPool = InitMiniMemPool(mode, len(cfg.Workers), wr.debugLevel, wr.debugOn, window)

	common.Debug("Created a new worker instance", 0, wr.debugLevel, wr.debugOn)

	// initialize workerAddrList
	for i := 0; i < len(cfg.Workers); i++ {
		int32Name, _ := strconv.ParseInt(cfg.Workers[i].Name, 10, 32)
		wr.workerAddrList[int32(int32Name)] = cfg.Workers[i].Address
		wr.workerArrayIndex[int32(int32Name)] = i
		wr.outgoingWorkerWriterMutexs[i] = sync.Mutex{}
	}

	// initialize clientAddrList
	for i := 0; i < len(cfg.Clients); i++ {
		int32Name, _ := strconv.ParseInt(cfg.Clients[i].Name, 10, 32)
		wr.clientAddrList[int32(int32Name)] = cfg.Clients[i].Address
		wr.clientArrayIndex[int32(int32Name)] = i
		wr.outgoingClientWriterMutexs[i] = sync.Mutex{}
	}
	/*
		Register the rpcs
	*/
	wr.RegisterRPC(new(proto.ClientBatch), wr.messageCodes.ClientBatchRpc)
	wr.RegisterRPC(new(proto.Status), wr.messageCodes.StatusRPC)
	wr.RegisterRPC(new(proto.MemPoolMini), wr.messageCodes.MemPoolMiniRPC)
	wr.RegisterRPC(new(proto.MemPool), wr.messageCodes.MemPoolRPC)

	common.Debug("Registered RPCs in the table", 0, wr.debugLevel, wr.debugOn)

	// Set random seed
	rand.Seed(time.Now().UTC().UnixNano())

	pid := os.Getpid()
	fmt.Printf("initialized worker %v with process id: %v \n", name, pid)

	return &wr
}

/*
	Returns the replica associated with this worker
*/

func getDefaultReplica(assignment map[int32][]int32, name int32) int32 {
	for key, element := range assignment {
		for i := 0; i < len(element); i++ {
			if element[i] == name {
				return key
			}
		}
	}
	return -1
}

/*
	Fill the RPC table by assigning a unique id to each message type
*/

func (wr *Worker) RegisterRPC(msgObj proto.Serializable, code uint8) {
	wr.rpcTable[code] = &common.RPCPair{Code: code, Obj: msgObj}
}

/*
	get node type using name
*/

func (wr *Worker) getNodeType(id int32) string {
	if wr.defaultReplicaName == id {
		return "replica"
	}
	if _, ok := wr.workerAddrList[id]; ok {
		return "worker"
	}
	if _, ok := wr.clientAddrList[id]; ok {
		return "client"
	}
	return ""
}
