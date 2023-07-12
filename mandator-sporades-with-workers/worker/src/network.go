package src

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"with-workers/common"
	"with-workers/proto"
)

/*
	Given an int32 id, connect to that node
	nodeType should be one of client, worker and replica
*/

func (wr *Worker) ConnectToNode(id int32, address string, nodeType string) {

	common.Debug("Connecting to "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)

	var b [4]byte
	bs := b[:4]

	for true {
		conn, err := net.Dial("tcp", address)
		if err == nil {

			if nodeType == "client" {
				wr.outgoingClientWriters[wr.clientArrayIndex[id]] = bufio.NewWriter(conn)
				binary.LittleEndian.PutUint16(bs, uint16(wr.name))
				_, err := conn.Write(bs)
				if err != nil {
					common.Debug("Error while connecting to client "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)
					panic(err)
				}
			} else if nodeType == "worker" {
				wr.outgoingWorkerWriters[wr.workerArrayIndex[id]] = bufio.NewWriter(conn)
				binary.LittleEndian.PutUint16(bs, uint16(wr.name))
				_, err := conn.Write(bs)
				if err != nil {
					common.Debug("Error while connecting to worker "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)
					panic(err)
				}
			} else if nodeType == "replica" {
				wr.outgoingDefaultReplicaWriter = bufio.NewWriter(conn)
				binary.LittleEndian.PutUint16(bs, uint16(wr.name))
				_, err := conn.Write(bs)
				if err != nil {
					common.Debug("Error while connecting to default replica "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)
					panic(err)
				}
			} else {
				panic("Unknown node id")
			}
			common.Debug("Established outgoing connection to "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)
			break
		}
	}
}

/*
	connect to all workers and default replica on bootstrap
*/

func (wr *Worker) ConnectBootStrap() {
	for name, address := range wr.workerAddrList {
		wr.ConnectToNode(name, address, "worker")
	}
	wr.ConnectToNode(wr.defaultReplicaName, wr.defaultReplicaAddr, "replica")
}

/*
	listen to a given connection reader. Upon receiving any message, put it into the central buffer
*/

func (wr *Worker) connectionListener(reader *bufio.Reader, id int32) {

	var msgType uint8
	var err error = nil

	for true {

		if msgType, err = reader.ReadByte(); err != nil {
			common.Debug("Error while reading message code: connection broken from "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)
			common.Debug("Error code "+err.Error(), 0, wr.debugLevel, wr.debugOn)
			return
		}

		if rpair, present := wr.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				common.Debug("Error while unmarshalling from "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)
				return
			}
			wr.incomingChan <- &common.RPCPair{
				Code: msgType,
				Obj:  obj,
			}
			common.Debug("Pushed a message from "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)

		} else {
			common.Debug("Error received unknown message type from "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)
		}
	}
}

/*
	Listen on the worker port for new connections from default replica, all workers and clients
*/

func (wr *Worker) WaitForConnections() {
	go func() {
		var b [4]byte
		bs := b[:4]
		Listener, _ := net.Listen("tcp", wr.listenAddress)
		common.Debug("Listening to incoming connection in "+wr.listenAddress, 0, wr.debugLevel, wr.debugOn)

		for true {
			conn, err := Listener.Accept()
			if err != nil {
				common.Debug("Socket accept error", 0, wr.debugLevel, wr.debugOn)
				panic(err)
			}
			if _, err := io.ReadFull(conn, bs); err != nil {
				common.Debug("Connection read error when establishing incoming connections", 0, wr.debugLevel, wr.debugOn)
				panic(err)
			}
			id := int32(binary.LittleEndian.Uint16(bs))
			common.Debug("Received incoming connection from "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)
			nodeType := wr.getNodeType(id)
			if nodeType == "worker" {
				wr.incomingWorkerReaders[wr.workerArrayIndex[id]] = bufio.NewReader(conn)
				go wr.connectionListener(wr.incomingWorkerReaders[wr.workerArrayIndex[id]], id)
				common.Debug("Started listening to worker "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)
			} else if nodeType == "replica" {
				wr.incomingDefaultReplicaReader = bufio.NewReader(conn)
				go wr.connectionListener(wr.incomingDefaultReplicaReader, id)
				common.Debug("Started listening to default replica "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)
			} else if nodeType == "client" {
				wr.incomingClientReaders[wr.clientArrayIndex[id]] = bufio.NewReader(conn)
				go wr.connectionListener(wr.incomingClientReaders[wr.clientArrayIndex[id]], id)
				common.Debug("Started listening to client "+strconv.Itoa(int(id)), 0, wr.debugLevel, wr.debugOn)
				// dial back the client
				wr.ConnectToNode(id, wr.clientAddrList[id], "client")
			}
		}
	}()
}

/*
	This is the main execution thread that listens to all the incoming messages
	It listens to incoming messages from the incomingChan, and invokes the appropriate handler depending on the message type
*/

func (wr *Worker) Run() {

	for true {

		common.Debug("Checking channel..", 0, wr.debugLevel, wr.debugOn)
		replicaMessage := <-wr.incomingChan
		common.Debug("Received  message ", 0, wr.debugLevel, wr.debugOn)

		switch replicaMessage.Code {
		case wr.messageCodes.ClientBatchRpc:
			clientBatch := replicaMessage.Obj.(*proto.ClientBatch)
			common.Debug("Client batch from "+fmt.Sprintf("%#v", clientBatch.Sender), 0, wr.debugLevel, wr.debugOn)
			wr.handleClientBatch(clientBatch)
			break

		case wr.messageCodes.StatusRPC:
			statusMessage := replicaMessage.Obj.(*proto.Status)
			common.Debug("Status message from "+fmt.Sprintf("%#v", statusMessage.Sender), 0, wr.debugLevel, wr.debugOn)
			wr.handleStatus(statusMessage)
			break

		case wr.messageCodes.MemPoolMiniRPC:
			memPoolMiniMessage := replicaMessage.Obj.(*proto.MemPoolMini)
			common.Debug("MemPoolMini message from "+fmt.Sprintf("%#v", memPoolMiniMessage.Sender), 0, wr.debugLevel, wr.debugOn)
			wr.handleMemPoolMini(memPoolMiniMessage)
			break
		}
	}

}

