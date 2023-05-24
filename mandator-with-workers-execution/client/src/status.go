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
	fmt.Print("Status response from " + strconv.Itoa(int(response.Sender)) + " \n")
}

/*
	Send a status request to all the workers
*/

func (cl *Client) SendStatus(operationType int) {
	common.Debug("Sending status request to all workers", 0, cl.debugLevel, cl.debugOn)

	for name, _ := range cl.workerAddrList {

		statusRequest := proto.Status{
			Sender:   cl.clientName,
			Receiver: name,
			UniqueId: "",
			Type:     int32(operationType),
			Note:     "",
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
