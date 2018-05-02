package geth

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/ethereum"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
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

func (ci *ColdImporter) Execute(startingBlockNumber int64, endingBlockNumber int64) error {
	for i := startingBlockNumber; i <= endingBlockNumber; i++ {
		hash := ci.ethDB.GetBlockHash(i)

		blockId, err := ci.createBlocksAndTransactions(hash, i)
		if err != nil {
			if err == repositories.ErrBlockExists {
				continue
			}
			return err
		}
		err = ci.createReceiptsAndLogs(hash, i, blockId)
		if err != nil {
			return err
		}
	}
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
