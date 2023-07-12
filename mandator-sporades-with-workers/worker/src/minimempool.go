package src

import (
	"async-consensus/common"
	"async-consensus/proto"
	"math/rand"
	"strconv"
	"time"
)

/* this file contains all the worker related logic */

/*
	Defines the data structures specific to Mini Mem Blocks
*/

type MiniMemPool struct {
	miniMap                  common.MiniMessageStore //saves the mini blocks that are awaiting acks
	incomingBuffer           []*proto.ClientBatch    // saves client requests that need to be added to a future mini block
	lastTimeMiniBlockCreated time.Time
	lastSentMiniBlockIndex   []int // last sent miniblock index for each worker
	lastSeenAck              []int // the last seen ack from each worker
	indexCounter             int   // to create unique mini mem block ids
	mode                     int   // 1 if all to all broadcast, 2 if selective broadcast with back pressure
	window                   int   // window for the number of outstanding mini blocks awaiting acks
}

/*
	Initialize a new MiniMemPool
*/

func InitMiniMemPool(mode int, numWorkers int, debugLevel int, debugOn bool, window int) *MiniMemPool {
	mmp := MiniMemPool{
		miniMap:                  common.MiniMessageStore{},
		incomingBuffer:           make([]*proto.ClientBatch, 0),
		lastTimeMiniBlockCreated: time.Now(),
		lastSentMiniBlockIndex:   make([]int, numWorkers),
		lastSeenAck:              make([]int, numWorkers),
		indexCounter:             1,
		mode:                     mode,
		window:                   window,
	}

	for i := 0; i < numWorkers; i++ {
		mmp.lastSeenAck[i] = 0
		mmp.lastSentMiniBlockIndex[i] = 0
	}

	mmp.miniMap.Init(debugLevel, debugOn)
	return &mmp
}

/*
	Handler for the client batch messages
*/

func (wr *Worker) handleClientBatch(batch *proto.ClientBatch) {
	//Upon receiving a client batch, save it in an internal pending buffer
	wr.miniPool.incomingBuffer = append(wr.miniPool.incomingBuffer, batch)

	/*
		If the buffer size is greater than batch size, or
		batch time has passed, then create a new mem block, update the last sendTime
	*/

	if len(wr.miniPool.incomingBuffer) > wr.workerBatchSize ||
		(time.Now().Sub(wr.miniPool.lastTimeMiniBlockCreated).Microseconds() > int64(wr.workerBatchTime) &&
			len(wr.miniPool.incomingBuffer) > 0) {

		// create a new block

		newMiniBlock := proto.MemPoolMini{
			Sender:   wr.name,
			Receiver: 0, // we do not assign a receiver at this point
			UniqueId: strconv.Itoa(int(wr.name)) + "." + strconv.Itoa(wr.miniPool.indexCounter),
			Type:     1,
			Note:     "",
			Commands: wr.miniPool.convertToMiniBloclkClientBatchArray(wr.miniPool.incomingBuffer),
			Creator:  wr.name,
		}
		// increment the counter, reset the incoming buffer, and update time
		wr.miniPool.indexCounter++
		wr.miniPool.incomingBuffer = make([]*proto.ClientBatch, 0)
		wr.miniPool.lastTimeMiniBlockCreated = time.Now()

		// Save the mini block in the minimempool
		wr.miniPool.miniMap.Add(&newMiniBlock)
		wr.miniPool.miniMap.AddAck(newMiniBlock.UniqueId, wr.name)

		var workerAssignmentMinusSelf map[int32][]int32
		workerAssignmentMinusSelf = wr.removeSelfFromWorkerAssignment(wr.workerAssignment)

		// If the broadcast mode is 1, then send the mini block to a worker at each replica
		if wr.miniPool.mode == 1 {
			wr.sendMiniBlockToEveryone(&newMiniBlock, workerAssignmentMinusSelf)
		}
		// if the broadcast mode is 2, then create a healthy worker set, that includes all the workers
		// who have sent a miniblock ack for the last sent miniblock, but with a window
		// if the healthy worker set is less than a majority, fill it with workers who are most up to date
		// Send Mem-Pool-Mini-Mem-Block to the selected worker set

		if wr.miniPool.mode == 2 {
			wr.sendMiniBlockToBestMajority(&newMiniBlock, workerAssignmentMinusSelf)
		}
	}

}

/*
	A helper function to remove self from the worker assignment
*/

func (wr *Worker) removeSelfFromWorkerAssignment(assignment map[int32][]int32) map[int32][]int32 {
	returnMap := make(map[int32][]int32)
	for replica, workers := range assignment {
		if replica != wr.defaultReplicaName {
			returnMap[replica] = workers
		}
	}
	return returnMap
}

