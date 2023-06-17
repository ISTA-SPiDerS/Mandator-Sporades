package src

import (
	"fmt"
	"pipelined-sporades/common"
	"pipelined-sporades/proto"
	"strconv"
	"time"
)

// handler for new view messages

func (rp *Replica) handleConsensusNewViewMessage(message *proto.Pipelined_Sporades) bool {
	if message.V >= rp.consensus.vCurr {
		if !rp.consensus.isAsync {
			_, ok := rp.consensus.newViewMessages[message.V]
			if !ok {
				rp.consensus.newViewMessages[message.V] = make([]*proto.Pipelined_Sporades, 0)
			}
			// add the new view to the recieved new view messages
			rp.consensus.newViewMessages[message.V] = append(rp.consensus.newViewMessages[message.V], message)
			if rp.debugOn {
				rp.debug("added a new view to buffer "+fmt.Sprintf("%v", message), 0)
			}
			if len(rp.consensus.newViewMessages[message.V]) == (rp.numReplicas/2 + 1) {
				if rp.debugOn {
					rp.debug("received n-f number of new view messages in view "+fmt.Sprintf("%v", message.V), 0)
				}
				// if we have collected a majority of new view messages, set block_high to be the highest received block_high in the new view messages
				highest_block_high := rp.extractHighestRankedBlockHigh(rp.consensus.newViewMessages[message.V])
				if rp.hasGreaterRank(highest_block_high.V, highest_block_high.R, rp.consensus.blockHigh.V, rp.consensus.blockHigh.R) {
					rp.consensus.blockHigh = highest_block_high
					if rp.debugOn {
						rp.debug("added a new block to store "+fmt.Sprintf("%v", highest_block_high), 0)
					}
					rp.consensus.consensusPool.Add(highest_block_high)
					rp.consensus.vCurr = rp.consensus.blockHigh.V
					rp.consensus.rCurr = rp.consensus.blockHigh.R
					if rp.debugOn {
						rp.debug("updated block_high to "+fmt.Sprintf("%v", rp.consensus.blockHigh), 0)
					}
				}
				if rp.debugOn {
					rp.debug("invoking propose after becoming the leader in view "+fmt.Sprintf("%v", message.V), 0)
				}
				rp.consensus.pipelinedSoFar = 0
				rp.propose(true, true)
			}
			return true
		} else {
			if rp.debugOn {
				rp.debug("delayed processing new view message"+fmt.Sprintf("%v", message)+" because i am still in async mode", 0)
			}
			return false
		}
	} else {
		if rp.debugOn {
			rp.debug("received a new view for an older view, hence rejected ", 0)
		}
		return true // we do not need this message, so just remove it
	}
}

/*
	Handler for sync consensus propose messages
*/

