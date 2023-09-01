package src

import (
	"fmt"
	"pipelined-sporades/common"
	"pipelined-sporades/proto"
	"strconv"
	"strings"
)

/*
	handler for external consensus timeout messages
*/

func (rp *Replica) handleConsensusTimeout(message *proto.Pipelined_Sporades) bool {
	// if the view of the message is greater than or equals vCur and isAsync is false
	if message.V >= rp.consensus.vCurr {
		if !rp.consensus.isAsync {
			//  save the timeout message in the timeoutMessages
			_, ok := rp.consensus.timeoutMessages[message.V]
			if !ok {
				rp.consensus.timeoutMessages[message.V] = make([]*proto.Pipelined_Sporades, 0)
			}

			if rp.senderInSet(rp.consensus.timeoutMessages[message.V], message.Sender) {
				if rp.debugOn {
					rp.debug("duplicate messages", 0)
				}
				return true // though unsuccessful, we process the message
			}

			rp.consensus.timeoutMessages[message.V] = append(rp.consensus.timeoutMessages[message.V], message)

			timeouts, _ := rp.consensus.timeoutMessages[message.V]

			//	if n-f timeout messages for the v of message
			if len(timeouts) == rp.numReplicas/2+1 {
				//	set is Async to true
				rp.consensus.isAsync = true
				if rp.consensus.viewTimer != nil {
					rp.consensus.viewTimer.Cancel()
					rp.consensus.viewTimer = nil
				}
				rp.consensus.pipelinedRequests = 0
				if rp.debugOn {
					rp.debug("Entering view change in view "+strconv.Itoa(int(rp.consensus.vCurr)), 1)
				}
				//	Update block high to be the block in the timeout messages with the highest rank and my own block high
				tempBlockHigh := rp.extractHighestRankedBlockHigh(timeouts)
				if rp.hasGreaterRank(tempBlockHigh.V, tempBlockHigh.R, rp.consensus.blockHigh.V, rp.consensus.blockHigh.R) {
					//	set vCur, rCur to v, max(r cur , block high .r)
					rp.consensus.consensusPool.Add(tempBlockHigh)
					rp.consensus.blockHigh = rp.CloneMyStruct(tempBlockHigh)
					rp.consensus.rCurr = rp.consensus.blockHigh.R
				}
				rp.consensus.vCurr = message.V

				var batches []*proto.ClientBatch
				if len(rp.incomingRequests) <= rp.replicaBatchSize {
					batches = rp.incomingRequests
					rp.incomingRequests = make([]*proto.ClientBatch, 0)
				} else {
					batches = rp.incomingRequests[:rp.replicaBatchSize]
					rp.incomingRequests = rp.incomingRequests[rp.replicaBatchSize:]
				}

				commands := proto.ReplicaBatch{
					UniqueId: strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(int(rp.consensus.vCurr)) + "." + strconv.Itoa(int(rp.consensus.rCurr+1)) + "." + "f" + "." + strconv.Itoa(int(1)),
					Requests: batches,
					Sender:   int64(rp.name),
				}

				height := rp.consensus.blockHigh.R - rp.consensus.blockCommit.R

				//	Form a new height 1 fallback block B f1 =(cmnds, v cur , r cur +1, block high)
				newLevel1FallBackBlock := proto.Pipelined_Sporades_Block{
					// creator_name.v.r.type.level. type can be r (regular) or f (fallback), level can be 1,2 or -1 (for regular blocks)
					Id:       strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(int(rp.consensus.vCurr)) + "." + strconv.Itoa(int(rp.consensus.rCurr+1)) + "." + "f" + "." + strconv.Itoa(int(1)),
					V:        rp.consensus.vCurr,
					R:        rp.consensus.rCurr + 1,
					ParentId: rp.consensus.blockHigh.Id,
					Parent:   rp.makeNChain(rp.consensus.blockHigh, int(height)),
					Commands: &commands,
					Level:    1,
				}

				// save the new block in the store
				rp.consensus.consensusPool.Add(&newLevel1FallBackBlock)

				//	broadcast <propose-async, B f1>
				for name, _ := range rp.replicaAddrList {

					proposeAsyncMsg := proto.Pipelined_Sporades{
						Sender:      rp.name,
						Receiver:    name,
						UniqueId:    "",
						Type:        4,
						Note:        "",
						V:           rp.consensus.vCurr,
						R:           rp.consensus.rCurr + 1,
						BlockHigh:   nil,
						BlockNew:    &newLevel1FallBackBlock,
						BlockCommit: nil,
					}

					rpcPair := common.RPCPair{
						Code: rp.messageCodes.SporadesConsensus,
						Obj:  &proposeAsyncMsg,
					}

					rp.sendMessage(name, rpcPair)
					if rp.debugOn {
						rp.debug("Sent propose-async level 1 to "+strconv.Itoa(int(name)), 1)
					}
				}
			}
			return true
		} else {
			// process later
			if message.V > rp.consensus.vCurr {
				if rp.debugOn {
					rp.debug("Sent an internal timeout because it is for a future view change ", 2)
				}
				return false
			}
			return true // we discard the message
		}
	} else {
		if rp.debugOn {
			rp.debug("Rejected an external timeout message because I am in a higher view", 1)
		}
		return true // this message was processed
	}
}

