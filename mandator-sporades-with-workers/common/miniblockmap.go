package common

import (
	"with-workers/proto"
)

/*
	Implements a block store for the mini mem pool which is the store for the mini mem blocks
*/

/*
	A single element in the map is a mini block with its set of acks
*/

type Block struct {
	MessageBlock *proto.MemPoolMini
	acks         []int32 // contains the set of nodes who acknowledged the block, a simple counter might work, but having an array is extensible
}

/*
	Message store defines the underlying map
*/

type MiniMessageStore struct {
	MessageBlocks map[string]Block
	debugLevel    int
	debugOn       bool
}

/*
	allocate the map object
*/

func (ms *MiniMessageStore) Init(debugLevel int, debugOn bool) {
	ms.MessageBlocks = make(map[string]Block)
	ms.debugLevel = debugLevel
	ms.debugOn = debugOn
	Debug("Initialized a new mini message store", 0, ms.debugLevel, ms.debugOn)
}

/*
	Add a new mini block to the store if it is not already there
*/

func (ms *MiniMessageStore) Add(block *proto.MemPoolMini) {
	_, ok := ms.MessageBlocks[block.UniqueId]
	if !ok {
		ms.MessageBlocks[block.UniqueId] = Block{
			MessageBlock: block,
			acks:         make([]int32, 0),
		}
		Debug("Added a new mini block to store with id "+block.UniqueId, 0, ms.debugLevel, ms.debugOn)
	}

}

/*
	return an existing block
*/

func (ms *MiniMessageStore) Get(id string) (*proto.MemPoolMini, bool) {
	block, ok := ms.MessageBlocks[id]
	if !ok {
		Debug("Requested block does not exist, hence returning nil for id "+id, 0, ms.debugLevel, ms.debugOn)
		return nil, ok
	} else {
		Debug("Requested block exist, hence returning the mini block for id "+id, 0, ms.debugLevel, ms.debugOn)
		return block.MessageBlock, ok
	}
}

/*
	return the set of acks for a given block
*/

func (ms *MiniMessageStore) GetAcks(id string) []int32 {
	block, ok := ms.MessageBlocks[id]
	if ok {
		Debug("Requested block exists, hence returning the acks for id "+id, 0, ms.debugLevel, ms.debugOn)
		return block.acks
	}
	Debug("Requested block does not exist, hence returning nil acks for id "+id, 0, ms.debugLevel, ms.debugOn)
	return nil
}

/*
	Remove an element from the map
*/

func (ms *MiniMessageStore) Remove(id string) {
	delete(ms.MessageBlocks, id)
	Debug("Removed a mini block with id "+id, 0, ms.debugLevel, ms.debugOn)
}

/*
	add a new ack to the ack list of a block
*/

func (ms *MiniMessageStore) AddAck(id string, node int32) {
	block, ok := ms.MessageBlocks[id]
	if ok {
		tempAcks := block.acks
		tempBlock := block.MessageBlock
		tempAcks = append(tempAcks, node)
		ms.Remove(id)
		ms.MessageBlocks[id] = Block{
			MessageBlock: tempBlock,
			acks:         tempAcks,
		}
		Debug("Added an ack for mini block with id "+id, 0, ms.debugLevel, ms.debugOn)
	}
	Debug("Adding an ack failed for mini block with id "+id+" because the mini block is not found", 0, ms.debugLevel, ms.debugOn)
}
