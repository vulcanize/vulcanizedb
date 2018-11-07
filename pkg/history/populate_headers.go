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
