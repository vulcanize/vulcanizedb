package rep

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"math/big"
)

type IRepFetcher interface {
	FetchRepValue(header core.Header) (string, error)
}

type RepFetcher struct {
	chain core.BlockChain
}

func NewRepFetcher(chain core.BlockChain) RepFetcher {
	return RepFetcher{
		chain: chain,
	}
}

func (fetcher RepFetcher) FetchRepValue(header core.Header) (string, error) {
	blockNumber := big.NewInt(header.BlockNumber)
	logs, err := fetcher.chain.GetLogs(price_feeds.RepAddress, price_feeds.RepLogTopic0, blockNumber, blockNumber)
	return fetcher.getLogValue(logs, err)
}

func (fetcher RepFetcher) getLogValue(logs []core.Log, err error) (string, error) {
	if err != nil {
		return "", err
	}
	if len(logs) < 1 {
		return "", price_feeds.ErrNoMatchingLog
	}
	if len(logs) > 1 {
		return "", price_feeds.ErrMultipleLogs
	}
	return logs[0].Data, nil
}
