package repositories

import "github.com/8thlight/vulcanizedb/pkg/core"

type Repository interface {
	CreateBlock(block core.Block) error
	BlockCount() int
	FindBlockByNumber(blockNumber int64) *core.Block
	MaxBlockNumber() int64
	MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64) []int64
}
