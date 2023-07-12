package common

import (
	"async-consensus/proto"
	"fmt"
)

/*
	Implements a block store for the consensus blocks
*/

/*
	A single element in the map is an async-consensus-block with its set of acks
*/

type AsynConsensusBlock struct {
	ConsensusBlock *proto.AsyncConsensus_Block
	acks           []int32 // contains the set of nodes who acknowledged the consensus block, a simple counter might work, but having an array is extensible
}

/*
	Async consensus store defines the underlying map
*/

type AsyncConsensusStore struct {
	ConsensusBlocks map[string]AsynConsensusBlock
	debugLevel      int
	debugOn         bool
}

/*
	allocate the map object
*/

func (ms *AsyncConsensusStore) Init(debugLevel int, debugOn bool) {
	ms.ConsensusBlocks = make(map[string]AsynConsensusBlock)
	ms.debugLevel = debugLevel
	ms.debugOn = debugOn
	Debug("Initialized a new consensus message store", 0, ms.debugLevel, ms.debugOn)
}

/*
	Add a new consensus block to the store if it is not already there
*/

func (ms *AsyncConsensusStore) Add(block *proto.AsyncConsensus_Block) {
	if block.R > 2 && block.Parent == nil {
		panic("Error Nil parent found in " + fmt.Sprintf("%v", block))
	}
	_, ok := ms.ConsensusBlocks[block.Id]
	if !ok {
		ms.ConsensusBlocks[block.Id] = AsynConsensusBlock{
			ConsensusBlock: block,
			acks:           make([]int32, 0),
		}
		Debug("Added a new consensus block to consensus store with  "+fmt.Sprintf("%v", block.Id), 0, ms.debugLevel, ms.debugOn)
	}
}

/*
	return an existing consensus block
*/

func (ms *AsyncConsensusStore) Get(id string) (*proto.AsyncConsensus_Block, bool) {
	block, ok := ms.ConsensusBlocks[id]
	if !ok {
		Debug("Requested consensus block does not exist, hence returning nil for id "+id, 0, ms.debugLevel, ms.debugOn)
		return nil, ok
	} else {
		Debug("Requested consensus block exists, hence returning the block for id "+id, 0, ms.debugLevel, ms.debugOn)
		if block.ConsensusBlock.R > 2 && block.ConsensusBlock.Parent == nil {
			panic("Error Nil parent found in " + fmt.Sprintf("%v", block.ConsensusBlock))
		}
		return block.ConsensusBlock, ok
	}
}

/*
	return the set of acks for a given consensus block
*/

func (ms *AsyncConsensusStore) GetAcks(id string) []int32 {
	block, ok := ms.ConsensusBlocks[id]
	if ok {
		Debug("Requested mem block exists, hence returning the acks for id "+id, 0, ms.debugLevel, ms.debugOn)
		return block.acks
	}
	Debug("Requested mem block does not exist, hence returning nil acks for id "+id, 0, ms.debugLevel, ms.debugOn)
	return nil
}

/*
	Remove an element from the map
*/

func (ms *AsyncConsensusStore) Remove(id string) {
	delete(ms.ConsensusBlocks, id)
	Debug("Removed a consensus block with id "+id, 0, ms.debugLevel, ms.debugOn)
}

/*
	add a new ack to the ack list of a consensus block
*/

func (ms *AsyncConsensusStore) AddAck(id string, node int32) {
	block, ok := ms.ConsensusBlocks[id]
	if ok {
		tempAcks := block.acks
		tempBlock := block.ConsensusBlock
		tempAcks = append(tempAcks, node)
		ms.Remove(id)
		ms.ConsensusBlocks[id] = AsynConsensusBlock{
			ConsensusBlock: tempBlock,
			acks:           tempAcks,
		}
		Debug("Added an ack for consensus block with id "+id, 0, ms.debugLevel, ms.debugOn)
	}
	Debug("Adding an ack failed for consensus block with id "+id+" because the consensus block does not exist", 0, ms.debugLevel, ms.debugOn)
}
