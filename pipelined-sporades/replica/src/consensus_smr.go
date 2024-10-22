package src

import (
	"fmt"
	"pipelined-sporades/common"
	"pipelined-sporades/proto"
	"strconv"
	"strings"
	"time"
)

// upon receiving a client batch, add the batches to the buffer

func (rp *Replica) handleClientBatch(batch *proto.ClientBatch) {
	if rp.logPrinted {
		return
	}

	rp.incomingRequests = append(rp.incomingRequests, batch)
	if rp.debugOn {
		rp.debug("put incoming client batch to buffer: "+fmt.Sprintf("%v", batch), 0)
	}
	rp.propose(false, false)
}

/*
	Commit all the entries from the last committed block +1 to blockCommit
*/

func (rp *Replica) updateSMR() {

	head := rp.consensus.blockCommit // the last block to commit
	if head == nil {
		return
	}
	//toCommit contains all the blocks from the lastCommittedBlock (not including) to rp.blockCommit (included)
	toCommit := make([]*proto.Pipelined_Sporades_Block, 0)

	for head.Id != rp.consensus.lastCommittedBlock.Id {

		//	toCommit.append(head)
		toCommit = append([]*proto.Pipelined_Sporades_Block{head}, toCommit...)

		//	head = head.parent
		parent_id := head.ParentId
		//	genesis block doesn't have a parent
		if parent_id == "genesis-block" {
			break
		}
		if len(strings.Trim(parent_id, " ")) == 0 {
			panic("empty parent id found in the block " + head.Id)
		}
		//	if parent_id is in the consensus pool
		headBlock, ok := rp.consensus.consensusPool.Get(parent_id)
		if ok {
			head = headBlock
		} else {
			// request the consensus block from some random peer
			if rp.debugOn {
				rp.debug("Consensus block "+parent_id+" does not exist, sending an external request", 1)
			}
			rp.sendExternalConsensusRequest(parent_id)
			return // because we are missing the blocks in the history
		}

	}

	// if the code comes to this place, then we have all the consensus blocks from lastCommittedBlock (not included) to blockCommit (included)

	//if there is nothing to commit
	if len(toCommit) == 0 {
		if rp.debugOn {
			rp.debug("There is nothing to commit", 0)
		}
		return
	}

	// we have every block needed to commit
	for i := 0; i < len(toCommit); i++ {
		nextBlockToCommit := toCommit[i] // toCommit[i] is the next block to be committed
		if nextBlockToCommit.Commands == nil {
			panic("should not happen, commands field nil in " + fmt.Sprintf("%v", nextBlockToCommit))
		}

		clientBatches := nextBlockToCommit.Commands.Requests
		responses := rp.updateApplicationLogic(clientBatches)
		rp.consensus.lastCommittedBlock = nextBlockToCommit
		if rp.debugOn {
			senders := make([]int64, 0)
			for b := 0; b < len(responses); b++ {
				senders = append(senders, responses[b].Sender)
			}
			rp.debug("Committed block "+rp.consensus.lastCommittedBlock.Id+" at time "+fmt.Sprintf(" %v", time.Now().Sub(rp.consensus.startTime))+" with "+strconv.Itoa(len(clientBatches))+" client batches from "+fmt.Sprintf("%v", senders), 42)
		}

		rp.decrementPipelined()
		if rp.debugOn {
			rp.debug("pipeline length "+strconv.Itoa(rp.consensus.pipelinedRequests), 0)
		}

		rp.sendClientResponses(responses)
		rp.removeDecidedItemsFromFutureProposals(clientBatches)
	}

}

// remove the items in array from incomingRequests

func (rp *Replica) removeDecidedItemsFromFutureProposals(batches []*proto.ClientBatch) {
	array := make([]string, 0)

	for i := 0; i < len(batches); i++ {
		array = append(array, batches[i].UniqueId)
	}

	// Create a set to store the elements of array
	set := make(map[string]bool)
	for _, elem := range array {
		set[elem] = true
	}

	// Remove elements from pr.incomingRequests if they exist in the set
	j := 0
	for i := 0; i < len(rp.incomingRequests); i++ {
		if !set[rp.incomingRequests[i].UniqueId] {
			rp.incomingRequests[j] = rp.incomingRequests[i]
			j++
		}
	}

	// Truncate pr.toBeProposed to remove the remaining elements
	rp.incomingRequests = rp.incomingRequests[:j]
}

/*
	A generic application handler for processing a client batches, and sending responses
*/

func (rp *Replica) updateApplicationLogic(commands []*proto.ClientBatch) []*proto.ClientBatch {
	return rp.state.Execute(commands)
}

/*
	send the responses to client
*/

func (rp *Replica) sendClientResponses(commands []*proto.ClientBatch) {

	for i := 0; i < len(commands); i++ {
		if commands[i].Sender == -1 {
			continue
		}
		// send the response back to the client
		resClientBatch := proto.ClientBatch{
			UniqueId: commands[i].UniqueId,
			Requests: commands[i].Requests,
			Sender:   commands[i].Sender,
		}

		rpcPair := common.RPCPair{
			Code: rp.messageCodes.ClientBatchRpc,
			Obj:  &resClientBatch,
		}

		rp.sendMessage(int32(resClientBatch.Sender), rpcPair)
		if rp.debugOn {
			rp.debug("sent a client response to "+strconv.Itoa(int(resClientBatch.Sender)), -1)
		}
	}
}
