package fakes

import (
	"github.com/8thlight/vulcanizedb/core"
)

type Blockchain struct {
	observers []core.BlockchainObserver
}

func (blockchain *Blockchain) RegisterObserver(observer core.BlockchainObserver) {
	blockchain.observers = append(blockchain.observers, observer)
}

func (blockchain *Blockchain) AddBlock(block core.Block) {
	for _, observer := range blockchain.observers {
		observer.NotifyBlockAdded(block)
	}
}
