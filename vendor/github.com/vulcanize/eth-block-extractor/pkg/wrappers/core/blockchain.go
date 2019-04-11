package core

import (
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

type GethCoreBlockChain interface {
	BlockChain() *core.BlockChain
	Config() *params.ChainConfig
	Engine() consensus.Engine
}

type BlockChain struct {
	blockChain *core.BlockChain
}

func NewBlockChain(databaseConnection ethdb.Database) (*BlockChain, error) {
	blockchain, err := core.NewBlockChain(databaseConnection, nil, params.MainnetChainConfig, ethash.NewFaker(), vm.Config{})
	if err != nil {
		return nil, err
	}
	return &BlockChain{blockChain: blockchain}, nil
}

func (sb *BlockChain) BlockChain() *core.BlockChain {
	return sb.blockChain
}

func (sb *BlockChain) Config() *params.ChainConfig {
	return sb.blockChain.Config()
}

func (sb *BlockChain) Engine() consensus.Engine {
	return sb.blockChain.Engine()
}
