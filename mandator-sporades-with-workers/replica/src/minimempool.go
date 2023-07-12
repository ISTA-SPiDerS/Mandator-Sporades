package src

import (
	"strconv"
	"with-workers/common"
	"with-workers/proto"
)

/*
	Handler for mempoolmini messages
*/
func (rp *Replica) handleMemPoolMini(message *proto.MemPoolMini) {
	if message.Type == 3 {

		// Mem-Pool-Mini-Mem-Block-Internal-Send 3
		// the worker sends an internal send to its default replica upon receiving a new minimemblock from one of other workers
		// Upon receiving a Mem-Pool-Mini-Mem-Block-Internal-Send 3, the replica should save it in its mini mem pool
		rp.memPool.miniMap.Add(message)
		// and then send an Mem-Pool-Mini-Mem-Block-Internal-Send-Ack 4 to the sender

		minimempoolinternalack := proto.MemPoolMini{
			Sender:   rp.name,
			Receiver: message.Sender,
			UniqueId: message.UniqueId,
			Type:     4,
			Note:     message.Note,
			Commands: nil,
			Creator:  message.Creator,
		}

		rpcPair := common.RPCPair{
			Code: rp.messageCodes.MemPoolMiniRPC,
			Obj:  &minimempoolinternalack,
		}

		rp.sendMessage(message.Sender, rpcPair)
		common.Debug("Sent Mini Mem Pool internal send ack message with type 4 to "+strconv.Itoa(int(message.Sender)), 0, rp.debugLevel, rp.debugOn)
	}
	if message.Type == 5 {

		//Mem-Pool-Mini-Mem-Block-External-Request 5
		// an external replica is asking for a mini mem block.
		// if the mini mem block is in the mini store, reply with a Mem-Pool-Mini-Mem-Block-External-Response 6
		miniBlock, ok := rp.memPool.miniMap.Get(message.UniqueId)
		if ok {
			minimempoolexternalResponse := proto.MemPoolMini{
				Sender:   rp.name,
				Receiver: message.Sender,
				UniqueId: miniBlock.UniqueId,
				Type:     6,
				Note:     miniBlock.Note,
				Commands: miniBlock.Commands,
				Creator:  miniBlock.Creator,
			}

			rpcPair := common.RPCPair{
				Code: rp.messageCodes.MemPoolMiniRPC,
				Obj:  &minimempoolexternalResponse,
			}

			rp.sendMessage(message.Sender, rpcPair)
			common.Debug("Sent Mini Mem Pool external response message with type 6 to "+strconv.Itoa(int(message.Sender)), 1, rp.debugLevel, rp.debugOn)
		} else {
			common.Debug("Received an external request for mini block id  "+message.UniqueId+" from "+strconv.Itoa(int(message.Sender))+" but i don't have it", 1, rp.debugLevel, rp.debugOn)
		}
	}
	if message.Type == 6 {
		//Mem-Pool-Mini-Mem-Block-External-Response 6
		// this a resposne to one of my previous requests, save it in the mini mem pool
		rp.memPool.miniMap.Add(message)
		common.Debug("Saved the mini block in mini store after receiving from external replica "+strconv.Itoa(int(message.Sender)), 1, rp.debugLevel, rp.debugOn)
	}

	if message.Type == 7 {
		// Mem-Pool-Mini-Mem-Block-Confirm 7
		// this indicates that the resulting mini mem block has been received by a majority of replicas
		// add this this block to mini map
		rp.memPool.miniMap.Add(message)
		// add hash of mini block to the incoming channel
		rp.memPool.incomingBuffer = append(rp.memPool.incomingBuffer, message.UniqueId)
		// send an Mem-Pool-Mini-Mem-Block-Internal-Send-Ack 4 to the sender
		minimempoolinternalack := proto.MemPoolMini{
			Sender:   rp.name,
			Receiver: message.Sender,
			UniqueId: message.UniqueId,
			Type:     4,
			Note:     message.Note,
			Commands: nil,
			Creator:  message.Creator,
		}

		rpcPair := common.RPCPair{
			Code: rp.messageCodes.MemPoolMiniRPC,
			Obj:  &minimempoolinternalack,
		}

		rp.sendMessage(message.Sender, rpcPair)
		common.Debug("Sent Mini Mem Pool internal send ack message with type 4 to "+strconv.Itoa(int(message.Sender)), 0, rp.debugLevel, rp.debugOn)

		// create a new mem block if conditions are satisfied
		rp.createNewMemBlock()
	}
}

/*
	This method is invoked when the replica needs mini block to commit the block
	Randomly selects a replica and sends a Mem-Pool-Mini-Mem-Block-External-Request 5
*/

func (rp *Replica) sendExternalMiniBlockRequest(id string) {

	randomReplica := common.Get_Some_Node(rp.replicaArrayIndex)

	for randomReplica == rp.name {
		randomReplica = common.Get_Some_Node(rp.replicaArrayIndex)
	}

	externalMiniBlockRequest := proto.MemPoolMini{
		Sender:   rp.name,
		Receiver: randomReplica,
		UniqueId: id,
		Type:     5,
		Note:     "",
		Commands: nil,
		Creator:  -1, // not important in this case
	}

	rpcPair := common.RPCPair{
		Code: rp.messageCodes.MemPoolMiniRPC,
		Obj:  &externalMiniBlockRequest,
	}

	rp.sendMessage(randomReplica, rpcPair)
	common.Debug("Sent Mini Mem Pool external request message with type 5 to "+strconv.Itoa(int(randomReplica)), 1, rp.debugLevel, rp.debugOn)
}

/*
	This method is invoked when the replica wants to send a Mem-Pool-Mini-Mem-Block-Client-Response 8 to
	a client
*/

func (rp *Replica) sendMiniMemPoolClientResponse(message *proto.MemPoolMini) {
	message.Sender = rp.name
	message.Type = 8

	rpcPair := common.RPCPair{
		Code: rp.messageCodes.MemPoolMiniRPC,
		Obj:  message,
	}

	rp.sendMessage(message.Receiver, rpcPair)
	common.Debug("Sent Mini Mem Pool client response with type 8 to "+strconv.Itoa(int(message.Receiver)), 0, rp.debugLevel, rp.debugOn)
}
