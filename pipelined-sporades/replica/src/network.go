package src

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"pipelined-sporades/common"
	"pipelined-sporades/proto"
	"strconv"
	"time"
)

/*
	listen on the replica port for new connections from clients, and all replicas
*/

func (rp *Replica) WaitForConnections() {
	go func() {
		var b [4]byte
		bs := b[:4]
		Listener, _ := net.Listen("tcp", rp.listenAddress)
		if rp.debugOn {
			rp.debug("Listening to incoming connections on "+rp.listenAddress, 0)
		}
		for true {
			conn, err := Listener.Accept()
			if err != nil {
				panic(err.Error())
			}
			if _, err := io.ReadFull(conn, bs); err != nil {
				panic(err.Error())
			}
			id := int32(binary.LittleEndian.Uint16(bs))
			if rp.debugOn {
				rp.debug("Received incoming connection from "+strconv.Itoa(int(id)), 0)
			}
			nodeType := rp.getNodeType(id)
			if nodeType == "client" {
				rp.incomingClientReaders[id] = bufio.NewReader(conn)
				go rp.connectionListener(rp.incomingClientReaders[id], id)
				if rp.debugOn {
					rp.debug("Started listening to client "+strconv.Itoa(int(id)), 0)
				}
				rp.ConnectToNode(id, rp.clientAddrList[id], "client")

			} else if nodeType == "replica" {
				rp.incomingReplicaReaders[id] = bufio.NewReader(conn)
				go rp.connectionListener(rp.incomingReplicaReaders[id], id)
				if rp.debugOn {
					rp.debug("Started listening to replica "+strconv.Itoa(int(id)), 0)
				}
			} else {
				panic("should not happen")
			}
		}
	}()
}

/*
	listen to a given connection reader. Upon receiving any message, put it into the central incoming buffer
*/

func (rp *Replica) connectionListener(reader *bufio.Reader, id int32) {

	var msgType uint8
	var err error = nil

	for true {

		if msgType, err = reader.ReadByte(); err != nil {
			if rp.debugOn {
				rp.debug("Error while reading message code: connection broken from "+strconv.Itoa(int(id))+fmt.Sprintf(" %v", err), 0)
			}
			return
		}

		if rpair, present := rp.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				if rp.debugOn {
					rp.debug("Error while unmarshalling from "+strconv.Itoa(int(id))+fmt.Sprintf(" %v", err), 0)
				}
				return
			}
			rp.incomingChan <- &common.RPCPair{
				Code: msgType,
				Obj:  obj,
			}
			if rp.debugOn {
				rp.debug("Pushed a message from "+strconv.Itoa(int(id)), -1)
			}
		} else {
			if rp.debugOn {
				rp.debug("Error received unknown message type from "+strconv.Itoa(int(id)), 0)
			}
			return
		}
	}
}

/*
	given an int32 id, connect to that node
	nodeType should be one of client or replica
*/

func (rp *Replica) ConnectToNode(id int32, address string, nodeType string) {

	if rp.debugOn {
		rp.debug("Connecting to "+strconv.Itoa(int(id)), 0)
	}

	var b [4]byte
	bs := b[:4]
	counter := 0
	for counter < 1000000 {
		counter++
		conn, err := net.Dial("tcp", address)
		if err == nil {

			if nodeType == "client" {
				rp.outgoingClientWriters[id] = bufio.NewWriter(conn)
				binary.LittleEndian.PutUint16(bs, uint16(rp.name))
				_, err := conn.Write(bs)
				if err != nil {
					panic("Error while connecting to client " + strconv.Itoa(int(id)) + err.Error())
				}
			} else if nodeType == "replica" {
				rp.outgoingReplicaWriters[id] = bufio.NewWriter(conn)
				binary.LittleEndian.PutUint16(bs, uint16(rp.name))
				_, err := conn.Write(bs)
				if err != nil {
					panic("Error while connecting to replica " + strconv.Itoa(int(id)) + err.Error())
				}
			} else {
				panic("Unknown node type")
			}
			if rp.debugOn {
				rp.debug("Established outgoing connection to "+strconv.Itoa(int(id)), 0)
			}
			break
		} else {
			if counter == 1000000 {
				panic(fmt.Sprintf("node %v cannot be reached to establish a connection "+err.Error(), id))
			}
		}
	}
}

/*
	Connect to all replicas on bootstrap
*/

func (rp *Replica) ConnectBootStrap() {

	for name, address := range rp.replicaAddrList {
		rp.ConnectToNode(name, address, "replica")
	}
}

/*
	this is the main execution thread that listens to all the incoming messages
	It listens to incoming messages from the incomingChan, and invokes the appropriate handler depending on the message type
*/

func (rp *Replica) Run() {
	for true {
		select {
		case replicaMessage := <-rp.incomingChan:

			if rp.debugOn {
				rp.debug("Received replica message", -1)
			}

			switch replicaMessage.Code {

			case rp.messageCodes.StatusRPC:
				statusMessage := replicaMessage.Obj.(*proto.Status)
				if rp.debugOn {
					rp.debug("Status message from "+fmt.Sprintf("%#v", statusMessage.Sender), -1)
				}
				rp.handleStatus(statusMessage)
				break

			case rp.messageCodes.ClientBatchRpc:
				clientBatch := replicaMessage.Obj.(*proto.ClientBatch)
				if rp.debugOn {
					rp.debug("Client batch message from "+fmt.Sprintf("%#v", clientBatch.Sender), 0)
				}
				rp.handleClientBatch(clientBatch)
				break

			case rp.messageCodes.SporadesConsensus:
				sporadesConsensusMessage := replicaMessage.Obj.(*proto.Pipelined_Sporades)
				if rp.debugOn {
					rp.debug("Sporades consensus message from "+fmt.Sprintf("%#v", sporadesConsensusMessage.Sender), -1)
				}
				rp.handleSporadesConsensus(sporadesConsensusMessage)
				break

			}
			break

		}
	}
}