/*
	Write a message to the wire, first the message type is written and then the actual message
*/

func (wr *Worker) internalSendMessage(peer int32, rpcPair *common.RPCPair) {
	peerType := wr.getNodeType(peer)
	if peerType == "client" {
		w := wr.outgoingClientWriters[wr.clientArrayIndex[peer]]
		if w == nil {
			return
		}
		wr.outgoingClientWriterMutexs[wr.clientArrayIndex[peer]].Lock()
		err := w.WriteByte(rpcPair.Code)
		if err != nil {
			common.Debug("Error writing message code byte:"+err.Error(), 0, wr.debugLevel, wr.debugOn)
			wr.outgoingClientWriterMutexs[wr.clientArrayIndex[peer]].Unlock()
			return
		}
		err = rpcPair.Obj.Marshal(w)
		if err != nil {
			common.Debug("Error while marshalling:"+err.Error(), 0, wr.debugLevel, wr.debugOn)
			wr.outgoingClientWriterMutexs[wr.clientArrayIndex[peer]].Unlock()
			return
		}
		err = w.Flush()
		if err != nil {
			common.Debug("Error while flushing:"+err.Error(), 0, wr.debugLevel, wr.debugOn)
			wr.outgoingClientWriterMutexs[wr.clientArrayIndex[peer]].Unlock()
			return
		}
		wr.outgoingClientWriterMutexs[wr.clientArrayIndex[peer]].Unlock()
		common.Debug("Internal sent message to "+strconv.Itoa(int(peer)), 0, wr.debugLevel, wr.debugOn)

	} else if peerType == "replica" {
		w := wr.outgoingDefaultReplicaWriter
		if w == nil {
			return
		}
		wr.outgoingDefaultReplicaWriterMutex.Lock()
		err := w.WriteByte(rpcPair.Code)
		if err != nil {
			common.Debug("Error writing message code byte:"+err.Error(), 0, wr.debugLevel, wr.debugOn)
			wr.outgoingDefaultReplicaWriterMutex.Unlock()
			return
		}
		err = rpcPair.Obj.Marshal(w)
		if err != nil {
			common.Debug("Error while marshalling:"+err.Error(), 0, wr.debugLevel, wr.debugOn)
			wr.outgoingDefaultReplicaWriterMutex.Unlock()
			return
		}
		err = w.Flush()
		if err != nil {
			common.Debug("Error while flushing:"+err.Error(), 0, wr.debugLevel, wr.debugOn)
			wr.outgoingDefaultReplicaWriterMutex.Unlock()
			return
		}
		wr.outgoingDefaultReplicaWriterMutex.Unlock()
		common.Debug("Internal sent message to "+strconv.Itoa(int(peer)), 0, wr.debugLevel, wr.debugOn)

	} else if peerType == "worker" {
		w := wr.outgoingWorkerWriters[wr.workerArrayIndex[peer]]
		if w == nil {
			return
		}
		wr.outgoingWorkerWriterMutexs[wr.workerArrayIndex[peer]].Lock()
		err := w.WriteByte(rpcPair.Code)
		if err != nil {
			common.Debug("Error writing message code byte:"+err.Error(), 0, wr.debugLevel, wr.debugOn)
			wr.outgoingWorkerWriterMutexs[wr.workerArrayIndex[peer]].Unlock()
			return
		}
		err = rpcPair.Obj.Marshal(w)
		if err != nil {
			common.Debug("Error while marshalling:"+err.Error(), 0, wr.debugLevel, wr.debugOn)
			wr.outgoingWorkerWriterMutexs[wr.workerArrayIndex[peer]].Unlock()
			return
		}
		err = w.Flush()
		if err != nil {
			common.Debug("Error while flushing:"+err.Error(), 0, wr.debugLevel, wr.debugOn)
			wr.outgoingWorkerWriterMutexs[wr.workerArrayIndex[peer]].Unlock()
			return
		}
		wr.outgoingWorkerWriterMutexs[wr.workerArrayIndex[peer]].Unlock()
		common.Debug("Internal sent message to "+strconv.Itoa(int(peer)), 0, wr.debugLevel, wr.debugOn)
	} else {
		panic("Unknown id")
	}
}

/*
	A set of threads that manages outgoing messages
*/

func (wr *Worker) StartOutgoingLinks() {
	for i := 0; i < numOutgoingThreads; i++ {
		go func() {
			for true {
				outgoingMessage := <-wr.outgoingMessageChan
				wr.internalSendMessage(outgoingMessage.Peer, outgoingMessage.RpcPair)
				common.Debug("Invoked internal sent to replica "+strconv.Itoa(int(outgoingMessage.Peer)), 0, wr.debugLevel, wr.debugOn)
			}
		}()
	}
}

/*
	Add a new out-going message to the outgoing channel
*/

func (wr *Worker) sendMessage(peer int32, rpcPair common.RPCPair) {
	wr.outgoingMessageChan <- &common.OutgoingRPC{
		RpcPair: &rpcPair,
		Peer:    peer,
	}
	common.Debug("Added RPC pair to outgoing channel to peer "+strconv.Itoa(int(peer)), 0, wr.debugLevel, wr.debugOn)
}
