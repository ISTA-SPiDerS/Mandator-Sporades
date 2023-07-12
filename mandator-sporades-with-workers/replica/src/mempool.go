package src

import (
	"fmt"
	"strconv"
	"time"
	"with-workers/common"
	"with-workers/proto"
)

/*
	Defines the data structures specific to Mem Blocks and mini mem blocks
*/
type MemPool struct {
	miniMap              common.MiniMessageStore //saves the mini blocks
	blockMap             common.MessageStore     //saves the blocks
	incomingBuffer       []string                // saves confirmed mini blocks that need to be added to a future  mem block
	lastTimeBlockCreated time.Time               // last time a mem block was created
	lastSentBlock        []int                   // the last sent block to each replica
	lastSeenAck          []int                   // the last seen ack from each replica
	indexCounter         int                     // to create unique block ids
	mode                 int                     // 1 if all to all broadcast, 2 if selective broadcast with back pressure
	window               int                     // window for the number of outstanding blocks awaiting acks

	lastCompletedRounds []int //array of N elements (N= number of replica) that keeps track of the last block for which
	// at least n-f block-acks were collected, for each replica
	awaitingAcks bool //states whether this replica is waiting for block-acks
}

/*
	Initialize a new MemPool
*/

func InitMemPool(mode int, numReplicas int, debugLevel int, debugOn bool, window int) *MemPool {
	mmp := MemPool{
		miniMap:              common.MiniMessageStore{},
		blockMap:             common.MessageStore{},
		incomingBuffer:       make([]string, 0),
		lastTimeBlockCreated: time.Now(),
		lastSeenAck:          make([]int, numReplicas),
		lastSentBlock:        make([]int, numReplicas),
		indexCounter:         1,
		mode:                 mode,
		window:               window,
		lastCompletedRounds:  make([]int, numReplicas),
		awaitingAcks:         false,
	}

	for i := 0; i < numReplicas; i++ {
		mmp.lastSeenAck[i] = 0
		mmp.lastSentBlock[i] = 0
		mmp.lastCompletedRounds[i] = 0
	}

	mmp.miniMap.Init(debugLevel, debugOn)
	mmp.blockMap.Init(debugLevel, debugOn)
	return &mmp
}

/*
	Handler for MemPool Messages
*/
func (rp *Replica) handleMemPool(message *proto.MemPool) {

	if message.Type == 1 {
		//Mem-Pool-Mem-Block 1
		// save the mem block
		rp.memPool.blockMap.Add(message)
		// Set lastCompletedRounds[sender] to the parent blocks round
		parentId := message.ParentBlockId
		node, sequence := common.ExtractSequenceNumber(parentId)
		if node != rp.name {
			rp.memPool.lastCompletedRounds[rp.replicaArrayIndex[node]] = sequence
		}
		common.Debug("Last Completed Rounds is "+fmt.Sprintf("%v", rp.memPool.lastCompletedRounds), 0, rp.debugLevel, rp.debugOn)
		// send Mem-Pool-Mem-Block-Ack 2 to the sender
		memPoolAck := proto.MemPool{
			Sender:        rp.name,
			Receiver:      message.Sender,
			UniqueId:      message.UniqueId,
			Type:          2,
			Note:          message.Note,
			Minimemblocks: nil,
			RoundNumber:   message.RoundNumber,
			ParentBlockId: message.ParentBlockId,
			Creator:       message.Creator,
		}

		rpcPair := common.RPCPair{
			Code: rp.messageCodes.MemPoolRPC,
			Obj:  &memPoolAck,
		}

		rp.sendMessage(message.Sender, rpcPair)
		common.Debug("Sent Mem Pool Ack message with type 2 to "+strconv.Itoa(int(message.Sender)), 0, rp.debugLevel, rp.debugOn)

	} else if message.Type == 2 {
		//update rp.memPool.lastSeenAck
		//Mem-Pool-Mem-Block-Ack 2
		// if the ack is for the last block I sent
		//		save the ack in the pool
		// 		if n-f acks are received for my last created block and awaitingAcks is true
		//			set AwaitingAcks to false
		//			set lastCompletedRound[self]++
		_, sequence := common.ExtractSequenceNumber(message.UniqueId)
		if sequence > rp.memPool.lastSeenAck[rp.replicaArrayIndex[message.Sender]] {
			rp.memPool.lastSeenAck[rp.replicaArrayIndex[message.Sender]] = sequence
		}
		if message.UniqueId == strconv.Itoa(int(rp.name))+"."+strconv.Itoa(rp.memPool.indexCounter-1) && rp.memPool.awaitingAcks == true {
			rp.memPool.blockMap.AddAck(message.UniqueId, message.Sender)
			acks := rp.memPool.blockMap.GetAcks(message.UniqueId)
			if acks != nil && len(acks) == len(rp.replicaAddrList)/2+1 {
				rp.memPool.awaitingAcks = false
				rp.memPool.lastCompletedRounds[rp.replicaArrayIndex[rp.name]]++
				common.Debug("Received n-f acks for the block "+message.UniqueId, 0, rp.debugLevel, rp.debugOn)
				// for testing purposes of the mem pool, send a dummy reponse to the client
				//rp.sendDummyResponse(message.UniqueId)
			}

		}
	} else if message.Type == 3 {
		// Mem-Pool-Mem-Block-Request 3
		// if the request mem block exists, send a Mem-Pool-Mem-Block-Response 4 to the sender
		block, ok := rp.memPool.blockMap.Get(message.UniqueId)
		if ok {
			memPoolResponse := proto.MemPool{
				Sender:        rp.name,
				Receiver:      message.Sender,
				UniqueId:      block.UniqueId,
				Type:          4,
				Note:          block.Note,
				Minimemblocks: block.Minimemblocks,
				RoundNumber:   block.RoundNumber,
				ParentBlockId: block.ParentBlockId,
				Creator:       block.Creator,
			}

			rpcPair := common.RPCPair{
				Code: rp.messageCodes.MemPoolRPC,
				Obj:  &memPoolResponse,
			}

			rp.sendMessage(message.Sender, rpcPair)
			common.Debug("Sent Mem Pool response message with type 4 to "+strconv.Itoa(int(message.Sender)), 0, rp.debugLevel, rp.debugOn)
		}

	} else if message.Type == 4 {
		// Mem-Pool-Mem-Block-Response 4
		// save the block in the store
		rp.memPool.blockMap.Add(message)
		common.Debug("Saved a mem block as an explicit response from "+strconv.Itoa(int(message.Sender)), 0, rp.debugLevel, rp.debugOn)
	}

}

