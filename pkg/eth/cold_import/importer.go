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

package cold_import

import (
	"github.com/vulcanize/vulcanizedb/pkg/eth/converters/common"
	"github.com/vulcanize/vulcanizedb/pkg/eth/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/eth/datastore/ethereum"
)

type ColdImporter struct {
	blockRepository   datastore.BlockRepository
	converter         common.BlockConverter
	ethDB             ethereum.Database
	receiptRepository datastore.FullSyncReceiptRepository
}

func NewColdImporter(ethDB ethereum.Database, blockRepository datastore.BlockRepository, receiptRepository datastore.FullSyncReceiptRepository, converter common.BlockConverter) *ColdImporter {
	return &ColdImporter{
		blockRepository:   blockRepository,
		converter:         converter,
		ethDB:             ethDB,
		receiptRepository: receiptRepository,
	}
}

func (ci *ColdImporter) Execute(startingBlockNumber int64, endingBlockNumber int64, nodeID string) error {
	missingBlocks := ci.blockRepository.MissingBlockNumbers(startingBlockNumber, endingBlockNumber, nodeID)
	for _, n := range missingBlocks {
		hash := ci.ethDB.GetBlockHash(n)

		blockID, err := ci.createBlocksAndTransactions(hash, n)
		if err != nil {
			return err
		}
		err = ci.createReceiptsAndLogs(hash, n, blockID)
		if err != nil {
			return err
		}
	}
	ci.blockRepository.SetBlocksStatus(endingBlockNumber)
	return nil
}

func (ci *ColdImporter) createBlocksAndTransactions(hash []byte, i int64) (int64, error) {
	block := ci.ethDB.GetBlock(hash, i)
	coreBlock, err := ci.converter.ToCoreBlock(block)
	if err != nil {
		return 0, err
	}
	return ci.blockRepository.CreateOrUpdateBlock(coreBlock)
}

func (ci *ColdImporter) createReceiptsAndLogs(hash []byte, number int64, blockID int64) error {
	receipts := ci.ethDB.GetBlockReceipts(hash, number)
	coreReceipts, err := common.ToCoreReceipts(receipts)
	if err != nil {
		return err
	}
	return ci.receiptRepository.CreateReceiptsAndLogs(blockID, coreReceipts)
}