func (rp *Replica) handleConsensusProposeSync(message *proto.Pipelined_Sporades) bool {

	// if the new block has a rank greater v_cur, r_cur

	if rp.hasGreaterRank(message.BlockNew.V, message.BlockNew.R, rp.consensus.vCurr, rp.consensus.rCurr) {
		if !rp.consensus.isAsync {
			// cancel the timer
			if rp.consensus.viewTimer != nil {
				rp.consensus.viewTimer.Cancel()
				rp.consensus.viewTimer = nil
			}
			//	Save the new block in the consensus pool
			rp.consensus.consensusPool.Add(message.BlockNew)
			//	update v_cur, r_cur to that of v,r of the block
			rp.consensus.vCurr = message.BlockNew.V
			rp.consensus.rCurr = message.BlockNew.R
			// set block high to the new block
			rp.consensus.blockHigh = message.BlockNew
			// set block commit to block commit, and call updateSMR()
			if rp.hasGreaterRank(message.BlockCommit.V, message.BlockCommit.R, rp.consensus.blockCommit.V, rp.consensus.blockCommit.R) &&
				rp.hasGreaterRank(message.BlockCommit.V, message.BlockCommit.R, rp.consensus.lastCommittedBlock.V, rp.consensus.lastCommittedBlock.R) {
				rp.consensus.blockCommit = message.BlockCommit
				rp.consensus.consensusPool.Add(message.BlockCommit)
				rp.updateSMR()
			}
			// 	send <vote, v cur , r cur , block high > to Vcur leader
			nextLeader := rp.getLeader(rp.consensus.vCurr)
			if rp.debugOn {
				rp.debug("Sending sync vote to "+strconv.Itoa(int(nextLeader)), 0)
			}

			vote_block_high, err := CloneMyStruct(rp.consensus.blockHigh)
			if err != nil {
				panic(err.Error())
			}
			vote_block_high.Commands.Requests = nil
			voteMsg := proto.Pipelined_Sporades{
				Sender:      rp.name,
				Receiver:    nextLeader,
				UniqueId:    "",
				Type:        2,
				Note:        "",
				V:           rp.consensus.vCurr,
				R:           rp.consensus.rCurr,
				BlockHigh:   vote_block_high,
				BlockNew:    nil,
				BlockCommit: nil,
			}

			rpcPair := common.RPCPair{
				Code: rp.messageCodes.SporadesConsensus,
				Obj:  &voteMsg,
			}

			rp.sendMessage(nextLeader, rpcPair)
			if rp.debugOn {
				rp.debug("Sent sync vote to "+strconv.Itoa(int(nextLeader)), 0)
			}
			// start the timeout
			rp.setViewTimer()
			return true
		} else {
			if message.V > rp.consensus.vCurr {
				if rp.debugOn {
					rp.debug("cannot process the propose message "+fmt.Sprintf("%v", message)+" because I still haven't changed my mode to sync ", 0)
				}
				return false
			} else {
				if rp.debugOn {
					rp.debug("dismissed old propose sync message "+fmt.Sprintf("%v", message), 0)
				}
				return true // discard this message
			}
		}
	} else {
		if rp.debugOn {
			rp.debug("Rejected a propose sync message because its for a previous rank of "+fmt.Sprintf("view: %v, round: %v", message.V, message.R)+" where as I am in "+fmt.Sprintf("view: %v, round: %v", rp.consensus.vCurr, rp.consensus.rCurr)+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 0)
		}
		return true // we discard this, but we successfully processed the message
	}
}

/*
	Upon a timeout, broadcast a timeout message with vcur
*/

func (rp *Replica) handleConsensusInternalTimeout(message *proto.Pipelined_Sporades) bool {
	// this is triggered by the timeout
	rp.consensus.viewTimer = nil
	if rp.hasGreaterThanOrEqualRank(message.V, message.R, rp.consensus.vCurr, rp.consensus.rCurr) {
		if !rp.consensus.isAsync {

			height := rp.consensus.blockHigh.R - rp.consensus.blockCommit.R

			// broadcast <timeout, v cur , r cur , block high >

			for name, _ := range rp.replicaAddrList {

				timeoutMsg := proto.Pipelined_Sporades{
					Sender:      rp.name,
					Receiver:    name,
					UniqueId:    "",
					Type:        3,
					Note:        "",
					V:           rp.consensus.vCurr,
					R:           rp.consensus.rCurr,
					BlockHigh:   rp.makeNChain(rp.consensus.blockHigh, int(height)),
					BlockNew:    nil,
					BlockCommit: nil,
				}

				rpcPair := common.RPCPair{
					Code: rp.messageCodes.SporadesConsensus,
					Obj:  &timeoutMsg,
				}

				rp.sendMessage(name, rpcPair)
				if rp.debugOn {
					rp.debug("Sent timeout to "+strconv.Itoa(int(name)), 0)
				}
			}

			return true
		} else {
			if rp.debugOn {
				rp.debug("Rejected an an internal timeout notification because I am already in asycn I am in "+fmt.Sprintf("view: %v, round: %v", rp.consensus.vCurr, rp.consensus.rCurr)+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 0)
			}
			return true // we are already in the async path
		}
	} else {
		if rp.debugOn {
			rp.debug("Rejected an internal timeout message because its for a previous rank of "+fmt.Sprintf("view: %v, round: %v", message.V, message.R)+" where as I am in "+fmt.Sprintf("view: %v, round: %v", rp.consensus.vCurr, rp.consensus.rCurr)+" at time "+fmt.Sprintf("%v", time.Now().Sub(rp.consensus.startTime)), 0)
		}
		return true // this is an old message
	}
}

