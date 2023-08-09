package src

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"pipelined-sporades/common"
	"pipelined-sporades/proto"
	"strconv"
	"time"
)

/*
	SporadesConsensus stores the data structures for the consensus
*/

type SporadesConsensus struct {
	vCurr       int32                           // current view number
	rCurr       int32                           // current round number
	blockHigh   *proto.Pipelined_Sporades_Block // the block with the highest round number in which the node received a propose message
	blockCommit *proto.Pipelined_Sporades_Block // the last block to commit
	isAsync     bool                            // true if in fallback mode, and false if in sync mode

	bFall           map[string][]string                    // for storing the fallback block ids: the key is view.level, value is the array of fallback block ids
	consensusPool   *AsyncConsensusStore                   // an id of a consensus block is creator_name.v.r.type.level. type can be r (regular) or f (fallback), level can be 1,2 or -1 (for regular blocks)
	voteReplies     map[string][]*proto.Pipelined_Sporades // stores the vote messages received in the sync path. The key is v.r, the value is the array of received votes
	newViewMessages map[int32][]*proto.Pipelined_Sporades  // stores the new view messages received for each view. The key is v, the value is the array of received new view messages
	timeoutMessages map[int32][]*proto.Pipelined_Sporades  // stores the timeout messages for each view, the key is the view number, value is the array of received timeout messages
	viewTimer       *common.TimerWithCancel                // a timer to set the view timeouts in the replicas

	lastCommittedBlock *proto.Pipelined_Sporades_Block // the last committed block

	randomness []int // predefined random leader node names for each view

	sentLevel2Block  map[int32]bool // records if level 2 fallback block was sent for view v; key is v
	startTime        time.Time      // time when the consensus was started
	lastProposedTime time.Time      //time when last proposed

	orderedMessages [][]*proto.Pipelined_Sporades // pending Sporades  messages to be processed from each replica

	pipelinedRequests int

	sentFirstProposal map[int32]bool // did we send the first proposal upon starting a new view
}

// increment pipelinedRequests

func (rp *Replica) incrementPipelined() {
	rp.consensus.pipelinedRequests++
}

// decrement pipelinedRequests

func (rp *Replica) decrementPipelined() {
	if rp.consensus.pipelinedRequests > 0 {
		rp.consensus.pipelinedRequests--
	}
}

/*
	Init Async Consensus Data Structs
*/

func InitAsyncConsensus(debugLevel int, debugOn bool, numReplicas int) *SporadesConsensus {

	blockHigh := &proto.Pipelined_Sporades_Block{
		Id:       "genesis-block",
		V:        0,
		R:        0,
		ParentId: "",
		Parent:   nil,
		Commands: &proto.ReplicaBatch{
			UniqueId: "nil",
			Requests: make([]*proto.ClientBatch, 0),
			Sender:   -1,
		},
		Level: -1,
	}

	asyncConsensus := SporadesConsensus{
		vCurr:           0,
		rCurr:           0,
		blockHigh:       blockHigh,
		blockCommit:     blockHigh,
		isAsync:         false,
		bFall:           make(map[string][]string),
		consensusPool:   &AsyncConsensusStore{},
		voteReplies:     make(map[string][]*proto.Pipelined_Sporades),
		newViewMessages: make(map[int32][]*proto.Pipelined_Sporades),
		timeoutMessages: make(map[int32][]*proto.Pipelined_Sporades),

		lastCommittedBlock: blockHigh,
		randomness:         make([]int, 0),
		sentLevel2Block:    make(map[int32]bool),
		startTime:          time.Now(),
		orderedMessages:    make([][]*proto.Pipelined_Sporades, numReplicas),
		pipelinedRequests:  0,
		sentFirstProposal:  make(map[int32]bool),
	}

	// initialize the consensus pool
	asyncConsensus.consensusPool.Init(debugLevel, debugOn)
	// add genesis block to the AsyncConsensusStore
	asyncConsensus.consensusPool.Add(blockHigh)

	//initialize randomness and sentLevel2Block; assumes 1000000 max view changes, increase the 1000000 if needed
	s2 := rand.NewSource(42)
	r2 := rand.New(s2)
	for i := 0; i < 1000000; i++ {
		asyncConsensus.randomness = append(asyncConsensus.randomness, (r2.Intn(42))%numReplicas+1)
		asyncConsensus.sentLevel2Block[int32(i)] = false
		asyncConsensus.sentFirstProposal[int32(i)] = false
	}

	asyncConsensus.randomness[10] = 1 // initial sync leader is 1

	// initialize ordered Messages

	for i := 0; i < numReplicas; i++ {
		asyncConsensus.orderedMessages[i] = make([]*proto.Pipelined_Sporades, 0)
	}

	return &asyncConsensus
}

