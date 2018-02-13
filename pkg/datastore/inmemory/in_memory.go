package inmemory

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
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
