package src

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
	"with-workers/common"
	"with-workers/proto"
)

/*
	Given an int32 id, connect to that node
	nodeType should be one of worker and replica
*/

func (rp *Replica) ConnectToNode(id int32, address string, nodeType string) {

	common.Debug("Connecting to "+strconv.Itoa(int(id)), 0, rp.debugLevel, rp.debugOn)

	var b [4]byte
	bs := b[:4]

	for true {
		conn, err := net.Dial("tcp", address)
		if err == nil {

			if nodeType == "worker" {
				rp.outgoingWorkerWriters[rp.workerArrayIndex[id]] = bufio.NewWriter(conn)
				binary.LittleEndian.PutUint16(bs, uint16(rp.name))
				_, err := conn.Write(bs)
				if err != nil {
					common.Debug("Error while connecting to worker "+strconv.Itoa(int(id)), 0, rp.debugLevel, rp.debugOn)
					panic(err)
				}
			} else if nodeType == "replica" {
				rp.outgoingReplicaWriters[rp.replicaArrayIndex[id]] = bufio.NewWriter(conn)
				binary.LittleEndian.PutUint16(bs, uint16(rp.name))
				_, err := conn.Write(bs)
				if err != nil {
					common.Debug("Error while connecting to replica "+strconv.Itoa(int(id)), 0, rp.debugLevel, rp.debugOn)
					panic(err)
				}
			} else {
				panic("Unknown node id")
			}
			common.Debug("Established outgoing connection to "+strconv.Itoa(int(id)), 0, rp.debugLevel, rp.debugOn)
			break
		}
	}
}

/*
	Connect to all designated workers and all replicas on bootstrap
*/

func (rp *Replica) ConnectBootStrap() {

	for name, address := range rp.workerAddrList {
		rp.ConnectToNode(name, address, "worker")
	}

	for name, address := range rp.replicaAddrList {
		rp.ConnectToNode(name, address, "replica")
	}
}

/*
	listen to a given connection reader. Upon receiving any message, put it into the central incoming buffer
*/

func (rp *Replica) connectionListener(reader *bufio.Reader, id int32) {

	var msgType uint8
	var err error = nil

	for true {

		if msgType, err = reader.ReadByte(); err != nil {
			common.Debug("Error while reading message code: connection broken from "+strconv.Itoa(int(id))+fmt.Sprintf(" %v", err), 0, rp.debugLevel, rp.debugOn)
			return
		}

		if rpair, present := rp.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				common.Debug("Error while unmarshalling from "+strconv.Itoa(int(id))+fmt.Sprintf(" %v", err), 0, rp.debugLevel, rp.debugOn)
				return
			}
			rp.incomingChan <- &common.RPCPair{
				Code: msgType,
				Obj:  obj,
			}
			common.Debug("Pushed a message from "+strconv.Itoa(int(id)), 0, rp.debugLevel, rp.debugOn)

		} else {
			common.Debug("Error received unknown message type from "+strconv.Itoa(int(id)), 0, rp.debugLevel, rp.debugOn)
		}
	}
}

/*
	Listen on the replica port for new connections from default workers, and all replicas
*/

func (rp *Replica) WaitForConnections() {
	go func() {
		var b [4]byte
		bs := b[:4]
		Listener, _ := net.Listen("tcp", rp.listenAddress)
		common.Debug("Listening to incoming connection in "+rp.listenAddress, 0, rp.debugLevel, rp.debugOn)

		for true {
			conn, err := Listener.Accept()
			if err != nil {
				common.Debug("Socket accept error", 0, rp.debugLevel, rp.debugOn)
				panic(err)
			}
			if _, err := io.ReadFull(conn, bs); err != nil {
				common.Debug("Connection read error when establishing incoming connections", 0, rp.debugLevel, rp.debugOn)
				panic(err)
			}
			id := int32(binary.LittleEndian.Uint16(bs))
			common.Debug("Received incoming connection from "+strconv.Itoa(int(id)), 0, rp.debugLevel, rp.debugOn)
			nodeType := rp.getNodeType(id)
			if nodeType == "worker" {
				rp.incomingWorkerReaders[rp.workerArrayIndex[id]] = bufio.NewReader(conn)
				go rp.connectionListener(rp.incomingWorkerReaders[rp.workerArrayIndex[id]], id)
				common.Debug("Started listening to worker "+strconv.Itoa(int(id)), 0, rp.debugLevel, rp.debugOn)
			} else if nodeType == "replica" {
				rp.incomingReplicaReaders[rp.replicaArrayIndex[id]] = bufio.NewReader(conn)
				go rp.connectionListener(rp.incomingReplicaReaders[rp.replicaArrayIndex[id]], id)
				common.Debug("Started listening to replica "+strconv.Itoa(int(id)), 0, rp.debugLevel, rp.debugOn)
			}
		}
	}()
}

/*
	This is the main execution thread that listens to all the incoming messages
	It listens to incoming messages from the incomingChan, and invokes the appropriate handler depending on the message type
*/