/*
	Returns the sync leader for the view; view leaders are pre-defined
*/

func (rp *Replica) getLeader(view int32) int32 {
	leaderIndex := rp.consensus.randomness[view+10]
	return int32(leaderIndex)
}

/*
	send the initial consensus new view message for the genesis block
*/

func (rp *Replica) sendGenesisConsensusNewView() {
	rp.consensus.startTime = time.Now()

	//send <new-view, vcur , rcur , blockhigh > to leader
	nextLeader := rp.getLeader(rp.consensus.vCurr)

	genesisBlock, ok := rp.consensus.consensusPool.Get("genesis-block")

	if !ok {
		panic("Genesis consensus block not found")
	}

	bootStrapNewView := proto.Pipelined_Sporades{
		Sender:      rp.name,
		Receiver:    nextLeader,
		UniqueId:    "",
		Type:        10,
		Note:        "",
		V:           rp.consensus.vCurr,
		R:           rp.consensus.rCurr,
		BlockHigh:   genesisBlock,
		BlockNew:    nil,
		BlockCommit: nil,
	}

	rpcPair := common.RPCPair{
		Code: rp.messageCodes.SporadesConsensus,
		Obj:  &bootStrapNewView,
	}

	rp.sendMessage(nextLeader, rpcPair)
	if rp.debugOn {
		rp.debug("Sent boot strap new view vote to "+strconv.Itoa(int(nextLeader)), 0)
	}
	// start the timeout
	rp.setViewTimer()
}

/*
	Sets a timer, which once timeout will broadcast a timeout message
*/

func (rp *Replica) setViewTimer() {

	rp.consensus.viewTimer = common.NewTimerWithCancel(time.Duration(rp.viewTimeout) * time.Microsecond)
	vCurr := rp.consensus.vCurr
	rCurr := rp.consensus.rCurr
	rp.consensus.viewTimer.SetTimeoutFunction(func() {
		// this function runs in a separate thread, hence we do not send timeout message in this function, instead send a timeout-internal signal
		internalTimeoutNotification := proto.Pipelined_Sporades{
			Sender:      rp.name,
			Receiver:    rp.name,
			UniqueId:    "",
			Type:        6,
			Note:        "",
			V:           vCurr,
			R:           rCurr,
			BlockHigh:   nil,
			BlockNew:    nil,
			BlockCommit: nil,
		}

		rp.incomingChan <- &common.RPCPair{
			Code: rp.messageCodes.SporadesConsensus,
			Obj:  &internalTimeoutNotification,
		}

		if rp.debugOn {
			rp.debug("Sent an internal timeout notification "+fmt.Sprintf("%v", internalTimeoutNotification), 0)
		}
	})
	rp.consensus.viewTimer.Start()
}

/*
	Handler for the sporades consensus messages
*/