/*
	propose a batch of client commands in the sync path
*/

func (rp *Replica) propose(sendHistory bool, immediate bool) {
	if rp.consensus.isAsync {
		return // we are in the async path
	}

	if rp.name != rp.getLeader(rp.consensus.vCurr) {
		return // i am not the leader
	}

	if rp.consensus.pipelinedSoFar > rp.pipelineLength {
		if rp.debugOn {
			rp.debug("did not propose because pipeline length full with outstanding proposals "+strconv.Itoa(rp.consensus.pipelinedSoFar), 0)
		}
		return
	}

	if len(rp.incomingRequests) >= rp.replicaBatchSize || (time.Now().Sub(rp.consensus.lastProposedTime).Microseconds() > int64(rp.replicaBatchTime)) || immediate {

		if rp.debugOn {
			rp.debug("proposing a new batch in the sync path", 0)
		}

		var batches []*proto.ClientBatch
		if len(rp.incomingRequests) <= rp.replicaBatchSize {
			batches = rp.incomingRequests
			rp.incomingRequests = make([]*proto.ClientBatch, 0)
		} else {
			batches = rp.incomingRequests[:rp.replicaBatchSize]
			rp.incomingRequests = rp.incomingRequests[rp.replicaBatchSize:]
		}

		commands := &proto.ReplicaBatch{
			UniqueId: strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(int(rp.consensus.vCurr)) + "." + strconv.Itoa(int(rp.consensus.rCurr+1)) + "." + "r" + "." + strconv.Itoa(int(-1)),
			Requests: batches,
			Sender:   int64(rp.name),
		}

		var newBlock proto.Pipelined_Sporades_Block

		height := rp.consensus.blockHigh.R - rp.consensus.blockCommit.R

		if sendHistory {
			newBlock = proto.Pipelined_Sporades_Block{
				// creator_name.v.r.type.level. type can be r (regular) or f (fallback), level can be 1,2 or -1 (for regular blocks)
				Id:       strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(int(rp.consensus.vCurr)) + "." + strconv.Itoa(int(rp.consensus.rCurr+1)) + "." + "r" + "." + strconv.Itoa(int(-1)),
				V:        rp.consensus.vCurr,
				R:        rp.consensus.rCurr + 1,
				ParentId: rp.consensus.blockHigh.Id,
				Parent:   rp.makeNChain(rp.consensus.blockHigh, int(height)),
				Commands: commands,
				Level:    -1,
			}
		} else {
			newBlock = proto.Pipelined_Sporades_Block{
				// creator_name.v.r.type.level. type can be r (regular) or f (fallback), level can be 1,2 or -1 (for regular blocks)
				Id:       strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(int(rp.consensus.vCurr)) + "." + strconv.Itoa(int(rp.consensus.rCurr+1)) + "." + "r" + "." + strconv.Itoa(int(-1)),
				V:        rp.consensus.vCurr,
				R:        rp.consensus.rCurr + 1,
				ParentId: rp.consensus.blockHigh.Id,
				Parent:   nil,
				Commands: commands,
				Level:    -1,
			}
		}

		// add the new block to pool
		rp.consensus.consensusPool.Add(&newBlock)

		//	broadcast <propose, new block, block commit >
		if rp.debugOn {
			rp.debug("broadcasting propose type 1 "+fmt.Sprintf("%v", newBlock), 0)
		}
		for name, _ := range rp.replicaAddrList {
			if name == rp.name {
				continue // we do not send to self in pipelined version
			}
			proposeMsg := proto.Pipelined_Sporades{
				Sender:      rp.name,
				Receiver:    name,
				UniqueId:    "",
				Type:        1,
				Note:        "",
				V:           rp.consensus.vCurr,
				R:           rp.consensus.rCurr + 1,
				BlockHigh:   nil,
				BlockNew:    &newBlock,
				BlockCommit: rp.consensus.blockCommit,
			}

			rpcPair := common.RPCPair{
				Code: rp.messageCodes.SporadesConsensus,
				Obj:  &proposeMsg,
			}

			rp.sendMessage(name, rpcPair)
		}
		if rp.debugOn {
			rp.debug("broadcast propose type 1 ", 0)
		}

		rp.consensus.lastProposedTime = time.Now()
		rp.consensus.pipelinedSoFar++
		if rp.debugOn {
			rp.debug("pipeline length "+strconv.Itoa(rp.consensus.pipelinedSoFar), 0)
		}

		// update rank

		rp.consensus.vCurr = newBlock.V
		rp.consensus.rCurr = newBlock.R

		rp.consensus.blockHigh = &newBlock

		// send a vote to self
		if rp.debugOn {
			rp.debug("Sending sync vote to self", 0)
		}

		vote_block_high, err := CloneMyStruct(rp.consensus.blockHigh)
		if err != nil {
			panic(err)
		}
		vote_block_high.Parent = nil

		voteMsg := proto.Pipelined_Sporades{
			Sender:      rp.name,
			Receiver:    rp.name,
			UniqueId:    "",
			Type:        2,
			Note:        "",
			V:           rp.consensus.vCurr,
			R:           rp.consensus.rCurr,
			BlockHigh:   vote_block_high,
			BlockNew:    nil,
			BlockCommit: nil,
		}

		rpcPair := common.RPCPair{
			Code: rp.messageCodes.SporadesConsensus,
			Obj:  &voteMsg,
		}

		rp.sendMessage(rp.name, rpcPair)
		if rp.debugOn {
			rp.debug("Sent self vote", 0)
		}

		// start the timeout
		rp.setViewTimer()

	}

}

