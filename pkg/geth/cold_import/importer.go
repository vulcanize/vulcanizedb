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

package cold_import

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/ethereum"
	"github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
)

type ColdImporter struct {
	blockRepository   datastore.BlockRepository
	converter         common.BlockConverter
	ethDB             ethereum.Database
	receiptRepository datastore.ReceiptRepository
}

func NewColdImporter(ethDB ethereum.Database, blockRepository datastore.BlockRepository, receiptRepository datastore.ReceiptRepository, converter common.BlockConverter) *ColdImporter {
	return &ColdImporter{
		blockRepository:   blockRepository,
		converter:         converter,
		ethDB:             ethDB,
		receiptRepository: receiptRepository,
	}
}

func (ci *ColdImporter) Execute(startingBlockNumber int64, endingBlockNumber int64, nodeId string) error {
	missingBlocks := ci.blockRepository.MissingBlockNumbers(startingBlockNumber, endingBlockNumber, nodeId)
	for _, n := range missingBlocks {
		hash := ci.ethDB.GetBlockHash(n)

		blockId, err := ci.createBlocksAndTransactions(hash, n)
		if err != nil {
			return err
		}
		err = ci.createReceiptsAndLogs(hash, n, blockId)
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

func (ci *ColdImporter) createReceiptsAndLogs(hash []byte, number int64, blockId int64) error {
	receipts := ci.ethDB.GetBlockReceipts(hash, number)
	coreReceipts := common.ToCoreReceipts(receipts)
	return ci.receiptRepository.CreateReceiptsAndLogs(blockId, coreReceipts)
}
