package src

import (
	"async-consensus/common"
	"async-consensus/proto"
	"fmt"
	"strconv"
	"time"
)

func (wr *Worker) handleStatus(message *proto.Status) {
	fmt.Print("Status from " + strconv.Itoa(int(message.Sender)) + " \n")
	if message.Type == 1 {
		if wr.serverStarted == false {
			wr.serverStarted = true
			wr.ConnectBootStrap()
		}

	}

	time.Sleep(2 * time.Second)

	// if the sender is client, then forward the status to replica
	if wr.getNodeType(message.Sender) == "client" {
		common.Debug("Sending status to default replica", 0, wr.debugLevel, wr.debugOn)

		statusMessage := proto.Status{
			Sender:   wr.name,
			Receiver: wr.defaultReplicaName,
			UniqueId: message.UniqueId,
			Type:     message.Type,
			Note:     message.Note,
		}

		rpcPair := common.RPCPair{
			Code: wr.messageCodes.StatusRPC,
			Obj:  &statusMessage,
		}

		wr.sendMessage(wr.defaultReplicaName, rpcPair)
		common.Debug("Sent status to default replica", 0, wr.debugLevel, wr.debugOn)

	}

	// if the sender is the default replica, forward the status to client
	if wr.getNodeType(message.Sender) == "replica" {
		common.Debug("Sending status request to all clients", 0, wr.debugLevel, wr.debugOn)
		// note that worker might not have connections to all clients, in this case this message will be dropped

		for name, _ := range wr.clientAddrList {

			statusMessage := proto.Status{
				Sender:   wr.name,
				Receiver: name,
				UniqueId: message.UniqueId,
				Type:     message.Type,
				Note:     message.Note,
			}

			rpcPair := common.RPCPair{
				Code: wr.messageCodes.StatusRPC,
				Obj:  &statusMessage,
			}

			wr.sendMessage(name, rpcPair)
			common.Debug("Sent status to "+strconv.Itoa(int(name)), 0, wr.debugLevel, wr.debugOn)
		}
	}
}
