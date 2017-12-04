package core

import "math/big"

type Blockchain interface {
	GetBlockByNumber(blockNumber int64) Block
	SubscribeToBlocks(blocks chan Block)
	StartListening()
	StopListening()
	GetContract(contractHash string) (Contract, error)
	GetAttribute(contract Contract, attributeName string, blockNumber *big.Int) (interface{}, error)
}
