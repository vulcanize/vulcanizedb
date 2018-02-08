package inmemory

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

type Blocks struct {
	*InMemory
}

func (blocks *Blocks) CreateOrUpdateBlock(block core.Block) error {
	blocks.CreateOrUpdateBlockCallCount++
	blocks.blocks[block.Number] = block
	for _, transaction := range block.Transactions {
		blocks.receipts[transaction.Hash] = transaction.Receipt
		blocks.logs[transaction.TxHash] = transaction.Logs
	}
	return nil
}

func (repository *Blocks) GetBlock(blockNumber int64) (core.Block, error) {
	if block, ok := repository.blocks[blockNumber]; ok {
		return block, nil
	}
	return core.Block{}, repositories.ErrBlockDoesNotExist(blockNumber)
}

func (repository *Blocks) MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64) []int64 {
	missingNumbers := []int64{}
	for blockNumber := int64(startingBlockNumber); blockNumber <= endingBlockNumber; blockNumber++ {
		if _, ok := repository.blocks[blockNumber]; !ok {
			missingNumbers = append(missingNumbers, blockNumber)
		}
	}
	return missingNumbers
}

func (repository *Blocks) SetBlocksStatus(chainHead int64) {
	for key, block := range repository.blocks {
		if key < (chainHead - blocksFromHeadBeforeFinal) {
			tmp := block
			tmp.IsFinal = true
			repository.blocks[key] = tmp
		}
	}
}

func (repository *Blocks) BlockCount() int {
	return len(repository.blocks)
}
