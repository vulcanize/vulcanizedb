package history

import (
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

func PopulateBlocks(blockchain core.Blockchain, repository repositories.Repository, startingBlockNumber int64) int {
	blockNumbers := repository.MissingBlockNumbers(startingBlockNumber, repository.MaxBlockNumber())
	for _, blockNumber := range blockNumbers {
		block := blockchain.GetBlockByNumber(blockNumber)
		repository.CreateBlock(block)
	}
	return len(blockNumbers)
}
