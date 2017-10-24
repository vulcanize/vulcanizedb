package fakes

import (
	"github.com/8thlight/vulcanizedb/core"
)

type BlockchainObserver struct {
	wasToldBlockAdded bool
	blocks            []core.Block
}

func (blockchainObserver *BlockchainObserver) WasToldBlockAdded() bool {
	return blockchainObserver.wasToldBlockAdded
}

func (blockchainObserver *BlockchainObserver) NotifyBlockAdded(block core.Block) {
	blockchainObserver.blocks = append(blockchainObserver.blocks, block)
	blockchainObserver.wasToldBlockAdded = true
}

func (observer *BlockchainObserver) LastAddedBlock() core.Block {
	return observer.blocks[len(observer.blocks)-1]
}