/*
	creates a new Mem block if the conditions for creating a new block is satisfied
		condition 1: incoming buffer is full || maximum time is passed
		condition 2: awaitingAcks is false
*/

func (rp *Replica) createNewMemBlock() {
	// if (the channel is full or if the batch time is passed) and isAwaiting is false
	if (len(rp.memPool.incomingBuffer) > rp.replicaBatchSize || (time.Now().Sub(rp.memPool.lastTimeBlockCreated).Microseconds() > int64(rp.replicaBatchTime) &&
		len(rp.memPool.incomingBuffer) > 0)) && rp.memPool.awaitingAcks == false {

		common.Debug("Creating a new mem block with  "+strconv.Itoa(len(rp.memPool.incomingBuffer))+" mini blocks", 0, rp.debugLevel, rp.debugOn)

		// 		create a new Mem block
		bParentId := strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(rp.memPool.indexCounter-1) // because we always increase the index counter upon creating a new block

		acks := rp.memPool.blockMap.GetAcks(bParentId)

		if acks == nil || len(acks) < len(rp.replicaAddrList)/2+1 {
			return
		}

		newMemBlock := proto.MemPool{
			Sender:        rp.name,
			Receiver:      0, // we don't assign a receiver for now
			UniqueId:      strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(rp.memPool.indexCounter),
			Type:          1,
			Note:          "",
			Minimemblocks: rp.convertToSingleMiniMemBlockArray(rp.memPool.incomingBuffer),
			RoundNumber:   int64(rp.memPool.indexCounter),
			ParentBlockId: bParentId,
			Creator:       rp.name,
		}

		// set awaitingAcks to true so that until this block gets n-f acks, no future mem blocks are created
		rp.memPool.awaitingAcks = true
		//	save the mem block in mem store
		rp.memPool.blockMap.Add(&newMemBlock)

		// create a replica array which contains the names of each replica

		replicas := make([]int32, len(rp.replicaAddrList))

		i := 0
		for name, _ := range rp.replicaAddrList {
			replicas[i] = name
			i++
		}

		//	if the mode is 1
		//		send the mem block to all replicas
		if rp.mode == 1 {
			rp.sendMemBlockToEveryone(&newMemBlock, replicas)
		}
		//	else if the mode is 2
		//		send the mem block to the best healthy replicas
		if rp.mode == 2 {
			rp.sendBlockToBestMajority(&newMemBlock, replicas)
		}

		//		increment the counter, reset the incoming buffer, and update time
		rp.memPool.indexCounter++
		rp.memPool.incomingBuffer = make([]string, 0)
		rp.memPool.lastTimeBlockCreated = time.Now()
	}
}

/*
	A helper function that converts between string array and MemPool_SingleMiniMemBlock
*/

func (rp *Replica) convertToSingleMiniMemBlockArray(buffer []string) []*proto.MemPool_SingleMiniMemBlock {
	returnArray := make([]*proto.MemPool_SingleMiniMemBlock, 0)

	for i := 0; i < len(buffer); i++ {
		returnArray = append(returnArray, &proto.MemPool_SingleMiniMemBlock{
			UniqueId: buffer[i],
			Creator:  rp.name,
		})
	}

	return returnArray
}

