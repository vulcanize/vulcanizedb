package history

import (
	"log"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

func PopulateMissingBlocks(blockchain core.Blockchain, blockRepository repositories.BlockRepository, startingBlockNumber int64) int {
	lastBlock := blockchain.LastBlock().Int64()
	blockRange := blockRepository.MissingBlockNumbers(startingBlockNumber, lastBlock-1)
	log.SetPrefix("")
	log.Printf("Backfilling %d blocks\n\n", len(blockRange))
	RetrieveAndUpdateBlocks(blockchain, blockRepository, blockRange)
	return len(blockRange)
}

func RetrieveAndUpdateBlocks(blockchain core.Blockchain, blockRepository repositories.BlockRepository, blockNumbers []int64) int {
	for _, blockNumber := range blockNumbers {
		block := blockchain.GetBlockByNumber(blockNumber)
		blockRepository.CreateOrUpdateBlock(block)
	}
	return len(blockNumbers)
}
