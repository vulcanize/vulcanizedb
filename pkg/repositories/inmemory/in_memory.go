package inmemory

import (
	"fmt"

	"errors"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

const (
	blocksFromHeadBeforeFinal = 20
)

type InMemory struct {
	blocks                       map[int64]core.Block
	receipts                     map[string]core.Receipt
	contracts                    map[string]core.Contract
	logs                         map[string][]core.Log
	logFilters                   map[string]filters.LogFilter
	CreateOrUpdateBlockCallCount int
}

func (repository *InMemory) AddFilter(filter filters.LogFilter) error {
	key := filter.Name
	if _, ok := repository.logFilters[key]; ok || key == "" {
		return errors.New("filter name not unique")
	}
	repository.logFilters[key] = filter
	return nil
}

func NewInMemory() *InMemory {
	return &InMemory{
		CreateOrUpdateBlockCallCount: 0,
		blocks:     make(map[int64]core.Block),
		receipts:   make(map[string]core.Receipt),
		contracts:  make(map[string]core.Contract),
		logs:       make(map[string][]core.Log),
		logFilters: make(map[string]filters.LogFilter),
	}
}

func (repository *InMemory) FindReceipt(txHash string) (core.Receipt, error) {
	if receipt, ok := repository.receipts[txHash]; ok {
		return receipt, nil
	}
	return core.Receipt{}, repositories.ErrReceiptDoesNotExist(txHash)
}

func (repository *InMemory) SetBlocksStatus(chainHead int64) {
	for key, block := range repository.blocks {
		if key < (chainHead - blocksFromHeadBeforeFinal) {
			tmp := block
			tmp.IsFinal = true
			repository.blocks[key] = tmp
		}
	}
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
	repository.contracts[contract.Hash] = contract
	return nil
}

func (repository *InMemory) ContractExists(contractHash string) bool {
	_, present := repository.contracts[contractHash]
	return present
}

func (repository *InMemory) FindContract(contractHash string) (core.Contract, error) {
	contract, ok := repository.contracts[contractHash]
	if !ok {
		return core.Contract{}, repositories.ErrContractDoesNotExist(contractHash)
	}
	for _, block := range repository.blocks {
		for _, transaction := range block.Transactions {
			if transaction.To == contractHash {
				contract.Transactions = append(contract.Transactions, transaction)
			}
		}
	}
	return contract, nil
}

func (repository *InMemory) MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64) []int64 {
	missingNumbers := []int64{}
	for blockNumber := int64(startingBlockNumber); blockNumber <= endingBlockNumber; blockNumber++ {
		if _, ok := repository.blocks[blockNumber]; !ok {
			missingNumbers = append(missingNumbers, blockNumber)
		}
	}
	return missingNumbers
}

func (repository *InMemory) CreateOrUpdateBlock(block core.Block) error {
	repository.CreateOrUpdateBlockCallCount++
	repository.blocks[block.Number] = block
	for _, transaction := range block.Transactions {
		repository.receipts[transaction.Hash] = transaction.Receipt
		repository.logs[transaction.TxHash] = transaction.Logs
	}
	return nil
}

func (repository *InMemory) BlockCount() int {
	return len(repository.blocks)
}

func (repository *InMemory) FindBlockByNumber(blockNumber int64) (core.Block, error) {
	if block, ok := repository.blocks[blockNumber]; ok {
		return block, nil
	}
	return core.Block{}, repositories.ErrBlockDoesNotExist(blockNumber)
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
