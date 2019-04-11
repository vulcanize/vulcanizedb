package core

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
)

type MockValidator struct {
	passedBlock       *types.Block
	passedParentBlock *types.Block
	passedStateDB     *state.StateDB
	passedReceipts    types.Receipts
	passedUsedGas     uint64
	returnErr         error
}

func NewMockValidator() *MockValidator {
	return &MockValidator{
		passedBlock:       nil,
		passedParentBlock: nil,
		passedStateDB:     nil,
		passedReceipts:    nil,
		passedUsedGas:     0,
		returnErr:         nil,
	}
}

func (mv *MockValidator) SetReturnErr(err error) {
	mv.returnErr = err
}

func (mv *MockValidator) ValidateState(block, parent *types.Block, state *state.StateDB, receipts types.Receipts, usedGas uint64) error {
	mv.passedBlock = block
	mv.passedParentBlock = parent
	mv.passedStateDB = state
	mv.passedReceipts = receipts
	mv.passedUsedGas = usedGas
	return mv.returnErr
}

func (mv *MockValidator) AssertValidateStateCalledWith(block, parent *types.Block, stateDB *state.StateDB, receipts types.Receipts, usedGas uint64) {
	Expect(mv.passedBlock).To(Equal(block))
	Expect(mv.passedParentBlock).To(Equal(parent))
	Expect(mv.passedStateDB).To(Equal(stateDB))
	Expect(mv.passedReceipts).To(Equal(receipts))
	Expect(mv.passedUsedGas).To(Equal(usedGas))
}
