package src

import (
	"async-consensus/benchmark"
	"async-consensus/common"
	"async-consensus/configuration"
	"async-consensus/proto"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

/*
	This file defines the Replica struct and the new method that is invoked when creating a new replica
*/

type Replica struct {
	name          int32  // unique replica identifier as defined in the configuration.yml
	listenAddress string // TCP address to which the replica listens to new incoming TCP connections

	numReplicas int

	workerAssignment map[int32][]int32 // maps replica to its workers

	workerAddrList             map[int32]string // map with the IP:port address of every designated worker node
	workerArrayIndex           map[int32]int    // map each worker name to a unique index, this is required because the identifiers can be arbitary
	incomingWorkerReaders      []*bufio.Reader  // socket readers for each worker
	outgoingWorkerWriters      []*bufio.Writer  // socket writer for each worker
	outgoingWorkerWriterMutexs []sync.Mutex     // for mutual exclusion for each buffio.writer outgoingWorkerWriters

	replicaAddrList             map[int32]string // map with the IP:port address of every replica node
	replicaArrayIndex           map[int32]int    // map each replica name to a unique index, this is required because the identifiers can be arbitary
	incomingReplicaReaders      []*bufio.Reader  // socket readers for each replica
	outgoingReplicaWriters      []*bufio.Writer  // socket writer for each replica
	outgoingReplicaWriterMutexs []sync.Mutex     // for mutual exclusion for each buffio.writer outgoingReplicaWriters

	rpcTable     map[uint8]*common.RPCPair // map each RPC type (message type) to its unique number
	messageCodes proto.MessageCode

	incomingChan chan *common.RPCPair // used to collect all the incoming messages

	logFilePath string // the path to write the log, used for sanity checks

	replicaBatchSize int // maximum replica side batch size
	replicaBatchTime int // maximum replica side batch time in micro seconds

	outgoingMessageChan chan *common.OutgoingRPC // buffer for messages that are written to the wire

	debugOn    bool // if turned on, the debug messages will be printed on the console
	debugLevel int  // current debug level

	serverStarted bool // to bootstrap

	mode           int             // 1 for all to all broadcast and 2 for selective broadcast
	memPool        *MemPool        // mempool
	asyncConsensus *AsyncConsensus // async consensus data structs
	paxosConsensus *Paxos          // Paxos consensus data structs

	consensusStarted bool // to send the initial vote messages for genesis block
	viewTimeout      int  // view change timeout in micro seconds

	logPrinted bool // to check if log was printed before

	asyncBatchTime int // delay for consensus messages

	consAlgo string // async/paxos

	benchmarkMode int                  // 0 for resident K/V store, 1 for redis
	state         *benchmark.Benchmark // k/v store
}

const numOutgoingThreads = 100       // number of wire writers: since the I/O writing is expensive we delegate that task to a thread pool and separate from the critical path todo this number should be tuned
const incomingBufferSize = 100000000 // the size of the buffer which receives all the incoming messages
const outgoingBufferSize = 100000000 // size of the buffer that collects messages to be written to the wire

/*
	Instantiate a new replica instance, allocate the buffers
*/

func New(name int32, cfg *configuration.InstanceConfig, workerAssignment map[int32][]int32, logFilePath string, replicaBatchSize int, replicaBatchTime int, debugOn bool, mode int, debugLevel int, viewTimeout int, window int, asyncBatchTime int, consAlgo string, benchmarkMode int, keyLen int, valLen int) *Replica {
	rp := Replica{
		name:          name,
		listenAddress: common.GetAddress(cfg.Peers, name),
		numReplicas:   len(cfg.Peers),

		workerAssignment: workerAssignment,

		workerAddrList:             make(map[int32]string),
		workerArrayIndex:           make(map[int32]int),
		incomingWorkerReaders:      make([]*bufio.Reader, len(workerAssignment[name])),
		outgoingWorkerWriters:      make([]*bufio.Writer, len(workerAssignment[name])),
		outgoingWorkerWriterMutexs: make([]sync.Mutex, len(workerAssignment[name])),

		replicaAddrList:             make(map[int32]string),
		replicaArrayIndex:           make(map[int32]int),
		incomingReplicaReaders:      make([]*bufio.Reader, len(cfg.Peers)),
		outgoingReplicaWriters:      make([]*bufio.Writer, len(cfg.Peers)),
		outgoingReplicaWriterMutexs: make([]sync.Mutex, len(cfg.Peers)),

		rpcTable:     make(map[uint8]*common.RPCPair),
		messageCodes: proto.GetRPCCodes(),

		incomingChan: make(chan *common.RPCPair, incomingBufferSize),

		logFilePath: logFilePath,

		replicaBatchSize: replicaBatchSize,
		replicaBatchTime: replicaBatchTime,

		outgoingMessageChan: make(chan *common.OutgoingRPC, outgoingBufferSize),
		debugOn:             debugOn,
		debugLevel:          debugLevel,
		serverStarted:       false,
		mode:                mode,
		consensusStarted:    false,
		viewTimeout:         viewTimeout,
		logPrinted:          false,
		asyncBatchTime:      asyncBatchTime,
		consAlgo:            consAlgo,
		benchmarkMode:       benchmarkMode,
		state:               benchmark.Init(benchmarkMode, name, keyLen, valLen),
	}

	// init mem pool and consensus pools
	rp.memPool = InitMemPool(mode, len(cfg.Peers), rp.debugLevel, rp.debugOn, window)

	// add genesis mem block with round 0 to mem pool
	rp.memPool.blockMap.Add(&proto.MemPool{
		Sender:        name,
		Receiver:      0,
		UniqueId:      strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(0),
		Type:          0,
		Note:          "",
		Minimemblocks: make([]*proto.MemPool_SingleMiniMemBlock, 0),
		RoundNumber:   0,
		ParentBlockId: "",
		Creator:       name,
	})

	rp.asyncConsensus = InitAsyncConsensus(debugLevel, debugOn, len(cfg.Peers))
	rp.paxosConsensus = InitPaxosConsensus(len(cfg.Peers))

	// add n acks to the genesis mem block

	for i := 0; i < len(cfg.Peers); i++ {
		peer, _ := strconv.Atoi(cfg.Peers[i].Name)
		rp.memPool.blockMap.AddAck(strconv.Itoa(int(rp.name))+"."+strconv.Itoa(0), int32(peer))
	}

	common.Debug("Created a new replica instance", 0, rp.debugLevel, rp.debugOn)

	// initialize workerAddrList
	for i := 0; i < len(workerAssignment[name]); i++ {
		int32Name := workerAssignment[name][i]
		rp.workerAddrList[int32Name] = common.GetAddress(cfg.Workers, int32Name)
		rp.workerArrayIndex[int32Name] = i
		rp.outgoingWorkerWriterMutexs[i] = sync.Mutex{}
	}

	// initialize replicaAddrList
	for i := 0; i < len(cfg.Peers); i++ {
		int32Name, _ := strconv.ParseInt(cfg.Peers[i].Name, 10, 32)
		rp.replicaAddrList[int32(int32Name)] = cfg.Peers[i].Address
		rp.replicaArrayIndex[int32(int32Name)] = i
		rp.outgoingReplicaWriterMutexs[i] = sync.Mutex{}
	}

	/*
		Register the rpcs
	*/
	rp.RegisterRPC(new(proto.ClientBatch), rp.messageCodes.ClientBatchRpc)
	rp.RegisterRPC(new(proto.Status), rp.messageCodes.StatusRPC)
	rp.RegisterRPC(new(proto.MemPoolMini), rp.messageCodes.MemPoolMiniRPC)
	rp.RegisterRPC(new(proto.MemPool), rp.messageCodes.MemPoolRPC)
	rp.RegisterRPC(new(proto.AsyncConsensus), rp.messageCodes.AsyncConsensus)
	rp.RegisterRPC(new(proto.PaxosConsensus), rp.messageCodes.PaxosConsensus)

	common.Debug("Registered RPCs in the table", 0, rp.debugLevel, rp.debugOn)

	pid := os.Getpid()
	fmt.Printf("initialized %v replica %v with process id: %v \n", consAlgo, name, pid)

	return &rp
}

/*
	Fill the RPC table by assigning a unique id to each message type
*/

func (rp *Replica) RegisterRPC(msgObj proto.Serializable, code uint8) {
	rp.rpcTable[code] = &common.RPCPair{Code: code, Obj: msgObj}
}

/*
	Given an id, return the node type: worker/replica
*/

func (rp *Replica) getNodeType(id int32) string {
	if _, ok := rp.workerAddrList[id]; ok {
		return "worker"
	}
	if _, ok := rp.replicaAddrList[id]; ok {
		return "replica"
	}
	return ""
}

/*
	A helper function to check if I am still alive, and to print the channel lengths
*/

func (rp *Replica) livenessDebug() {
	go func() {
		for true {
			time.Sleep(2 * time.Second)
			common.Debug(fmt.Sprintf("\n \n I, %v hereby claim that I am alive with incoming channel length of %v and outgoing channel length of %v\n \n",
				rp.name, len(rp.incomingChan), len(rp.outgoingMessageChan)), 0, rp.debugLevel, rp.debugOn)
		}
	}()

}
