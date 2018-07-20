package history

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"log"
)

func PopulateMissingHeaders(blockchain core.BlockChain, headerRepository datastore.HeaderRepository, startingBlockNumber int64) int {
	lastBlock := blockchain.LastBlock().Int64()
	blockRange := headerRepository.MissingBlockNumbers(startingBlockNumber, lastBlock, blockchain.Node().ID)
	log.SetPrefix("")
	log.Printf("Backfilling %d blocks\n\n", len(blockRange))
	RetrieveAndUpdateHeaders(blockchain, headerRepository, blockRange)
	return len(blockRange)
}

func RetrieveAndUpdateHeaders(blockchain core.BlockChain, headerRepository datastore.HeaderRepository, blockNumbers []int64) int {
	for _, blockNumber := range blockNumbers {
		header, err := blockchain.GetHeaderByNumber(blockNumber)
		if err != nil {
			log.Printf("failed to retrieve block number: %d\n", blockNumber)
			return 0
		}
		// TODO: handle possible error here
		headerRepository.CreateOrUpdateHeader(header)
	}
	return len(blockNumbers)
}
