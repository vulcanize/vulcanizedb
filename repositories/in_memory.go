package repositories

import (
	"github.com/8thlight/vulcanizedb/core"
)

type InMemory struct {
	blocks map[int64]*core.Block
}

func NewInMemory() *InMemory {
	return &InMemory{
		blocks: make(map[int64]*core.Block),
	}
}

func (repository *InMemory) CreateBlock(block core.Block) {
	repository.blocks[block.Number] = &block
}

func (repository *InMemory) BlockCount() int {
	return len(repository.blocks)
}

func (repository *InMemory) FindBlockByNumber(blockNumber int64) *core.Block {
	return repository.blocks[blockNumber]
}
