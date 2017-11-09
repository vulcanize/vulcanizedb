package repositories

import (
	"github.com/8thlight/vulcanizedb/pkg/core"
)

type InMemory struct {
	blocks    map[int64]*core.Block
	contracts map[string]*core.Contract
}

func (repository *InMemory) CreateContract(contract core.Contract) error {
	repository.contracts[contract.Hash] = &contract
	return nil
}

func (repository *InMemory) IsWatchedContract(contractHash string) bool {
	_, present := repository.contracts[contractHash]
	return present
}

func (repository *InMemory) MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64) []int64 {
	missingNumbers := []int64{}
	for blockNumber := int64(startingBlockNumber); blockNumber <= endingBlockNumber; blockNumber++ {
		if repository.blocks[blockNumber] == nil {
			missingNumbers = append(missingNumbers, blockNumber)
		}
	}
	return missingNumbers
}

func NewInMemory() *InMemory {
	return &InMemory{
		blocks:    make(map[int64]*core.Block),
		contracts: make(map[string]*core.Contract),
	}
}

func (repository *InMemory) CreateBlock(block core.Block) error {
	repository.blocks[block.Number] = &block
	return nil
}

func (repository *InMemory) BlockCount() int {
	return len(repository.blocks)
}

func (repository *InMemory) FindBlockByNumber(blockNumber int64) *core.Block {
	return repository.blocks[blockNumber]
}

func (repository *InMemory) MaxBlockNumber() int64 {
	highestBlockNumber := int64(-1)
	for key := range repository.blocks {
		if key > highestBlockNumber {
			highestBlockNumber = key
		}
	}
	return highestBlockNumber
}
