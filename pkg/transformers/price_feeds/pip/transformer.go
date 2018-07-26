package pip

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
)

type PipTransformer struct {
	fetcher    IPipFetcher
	repository IPipRepository
}

func NewPipTransformer(chain core.BlockChain, db *postgres.DB, contractAddress string) PipTransformer {
	fetcher := NewPipFetcher(chain, contractAddress)
	repository := NewPipRepository(db)
	return PipTransformer{
		fetcher:    fetcher,
		repository: repository,
	}
}

func (transformer PipTransformer) Execute(header core.Header, headerID int64) error {
	value, err := transformer.fetcher.FetchPipValue(header)
	if err != nil {
		if err == price_feeds.ErrNoMatchingLog {
			return nil
		}
		return err
	}
	pip := getPip(value, header, headerID)
	return transformer.repository.CreatePip(pip)
}

func getPip(logValue string, header core.Header, headerID int64) price_feeds.PriceUpdate {
	valueInUSD := price_feeds.Convert("wad", logValue, 15)
	pep := price_feeds.PriceUpdate{
		BlockNumber: header.BlockNumber,
		HeaderID:    headerID,
		UsdValue:    valueInUSD,
	}
	return pep
}
