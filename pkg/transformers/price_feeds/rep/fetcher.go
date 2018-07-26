package rep

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"math/big"
)

type IRepFetcher interface {
	FetchRepValue(header core.Header) (string, error)
}

type RepFetcher struct {
	chain           core.BlockChain
	contractAddress string
}

func NewRepFetcher(chain core.BlockChain, contractAddress string) RepFetcher {
	return RepFetcher{
		chain:           chain,
		contractAddress: contractAddress,
	}
}

func (fetcher RepFetcher) FetchRepValue(header core.Header) (string, error) {
	blockNumber := big.NewInt(header.BlockNumber)
	query := ethereum.FilterQuery{
		FromBlock: blockNumber,
		ToBlock:   blockNumber,
		Addresses: []common.Address{common.HexToAddress(fetcher.contractAddress)},
		Topics:    [][]common.Hash{{common.HexToHash(price_feeds.RepLogTopic0)}},
	}
	logs, err := fetcher.chain.GetEthLogsWithCustomQuery(query)
	return fetcher.getLogValue(logs, err)
}

func (fetcher RepFetcher) getLogValue(logs []types.Log, err error) (string, error) {
	if err != nil {
		return "", err
	}
	if len(logs) < 1 {
		return "", price_feeds.ErrNoMatchingLog
	}
	if len(logs) > 1 {
		return "", price_feeds.ErrMultipleLogs
	}
	return string(logs[0].Data), nil
}
