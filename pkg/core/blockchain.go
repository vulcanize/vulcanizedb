package core

import "math/big"

type Blockchain interface {
	GetBlockByNumber(blockNumber int64) Block
	Node() Node
	SubscribeToBlocks(blocks chan Block)
	StartListening()
	StopListening()
	GetAttributes(contract Contract) (ContractAttributes, error)
	GetAttribute(contract Contract, attributeName string, blockNumber *big.Int) (interface{}, error)
}
