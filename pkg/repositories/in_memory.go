package repositories

import (
	"fmt"

	"github.com/8thlight/vulcanizedb/pkg/core"
)

type InMemory struct {
	blocks    map[int64]*core.Block
	contracts map[string]*core.Contract
	logs      map[string][]core.Log
}

func (repository *InMemory) CreateLogs(logs []core.Log) error {
	for _, log := range logs {
		key := fmt.Sprintf("%s%s", log.BlockNumber, log.Index)
		var logs []core.Log
		repository.logs[key] = append(logs, log)
	}
	return nil
}

func (repository *InMemory) FindLogs(address string, blockNumber int64) []core.Log {
	var matchingLogs []core.Log
	for _, logs := range repository.logs {
		for _, log := range logs {
			if log.Address == address && log.BlockNumber == blockNumber {
				matchingLogs = append(matchingLogs, log)
			}
		}
	}
	return matchingLogs
}

func (repository *InMemory) CreateContract(contract core.Contract) error {
	repository.contracts[contract.Hash] = &contract
	return nil
}

func (repository *InMemory) ContractExists(contractHash string) bool {
	_, present := repository.contracts[contractHash]
	return present
}

func (repository *InMemory) FindContract(contractHash string) *core.Contract {
	contract, ok := repository.contracts[contractHash]
	if !ok {
		return nil
	}
	for _, block := range repository.blocks {
		for _, transaction := range block.Transactions {
			if transaction.To == contractHash {
				contract.Transactions = append(contract.Transactions, transaction)
			}
		}
	}
	return contract
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
		logs:      make(map[string][]core.Log),
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
