package fakes

import (
	"github.com/8thlight/vulcanizedb/core"
)

type BlockchainObserver struct {
	wasToldBlockAdded bool
	LastAddedBlock    core.Block
}

func (blockchainObserver *BlockchainObserver) WasToldBlockAdded() bool {
	return blockchainObserver.wasToldBlockAdded
}

func (blockchainObserver *BlockchainObserver) NotifyBlockAdded(block core.Block) {
	blockchainObserver.LastAddedBlock = block
	blockchainObserver.wasToldBlockAdded = true
}
