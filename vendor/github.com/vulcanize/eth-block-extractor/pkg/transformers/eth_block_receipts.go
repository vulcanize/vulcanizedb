package transformers

import (
	"github.com/vulcanize/eth-block-extractor/pkg/db"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"log"
)

type EthBlockReceiptTransformer struct {
	database  db.Database
	publisher ipfs.Publisher
}

func NewEthBlockReceiptTransformer(database db.Database, publisher ipfs.Publisher) *EthBlockReceiptTransformer {
	return &EthBlockReceiptTransformer{
		database:  database,
		publisher: publisher,
	}
}

func (transformer EthBlockReceiptTransformer) Execute(startingBlockNumber int64, endingBlockNumber int64) error {
	if endingBlockNumber < startingBlockNumber {
		return ErrInvalidRange
	}
	for i := startingBlockNumber; i <= endingBlockNumber; i++ {
		receipts := transformer.database.GetBlockReceipts(i)
		cids, err := transformer.publisher.Write(receipts)
		if err != nil {
			return err
		}
		log.Println("Generated IPLDs: ", cids)
	}
	return nil
}
