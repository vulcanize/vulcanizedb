package core

import "math/big"

type Blockchain interface {
	ContractDataFetcher
	GetBlockByNumber(blockNumber int64) (Block, error)
	GetLogs(contract Contract, startingBlockNumber *big.Int, endingBlockNumber *big.Int) ([]Log, error)
	LastBlock() *big.Int
	Node() Node
}

type ContractDataFetcher interface {
	FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error
}
