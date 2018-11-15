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
)

func PopulateMissingBlocks(blockchain core.BlockChain, blockRepository datastore.BlockRepository, startingBlockNumber int64) int {
	lastBlock := blockchain.LastBlock().Int64()
	blockRange := blockRepository.MissingBlockNumbers(startingBlockNumber, lastBlock, blockchain.Node().ID)
	log.SetPrefix("")
	log.Printf("Backfilling %d blocks\n\n", len(blockRange))
	RetrieveAndUpdateBlocks(blockchain, blockRepository, blockRange)
	return len(blockRange)
}

func RetrieveAndUpdateBlocks(blockchain core.BlockChain, blockRepository datastore.BlockRepository, blockNumbers []int64) int {
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
