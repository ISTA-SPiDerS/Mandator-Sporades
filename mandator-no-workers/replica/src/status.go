package src

import (
	"async-consensus/common"
	"async-consensus/proto"
	"fmt"
	"time"
)

/*
	Handler for status message
		1. Invoke bootstrap or printlog depending on the operation type
		2. Send a response back to the sender
*/

func (rp *Replica) handleStatus(message *proto.Status) {
	fmt.Print("Status  " + fmt.Sprintf("%v", message) + " \n")
	if message.Type == 1 {
		if rp.serverStarted == false {
			rp.serverStarted = true
			rp.ConnectBootStrap()
			time.Sleep(2 * time.Second)
		}
	} else if message.Type == 2 {
		if rp.logPrinted == false {
			rp.logPrinted = true
			fmt.Printf(fmt.Sprintf("Last completed Mem Blocks %v\n", rp.memPool.lastCompletedRounds))
			// empty the incoming channel
			go func() {
				for true {
					_ = <-rp.incomingChan
				}
			}()
			rp.printLogMemPool() // this is for the mem pool testing purposes
			if rp.consAlgo == "async" {
				fmt.Printf(fmt.Sprintf("last committed consensus block %v at time %v\n", rp.asyncConsensus.lastCommittedBlock, rp.asyncConsensus.lastCommittedTime.Sub(rp.asyncConsensus.startTime)))
				rp.printLogConsensus() // this is for consensus testing purposes
			} else if rp.consAlgo == "paxos" {
				fmt.Printf(fmt.Sprintf("last committed consensus instance %v at time %v\n", rp.paxosConsensus.lastCommittedLogIndex, rp.paxosConsensus.lastCommittedTime.Sub(rp.paxosConsensus.startTime)))
				fmt.Printf(fmt.Sprintf("last decided consensus instance %v\n", rp.paxosConsensus.lastDecidedLogIndex))
				fmt.Printf(fmt.Sprintf("View %v \n", rp.paxosConsensus.view))
				rp.printPaxosLogConsensus() // this is for consensus testing purposes
			}

		}
	} else if message.Type == 3 {
		if rp.consensusStarted == false {
			rp.consensusStarted = true
			if rp.consAlgo == "async" {
				rp.sendGenesisConsensusVote()
			} else if rp.consAlgo == "paxos" {
				rp.paxosConsensus.startTime = time.Now()
				rp.paxosConsensus.lastCommittedTime = time.Now()
				initLeader := rp.name
				for name, index := range rp.replicaArrayIndex {
					if index == 0 {
						initLeader = name
					}
				}
				if rp.name == initLeader {
					rp.sendPrepare()
				}
			}
		}
	}

	common.Debug("Sending status reply ", 0, rp.debugLevel, rp.debugOn)

	statusMessage := proto.Status{
		Type: message.Type,
		Note: message.Note,
	}

	rpcPair := common.RPCPair{
		Code: rp.messageCodes.StatusRPC,
		Obj:  &statusMessage,
	}

	rp.sendMessage(int32(message.Sender), rpcPair)
	common.Debug("Sent status ", 0, rp.debugLevel, rp.debugOn)

}
