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
 defines the Replica struct and the new method that is invoked when creating a new replica
*/

type Replica struct {
	name          int32  // unique replica identifier as defined in the configuration.yml
	listenAddress string // TCP address to which the replica listens to new incoming TCP connections

	numReplicas int

	clientAddrList             map[int32]string // map with the IP:port address of every client
	clientArrayIndex           map[int32]int    // map each client name to a unique index, this is required because the identifiers can be arbitary
	incomingClientReaders      []*bufio.Reader  // socket readers for each client
	outgoingClientWriters      []*bufio.Writer  // socket writer for each client
	outgoingClientWriterMutexs []sync.Mutex     // for mutual exclusion for each buffio.writer outgoingClientWriters

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
	instantiate a new replica instance, allocate the buffers
*/

func New(name int32, cfg *configuration.InstanceConfig, logFilePath string, replicaBatchSize int, replicaBatchTime int, debugOn bool, mode int, debugLevel int, viewTimeout int, window int, asyncBatchTime int, consAlgo string, benchmarkMode int, keyLen int, valLen int) *Replica {
	rp := Replica{
		name:          name,
		listenAddress: common.GetAddress(cfg.Peers, name),
		numReplicas:   len(cfg.Peers),

		clientAddrList:             make(map[int32]string),
		clientArrayIndex:           make(map[int32]int),
		incomingClientReaders:      make([]*bufio.Reader, len(cfg.Clients)),
		outgoingClientWriters:      make([]*bufio.Writer, len(cfg.Clients)),
		outgoingClientWriterMutexs: make([]sync.Mutex, len(cfg.Clients)),

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

	common.Debug("Created a new replica instance", 0, rp.debugLevel, rp.debugOn)

	// initialize clientAddrList
	for i := 0; i < len(cfg.Clients); i++ {
		int32Name, _ := strconv.ParseInt(cfg.Clients[i].Name, 10, 32)
		rp.clientAddrList[int32(int32Name)] = cfg.Clients[i].Address
		rp.clientArrayIndex[int32(int32Name)] = i
		rp.outgoingClientWriterMutexs[i] = sync.Mutex{}
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
	rp.RegisterRPC(new(proto.MemPool), rp.messageCodes.MemPoolRPC)
	rp.RegisterRPC(new(proto.AsyncConsensus), rp.messageCodes.AsyncConsensus)
	rp.RegisterRPC(new(proto.PaxosConsensus), rp.messageCodes.PaxosConsensus)

	common.Debug("Registered RPCs in the table", 0, rp.debugLevel, rp.debugOn)

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
	A helper function to check if I am still alive, and to print the channel lengths
*/

func (rp *Replica) livenessDebug() {
	go func() {
		for true {
			time.Sleep(2 * time.Second)
			common.Debug(fmt.Sprintf("\n \n %v is alive with incoming channel length of %v and outgoing channel length of %v\n \n",
				rp.name, len(rp.incomingChan), len(rp.outgoingMessageChan)), 0, rp.debugLevel, rp.debugOn)
		}
	}()

}
