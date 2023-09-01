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

	finished         bool // to finish consensus
	networkbatchTime int  // network level batchTime in "milli" seconds

	rejectedCount int // number of messages rejected because self is still not updated to a rank

	isAsynchronousSimulation bool
	asyncSimTimeout          int
	asynchronousReplicas     map[int][]int // for each time based epoch, the minority replicas that are attacked
	timeEpochSize            int           // how many ms for a given time epoch
}

const incomingBufferSize = 1000000 // the size of the buffer which receives all the incoming messages
const outgoingBufferSize = 1000000 // size of the buffer that collects messages to be written to the wire

/*
	instantiate a new replica instance, allocate the buffers
*/

func New(name int32, cfg *configuration.InstanceConfig, logFilePath string, replicaBatchSize int, replicaBatchTime int, debugOn bool, debugLevel int, viewTimeout int, benchmarkMode int, keyLen int, valLen int, pipelineLength int, networkbatchTime int, isAsyncSim bool, asyncSimTimeout int, timeEpochSize int) *Replica {
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

		benchmarkMode:            benchmarkMode,
		state:                    Init(benchmarkMode, name, keyLen, valLen),
		incomingRequests:         make([]*proto.ClientBatch, 0),
		pipelineLength:           pipelineLength,
		finished:                 false,
		networkbatchTime:         networkbatchTime,
		rejectedCount:            0,
		isAsynchronousSimulation: isAsyncSim,
		asyncSimTimeout:          asyncSimTimeout,
		asynchronousReplicas:     make(map[int][]int),
		timeEpochSize:            timeEpochSize,
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

	if rp.debugOn {
		rp.debug("Registered RPCs in the table", 0)
	}
	rp.consensus = InitAsyncConsensus(debugLevel, debugOn, rp.numReplicas)

	if rp.isAsynchronousSimulation {
		// initialize the attack replicas for each time epoch, we assume a total number of time of the run to be 10 minutes just for convenience, but this does not affect the correctness
		numEpochs := 10 * 60 * 1000 / rp.timeEpochSize
		s2 := rand.NewSource(39)
		r2 := rand.New(s2)

		for i := 0; i < numEpochs; i++ {
			rp.asynchronousReplicas[i] = []int{}
			for j := 0; j < rp.numReplicas/2; j++ {
				newReplica := r2.Intn(39)%rp.numReplicas + 1
				for rp.inArray(rp.asynchronousReplicas[i], newReplica) {
					newReplica = r2.Intn(39)%rp.numReplicas + 1
				}
				rp.asynchronousReplicas[i] = append(rp.asynchronousReplicas[i], newReplica)
			}
		}

		if rp.debugOn {
			rp.debug(fmt.Sprintf("set of attacked nodes %v ", rp.asynchronousReplicas), 0)
		}
	}

	pid := os.Getpid()
	fmt.Printf("--Initialized replica %v with process id: %v \n", name, pid)

	return &rp
}

/*
	checks if replica is in ints
*/

func (rp *Replica) inArray(ints []int, replica int) bool {
	for i := 0; i < len(ints); i++ {
		if ints[i] == replica {
			return true
		}
	}
	return false
}

/*
	checks if self is in the set of attacked nodes for this replica in this time epoch
*/

func (rp *Replica) amIAttacked(epoch int) bool {
	attackedNodes := rp.asynchronousReplicas[epoch]
	for i := 0; i < len(attackedNodes); i++ {
		if rp.name == int32(attackedNodes[i]) {
			return true
		}
	}
	return false
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
