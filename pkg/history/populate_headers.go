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
	"log"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

func PopulateMissingHeaders(blockchain core.BlockChain, headerRepository datastore.HeaderRepository, startingBlockNumber int64) (int, error) {
	lastBlock := blockchain.LastBlock().Int64()
	blockRange := headerRepository.MissingBlockNumbers(startingBlockNumber, lastBlock, blockchain.Node().ID)
	log.SetPrefix("")
	log.Printf("Backfilling %d blocks\n\n", len(blockRange))
	_, err := RetrieveAndUpdateHeaders(blockchain, headerRepository, blockRange)
	if err != nil {
		return 0, err
	}
	return len(blockRange), nil
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
