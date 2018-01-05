package core

import "math/big"

type Blockchain interface {
	GetBlockByNumber(blockNumber int64) Block
	LastBlock() *big.Int
	Node() Node
	GetAttributes(contract Contract) (ContractAttributes, error)
	GetAttribute(contract Contract, attributeName string, blockNumber *big.Int) (interface{}, error)
	GetLogs(contract Contract, startingBlockNumber *big.Int, endingBlockNumber *big.Int) ([]Log, error)
}
