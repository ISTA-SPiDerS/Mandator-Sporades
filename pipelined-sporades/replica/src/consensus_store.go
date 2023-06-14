package src

import (
	"fmt"
	"pipelined-sporades/proto"
)

/*
	Implements a block store for the consensus blocks
*/

/*
	A single element in the map is an async-consensus-block with its set of acks
*/

type AsynConsensusBlock struct {
	ConsensusBlock *proto.Pipelined_Sporades_Block
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
	print debug messages
*/

func (ms *AsyncConsensusStore) debug(message string, level int) {
	if ms.debugOn {
		if ms.debugLevel <= level {
			fmt.Print(message)
		}
	}
}

/*
	allocate the map object
*/

func (ms *AsyncConsensusStore) Init(debugLevel int, debugOn bool) {
	ms.ConsensusBlocks = make(map[string]AsynConsensusBlock)
	ms.debugLevel = debugLevel
	ms.debugOn = debugOn
	if ms.debugOn {
		ms.debug("Initialized a new consensus message store", 0)
	}
}

/*
	Add a new consensus block to the store if it is not already there
*/

func (ms *AsyncConsensusStore) Add(block *proto.Pipelined_Sporades_Block) {

	if block.Parent != nil {
		ms.Add(block.Parent)
	}

	block.Parent = nil
	_, ok := ms.ConsensusBlocks[block.Id]
	if !ok {
		ms.ConsensusBlocks[block.Id] = AsynConsensusBlock{
			ConsensusBlock: block,
			acks:           make([]int32, 0),
		}
		if ms.debugOn {
			ms.debug("Added a new consensus block to consensus store with  "+fmt.Sprintf("%v", block.Id), 0)
		}
	}
}

/*
	return an existing consensus block
*/

func (ms *AsyncConsensusStore) Get(id string) (*proto.Pipelined_Sporades_Block, bool) {
	block, ok := ms.ConsensusBlocks[id]
	if !ok {
		if ms.debugOn {
			ms.debug("Requested consensus block does not exist, hence returning nil for id "+id, 0)
		}
		return nil, ok
	} else {
		if ms.debugOn {
			ms.debug("Requested consensus block exists, hence returning the block for id "+id, 0)
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
		if ms.debugOn {
			ms.debug("Requested block exists, hence returning the acks for id "+id, 0)
		}
		return block.acks
	}
	if ms.debugOn {
		ms.debug("Requested consensus block does not exist, hence returning nil acks for id "+id, 0)
	}
	return nil
}

/*
	Remove an element from the map
*/

func (ms *AsyncConsensusStore) Remove(id string) {
	delete(ms.ConsensusBlocks, id)
	if ms.debugOn {
		ms.debug("Removed a consensus block with id "+id, 0)
	}
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
		if ms.debugOn {
			ms.debug("Added an ack for consensus block with id "+id, 0)
		}
	} else {
		if ms.debugOn {
			ms.debug("Adding an ack failed for consensus block with id "+id+" because the consensus block does not exist", 0)
		}
	}
}