/*
	A helper function to convert client batch array in two different proto definitions
*/

func (p MiniMemPool) convertToMiniBloclkClientBatchArray(buffer []*proto.ClientBatch) []*proto.MemPoolMini_ClientBatch {
	returnArray := make([]*proto.MemPoolMini_ClientBatch, 0)
	for i := 0; i < len(buffer); i++ {
		returnArray = append(returnArray, &proto.MemPoolMini_ClientBatch{
			Commands: p.convertToMiniMemBlockClientbatchArray(buffer[i].Requests),
			Id:       buffer[i].UniqueId,
			Creator:  buffer[i].Sender,
		})
	}
	return returnArray
}

/*
	A helper function to convert client requests array in two different proto definitions
*/

func (p MiniMemPool) convertToMiniMemBlockClientbatchArray(requests []*proto.ClientBatch_SingleOperation) []*proto.MemPoolMini_SingleOperation {
	returnArray := make([]*proto.MemPoolMini_SingleOperation, 0)
	for i := 0; i < len(requests); i++ {
		returnArray = append(returnArray, &proto.MemPoolMini_SingleOperation{
			UniqueId: requests[i].Id,
			Command:  requests[i].Command,
		})
	}
	return returnArray
}

/*
	Send a mini block to 1 worker in each replica
*/

func (wr *Worker) sendMiniBlockToEveryone(m *proto.MemPoolMini, workerAssignment map[int32][]int32) {
	for _, workers := range workerAssignment {
		randWorkerIndex := rand.Intn(len(workers))
		randWorker := workers[randWorkerIndex]

		memPoolMini := proto.MemPoolMini{
			Sender:   wr.name,
			Receiver: randWorker,
			UniqueId: m.UniqueId,
			Type:     m.Type,
			Note:     m.Note,
			Commands: m.Commands,
			Creator:  m.Creator,
		}

		rpcPair := common.RPCPair{
			Code: wr.messageCodes.MemPoolMiniRPC,
			Obj:  &memPoolMini,
		}

		wr.sendMessage(randWorker, rpcPair)

		_, wr.miniPool.lastSentMiniBlockIndex[wr.workerArrayIndex[randWorker]] = common.ExtractSequenceNumber(m.UniqueId)

		common.Debug("Sent Mini Mem Pool message with type 1 to "+strconv.Itoa(int(randWorker)), 0, wr.debugLevel, wr.debugOn)
	}
}

/*
	Back pressure based broadcasting where the sender sends only to healthy workers
	Healthy workers are the ones who sent an acknoledgements for the previously sent mini blocks
	To balance the tradeoff between perfect back pressure and performance, we use a window
*/

func (wr *Worker) sendMiniBlockToBestMajority(m *proto.MemPoolMini, workerAssignment map[int32][]int32) {
	healthyWorkers := make(map[int32][]int32) // a map that has healthy workers
	healthyCount := 0
	threshold := -10

	for healthyCount < len(workerAssignment)/2 {
		threshold += 10
		healthyWorkers = make(map[int32][]int32)
		healthyCount = 0
		for replica, workers := range workerAssignment {
			for i := 0; i < len(workers); i++ {
				if wr.miniPool.lastSentMiniBlockIndex[wr.workerArrayIndex[workers[i]]]-wr.miniPool.lastSeenAck[wr.workerArrayIndex[workers[i]]] < wr.miniPool.window+threshold {
					healthyWorkers[replica] = []int32{workers[i]}
					healthyCount++
					break
				}
			}
		}
	}
	if len(healthyWorkers) != len(workerAssignment) {
		common.Debug("Selected a healthy worker set when threshold is "+strconv.Itoa(int(threshold))+" and the number of cohorts is "+strconv.Itoa(len(healthyWorkers)), 1, wr.debugLevel, wr.debugOn)
	}
	wr.sendMiniBlockToEveryone(m, healthyWorkers)
}

/*
	Handler for mempool mini messages
*/

