package core

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
)

type GethBlockValidator interface {
	ValidateState(block, parent *types.Block, state *state.StateDB, receipts types.Receipts, usedGas uint64) error
}

type BlockValidator struct {
	validator core.Validator
}

func NewBlockValidator(blockChain BlockChain) *BlockValidator {
	validator := core.NewBlockValidator(blockChain.Config(), blockChain.BlockChain(), blockChain.Engine())
	return &BlockValidator{validator: validator}
}

func (sv *BlockValidator) ValidateState(block, parent *types.Block, state *state.StateDB, receipts types.Receipts, usedGas uint64) error {
	return sv.validator.ValidateState(block, parent, state, receipts, usedGas)
}
