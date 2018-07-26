package pep

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
)

type PepTransformer struct {
	fetcher    IPepFetcher
	repository IPepRepository
}

func NewPepTransformer(chain core.BlockChain, db *postgres.DB, contractAddress string) PepTransformer {
	fetcher := NewPepFetcher(chain, contractAddress)
	repository := NewPepRepository(db)
	return PepTransformer{
		fetcher:    fetcher,
		repository: repository,
	}
}

func (transformer PepTransformer) Execute(header core.Header, headerID int64) error {
	logValue, err := transformer.fetcher.FetchPepValue(header)
	if err != nil {
		if err == price_feeds.ErrNoMatchingLog {
			return nil
		}
		return err
	}
	pep := getPep(logValue, header, headerID)
	return transformer.repository.CreatePep(pep)
}

func getPep(logValue string, header core.Header, headerID int64) price_feeds.PriceUpdate {
	valueInUSD := price_feeds.Convert("wad", logValue, 15)
	pep := price_feeds.PriceUpdate{
		BlockNumber: header.BlockNumber,
		HeaderID:    headerID,
		UsdValue:    valueInUSD,
	}
	return pep
}
