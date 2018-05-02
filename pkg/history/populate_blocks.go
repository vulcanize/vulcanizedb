package history

import (
	"log"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
)

func PopulateMissingBlocks(blockchain core.Blockchain, blockRepository datastore.BlockRepository, startingBlockNumber int64) int {
	lastBlock := blockchain.LastBlock().Int64()
	blockRange := blockRepository.MissingBlockNumbers(startingBlockNumber, lastBlock-1)
	log.SetPrefix("")
	log.Printf("Backfilling %d blocks\n\n", len(blockRange))
	RetrieveAndUpdateBlocks(blockchain, blockRepository, blockRange)
	return len(blockRange)
}

func RetrieveAndUpdateBlocks(blockchain core.Blockchain, blockRepository datastore.BlockRepository, blockNumbers []int64) int {
	for _, blockNumber := range blockNumbers {
		block, err := blockchain.GetBlockByNumber(blockNumber)
		if err != nil {
			log.Printf("failed to retrieve block number: %d\n", blockNumber)
			return 0
		}
		// TODO: handle possible error here
		blockRepository.CreateOrUpdateBlock(block)
	}
	return len(blockNumbers)
}
