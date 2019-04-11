package core

import (
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

type GethStateProcessor interface {
	Process(block *types.Block, statedb *state.StateDB) (types.Receipts, []*types.Log, uint64, error)
}

type StateProcessor struct {
	processor *core.StateProcessor
}

func NewStateProcessor(blockChain BlockChain) *StateProcessor {
	processor := core.NewStateProcessor(params.MainnetChainConfig, blockChain.BlockChain(), ethash.NewFaker())
	return &StateProcessor{processor: processor}
}

func (sp *StateProcessor) Process(block *types.Block, statedb *state.StateDB) (types.Receipts, []*types.Log, uint64, error) {
	return sp.processor.Process(block, statedb, vm.Config{})
}
