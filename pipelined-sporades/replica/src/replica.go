package src

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"pipelined-sporades/common"
	"pipelined-sporades/configuration"
	"pipelined-sporades/proto"
	"strconv"
	"sync"
	"time"
)

/*
	defines the Replica struct and the new method that is invoked when creating a new replica
*/

type Replica struct {
	name          int32  // unique replica identifier as defined in the configuration.yml
	listenAddress string // TCP address to which the replica listens to new incoming TCP connections

	numReplicas int

	clientAddrList             map[int32]string         // map with the IP:port address of every client
	incomingClientReaders      map[int32]*bufio.Reader  // socket readers for each client
	outgoingClientWriters      map[int32]*bufio.Writer  // socket writer for each client
	outgoingClientWriterMutexs map[int32]*sync.Mutex    // for mutual exclusion for each buffio.writer outgoingClientWriters
	outgoingClientMessageChan  chan *common.OutgoingRPC // buffer for messages that are written to clients

	replicaAddrList             map[int32]string           // map with the IP:port address of every replica node
	incomingReplicaReaders      map[int32]*bufio.Reader    // socket readers for each replica
	outgoingReplicaWriters      map[int32]*bufio.Writer    // socket writer for each replica
	outgoingReplicaWriterMutexs map[int32]*sync.Mutex      // for mutual exclusion for each buffio.writer outgoingReplicaWriters
	outgoingReplicaMessageChans []chan *common.OutgoingRPC // buffer for messages that are written to replicas

	rpcTable     map[uint8]*common.RPCPair // map each RPC type (message type) to its unique number
	messageCodes proto.MessageCode

	incomingChan chan *common.RPCPair // used to collect all the incoming messages

	logFilePath string // the path to write the log, used for sanity checks

	replicaBatchSize int // maximum replica side batch size
	replicaBatchTime int // maximum replica side batch time in micro seconds

	debugOn    bool // if turned on, the debug messages will be printed on the console
	debugLevel int  // current debug level

	serverStarted bool // to bootstrap

	consensus        *SporadesConsensus
	consensusStarted bool
	viewTimeout      int // view change timeout in micro seconds

	logPrinted bool // to check if log was printed before

	benchmarkMode int        // 0 for resident K/V store, 1 for redis
	state         *Benchmark // k/v store

	incomingRequests []*proto.ClientBatch
	pipelineLength   int

	lastProposedTime time.Time

	finished bool // to finish consensus
}

const incomingBufferSize = 1000000 // the size of the buffer which receives all the incoming messages
const outgoingBufferSize = 1000000 // size of the buffer that collects messages to be written to the wire

/*
	instantiate a new replica instance, allocate the buffers
*/

func New(name int32, cfg *configuration.InstanceConfig, logFilePath string, replicaBatchSize int, replicaBatchTime int, debugOn bool, debugLevel int, viewTimeout int, benchmarkMode int, keyLen int, valLen int, pipelineLength int) *Replica {
	rp := Replica{
		name:          name,
		listenAddress: common.GetAddress(cfg.Peers, name),
		numReplicas:   len(cfg.Peers),

		clientAddrList:             make(map[int32]string),
		incomingClientReaders:      make(map[int32]*bufio.Reader),
		outgoingClientWriters:      make(map[int32]*bufio.Writer),
		outgoingClientWriterMutexs: make(map[int32]*sync.Mutex),
		outgoingClientMessageChan:  make(chan *common.OutgoingRPC, outgoingBufferSize),

		replicaAddrList:             make(map[int32]string),
		incomingReplicaReaders:      make(map[int32]*bufio.Reader),
		outgoingReplicaWriters:      make(map[int32]*bufio.Writer),
		outgoingReplicaWriterMutexs: make(map[int32]*sync.Mutex),
		outgoingReplicaMessageChans: make([]chan *common.OutgoingRPC, len(cfg.Peers)),

		rpcTable:     make(map[uint8]*common.RPCPair),
		messageCodes: proto.GetRPCCodes(),

		incomingChan: make(chan *common.RPCPair, incomingBufferSize),

		logFilePath: logFilePath,

		replicaBatchSize: replicaBatchSize,
		replicaBatchTime: replicaBatchTime,

		debugOn:    debugOn,
		debugLevel: debugLevel,

		serverStarted:    false,
		consensusStarted: false,
		viewTimeout:      viewTimeout,
		logPrinted:       false,

		benchmarkMode:    benchmarkMode,
		state:            Init(benchmarkMode, name, keyLen, valLen),
		incomingRequests: make([]*proto.ClientBatch, 0),
		pipelineLength:   pipelineLength,
		lastProposedTime: time.Now(),
		finished:         false,
	}

	// initialize clientAddrList
	for i := 0; i < len(cfg.Clients); i++ {
		int32Name, _ := strconv.ParseInt(cfg.Clients[i].Name, 10, 32)
		rp.clientAddrList[int32(int32Name)] = cfg.Clients[i].Address
		rp.outgoingClientWriterMutexs[int32(int32Name)] = &sync.Mutex{}
	}

	// initialize replicaAddrList
	for i := 0; i < len(cfg.Peers); i++ {
		int32Name, _ := strconv.ParseInt(cfg.Peers[i].Name, 10, 32)
		rp.replicaAddrList[int32(int32Name)] = cfg.Peers[i].Address
		rp.outgoingReplicaWriterMutexs[int32(int32Name)] = &sync.Mutex{}
		rp.outgoingReplicaMessageChans[i] = make(chan *common.OutgoingRPC, outgoingBufferSize)
	}

	/*
		Register the rpcs
	*/
	rp.RegisterRPC(new(proto.ClientBatch), rp.messageCodes.ClientBatchRpc)
	rp.RegisterRPC(new(proto.Status), rp.messageCodes.StatusRPC)
	rp.RegisterRPC(new(proto.Pipelined_Sporades), rp.messageCodes.SporadesConsensus)

	rand.Seed(time.Now().UnixNano() + int64(rp.name))

	if rp.debugOn {
		rp.debug("Registered RPCs in the table", 0)
	}
	rp.consensus = InitAsyncConsensus(debugLevel, debugOn, rp.numReplicas)

	pid := os.Getpid()
	fmt.Printf("--Initialized replica %v with process id: %v \n", name, pid)

	return &rp
}

/*
	Fill the RPC table by assigning a unique id to each message type
*/

func (rp *Replica) RegisterRPC(msgObj proto.Serializable, code uint8) {
	rp.rpcTable[code] = &common.RPCPair{Code: code, Obj: msgObj}
}

/*
	given an id, return the node type: client/replica
*/

func (rp *Replica) getNodeType(id int32) string {
	if _, ok := rp.clientAddrList[id]; ok {
		return "client"
	}
	if _, ok := rp.replicaAddrList[id]; ok {
		return "replica"
	}
	panic("should not happen")
}

//debug printing

func (rp *Replica) debug(s string, i int) {
	if rp.debugOn && i >= rp.debugLevel {
		fmt.Print(s + "\n")
	}
}