/*
	Handler for async propose fallback messages
*/

func (rp *Replica) handleConsensusProposeAsync(message *proto.Pipelined_Sporades) bool {
	// if the v of message is equal to v cur
	if message.V == rp.consensus.vCurr {
		//if is async is true
		if rp.consensus.isAsync == true {
			if rp.hasGreaterRank(message.BlockNew.V, message.BlockNew.R, rp.consensus.vCurr, rp.consensus.rCurr) {
				// save the block in the store
				rp.consensus.consensusPool.Add(message.BlockNew)
				vote_block_new := rp.CloneMyStruct(message.BlockNew)
				vote_block_new.Commands = nil
				// send <vote-async, B, h> to p
				voteAsyncMsg := proto.Pipelined_Sporades{
					Sender:      rp.name,
					Receiver:    message.Sender,
					UniqueId:    message.BlockNew.Id,
					Type:        5,
					Note:        "",
					V:           rp.consensus.vCurr,
					R:           rp.consensus.rCurr,
					BlockHigh:   nil,
					BlockNew:    vote_block_new,
					BlockCommit: nil,
				}

				rpcPair := common.RPCPair{
					Code: rp.messageCodes.SporadesConsensus,
					Obj:  &voteAsyncMsg,
				}

				rp.sendMessage(message.Sender, rpcPair)
				if rp.debugOn {
					rp.debug("Sent vote-async 5 to "+strconv.Itoa(int(message.Sender)), 1)
				}
				// if h == 2 set B fall [p] to B, and adapt level 1 block unless I have already sent level 2 block
				if message.BlockNew.Level == 2 {

					// save the level 2 block in b_fall

					key := strconv.Itoa(int(message.BlockNew.V)) + "." + strconv.Itoa(int(message.BlockNew.Level))
					_, ok := rp.consensus.bFall[key]
					if !ok {
						rp.consensus.bFall[key] = make([]string, 0)
					}
					duplicate := false
					for b := 0; b < len(rp.consensus.bFall[key]); b++ {
						if rp.consensus.bFall[key][b] == message.BlockNew.Id {
							duplicate = true // a duplicate
							break
						}
					}

					if !duplicate {
						rp.consensus.bFall[key] = append(rp.consensus.bFall[key], message.BlockNew.Id)
					}

					// if I still haven't sent a level 2 block, adapt the level 1 block, and send a level 2 block

					if rp.consensus.sentLevel2Block[message.V] == false {
						level1Block := message.BlockNew.Parent
						if level1Block != nil {
							if rp.debugOn {
								rp.debug("Adopted level 1 fallback block from  "+strconv.Itoa(int(message.Sender)), 1)
							}
							rp.consensus.consensusPool.Add(level1Block)

							var batches []*proto.ClientBatch
							if len(rp.incomingRequests) <= rp.replicaBatchSize {
								batches = rp.incomingRequests
								rp.incomingRequests = make([]*proto.ClientBatch, 0)
							} else {
								batches = rp.incomingRequests[:rp.replicaBatchSize]
								rp.incomingRequests = rp.incomingRequests[rp.replicaBatchSize:]
							}

							commands := proto.ReplicaBatch{
								UniqueId: strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(int(rp.consensus.vCurr)) + "." + strconv.Itoa(int(level1Block.R+1)) + "." + "f" + "." + strconv.Itoa(int(2)),
								Requests: batches,
								Sender:   int64(rp.name),
							}

							// create height 2 block and broadcast
							newLevel2FallBackBlock := proto.Pipelined_Sporades_Block{
								// creator_name.v.r.type.level. type can be r (regular) or f (fallback), level can be 1,2 or -1 (for regular blocks)
								Id:       strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(int(rp.consensus.vCurr)) + "." + strconv.Itoa(int(level1Block.R+1)) + "." + "f" + "." + strconv.Itoa(int(2)),
								V:        rp.consensus.vCurr,
								R:        level1Block.R + 1,
								ParentId: level1Block.Id,
								Parent:   level1Block,
								Commands: &commands,
								Level:    2,
							}

							// save the new block in the store
							rp.consensus.consensusPool.Add(&newLevel2FallBackBlock)

							//	broadcast <propose-async, B f2>

							for name, _ := range rp.replicaAddrList {

								proposeAsyncLevel2Msg := proto.Pipelined_Sporades{
									Sender:      rp.name,
									Receiver:    name,
									UniqueId:    "",
									Type:        4,
									Note:        "",
									V:           rp.consensus.vCurr,
									R:           level1Block.R + 1,
									BlockHigh:   nil,
									BlockNew:    &newLevel2FallBackBlock,
									BlockCommit: nil,
								}

								rpcPair := common.RPCPair{
									Code: rp.messageCodes.SporadesConsensus,
									Obj:  &proposeAsyncLevel2Msg,
								}

								rp.sendMessage(name, rpcPair)
								if rp.debugOn {
									rp.debug("Sent adopted and extended propose-async level 2 to "+strconv.Itoa(int(name)), 1)
								}
							}
							rp.consensus.sentLevel2Block[message.V] = true
						}
					}

				}
			}
			return true
		} else {
			// given that there is level 1 fallback block, eventually i should also get it, so save this incoming message to process later
			if rp.debugOn {
				rp.debug("Sent an internal fallback propose level 1/2 block because I still haven't received n-f timeouts, my rank is "+fmt.Sprintf("v: %v, r:%v and the message rank is v:%v, r:%v", rp.consensus.vCurr, rp.consensus.rCurr, message.V, message.R), 1)
			}
			return false
		}
	} else if message.V > rp.consensus.vCurr {
		if rp.debugOn {
			rp.debug("Sent an internal propose async level 1/2 block because I still haven't reached the view ", 1)
		}
		return false
	} else {
		return true // we do not care about old view change messages
	}
}