func (rp *Replica) handleSporadesConsensus(messageNew *proto.Pipelined_Sporades) {

	if rp.logPrinted {
		return
	}

	if messageNew.Type == 7 {
		if rp.debugOn {
			rp.debug("Received a consensus-external-request message from "+strconv.Itoa(int(messageNew.Sender))+
				" for view "+strconv.Itoa(int(messageNew.V))+" for round "+strconv.Itoa(int(messageNew.R))+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 0)
		}
		rp.handleConsensusExternalRequest(messageNew)
		return
	} else if messageNew.Type == 8 {
		if rp.debugOn {
			rp.debug("Received a consensus-external-response message from "+strconv.Itoa(int(messageNew.Sender))+
				" for view "+strconv.Itoa(int(messageNew.V))+" for round "+strconv.Itoa(int(messageNew.R))+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 0)
		}
		rp.handleConsensusExternalResponseMessage(messageNew)
		return
	} else {

		// add this message to the tail of the orderedMessages buffer of the sender

		rp.consensus.orderedMessages[messageNew.Sender-1] = append(rp.consensus.orderedMessages[messageNew.Sender-1], messageNew)

		output := true

		for output && len(rp.consensus.orderedMessages[messageNew.Sender-1]) > 0 {

			message := rp.consensus.orderedMessages[messageNew.Sender-1][0]

			if message.Type == 1 {
				if rp.debugOn {
					rp.debug("Received a propose message from "+strconv.Itoa(int(message.Sender))+
						" for view "+strconv.Itoa(int(message.V))+" for round "+strconv.Itoa(int(message.R))+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 1)
				}
				output = rp.handleConsensusProposeSync(message)

			} else if message.Type == 2 {
				if rp.debugOn {
					rp.debug("Received a vote message from "+strconv.Itoa(int(message.Sender))+
						" for view "+strconv.Itoa(int(message.V))+" for round "+strconv.Itoa(int(message.R))+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 1)
				}
				output = rp.handleConsensusVoteSync(message)

			} else if message.Type == 3 {
				if rp.debugOn {
					rp.debug("Received a timeout message from "+strconv.Itoa(int(message.Sender))+
						" for view "+strconv.Itoa(int(message.V))+" for round "+strconv.Itoa(int(message.R))+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 1)
				}
				output = rp.handleConsensusTimeout(message)

			} else if message.Type == 4 {
				if rp.debugOn {
					rp.debug("Received a propose-async message from "+strconv.Itoa(int(message.Sender))+
						" for view "+strconv.Itoa(int(message.V))+" for round "+strconv.Itoa(int(message.R))+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 1)
				}
				output = rp.handleConsensusProposeAsync(message)

			} else if message.Type == 5 {
				if rp.debugOn {
					rp.debug("Received a vote-async message from "+strconv.Itoa(int(message.Sender))+
						" for view "+strconv.Itoa(int(message.V))+" for round "+strconv.Itoa(int(message.R))+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 1)
				}
				output = rp.handleConsensusAsyncVote(message)

			} else if message.Type == 6 {
				if rp.debugOn {
					rp.debug("Received a timeout-internal message from "+strconv.Itoa(int(message.Sender))+
						" for view "+strconv.Itoa(int(message.V))+" for round "+strconv.Itoa(int(message.R))+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 1)
				}
				output = rp.handleConsensusInternalTimeout(message)
			} else if message.Type == 9 {
				if rp.debugOn {
					rp.debug("Received a async fallback-complete message from "+strconv.Itoa(int(message.Sender))+
						" for view "+strconv.Itoa(int(message.V))+" for round "+strconv.Itoa(int(message.R))+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 1)
				}
				output = rp.handleConsensusFallbackCompleteMessage(message)
			} else if message.Type == 10 {
				if rp.debugOn {
					rp.debug("Received a new view message from "+strconv.Itoa(int(message.Sender))+
						" for view "+strconv.Itoa(int(message.V))+" for round "+strconv.Itoa(int(message.R))+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 1)
				}
				output = rp.handleConsensusNewViewMessage(message)
			}

			if output {
				// the message was processed correctly
				rp.consensus.orderedMessages[messageNew.Sender-1][0] = nil
				if len(rp.consensus.orderedMessages[messageNew.Sender-1]) > 1 {
					rp.consensus.orderedMessages[messageNew.Sender-1] = rp.consensus.orderedMessages[messageNew.Sender-1][1:]
				} else {
					rp.consensus.orderedMessages[messageNew.Sender-1] = make([]*proto.Pipelined_Sporades, 0)
				}
			} else {
				if rp.debugOn {
					rp.debug("message "+fmt.Sprintf("%v", message)+" not processed, saved for future", 1)
					rp.rejectedCount++
					if rp.rejectedCount >= 10000 {
						rp.debug("somewhere wrong", 0)
					}
				}
			}
		}
	}
}

/*
	Util function to compare the ranks; checks if v1;r1 rank higher than v2,r2
*/

func (rp *Replica) hasGreaterRank(v1 int32, r1 int32, v2 int32, r2 int32) bool {
	if v1 > v2 {
		return true
	}
	if v1 == v2 && r1 > r2 {
		return true
	}
	return false
}

/*
	Util function to compare the ranks; checks if v1,r1 >= v2,r2
*/

