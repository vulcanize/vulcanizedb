package repositories

import (
	"errors"
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

var ErrBlockDoesNotExist = func(blockNumber int64) error {
	return errors.New(fmt.Sprintf("Block number %d does not exist", blockNumber))
}

type BlockRepository interface {
	CreateOrUpdateBlock(block core.Block) error
	GetBlock(blockNumber int64) (core.Block, error)
	MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64) []int64
	SetBlocksStatus(chainHead int64)
}

var ErrContractDoesNotExist = func(contractHash string) error {
	return errors.New(fmt.Sprintf("Contract %v does not exist", contractHash))
}

type ContractRepository interface {
	CreateContract(contract core.Contract) error
	GetContract(contractHash string) (core.Contract, error)
	ContractExists(contractHash string) bool
}

var ErrFilterDoesNotExist = func(name string) error {
	return errors.New(fmt.Sprintf("filter %s does not exist", name))
}

type FilterRepository interface {
	CreateFilter(filter filters.LogFilter) error
	GetFilter(name string) (filters.LogFilter, error)
}

type LogRepository interface {
	CreateLogs(logs []core.Log) error
	GetLogs(address string, blockNumber int64) []core.Log
}

var ErrReceiptDoesNotExist = func(txHash string) error {
	return errors.New(fmt.Sprintf("Receipt for tx: %v does not exist", txHash))
}

type ReceiptRepository interface {
	GetReceipt(txHash string) (core.Receipt, error)
}

type WatchedEventRepository interface {
	GetWatchedEvents(name string) ([]*core.WatchedEvent, error)
}
