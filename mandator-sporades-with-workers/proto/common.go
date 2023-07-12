package proto

import "io"

/*
	Each message sent over the network in Async-Consensus should implement this interface
	If a new message type needs to be added: first define it in a proto file, generate the go protobuf files using mage generate and then implement the three methods
*/

type Serializable interface {
	Marshal(io.Writer) error
	Unmarshal(io.Reader) error
	New() Serializable
}

/*
	A struct that allocates a unique uint8 for each message type. When you define a new proto message type, add the message to here
*/

type MessageCode struct {
	ClientBatchRpc uint8
	StatusRPC      uint8
	MemPoolMiniRPC uint8
	MemPoolRPC     uint8
	AsyncConsensus uint8
	PaxosConsensus uint8
}

/*
	A static function which assigns a unique uint8 to each message type. Update this function when you define new message types
*/

func GetRPCCodes() MessageCode {
	return MessageCode{
		ClientBatchRpc: 0,
		StatusRPC:      1,
		MemPoolMiniRPC: 2,
		MemPoolRPC:     3,
		AsyncConsensus: 4,
		PaxosConsensus: 5,
	}
}
