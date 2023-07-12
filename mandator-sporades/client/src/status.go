package src

import (
	"async-consensus/common"
	"async-consensus/proto"
	"fmt"
	"strconv"
	"time"
)

/*
	When a status response is received, print it to console
*/

func (cl *Client) handleClientStatusResponse(response *proto.Status) {
	fmt.Printf("Status response %v\n", response)
}

/*
	Send a status request to all the replicas
*/

func (cl *Client) SendStatus(operationType int) {
	common.Debug("Sending status request to all replicas", 0, cl.debugLevel, cl.debugOn)

	for name, _ := range cl.replicaAddrList {

		statusRequest := proto.Status{
			Type:   int32(operationType),
			Note:   "",
			Sender: int64(cl.clientName),
		}

		rpcPair := common.RPCPair{
			Code: cl.messageCodes.StatusRPC,
			Obj:  &statusRequest,
		}

		cl.sendMessage(name, rpcPair)
		common.Debug("Sent status to "+strconv.Itoa(int(name)), 0, cl.debugLevel, cl.debugOn)
	}
	time.Sleep(time.Duration(statusTimeout) * time.Second)
}
