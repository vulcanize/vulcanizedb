// VulcanizeDB
// Copyright © 2019 Vulcanize

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
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
)

func PopulateMissingBlocks(blockchain core.BlockChain, blockRepository datastore.BlockRepository, startingBlockNumber int64) (int, error) {
	lastBlock, err := blockchain.LastBlock()
	if err != nil {
		log.Error("PopulateMissingBlocks: error getting last block: ", err)
		return 0, err
	}
	blockRange := blockRepository.MissingBlockNumbers(startingBlockNumber, lastBlock.Int64(), blockchain.Node().ID)

	if len(blockRange) == 0 {
		return 0, nil
	}

	log.Debug(getBlockRangeString(blockRange))
	_, err = RetrieveAndUpdateBlocks(blockchain, blockRepository, blockRange)
	if err != nil {
		log.Error("PopulateMissingBlocks: error gettings/updating blocks: ", err)
		return 0, err
	}
	return len(blockRange), nil
}

func RetrieveAndUpdateBlocks(blockchain core.BlockChain, blockRepository datastore.BlockRepository, blockNumbers []int64) (int, error) {
	for _, blockNumber := range blockNumbers {
		block, err := blockchain.GetBlockByNumber(blockNumber)
		if err != nil {
			log.Error("RetrieveAndUpdateBlocks: error getting block: ", err)
			return 0, err
		}

		_, err = blockRepository.CreateOrUpdateBlock(block)
		if err != nil {
			log.Error("RetrieveAndUpdateBlocks: error creating/updating block: ", err)
			return 0, err
		}

	}
	return len(blockNumbers), nil
}

func getBlockRangeString(blockRange []int64) string {
	return fmt.Sprintf("Backfilling |%v| blocks", len(blockRange))
}