/*
	handler for consensus async vote messages
*/

func (rp *Replica) handleConsensusAsyncVote(message *proto.Pipelined_Sporades) bool {
	//	if the view is equal to current view and  isAsync is true
	if message.BlockNew.V == rp.consensus.vCurr && rp.consensus.isAsync {
		//	save the vote in the async vote store
		key := message.BlockNew.Id
		_, ok := rp.consensus.consensusPool.Get(key)
		if !ok {
			panic("Key " + key + " is not found in the block store, which is a fallback block. Triggering this error after receiving an async-vote")
		}

		acks := rp.consensus.consensusPool.GetAcks(key)
		found := false
		for a := 0; a < len(acks); a++ {
			if acks[a] == message.Sender {
				found = true
				break
			}
		}
		if !found {
			rp.consensus.consensusPool.AddAck(key, message.Sender)
		} else {
			return true
		}
		//	if there are n-f async votes for the block that I proposed
		acks = rp.consensus.consensusPool.GetAcks(key)
		if len(acks) == rp.numReplicas/2+1 {

			if message.BlockNew.Level == 1 && !rp.consensus.sentLevel2Block[rp.consensus.vCurr] {

				l1block, _ := rp.consensus.consensusPool.Get(key)
				rLevel1 := l1block.R

				var batches []*proto.ClientBatch
				if len(rp.incomingRequests) <= rp.replicaBatchSize {
					batches = rp.incomingRequests
					rp.incomingRequests = make([]*proto.ClientBatch, 0)
				} else {
					batches = rp.incomingRequests[:rp.replicaBatchSize]
					rp.incomingRequests = rp.incomingRequests[rp.replicaBatchSize:]
				}

				commands := proto.ReplicaBatch{
					UniqueId: strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(int(rp.consensus.vCurr)) + "." + strconv.Itoa(int(rLevel1+1)) + "." + "f" + "." + strconv.Itoa(int(2)),
					Requests: batches,
					Sender:   int64(rp.name),
				}

				//	Form a new height 2 fallback block B f2 =(cmnds, v cur , B.r+1, B, 2)
				newLevel2FallBackBlock := proto.Pipelined_Sporades_Block{
					Id:       strconv.Itoa(int(rp.name)) + "." + strconv.Itoa(int(rp.consensus.vCurr)) + "." + strconv.Itoa(int(rLevel1+1)) + "." + "f" + "." + strconv.Itoa(int(2)), //creator_name.v.r.type.level. type can be r (regular) or f (fallback), level can be 1,2 or -1 (for regular blocks)
					V:        rp.consensus.vCurr,
					R:        rLevel1 + 1,
					Parent:   l1block,
					ParentId: l1block.Id,
					Commands: &commands,
					Level:    2,
				}

				rp.consensus.consensusPool.Add(&newLevel2FallBackBlock)

				//	broadcast <propose-async, B f2 , self.id, 2>
				for name, _ := range rp.replicaAddrList {

					proposeAsyncLevel2Msg := proto.Pipelined_Sporades{
						Sender:      rp.name,
						Receiver:    name,
						UniqueId:    "",
						Type:        4,
						Note:        "",
						V:           rp.consensus.vCurr,
						R:           rLevel1 + 1,
						BlockHigh:   nil,
						BlockNew:    &newLevel2FallBackBlock,
						BlockCommit: nil,
					}

					rpcPair := common.RPCPair{
						Code: rp.messageCodes.SporadesConsensus,
						Obj:  &proposeAsyncLevel2Msg,
					}

					rp.sendMessage(name, rpcPair)
					if rp.debugOn {
						rp.debug("Sent propose-async level 2 to "+strconv.Itoa(int(name)), 1)
					}
				}

				rp.consensus.sentLevel2Block[rp.consensus.vCurr] = true

			} else if message.BlockNew.Level == 2 {
				level_2_block, _ := rp.consensus.consensusPool.Get(message.BlockNew.Id)
				// broadcast <fallback-complete, B, v cur , self.id>
				for name, _ := range rp.replicaAddrList {

					proposeAsyncFallbackComplete := proto.Pipelined_Sporades{
						Sender:      rp.name,
						Receiver:    name,
						UniqueId:    "",
						Type:        9,
						Note:        "",
						V:           rp.consensus.vCurr,
						R:           message.BlockNew.R,
						BlockHigh:   nil,
						BlockNew:    level_2_block,
						BlockCommit: nil,
					}

					rpcPair := common.RPCPair{
						Code: rp.messageCodes.SporadesConsensus,
						Obj:  &proposeAsyncFallbackComplete,
					}

					rp.sendMessage(name, rpcPair)
					if rp.debugOn {
						rp.debug("Sent propose-async-fallback-complete to "+strconv.Itoa(int(name)), 1)
					}
				}
			}
		}
		return true
	} else {
		if rp.debugOn {
			rp.debug("rejected a vote sync because i have moved on", 1)
		}
		return true // this message is obsolete
	}
}