/*
	Send a mem block to each replica in replicas
*/

func (rp *Replica) sendMemBlockToEveryone(m *proto.MemPool, replicas []int32) {

	for i := 0; i < len(replicas); i++ {
		replica := replicas[i]

		memPool := proto.MemPool{
			Sender:        rp.name,
			Receiver:      replica,
			UniqueId:      m.UniqueId,
			Type:          m.Type,
			Note:          m.Note,
			Minimemblocks: m.Minimemblocks,
			RoundNumber:   m.RoundNumber,
			ParentBlockId: m.ParentBlockId,
			Creator:       m.Creator,
		}

		rpcPair := common.RPCPair{
			Code: rp.messageCodes.MemPoolRPC,
			Obj:  &memPool,
		}

		rp.sendMessage(replica, rpcPair)
		_, sequence := common.ExtractSequenceNumber(m.UniqueId)
		rp.memPool.lastSentBlock[rp.replicaArrayIndex[replica]] = sequence

		common.Debug("Sent Mem Pool message with type 1 to "+strconv.Itoa(int(replica)), 0, rp.debugLevel, rp.debugOn)
	}
}

/*
	Back pressure based broadcasting where the sender sends only to healthy replicas
	Healthy replicas are the ones who sent an acknowledgements for the previously sent blocks
	To balance the tradeoff between perfect back pressure and performance, we use a window
*/

func (rp *Replica) sendBlockToBestMajority(m *proto.MemPool, replicas []int32) {
	healthyReplicas := make([]int32, 0) // an array that has the names of healthy replicas
	healthyCount := 0
	threshold := -10

	for healthyCount < len(replicas)/2+1 {
		threshold += 10
		healthyReplicas = make([]int32, 0)
		healthyCount = 0
		for i := 0; i < len(replicas); i++ {
			if rp.memPool.lastSentBlock[rp.replicaArrayIndex[replicas[i]]]-rp.memPool.lastSeenAck[rp.replicaArrayIndex[replicas[i]]] < rp.memPool.window+threshold {
				healthyReplicas = append(healthyReplicas, replicas[i])
				healthyCount++
			}

		}
	}
	if len(replicas) != len(healthyReplicas) {
		common.Debug("Selected a healthy replica set when threshold is "+strconv.Itoa(int(threshold))+" and the replicas are "+fmt.Sprintf("%v", healthyReplicas), 0, rp.debugLevel, rp.debugOn)
	}
	rp.sendMemBlockToEveryone(m, healthyReplicas)
}

/*
	This is only for testing the mem pool
	Upon receiving n-f block acks, the replica sends dummy response to the client through the worker
*/

func (rp *Replica) sendDummyResponse(id string) {
	// get the mempool block
	memPoolBlock, ok1 := rp.memPool.blockMap.Get(id)
	if ok1 {
		miniBlocksIds := memPoolBlock.Minimemblocks
		// for each mini block id, get the mini block
		for i := 0; i < len(miniBlocksIds); i++ {
			miniBlockI, ok2 := rp.memPool.miniMap.Get(miniBlocksIds[i].UniqueId)
			if ok2 {
				minimemPoolClientResponse := proto.MemPoolMini{
					Sender:   rp.name,
					Receiver: miniBlockI.Creator,
					UniqueId: miniBlockI.UniqueId,
					Type:     8,
					Note:     miniBlockI.Note,
					Commands: miniBlockI.Commands,
					Creator:  miniBlockI.Creator,
				}
				rp.sendMiniMemPoolClientResponse(&minimemPoolClientResponse)
			} else {
				panic("The mini mem pool for id " + miniBlocksIds[i].UniqueId + " was not found")
			}
		}

	} else {
		panic("The mem pool for id " + id + " was not found")
	}

	// for each mini block client batch, get the client batch
	// for each client
}

/*
	This method is invoked when the replica needs mem block to commit the block
	Randomly selects a replica and sends a Mem-Pool-Mem-Block-Request 3
*/

func (rp *Replica) sendExternalMemBlockRequest(id string) {

	randomReplica := common.Get_Some_Node(rp.replicaArrayIndex)

	for randomReplica == rp.name {
		randomReplica = common.Get_Some_Node(rp.replicaArrayIndex)
	}

	externalMemBlockRequest := proto.MemPool{
		Sender:        rp.name,
		Receiver:      randomReplica,
		UniqueId:      id,
		Type:          3,
		Note:          "",
		Minimemblocks: nil,
		RoundNumber:   -1,
		ParentBlockId: "",
		Creator:       -1,
	}

	rpcPair := common.RPCPair{
		Code: rp.messageCodes.MemPoolRPC,
		Obj:  &externalMemBlockRequest,
	}

	rp.sendMessage(randomReplica, rpcPair)
	common.Debug("Sent Mem Pool Mem block request message with type 3 to "+strconv.Itoa(int(randomReplica)), 0, rp.debugLevel, rp.debugOn)
}
