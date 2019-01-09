package history

import (
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

func PopulateMissingHeaders(blockchain core.BlockChain, headerRepository datastore.HeaderRepository, startingBlockNumber int64) (int, error) {
	lastBlock := blockchain.LastBlock().Int64()
	headerAlreadyExists, err := headerRepository.HeaderExists(lastBlock)
	if err != nil {
		return 0, err
	} else if headerAlreadyExists {
		return 0, nil
	}

	blockNumbers := headerRepository.MissingBlockNumbers(startingBlockNumber, lastBlock, blockchain.Node().ID)
	log.Printf("Backfilling %d blocks\n\n", len(blockNumbers))
	_, err = RetrieveAndUpdateHeaders(blockchain, headerRepository, blockNumbers)
	if err != nil {
		return 0, err
	}
	return len(blockNumbers), nil
}

func RetrieveAndUpdateHeaders(chain core.BlockChain, headerRepository datastore.HeaderRepository, blockNumbers []int64) (int, error) {
	headers, err := chain.GetHeaderByNumbers(blockNumbers)
	for _, header := range headers {
		_, err = headerRepository.CreateOrUpdateHeader(header)
		if err != nil {
			if err == repositories.ErrValidHeaderExists {
				continue
			}
			return 0, err
		}
	}

	return len(blockNumbers), nil
}
