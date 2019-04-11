package transformers

import (
	"log"

	"github.com/vulcanize/eth-block-extractor/pkg/db"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
)

type EthBlockHeaderTransformer struct {
	database  db.Database
	publisher ipfs.Publisher
}

func NewEthBlockHeaderTransformer(ethDB db.Database, publisher ipfs.Publisher) *EthBlockHeaderTransformer {
	return &EthBlockHeaderTransformer{database: ethDB, publisher: publisher}
}

func (t EthBlockHeaderTransformer) Execute(startingBlockNumber int64, endingBlockNumber int64) error {
	if endingBlockNumber < startingBlockNumber {
		return ErrInvalidRange
	}
	for i := startingBlockNumber; i <= endingBlockNumber; i++ {
		blockData := t.database.GetRawBlockHeaderByBlockNumber(i)
		output, err := t.publisher.Write(blockData)
		if err != nil {
			return NewExecuteError(PutIpldErr, err)
		}
		log.Printf("Created IPLD: %s", output)
	}
	return nil
}
