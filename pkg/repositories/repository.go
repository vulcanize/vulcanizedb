package repositories

import "github.com/8thlight/vulcanizedb/pkg/core"

type Repository interface {
	CreateBlock(block core.Block)
	BlockCount() int
	FindBlockByNumber(blockNumber int64) *core.Block
}
