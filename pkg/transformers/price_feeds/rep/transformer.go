package rep

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
)

type RepTransformer struct {
	fetcher    IRepFetcher
	repository IRepRepository
}

func NewRepTransformer(chain core.BlockChain, db *postgres.DB) RepTransformer {
	fetcher := NewRepFetcher(chain)
	repository := NewRepRepository(db)
	return RepTransformer{
		fetcher:    fetcher,
		repository: repository,
	}
}

func (transformer RepTransformer) Execute(header core.Header, headerID int64) error {
	logValue, err := transformer.fetcher.FetchRepValue(header)
	if err != nil {
		if err == price_feeds.ErrNoMatchingLog {
			return nil
		}
		return err
	}
	rep := getRep(logValue, header, headerID)
	return transformer.repository.CreateRep(rep)
}

func getRep(logValue string, header core.Header, headerID int64) price_feeds.PriceUpdate {
	valueInUSD := price_feeds.Convert("wad", logValue, 15)
	rep := price_feeds.PriceUpdate{
		BlockNumber: header.BlockNumber,
		HeaderID:    headerID,
		UsdValue:    valueInUSD,
	}
	return rep
}
