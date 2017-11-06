package observers

import (
	"github.com/8thlight/vulcanizedb/core"
	"github.com/8thlight/vulcanizedb/repositories"
)

type BlockchainDbObserver struct {
	repository repositories.Repository
}

func NewBlockchainDbObserver(repository repositories.Repository) BlockchainDbObserver {
	return BlockchainDbObserver{repository: repository}
}

func (observer BlockchainDbObserver) NotifyBlockAdded(block core.Block) {
	observer.repository.CreateBlock(block)
}
