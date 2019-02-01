// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
