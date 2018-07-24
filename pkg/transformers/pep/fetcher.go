package pep

import (
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"math/big"
)

var (
	ErrMultipleLogs  = errors.New("multiple matching logs found in block")
	ErrNoMatchingLog = errors.New("no matching log")
)

type IPepFetcher interface {
	FetchPepValue(header core.Header) (string, error)
}

type PepFetcher struct {
	blockChain core.BlockChain
}

func NewPepFetcher(chain core.BlockChain) PepFetcher {
	return PepFetcher{
		blockChain: chain,
	}
}

func (fetcher PepFetcher) FetchPepValue(header core.Header) (string, error) {
	blockNumber := big.NewInt(header.BlockNumber)
	query := ethereum.FilterQuery{
		FromBlock: blockNumber,
		ToBlock:   blockNumber,
		Addresses: []common.Address{common.HexToAddress(PepAddress)},
		Topics:    [][]common.Hash{{common.HexToHash(LogValueTopic0)}},
	}
	logs, err := fetcher.blockChain.GetEthLogsWithCustomQuery(query)
	return fetcher.getLogValue(logs, err)
}

func (fetcher PepFetcher) getLogValue(logs []types.Log, err error) (string, error) {
	if err != nil {
		return "", err
	}
	if len(logs) < 1 {
		return "", ErrNoMatchingLog
	}
	if len(logs) > 1 {
		return "", ErrMultipleLogs
	}
	return string(logs[0].Data), nil
}
