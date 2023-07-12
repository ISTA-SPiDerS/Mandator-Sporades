package src

import (
	"async-consensus/common"
	"async-consensus/proto"
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

/*
	Each client sends connection requests to all workers
*/

func (cl *Client) ConnectToWorkers() {

	common.Debug("Connecting to workers", 0, cl.debugLevel, cl.debugOn)

	var b [4]byte
	bs := b[:4]

	//connect to workers
	for name, address := range cl.workerAddrList {
		for true {
			conn, err := net.Dial("tcp", address)
			if err == nil {
				cl.outgoingWorkerWriters[cl.workerArrayIndex[name]] = bufio.NewWriter(conn)
				binary.LittleEndian.PutUint16(bs, uint16(cl.clientName))
				_, err := conn.Write(bs)
				if err != nil {
					common.Debug("Error while connecting to worker "+strconv.Itoa(int(name)), 0, cl.debugLevel, cl.debugOn)
					panic(err)
				}
				common.Debug("Established outgoing connection to "+strconv.Itoa(int(name)), 0, cl.debugLevel, cl.debugOn)
				break
			}
		}
	}
	common.Debug("Established outgoing connections to all workers", 0, cl.debugLevel, cl.debugOn)
}

/*
	Listen on the client port for new connections from workers
*/

func (cl *Client) WaitForConnections() {
	go func() {
		var b [4]byte
		bs := b[:4]
		Listener, _ := net.Listen("tcp", cl.clientListenAddress)
		common.Debug("Listening to incoming connection in "+cl.clientListenAddress, 0, cl.debugLevel, cl.debugOn)

		for true {
			conn, err := Listener.Accept()
			if err != nil {
				common.Debug("Socket accept error", 0, cl.debugLevel, cl.debugOn)
				panic(err)
			}
			if _, err := io.ReadFull(conn, bs); err != nil {
				common.Debug("Connection read error when establishing incoming connections", 0, cl.debugLevel, cl.debugOn)
				panic(err)
			}
			id := int32(binary.LittleEndian.Uint16(bs))
			common.Debug("Received incoming connection from "+strconv.Itoa(int(id)), 0, cl.debugLevel, cl.debugOn)

			cl.incomingWorkerReaders[cl.workerArrayIndex[id]] = bufio.NewReader(conn)
			go cl.connectionListener(cl.incomingWorkerReaders[cl.workerArrayIndex[id]], id)
			common.Debug("Started listening to "+strconv.Itoa(int(id)), 0, cl.debugLevel, cl.debugOn)

		}
	}()
}

/*
	listen to a given connection. Upon receiving any message, put it into the central buffer
*/

func (cl *Client) connectionListener(reader *bufio.Reader, id int32) {

	var msgType uint8
	var err error = nil

	for true {

		if msgType, err = reader.ReadByte(); err != nil {
			common.Debug("Error while reading message code: connection broken from "+strconv.Itoa(int(id)), 0, cl.debugLevel, cl.debugOn)
			return
		}

		if rpair, present := cl.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				common.Debug("Error while unmarshalling from "+strconv.Itoa(int(id)), 0, cl.debugLevel, cl.debugOn)
				return
			}
			cl.incomingChan <- &common.RPCPair{
				Code: msgType,
				Obj:  obj,
			}
			common.Debug("Pushed a message from "+strconv.Itoa(int(id)), 0, cl.debugLevel, cl.debugOn)

		} else {
			common.Debug("Error received unknown message type from "+strconv.Itoa(int(id)), 0, cl.debugLevel, cl.debugOn)
		}
	}
}

/*
	This is the main execution thread that listens to all the incoming messages
	It listens to incoming messages from the incomingChan, and invokes the appropriate handler depending on the message type
*/

func (cl *Client) Run() {
	go func() {
		for true {

			common.Debug("Checking channel..", 0, cl.debugLevel, cl.debugOn)
			replicaMessage := <-cl.incomingChan
			common.Debug("Received message", 0, cl.debugLevel, cl.debugOn)

			switch replicaMessage.Code {
			case cl.messageCodes.ClientBatchRpc:
				clientResponseBatch := replicaMessage.Obj.(*proto.ClientBatch)
				common.Debug("Client response batch from "+fmt.Sprintf("%#v", clientResponseBatch.Sender), 0, cl.debugLevel, cl.debugOn)
				cl.handleClientResponseBatch(clientResponseBatch)
				break

			case cl.messageCodes.StatusRPC:
				clientStatusResponse := replicaMessage.Obj.(*proto.Status)
				common.Debug("Client status from "+fmt.Sprintf("%#v", clientStatusResponse.Sender), 0, cl.debugLevel, cl.debugOn)
				cl.handleClientStatusResponse(clientStatusResponse)
				break
			}
		}
	}()
}

/*
	Write a message to the wire, first the message type is written and then the actual message
*/

func (cl *Client) internalSendMessage(peer int32, rpcPair *common.RPCPair) {
	w := cl.outgoingWorkerWriters[cl.workerArrayIndex[peer]]
	cl.outgoingWorkerWriterMutexs[cl.workerArrayIndex[peer]].Lock()
	err := w.WriteByte(rpcPair.Code)
	if err != nil {
		common.Debug("Error writing message code byte:"+err.Error(), 0, cl.debugLevel, cl.debugOn)
		cl.outgoingWorkerWriterMutexs[cl.workerArrayIndex[peer]].Unlock()
		return
	}
	err = rpcPair.Obj.Marshal(w)
	if err != nil {
		common.Debug("Error while marshalling:"+err.Error(), 0, cl.debugLevel, cl.debugOn)
		cl.outgoingWorkerWriterMutexs[cl.workerArrayIndex[peer]].Unlock()
		return
	}
	err = w.Flush()
	if err != nil {
		common.Debug("Error while flushing:"+err.Error(), 0, cl.debugLevel, cl.debugOn)
		cl.outgoingWorkerWriterMutexs[cl.workerArrayIndex[peer]].Unlock()
		return
	}
	cl.outgoingWorkerWriterMutexs[cl.workerArrayIndex[peer]].Unlock()
	common.Debug("Internal sent message to "+strconv.Itoa(int(peer)), 0, cl.debugLevel, cl.debugOn)
}

/*
	A set of threads that manages outgoing messages
*/

func (cl *Client) StartOutgoingLinks() {
	for i := 0; i < numOutgoingThreads; i++ {
		go func() {
			for true {
				outgoingMessage := <-cl.outgoingMessageChan
				cl.internalSendMessage(outgoingMessage.Peer, outgoingMessage.RpcPair)
				common.Debug("Invoked internal sent to replica "+strconv.Itoa(int(outgoingMessage.Peer)), 0, cl.debugLevel, cl.debugOn)
			}
		}()
	}
}

/*
	Add a new out-going message to the outgoing channel
*/

func (cl *Client) sendMessage(peer int32, rpcPair common.RPCPair) {
	cl.outgoingMessageChan <- &common.OutgoingRPC{
		RpcPair: &rpcPair,
		Peer:    peer,
	}
	common.Debug("Added RPC pair to outgoing channel to peer "+strconv.Itoa(int(peer)), 0, cl.debugLevel, cl.debugOn)
}