func (rp *Replica) hasGreaterThanOrEqualRank(v1 int32, r1 int32, v2 int32, r2 int32) bool {
	if v1 > v2 {
		return true
	}
	if v1 == v2 && r1 > r2 {
		return true
	}
	if v1 == v2 && r1 == r2 {
		return true
	}
	return false
}

/*
	Util function to extract the highest blockHigh from the received set of vote/timeout/new view messages
*/

func (rp *Replica) extractHighestRankedBlockHigh(messages []*proto.Pipelined_Sporades) *proto.Pipelined_Sporades_Block {
	highBlock := messages[0].BlockHigh
	for i := 0; i < len(messages); i++ {
		if rp.hasGreaterRank(messages[i].BlockHigh.V, messages[i].BlockHigh.R, highBlock.V, highBlock.R) {
			highBlock = messages[i].BlockHigh
		}
	}
	return highBlock
}

/*
	Util function to convert int[] to int32[]
*/

func (rp *Replica) convertToInt32Array(elements []int) []int32 {
	returnArray := make([]int32, 0)
	for i := 0; i < len(elements); i++ {
		returnArray = append(returnArray, int32(elements[i]))
	}
	return returnArray
}

/*
	This method is invoked when the replica needs a consensus block
	Randomly selects a replica and sends a 7-consensus-external-request
*/

func (rp *Replica) sendExternalConsensusRequest(id string) {

	randomReplica := int32(rand.Intn(rp.numReplicas) + 1)

	for randomReplica == rp.name {
		randomReplica = int32(rand.Intn(rp.numReplicas) + 1)
	}

	externalConsensusRequest := proto.Pipelined_Sporades{
		Sender:      rp.name,
		Receiver:    randomReplica,
		UniqueId:    id,
		Type:        7,
		Note:        "",
		V:           -1,
		R:           -1,
		BlockHigh:   nil,
		BlockNew:    nil,
		BlockCommit: nil,
	}

	rpcPair := common.RPCPair{
		Code: rp.messageCodes.SporadesConsensus,
		Obj:  &externalConsensusRequest,
	}

	rp.sendMessage(randomReplica, rpcPair)
	if rp.debugOn {
		rp.debug("Sent external consensus request message with type 7 to "+strconv.Itoa(int(randomReplica)), 1)
	}
}

/*
	Handler for external consensus request messages
*/

func (rp *Replica) handleConsensusExternalRequest(message *proto.Pipelined_Sporades) {
	// if the consensus block with id message.unique_id exists in the consensus store
	block, ok := rp.consensus.consensusPool.Get(message.UniqueId)
	if ok {
		// send an external consensus response 8 message to the sender
		externalConsensusResponse := proto.Pipelined_Sporades{
			Sender:      rp.name,
			Receiver:    message.Sender,
			UniqueId:    message.UniqueId,
			Type:        8,
			Note:        "",
			V:           -1,
			R:           -1,
			BlockHigh:   nil,
			BlockNew:    block,
			BlockCommit: nil,
		}

		rpcPair := common.RPCPair{
			Code: rp.messageCodes.SporadesConsensus,
			Obj:  &externalConsensusResponse,
		}

		rp.sendMessage(message.Sender, rpcPair)
		if rp.debugOn {
			rp.debug("Sent external consensus response message with type 8 to "+strconv.Itoa(int(message.Sender)), 1)
		}
	}
}

/*
	Handler for external consensus response messages
*/

func (rp *Replica) handleConsensusExternalResponseMessage(message *proto.Pipelined_Sporades) {
	rp.consensus.consensusPool.Add(message.BlockNew)
	if rp.debugOn {
		rp.debug("Added a consensus block "+message.BlockNew.Id+" from an external response", 1)
	}
	rp.updateSMR()
}

/*
	Make chain upto n: appends parent blocks upto n blocks in the history
*/

func (rp *Replica) makeNChain(blockOri *proto.Pipelined_Sporades_Block, n int) *proto.Pipelined_Sporades_Block {

	block, err := CloneMyStruct(blockOri)
	if err != nil {
		panic(err.Error())
	}
	head := block

	for n >= 0 {
		n--
		parent_id := block.ParentId
		if parent_id == "genesis-block" {
			return head
		}
		b, ok := rp.consensus.consensusPool.Get(parent_id)
		if !ok {
			return head
		}
		b, err = CloneMyStruct(b)
		if err != nil {
			panic(err.Error())
		}

		block.Parent = b
		block = block.Parent
	}
	return head

}

