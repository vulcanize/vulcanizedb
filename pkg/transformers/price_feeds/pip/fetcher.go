package pip

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"math/big"
)

type IPipFetcher interface {
	FetchPipValue(header core.Header) (string, error)
}

type PipFetcher struct {
	blockChain      core.BlockChain
	contractAddress string
}

func NewPipFetcher(chain core.BlockChain, contractAddress string) PipFetcher {
	return PipFetcher{
		blockChain:      chain,
		contractAddress: contractAddress,
	}
}

func (fetcher PipFetcher) FetchPipValue(header core.Header) (string, error) {
	blockNumber := big.NewInt(header.BlockNumber)
	query := ethereum.FilterQuery{
		FromBlock: blockNumber,
		ToBlock:   blockNumber,
		Addresses: []common.Address{common.HexToAddress(fetcher.contractAddress)},
		Topics:    [][]common.Hash{{common.HexToHash(price_feeds.PipLogTopic0)}},
	}
	logs, err := fetcher.blockChain.GetEthLogsWithCustomQuery(query)
	if err != nil {
		return "", err
	}
	if len(logs) > 0 {
		return fetcher.getLogValue(logs, err)
	}
	return "", price_feeds.ErrNoMatchingLog
}

func (fetcher PipFetcher) getLogValue(logs []types.Log, err error) (string, error) {
	var (
		ret0 = new([32]byte)
		ret1 = new(bool)
	)
	var r = &[]interface{}{
		ret0,
		ret1,
	}
	err = fetcher.blockChain.FetchContractData(price_feeds.PipMedianizerABI, fetcher.contractAddress, price_feeds.PeekMethodName, nil, r, int64(logs[0].BlockNumber))
	if err != nil {
		return "", err
	}
	result := newResult(*ret0, *ret1)
	return result.Value.String(), nil
}

type Value [32]byte

type Peek struct {
	Value
	OK bool
}

func (value Value) String() string {
	bi := big.NewInt(0).SetBytes(value[:])
	return bi.String()
}

func newResult(value [32]byte, ok bool) *Peek {
	return &Peek{Value: value, OK: ok}
}
