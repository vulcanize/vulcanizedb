package repositories

import (
	"errors"
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type Repository interface {
	BlockRepository
	ContractRepository
	LogsRepository
	ReceiptRepository
	FilterRepository
}

var ErrBlockDoesNotExist = func(blockNumber int64) error {
	return errors.New(fmt.Sprintf("Block number %d does not exist", blockNumber))
}

type BlockRepository interface {
	CreateOrUpdateBlock(block core.Block) error
	FindBlockByNumber(blockNumber int64) (core.Block, error)
	MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64) []int64
	SetBlocksStatus(chainHead int64)
}

var ErrContractDoesNotExist = func(contractHash string) error {
	return errors.New(fmt.Sprintf("Contract %v does not exist", contractHash))
}

type ContractRepository interface {
	CreateContract(contract core.Contract) error
	ContractExists(contractHash string) bool
	FindContract(contractHash string) (core.Contract, error)
}

type FilterRepository interface {
	AddFilter(filter filters.LogFilter) error
}

type LogsRepository interface {
	FindLogs(address string, blockNumber int64) []core.Log
	CreateLogs(logs []core.Log) error
}

var ErrReceiptDoesNotExist = func(txHash string) error {
	return errors.New(fmt.Sprintf("Receipt for tx: %v does not exist", txHash))
}

type ReceiptRepository interface {
	FindReceipt(txHash string) (core.Receipt, error)
}