/*
	Write a message to the wire, first the message type is written and then the actual message
*/

func (rp *Replica) internalSendMessage(peer int32, rpcPair *common.RPCPair) {
	peerType := rp.getNodeType(peer)
	if peerType == "replica" {

		if rpcPair.Code == rp.messageCodes.SporadesConsensus {
			time.Sleep(time.Duration(rp.asyncbatchTime) * time.Millisecond)
		}

		w := rp.outgoingReplicaWriters[peer]
		if w == nil {
			panic("replica not found")
		}
		rp.outgoingReplicaWriterMutexs[peer].Lock()
		err := w.WriteByte(rpcPair.Code)
		if err != nil {
			if rp.debugOn {
				rp.debug("Error writing message code byte:"+err.Error(), 0)
			}
			rp.outgoingReplicaWriterMutexs[peer].Unlock()
			return
		}
		err = rpcPair.Obj.Marshal(w)
		if err != nil {
			if rp.debugOn {
				rp.debug("Error while marshalling:"+err.Error(), 0)
			}
			rp.outgoingReplicaWriterMutexs[peer].Unlock()
			return
		}
		err = w.Flush()
		if err != nil {
			if rp.debugOn {
				rp.debug("Error while flushing:"+err.Error(), 0)
			}
			rp.outgoingReplicaWriterMutexs[peer].Unlock()
			return
		}
		rp.outgoingReplicaWriterMutexs[peer].Unlock()
		if rp.debugOn {
			rp.debug("Internal sent message to "+strconv.Itoa(int(peer)), -1)
		}
	} else if peerType == "client" {
		w := rp.outgoingClientWriters[peer]
		if w == nil {
			panic("client not found")
		}
		rp.outgoingClientWriterMutexs[peer].Lock()
		err := w.WriteByte(rpcPair.Code)
		if err != nil {
			if rp.debugOn {
				rp.debug("Error writing message code byte:"+err.Error(), 1)
			}
			rp.outgoingClientWriterMutexs[peer].Unlock()
			return
		}
		err = rpcPair.Obj.Marshal(w)
		if err != nil {
			if rp.debugOn {
				rp.debug("Error while marshalling:"+err.Error(), 1)
			}
			rp.outgoingClientWriterMutexs[peer].Unlock()
			return
		}
		err = w.Flush()
		if err != nil {
			if rp.debugOn {
				rp.debug("Error while flushing:"+err.Error(), 1)
			}
			rp.outgoingClientWriterMutexs[peer].Unlock()
			return
		}
		rp.outgoingClientWriterMutexs[peer].Unlock()
		if rp.debugOn {
			rp.debug("Internal sent message to "+strconv.Itoa(int(peer)), 1)
		}
	} else {
		panic("Unknown id from node name " + strconv.Itoa(int(peer)))
	}
}

/*
	A set of threads that manages outgoing messages
*/

func (rp *Replica) StartOutgoingLinks() {
	// for clients, we do not care about the order
	for i := 0; i < 100; i++ {
		go func() {
			for true {
				outgoingMessage := <-rp.outgoingClientMessageChan
				rp.internalSendMessage(outgoingMessage.Peer, outgoingMessage.RpcPair)
				if rp.debugOn {
					rp.debug("Invoked internal sent to "+strconv.Itoa(int(outgoingMessage.Peer)), 1)
				}
			}
		}()
	}
	for i := 0; i < rp.numReplicas; i++ {
		go func(peer int) {
			for true {
				outgoingMessage := <-rp.outgoingReplicaMessageChans[peer]
				rp.internalSendMessage(outgoingMessage.Peer, outgoingMessage.RpcPair)
				if rp.debugOn {
					rp.debug("Invoked internal sent to "+strconv.Itoa(int(outgoingMessage.Peer)), -1)
				}
			}
		}(i)
	}

}

/*
	Add a new out-going message to the correct outgoing channel
*/

func (rp *Replica) sendMessage(peer int32, rpcPair common.RPCPair) {

	peerType := rp.getNodeType(peer)
	if peerType == "replica" {
		rp.outgoingReplicaMessageChans[peer-1] <- &common.OutgoingRPC{
			RpcPair: &rpcPair,
			Peer:    peer,
		}
		if rp.debugOn {
			rp.debug("Added RPC pair to outgoing channel to peer "+strconv.Itoa(int(peer)), -1)
		}
	} else if peerType == "client" {
		rp.outgoingClientMessageChan <- &common.OutgoingRPC{
			RpcPair: &rpcPair,
			Peer:    peer,
		}
		if rp.debugOn {
			rp.debug("Added RPC pair to outgoing channel to client "+strconv.Itoa(int(peer)), 1)
		}
	} else {
		panic("unknown peer type")
	}

}
