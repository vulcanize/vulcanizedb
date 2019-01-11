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
		log.Error("Error in checking header in PopulateMissingHeaders: ", err)
		return 0, err
	} else if headerAlreadyExists {
		return 0, nil
	}

	blockNumbers, err := headerRepository.MissingBlockNumbers(startingBlockNumber, lastBlock, blockchain.Node().ID)
	if err != nil {
		log.Error("Error getting missing block numbers in PopulateMissingHeaders: ", err)
		return 0, err
	}

	log.Printf("Backfilling %d blocks\n\n", len(blockNumbers))
	_, err = RetrieveAndUpdateHeaders(blockchain, headerRepository, blockNumbers)
	if err != nil {
		return 0, err
	}
	return len(blockNumbers), nil
}

func RetrieveAndUpdateHeaders(chain core.BlockChain, headerRepository datastore.HeaderRepository, blockNumbers []int64) (int, error) {
	for _, blockNumber := range blockNumbers {
		header, err := chain.GetHeaderByNumber(blockNumber)
		if err != nil {
			return 0, err
		}
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