/*
	Deep copy a consensus block
*/

func CloneMyStructJson(orig *proto.Pipelined_Sporades_Block) (*proto.Pipelined_Sporades_Block, error) {
	origJSON, err := json.Marshal(orig)
	if err != nil {
		panic(err)
	}
	clone := proto.Pipelined_Sporades_Block{}
	if err = json.Unmarshal(origJSON, &clone); err != nil {
		return nil, err
	}
	return &clone, nil
}

/*
	manual copy a consensus block
*/

func CloneMyStruct(orig *proto.Pipelined_Sporades_Block) (*proto.Pipelined_Sporades_Block, error) {
	rtnBlc := &proto.Pipelined_Sporades_Block{
		Id:       orig.Id,
		V:        orig.V,
		R:        orig.R,
		ParentId: orig.ParentId,
		Parent:   nil,
		Commands: DuplicateCommands(orig.Commands),
		Level:    orig.Level,
	}

	return rtnBlc, nil
}

// copies the replica batch

func DuplicateCommands(commands *proto.ReplicaBatch) *proto.ReplicaBatch {
	returnBatch := &proto.ReplicaBatch{
		UniqueId: commands.UniqueId,
		Requests: DuplicateRequests(commands.Requests),
		Sender:   commands.Sender,
	}
	return returnBatch
}

// duplicate an array of client batches

func DuplicateRequests(requests []*proto.ClientBatch) []*proto.ClientBatch {
	returnArray := make([]*proto.ClientBatch, len(requests))
	for i := 0; i < len(requests); i++ {
		returnArray[i] = &proto.ClientBatch{
			UniqueId: requests[i].UniqueId,
			Requests: DuplicateClientRequests(requests[i].Requests),
			Sender:   requests[i].Sender,
		}
	}
	return returnArray
}

// duplicate client requests array

func DuplicateClientRequests(requests []*proto.SingleOperation) []*proto.SingleOperation {
	returnArray := make([]*proto.SingleOperation, len(requests))
	for i := 0; i < len(requests); i++ {
		returnArray[i] = &proto.SingleOperation{
			Command: requests[i].Command,
		}
	}
	return returnArray
}

/*
	Printing the replicated log for testing purposes
*/

func (rp *Replica) printLogConsensus() {
	f, err := os.Create(rp.logFilePath + strconv.Itoa(int(rp.name)) + "-consensus.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	head := rp.consensus.lastCommittedBlock // the last block to commit
	if head == nil {
		return
	}
	genesisBlock, ok := rp.consensus.consensusPool.Get("genesis-block")
	if !ok {
		panic("Genesis Block not found when printing the logs")
	}

	toPrint := make([]*proto.Pipelined_Sporades_Block, 0)

	for head.Id != genesisBlock.Id {

		toPrint = append([]*proto.Pipelined_Sporades_Block{head}, toPrint...)
		if rp.debugOn {
			rp.debug(fmt.Sprintf("printing: block %v\n", head), 22)
		}

		parent_id := head.ParentId
		if parent_id == "genesis-block" {
			break
		}

		parent, ok := rp.consensus.consensusPool.Get(parent_id)
		if !ok {
			panic("parent not found")
		}
		head = parent

	}
	fmt.Printf("Number of committed blocks: %v\n", len(toPrint))
	count := 0

	for i := 0; i < len(toPrint); i++ {
		for clientBatchIndex := 0; clientBatchIndex < len(toPrint[i].Commands.Requests); clientBatchIndex++ {
			clientBatch := toPrint[i].Commands.Requests[clientBatchIndex]
			clientBatchID := clientBatch.UniqueId
			clientBatchCommands := clientBatch.Requests
			for clientRequestIndex := 0; clientRequestIndex < len(clientBatchCommands); clientRequestIndex++ {
				clientRequest := clientBatchCommands[clientRequestIndex].Command
				_, _ = f.WriteString(toPrint[i].Id + "-" + clientBatchID + "-" + strconv.Itoa(clientRequestIndex) + ":" + clientRequest + "\n")
				count++
			}
		}
	}

	fmt.Print("Number of committed client requests: " + strconv.Itoa(count))

}
