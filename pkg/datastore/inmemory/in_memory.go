package inmemory

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

const (
	blocksFromHeadBeforeFinal = 20
)

type InMemory struct {
	CreateOrUpdateBlockCallCount int
	blocks                       map[int64]core.Block
	contracts                    map[string]core.Contract
	headers                      map[int64]core.Header
	logFilters                   map[string]filters.LogFilter
	logs                         map[string][]core.Log
	receipts                     map[string]core.Receipt
}

func NewInMemory() *InMemory {
	return &InMemory{
		CreateOrUpdateBlockCallCount: 0,
		blocks:     make(map[int64]core.Block),
		contracts:  make(map[string]core.Contract),
		headers:    make(map[int64]core.Header),
		logFilters: make(map[string]filters.LogFilter),
		logs:       make(map[string][]core.Log),
		receipts:   make(map[string]core.Receipt),
	}
}
