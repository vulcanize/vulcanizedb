package repositories

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

const (
	blocksFromHeadBeforeFinal = 20
)

type Repository interface {
	CreateOrUpdateBlock(block core.Block) error
	BlockCount() int
	FindBlockByNumber(blockNumber int64) (core.Block, error)
	MaxBlockNumber() int64
	MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64) []int64
	FindReceipt(txHash string) (core.Receipt, error)
	CreateContract(contract core.Contract) error
	ContractExists(contractHash string) bool
	FindContract(contractHash string) (core.Contract, error)
	CreateLogs(log []core.Log) error
	FindLogs(address string, blockNumber int64) []core.Log
	SetBlocksStatus(chainHead int64)
	AddFilter(filter filters.LogFilter) error
}
