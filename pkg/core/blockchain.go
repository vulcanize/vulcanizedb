package core

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type BlockChain interface {
	ContractDataFetcher
	GetBlockByNumber(blockNumber int64) (Block, error)
	GetHeaderByNumber(blockNumber int64) (Header, error)
	GetLogs(contract Contract, startingBlockNumber *big.Int, endingBlockNumber *big.Int) ([]Log, error)
	GetEthLogsWithCustomQuery(query ethereum.FilterQuery) ([]types.Log, error)
	LastBlock() *big.Int
	Node() Node
}

type ContractDataFetcher interface {
	FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error
}