func (wr *Worker) handleMemPoolMini(message *proto.MemPoolMini) {

	if message.Type == 1 {
		//     Mem-Pool-Mini-Mem-Block 1: when a remote worker sends a new block
		//     Save the mem block in the store
		//	   Send an ack to the sender, and send Mem-Pool-Mini-Mem-Block-Internal-Send 3 to default replica
		wr.miniPool.miniMap.Add(message)

		memPoolMiniAck := proto.MemPoolMini{
			Sender:   wr.name,
			Receiver: message.Sender,
			UniqueId: message.UniqueId,
			Type:     2,
			Note:     message.Note,
			Commands: nil,
			Creator:  message.Creator,
		}

		rpcPair := common.RPCPair{
			Code: wr.messageCodes.MemPoolMiniRPC,
			Obj:  &memPoolMiniAck,
		}

		wr.sendMessage(message.Sender, rpcPair)
		common.Debug("Sent Mini Mem Pool message ack  with type 2 to "+strconv.Itoa(int(message.Sender)), 0, wr.debugLevel, wr.debugOn)

		// send the miniblock to the default replica

		memPoolMiniInternalSend := proto.MemPoolMini{
			Sender:   wr.name,
			Receiver: wr.defaultReplicaName,
			UniqueId: message.UniqueId,
			Type:     3,
			Note:     message.Note,
			Commands: nil,
			Creator:  message.Creator,
		}

		rpcPair = common.RPCPair{
			Code: wr.messageCodes.MemPoolMiniRPC,
			Obj:  &memPoolMiniInternalSend,
		}

		wr.sendMessage(wr.defaultReplicaName, rpcPair)
		common.Debug("Sent Mini Mem Pool internal send  with type 3 to "+strconv.Itoa(int(wr.defaultReplicaName)), 0, wr.debugLevel, wr.debugOn)

	} else if message.Type == 2 {
		//	Mem-Pool-Mini-Mem-Block-Ack 2
		//	If the block exists in the MiniStore, then
		//		add the ack to the mini mem block
		//	   	if the number of acks is n-f, send a Mem-Pool-Mini-Mem-Block-Confirm 7 to the default replica
		//	Mark the last seen ack counter for the sender
		miniBlock, ok := wr.miniPool.miniMap.Get(message.UniqueId)
		if ok {
			wr.miniPool.miniMap.AddAck(message.UniqueId, message.Sender)
			acks := wr.miniPool.miniMap.GetAcks(message.UniqueId)
			if acks != nil {
				if len(acks) == len(wr.workerAssignment)/2+1 {

					memPoolMiniConfirm := proto.MemPoolMini{
						Sender:   wr.name,
						Receiver: wr.defaultReplicaName,
						UniqueId: miniBlock.UniqueId,
						Type:     7,
						Note:     miniBlock.Note,
						Commands: nil,
						Creator:  miniBlock.Creator,
					}

					rpcPair := common.RPCPair{
						Code: wr.messageCodes.MemPoolMiniRPC,
						Obj:  &memPoolMiniConfirm,
					}

					wr.sendMessage(wr.defaultReplicaName, rpcPair)
					common.Debug("Sent Mini Mem Pool confirm  with type 7 to "+strconv.Itoa(int(wr.defaultReplicaName)), 0, wr.debugLevel, wr.debugOn)
				}
			}
		}
		_, sequence := common.ExtractSequenceNumber(message.UniqueId)
		if sequence > wr.miniPool.lastSeenAck[wr.workerArrayIndex[message.Sender]] {
			wr.miniPool.lastSeenAck[wr.workerArrayIndex[message.Sender]] = sequence
		}

	} else if message.Type == 4 {
		//    Mem-Pool-Mini-Mem-Block-Internal-Send-Ack 4

	} else if message.Type == 8 {
		//  Mem-Pool-Mini-Mem-Block-Client-Response 8
		miniBlock, ok := wr.miniPool.miniMap.Get(message.UniqueId)
		if !ok {
			panic("should not happen")
		}
		message.Commands = miniBlock.Commands
		//  for each batch of requests, find the originating client, and send the ClientBatch
		for i := 0; i < len(message.Commands); i++ {
			// message.Commands[i] is a batch of client responses
			clientResponses := make([]*proto.ClientBatch_SingleOperation, 0)
			for j := 0; j < len(message.Commands[i].Commands); j++ {
				clientResponses = append(clientResponses, &proto.ClientBatch_SingleOperation{
					Id:      message.Commands[i].Commands[j].UniqueId,
					Command: message.Commands[i].Commands[j].Command,
				})
			}
			// send the clientResponses to the creator of the client batch
			clientBatchResponse := proto.ClientBatch{
				Sender:   wr.name,
				Receiver: message.Commands[i].Creator,
				UniqueId: message.Commands[i].Id,
				Type:     2,
				Note:     "",
				Requests: clientResponses,
			}

			rpcPair := common.RPCPair{
				Code: wr.messageCodes.ClientBatchRpc,
				Obj:  &clientBatchResponse,
			}

			wr.sendMessage(message.Commands[i].Creator, rpcPair)
			common.Debug("Sent client batch  with type 2 to "+strconv.Itoa(int(message.Commands[i].Creator)), 0, wr.debugLevel, wr.debugOn)
		}
	}

}
