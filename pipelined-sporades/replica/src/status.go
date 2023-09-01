package src

import (
	"fmt"
	"pipelined-sporades/common"
	"pipelined-sporades/proto"
	"time"
)

/*
	Handler for status message
		1. Invoke bootstrap / start consensus or printlog depending on the operation type
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

			// empty the incoming channel
			go func() {
				for true {
					_ = <-rp.incomingChan
				}
			}()

			rp.printLogConsensus() // this is for consensus testing purposes
		}
	} else if message.Type == 3 {
		if rp.consensusStarted == false {
			rp.consensusStarted = true
			rp.sendGenesisConsensusNewView()
			rp.consensus.startTime = time.Now()
			if rp.debugOn {
				rp.debug("started Sporades consensus", 0)
			}

			time.Sleep(time.Duration(2) * time.Second)

			// a dummy client request generator to avoid view timeouts due to absense of client requests
			go func() {
				for true {
					time.Sleep(time.Duration(rp.viewTimeout/2) * time.Microsecond)
					dummyClientBatch := proto.ClientBatch{
						UniqueId: "dummy",
						Requests: make([]*proto.SingleOperation, 0),
						Sender:   -1,
					}

					rpcPair := common.RPCPair{
						Code: rp.messageCodes.ClientBatchRpc,
						Obj:  &dummyClientBatch,
					}

					rp.incomingChan <- &rpcPair
				}
			}()
		}

	}

	statusMessage := proto.Status{
		Type: message.Type,
		Note: message.Note,
	}

	rpcPair := common.RPCPair{
		Code: rp.messageCodes.StatusRPC,
		Obj:  &statusMessage,
	}

	rp.sendMessage(int32(message.Sender), rpcPair)
	if rp.debugOn {
		rp.debug("Sent status response", 0)
	}
}
