package src

import (
	"bufio"
	"fmt"
	"mandator-sporades/common"
	"mandator-sporades/configuration"
	"mandator-sporades/proto"
	"os"
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

	clientAddrList             map[int32]string        // map with the IP:port address of every client
	incomingClientReaders      map[int32]*bufio.Reader // socket readers for each client
	outgoingClientWriters      map[int32]*bufio.Writer // socket writer for each client
	outgoingClientWriterMutexs map[int32]*sync.Mutex   // for mutual exclusion for each buffio.writer outgoingClientWriters

	replicaAddrList             map[int32]string        // map with the IP:port address of every replica node
	incomingReplicaReaders      map[int32]*bufio.Reader // socket readers for each replica
	outgoingReplicaWriters      map[int32]*bufio.Writer // socket writer for each replica
	outgoingReplicaWriterMutexs map[int32]*sync.Mutex   // for mutual exclusion for each buffio.writer outgoingReplicaWriters

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

	networkBatchTime int // delay for consensus propose messages ms

	consAlgo string // async/paxos

	benchmarkMode    int        // 0 for resident K/V store, 1 for redis
	state            *Benchmark // k/v store
	isAsync          bool
	asynchronousTime int // ms for simulating asynchronous timeouts
}

const numOutgoingThreads = 100       // number of wire writers: since the I/O writing is expensive we delegate that task to a thread pool and separate from the critical path
const incomingBufferSize = 100000000 // the size of the buffer which receives all the incoming messages
const outgoingBufferSize = 100000000 // size of the buffer that collects messages to be written to the wire

/*
	instantiate a new replica instance, allocate the buffers
*/

func New(name int32, cfg *configuration.InstanceConfig, logFilePath string, replicaBatchSize int, replicaBatchTime int, debugOn bool, mode int, debugLevel int, viewTimeout int, window int, networkBatchTime int, consAlgo string, benchmarkMode int, keyLen int, valLen int, isAsync bool, asyncSimTime int) *Replica {
	rp := Replica{
		name:          name,
		listenAddress: common.GetAddress(cfg.Peers, name),
		numReplicas:   len(cfg.Peers),

		clientAddrList:             make(map[int32]string),
		incomingClientReaders:      make(map[int32]*bufio.Reader),
		outgoingClientWriters:      make(map[int32]*bufio.Writer),
		outgoingClientWriterMutexs: make(map[int32]*sync.Mutex),

		replicaAddrList:             make(map[int32]string),
		incomingReplicaReaders:      make(map[int32]*bufio.Reader),
		outgoingReplicaWriters:      make(map[int32]*bufio.Writer),
		outgoingReplicaWriterMutexs: make(map[int32]*sync.Mutex),

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
		networkBatchTime:    networkBatchTime,
		consAlgo:            consAlgo,
		benchmarkMode:       benchmarkMode,
		state:               Init(benchmarkMode, name, keyLen, valLen),
		isAsync:             isAsync,
		asynchronousTime:    asyncSimTime,
	}

	// init mem pool
	rp.memPool = InitMemPool(mode, len(cfg.Peers), rp.debugLevel, rp.debugOn, window)
	rp.asyncConsensus = InitAsyncConsensus(debugLevel, debugOn, len(cfg.Peers))
	rp.paxosConsensus = InitPaxosConsensus(len(cfg.Peers))

	// add genesis mem block with round 0 to mem pool
	rp.memPool.blockMap.Add(&proto.MemPool{
		Sender:        name,
		UniqueId:      strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(0),
		Type:          0,
		Note:          "",
		ClientBatches: make([]*proto.ClientBatch, 0),
		RoundNumber:   0,
		ParentBlockId: "",
		Creator:       name,
	})

	// add n acks to the genesis mem block

	for i := 0; i < len(cfg.Peers); i++ {
		peer, _ := strconv.Atoi(cfg.Peers[i].Name)
		rp.memPool.blockMap.AddAck(strconv.Itoa(int(rp.name))+"."+strconv.Itoa(0), int32(peer))
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
	}

	/*
		Register the rpcs
	*/
	rp.RegisterRPC(new(proto.ClientBatch), rp.messageCodes.ClientBatchRpc)
	rp.RegisterRPC(new(proto.Status), rp.messageCodes.StatusRPC)
	rp.RegisterRPC(new(proto.MemPool), rp.messageCodes.MemPoolRPC)
	rp.RegisterRPC(new(proto.AsyncConsensus), rp.messageCodes.AsyncConsensus)
	rp.RegisterRPC(new(proto.PaxosConsensus), rp.messageCodes.PaxosConsensus)

	pid := os.Getpid()
	fmt.Printf("--Initialized %v replica %v with process id: %v \n", consAlgo, name, pid)

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
	return ""
}

/*
	debug print
*/

func (rp *Replica) debug(message string, level int) {
	if rp.debugOn && level >= rp.debugLevel {
		fmt.Print(message)
	}
}