/*
	Handler for async consensus fallback complete messages
*/

func (rp *Replica) handleConsensusFallbackCompleteMessage(message *proto.Pipelined_Sporades) bool {
	if message.BlockNew.V == rp.consensus.vCurr {
		if rp.consensus.isAsync == true {
			// save the level 2 confirmed block
			rp.consensus.consensusPool.Add(message.BlockNew)

			// add the id of the confirmed level 2 block to Bfall[v.3]
			key := strconv.Itoa(int(message.BlockNew.V)) + "." + strconv.Itoa(3)

			_, ok := rp.consensus.bFall[key]

			if !ok {
				rp.consensus.bFall[key] = make([]string, 0)
			}
			found := false
			for b := 0; b < len(rp.consensus.bFall[key]); b++ {
				if rp.consensus.bFall[key][b] == message.BlockNew.Id {
					found = true
					break
				}
			}
			if !found {
				rp.consensus.bFall[key] = append(rp.consensus.bFall[key], message.BlockNew.Id)
			}

			if len(rp.consensus.bFall[key]) == rp.numReplicas/2+1 {
				// received majority fallback complete messages

				leaderNode := rp.consensus.randomness[rp.consensus.vCurr] // l is the async leader of this view
				if rp.debugOn {
					rp.debug("Possible leader node for the view "+strconv.Itoa(int(rp.consensus.vCurr))+" is "+strconv.Itoa(leaderNode), 1)
				}
				// if height 2 block by leader exists in the first n-f fallback height 2 blocks received in the fallback messages
				height2ConfirmedLeaderBlockExists := false
				var height2ConfirmedLeaderBlock *proto.Pipelined_Sporades_Block
				for j := 0; j < len(rp.consensus.bFall[key]); j++ {
					creator := strings.Split(rp.consensus.bFall[key][j], ".")[0]
					if creator == strconv.Itoa(leaderNode) {
						height2ConfirmedLeaderBlockExists = true
						height2ConfirmedLeaderBlock, _ = rp.consensus.consensusPool.Get(rp.consensus.bFall[key][j])
						break
					}
				}
				if height2ConfirmedLeaderBlockExists {
					//	Set block high, block commit to height 2 block from l
					rp.consensus.blockHigh = rp.CloneMyStruct(height2ConfirmedLeaderBlock)
					rp.consensus.rCurr = rp.consensus.blockHigh.R
					if rp.hasGreaterRank(height2ConfirmedLeaderBlock.V, height2ConfirmedLeaderBlock.R, rp.consensus.blockCommit.V, rp.consensus.blockCommit.R) &&
						rp.hasGreaterRank(height2ConfirmedLeaderBlock.V, height2ConfirmedLeaderBlock.R, rp.consensus.lastCommittedBlock.V, rp.consensus.lastCommittedBlock.R) {
						rp.consensus.blockCommit = rp.CloneMyStruct(height2ConfirmedLeaderBlock)
						rp.updateSMR()
					}
					if rp.debugOn {
						rp.debug("Updated block commit in the async path for block "+rp.consensus.blockCommit.Id, 1)
					}
					//	Set v cur , r cur to rank(block high)
					rp.consensus.vCurr = message.V

				} else {
					if rp.debugOn {
						rp.debug("Leader node level 2 confirmed proposal does not exist for the view "+strconv.Itoa(int(rp.consensus.vCurr)), 1)
					}
					//else if height 2 block from the leader exists in the Bfall
					height2Key := strconv.Itoa(int(message.V)) + "." + strconv.Itoa(2)
					height2Blocks, ok := rp.consensus.bFall[height2Key]
					if ok {
						height2LeaderBlockExists := false
						var height2LeaderBlock *proto.Pipelined_Sporades_Block
						for k := 0; k < len(height2Blocks); k++ {
							creator := strings.Split(height2Blocks[k], ".")[0]
							if creator == strconv.Itoa(leaderNode) {
								height2LeaderBlockExists = true
								height2LeaderBlock, _ = rp.consensus.consensusPool.Get(height2Blocks[k])
								break
							}
						}
						if height2LeaderBlockExists {
							rp.consensus.blockHigh = rp.CloneMyStruct(height2LeaderBlock)
							rp.consensus.rCurr = rp.consensus.blockHigh.R
							rp.consensus.vCurr = message.V
							if rp.debugOn {
								rp.debug("Updated block high (not committed) in the async path for block "+rp.consensus.blockHigh.Id, 1)
							}
						} else {
							if rp.debugOn {
								rp.debug("Leader node level 2 async proposal does not exists for the view "+strconv.Itoa(int(rp.consensus.vCurr)), 1)
							}
						}
					}
				}
				//	Set v cur to v cur +1
				rp.consensus.vCurr++
				if rp.debugOn {
					rp.debug("Incremented the view to "+strconv.Itoa(int(rp.consensus.vCurr)), 1)
				}
				//	Set isAsync to false
				rp.consensus.isAsync = false
				// send <new-view, v cur block high > to L_Vcur

				nextLeader := rp.getLeader(rp.consensus.vCurr)

				newViewMsg := proto.Pipelined_Sporades{
					Sender:      rp.name,
					Receiver:    nextLeader,
					UniqueId:    "",
					Type:        10,
					Note:        "",
					V:           rp.consensus.vCurr,
					R:           rp.consensus.rCurr,
					BlockHigh:   rp.makeNChain(rp.consensus.blockHigh, int(rp.consensus.blockHigh.R-rp.consensus.blockCommit.R)),
					BlockNew:    nil,
					BlockCommit: nil,
				}

				rpcPair := common.RPCPair{
					Code: rp.messageCodes.SporadesConsensus,
					Obj:  &newViewMsg,
				}

				rp.sendMessage(nextLeader, rpcPair)
				if rp.debugOn {
					rp.debug("Exiting view change and sending new view to "+strconv.Itoa(int(nextLeader))+" after the view change", 1)
				}
				// start the timeout
				if rp.consensus.viewTimer != nil {
					rp.consensus.viewTimer.Cancel()
					rp.consensus.viewTimer = nil
				}
				rp.setViewTimer()
			}
			return true
		} else {
			if rp.debugOn {
				rp.debug("Sent an internal fallback-complete because I still haven't converted to async ", 1)
			}
			return false
		}
	} else if message.BlockNew.V > rp.consensus.vCurr {
		if rp.debugOn {
			rp.debug("Sent an internal fallback-complete because I still haven't reached the view ", 1)
		}
		return false
	} else {
		return true // this message is discarded
	}
}