func (rp *Replica) Run() {

	for true {

		common.Debug("Checking channel..", 0, rp.debugLevel, rp.debugOn)
		replicaMessage := <-rp.incomingChan
		common.Debug("Received replica message", 0, rp.debugLevel, rp.debugOn)

		switch replicaMessage.Code {

		case rp.messageCodes.StatusRPC:
			statusMessage := replicaMessage.Obj.(*proto.Status)
			common.Debug("Status message from "+fmt.Sprintf("%#v", statusMessage.Sender), 0, rp.debugLevel, rp.debugOn)
			rp.handleStatus(statusMessage)
			break

		case rp.messageCodes.MemPoolMiniRPC:
			memPoolMiniMessage := replicaMessage.Obj.(*proto.MemPoolMini)
			common.Debug("MemPoolMini message from "+fmt.Sprintf("%#v", memPoolMiniMessage.Sender), 0, rp.debugLevel, rp.debugOn)
			rp.handleMemPoolMini(memPoolMiniMessage)
			break

		case rp.messageCodes.MemPoolRPC:
			memPoolMessage := replicaMessage.Obj.(*proto.MemPool)
			common.Debug("Mem pool message from "+fmt.Sprintf("%#v", memPoolMessage.Sender), 0, rp.debugLevel, rp.debugOn)
			rp.handleMemPool(memPoolMessage)
			break

		case rp.messageCodes.AsyncConsensus:
			asyncConsensusMessage := replicaMessage.Obj.(*proto.AsyncConsensus)
			common.Debug("Async consensus message from "+fmt.Sprintf("%#v", asyncConsensusMessage.Sender), 0, rp.debugLevel, rp.debugOn)
			rp.handleAsyncConsensus(asyncConsensusMessage)
			break

		case rp.messageCodes.PaxosConsensus:
			paxosConsensusMessage := replicaMessage.Obj.(*proto.PaxosConsensus)
			common.Debug("Paxos consensus message from "+fmt.Sprintf("%#v", paxosConsensusMessage.Sender), 0, rp.debugLevel, rp.debugOn)
			rp.handlePaxosConsensus(paxosConsensusMessage)
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
		// we delay the consensus message for 20 milli second, to avoid the curse of small consensus message sizes
		if rpcPair.Code == rp.messageCodes.AsyncConsensus || rpcPair.Code == rp.messageCodes.PaxosConsensus {
			time.Sleep(time.Duration(rp.asyncBatchTime) * time.Millisecond)
		}
		w := rp.outgoingReplicaWriters[rp.replicaArrayIndex[peer]]
		if w == nil {
			panic("replica not found")
		}
		rp.outgoingReplicaWriterMutexs[rp.replicaArrayIndex[peer]].Lock()
		err := w.WriteByte(rpcPair.Code)
		if err != nil {
			common.Debug("Error writing message code byte:"+err.Error(), 0, rp.debugLevel, rp.debugOn)
			rp.outgoingReplicaWriterMutexs[rp.replicaArrayIndex[peer]].Unlock()
			return
		}
		err = rpcPair.Obj.Marshal(w)
		if err != nil {
			common.Debug("Error while marshalling:"+err.Error(), 0, rp.debugLevel, rp.debugOn)
			rp.outgoingReplicaWriterMutexs[rp.replicaArrayIndex[peer]].Unlock()
			return
		}
		err = w.Flush()
		if err != nil {
			common.Debug("Error while flushing:"+err.Error(), 0, rp.debugLevel, rp.debugOn)
			rp.outgoingReplicaWriterMutexs[rp.replicaArrayIndex[peer]].Unlock()
			return
		}
		rp.outgoingReplicaWriterMutexs[rp.replicaArrayIndex[peer]].Unlock()
		common.Debug("Internal sent message to "+strconv.Itoa(int(peer)), 0, rp.debugLevel, rp.debugOn)
	} else if peerType == "worker" {
		w := rp.outgoingWorkerWriters[rp.workerArrayIndex[peer]]
		if w == nil {
			panic("worker not found")
		}
		rp.outgoingWorkerWriterMutexs[rp.workerArrayIndex[peer]].Lock()
		err := w.WriteByte(rpcPair.Code)
		if err != nil {
			common.Debug("Error writing message code byte:"+err.Error(), 0, rp.debugLevel, rp.debugOn)
			rp.outgoingWorkerWriterMutexs[rp.workerArrayIndex[peer]].Unlock()
			return
		}
		err = rpcPair.Obj.Marshal(w)
		if err != nil {
			common.Debug("Error while marshalling:"+err.Error(), 0, rp.debugLevel, rp.debugOn)
			rp.outgoingWorkerWriterMutexs[rp.workerArrayIndex[peer]].Unlock()
			return
		}
		err = w.Flush()
		if err != nil {
			common.Debug("Error while flushing:"+err.Error(), 0, rp.debugLevel, rp.debugOn)
			rp.outgoingWorkerWriterMutexs[rp.workerArrayIndex[peer]].Unlock()
			return
		}
		rp.outgoingWorkerWriterMutexs[rp.workerArrayIndex[peer]].Unlock()
		common.Debug("Internal sent message to "+strconv.Itoa(int(peer)), 0, rp.debugLevel, rp.debugOn)
	} else {
		panic("Unknown id from node name " + strconv.Itoa(int(peer)))
	}
}

/*
	A set of threads that manages outgoing messages
*/

func (rp *Replica) StartOutgoingLinks() {
	for i := 0; i < numOutgoingThreads; i++ {
		go func() {
			for true {
				outgoingMessage := <-rp.outgoingMessageChan
				rp.internalSendMessage(outgoingMessage.Peer, outgoingMessage.RpcPair)
				common.Debug("Invoked internal sent to replica "+strconv.Itoa(int(outgoingMessage.Peer)), 0, rp.debugLevel, rp.debugOn)
			}
		}()
	}
}

/*
	Add a new out-going message to the outgoing channel
*/

func (rp *Replica) sendMessage(peer int32, rpcPair common.RPCPair) {
	rp.outgoingMessageChan <- &common.OutgoingRPC{
		RpcPair: &rpcPair,
		Peer:    peer,
	}
	common.Debug("Added RPC pair to outgoing channel to peer "+strconv.Itoa(int(peer)), 0, rp.debugLevel, rp.debugOn)
}