/*
	Handler for sync consensus vote messages
*/

func (rp *Replica) handleConsensusVoteSync(message *proto.Pipelined_Sporades) bool {
	if message.V == rp.consensus.vCurr {
		if !rp.consensus.isAsync {
			//	Save the vote in the vote replies:
			key := strconv.Itoa(int(message.V)) + "." + strconv.Itoa(int(message.R))
			//	if vote replies already has v.r key, then append to existing array, else create a new entry
			_, ok := rp.consensus.voteReplies[key]
			if !ok {
				rp.consensus.voteReplies[key] = make([]*proto.Pipelined_Sporades, 0)
			}
			rp.consensus.voteReplies[key] = append(rp.consensus.voteReplies[key], message)
			//	if for this v,r the array vote replies has n-f blocks
			votes, _ := rp.consensus.voteReplies[key]
			if len(votes) == rp.numReplicas/2+1 {

				//	If n-f received vote messages have the same block high and the rank of this received block high is (v,r):
				blockHigh, isSameHigh := rp.hasSameBlockHigh(votes)
				newBlockCommit := rp.consensus.blockCommit
				if isSameHigh && blockHigh.V == message.V && blockHigh.R == message.R {
					//	set block commit to this block high
					if rp.hasGreaterRank(blockHigh.V, blockHigh.R, rp.consensus.blockCommit.V, rp.consensus.blockCommit.R) &&
						rp.hasGreaterRank(blockHigh.V, blockHigh.R, rp.consensus.lastCommittedBlock.V, rp.consensus.lastCommittedBlock.R) {
						newBlockCommit = blockHigh
					}
				}

				savedBlock, ok := rp.consensus.consensusPool.Get(newBlockCommit.Id)

				if !ok {
					panic("voted block_commit does not appear in my store")
				}
				rp.consensus.blockCommit = savedBlock
				if rp.debugOn {
					rp.debug("sync leader updated block commit to "+fmt.Sprintf("%v", rp.consensus.blockCommit), 0)
				}
				rp.updateSMR()
				rp.propose(false, true)
			}
			return true
		} else {
			if rp.debugOn {
				rp.debug("Rejected a vote sync because i am in the async path of the same view ", 0)
			}
			return true
		}
	} else {
		if rp.debugOn {
			rp.debug("Rejected a sync vote message because its for a previous view", 0)
		}
		return true
	}
}

/*
	util method to check if all votes contain the same blockHigh
*/

func (rp *Replica) hasSameBlockHigh(votes []*proto.Pipelined_Sporades) (*proto.Pipelined_Sporades_Block, bool) {
	blockHigh := votes[0].BlockHigh
	for i := 1; i < len(votes); i++ {
		if blockHigh.Id != votes[i].BlockHigh.Id {
			return nil, false
		}
	}
	return blockHigh, true
}
