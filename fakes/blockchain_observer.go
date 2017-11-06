package fakes

import "github.com/8thlight/vulcanizedb/core"

type BlockchainObserver struct {
	CurrentBlocks []core.Block
	WasNotified   chan bool
}

func (observer *BlockchainObserver) LastBlock() core.Block {
	return observer.CurrentBlocks[len(observer.CurrentBlocks)-1]
}

func NewFakeBlockchainObserver() *BlockchainObserver {
	return &BlockchainObserver{
		WasNotified: make(chan bool),
	}
}

func (observer *BlockchainObserver) NotifyBlockAdded(block core.Block) {
	observer.CurrentBlocks = append(observer.CurrentBlocks, block)
	observer.WasNotified <- true
}
