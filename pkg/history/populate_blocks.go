package history

import (
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

type Window struct {
	LowerBound     int
	UpperBound     int
	MaxBlockNumber int
}

func (window Window) Size() int {
	return int(window.UpperBound - window.LowerBound)
}

func PopulateMissingBlocks(blockchain core.Blockchain, repository repositories.Repository, startingBlockNumber int64) int {
	blockRange := repository.MissingBlockNumbers(startingBlockNumber, repository.MaxBlockNumber())
	updateBlockRange(blockchain, repository, blockRange)
	return len(blockRange)
}

func UpdateBlocksWindow(blockchain core.Blockchain, repository repositories.Repository, windowSize int) Window {
	maxBlockNumber := repository.MaxBlockNumber()
	upperBound := repository.MaxBlockNumber() - int64(2)
	lowerBound := upperBound - int64(windowSize)
	blockRange := MakeRange(lowerBound, upperBound)
	updateBlockRange(blockchain, repository, blockRange)
	return Window{int(lowerBound), int(upperBound), int(maxBlockNumber)}
}

func updateBlockRange(blockchain core.Blockchain, repository repositories.Repository, blockNumbers []int64) int {
	for _, blockNumber := range blockNumbers {
		block := blockchain.GetBlockByNumber(blockNumber)
		repository.CreateOrUpdateBlock(block)
	}
	return len(blockNumbers)
}

func MakeRange(min, max int64) []int64 {
	a := make([]int64, max-min)
	for i := range a {
		a[i] = min + int64(i)
	}
	return a
}
