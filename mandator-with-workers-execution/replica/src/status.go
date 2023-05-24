package src

import (
	"async-consensus/common"
	"async-consensus/proto"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

/*
	Handler for status message
		1. Invoke bootstrap or printlog depending on the operation type
		2. Send a response back to the sender
*/

func (rp *Replica) handleStatus(message *proto.Status) {
	fmt.Print("Status from " + strconv.Itoa(int(message.Sender)) + " \n")
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

	common.Debug("Sending status reply to worker ", 0, rp.debugLevel, rp.debugOn)

	statusMessage := proto.Status{
		Sender:   rp.name,
		Receiver: message.Sender,
		UniqueId: message.UniqueId,
		Type:     message.Type,
		Note:     message.Note,
	}

	rpcPair := common.RPCPair{
		Code: rp.messageCodes.StatusRPC,
		Obj:  &statusMessage,
	}

	rp.sendMessage(message.Sender, rpcPair)
	common.Debug("Sent status ", 0, rp.debugLevel, rp.debugOn)

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

	head := rp.asyncConsensus.blockCommit // the last block to commit
	if head == nil {
		return
	}
	genesisBlock, ok := rp.asyncConsensus.consensusPool.Get("genesis-block")
	if !ok {
		panic("Genesis Block not found when printing the logs")
	}
	//toCommit = [] contains all the entries from the genesisBlock (not including) to rp.blockCommit (included)
	toCommit := make([]*proto.AsyncConsensus_Block, 0)

	for head.Id != genesisBlock.Id {
		//	toCommit.append(head)
		toCommit = append([]*proto.AsyncConsensus_Block{head}, toCommit...)
		//	head = head.parent
		head = head.Parent

		//	if head is in the consensus pool
		headBlock, ok := rp.asyncConsensus.consensusPool.Get(head.Id)
		if ok {
			head = headBlock
		} else {
			panic("Consensus block " + head.Id + " not found in the pool")
		}

	}

	lastCommittedMemPoolIndexes := make([]int, rp.numReplicas)
	for i := 0; i < rp.numReplicas; i++ {
		lastCommittedMemPoolIndexes[i] = 0
	}

	for i := 0; i < len(toCommit); i++ {
		nextBlockToCommit := toCommit[i] // toCommit[i] is the next block to be committed
		nextMemBlockLogPositionsToCommit := nextBlockToCommit.Commands

		// for each log position in nextMemBlockLogPositionsToCommit that corresponds to different replicas, check if the index is --
		// greater than the lastCommittedMemPoolIndexes
		for j := 0; j < rp.numReplicas; j++ {
			if int(nextMemBlockLogPositionsToCommit[j]) > lastCommittedMemPoolIndexes[j] {
				// there are new entries to commit for this index
				startMemPoolCounter := lastCommittedMemPoolIndexes[j] + 1
				lastMemPoolCounter := int(nextMemBlockLogPositionsToCommit[j])

				for k := startMemPoolCounter; k <= lastMemPoolCounter; k++ {
					memPoolName := strconv.Itoa(int(rp.getReplicaName(j))) + "." + strconv.Itoa(k)
					memBlock, _ := rp.memPool.blockMap.Get(memPoolName)
					miniBlockIDs := memBlock.Minimemblocks
					for l := 0; l < len(miniBlockIDs); l++ {
						miniBlockID := miniBlockIDs[l].UniqueId
						miniBlock, _ := rp.memPool.miniMap.Get(miniBlockID)
						clientBatches := miniBlock.Commands
						for clientBatchIndex := 0; clientBatchIndex < len(clientBatches); clientBatchIndex++ {
							clientBatch := clientBatches[clientBatchIndex]
							clientBatchID := clientBatch.Id
							clientBatchCommands := clientBatch.Commands
							for clientRequestIndex := 0; clientRequestIndex < len(clientBatchCommands); clientRequestIndex++ {
								clientRequestID := clientBatchCommands[clientRequestIndex].UniqueId
								_, _ = f.WriteString(nextBlockToCommit.Id + "-" + memPoolName + "-" + miniBlockID + "-" + clientBatchID + "-" + clientRequestID + ":" + clientBatchCommands[clientRequestIndex].Command + "\n")
							}
						}
					}
				}
				lastCommittedMemPoolIndexes[j] = lastMemPoolCounter
			}
		}
	}

}

/*
	Printing the mem store for debug purpose
*/

func (rp *Replica) printLogMemPool() {
	f, err := os.Create(rp.logFilePath + strconv.Itoa(int(rp.name)) + "-mem-pool.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for memBlockID, memBlock := range rp.memPool.blockMap.MessageBlocks {
		miniBlocks := memBlock.MessageBlock.Minimemblocks
		for miniBlockIndex := 0; miniBlockIndex < len(miniBlocks); miniBlockIndex++ {
			miniBlockId := miniBlocks[miniBlockIndex].UniqueId
			miniBlock, _ := rp.memPool.miniMap.Get(miniBlockId)
			clientBatches := miniBlock.Commands
			for clientBatchIndex := 0; clientBatchIndex < len(clientBatches); clientBatchIndex++ {
				clientBatch := clientBatches[clientBatchIndex]
				clientBatchID := clientBatch.Id
				clientBatchCommands := clientBatch.Commands
				for clientRequestIndex := 0; clientRequestIndex < len(clientBatchCommands); clientRequestIndex++ {
					clientRequestID := clientBatchCommands[clientRequestIndex].UniqueId
					_, _ = f.WriteString(memBlockID + "-" + miniBlockId + "-" + clientBatchID + "-" + clientRequestID + ":" + clientBatchCommands[clientRequestIndex].Command + "\n")
				}
			}
		}
	}

}
