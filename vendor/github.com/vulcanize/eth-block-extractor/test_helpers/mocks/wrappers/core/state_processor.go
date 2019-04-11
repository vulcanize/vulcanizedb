package core

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
)

type MockProcessor struct {
	passedBlock    *types.Block
	passedStateDB  *state.StateDB
	returnErr      error
	returnReceipts types.Receipts
	returnUsedGas  uint64
}

func NewMockProcessor() *MockProcessor {
	return &MockProcessor{
		passedBlock:    nil,
		passedStateDB:  nil,
		returnErr:      nil,
		returnReceipts: nil,
		returnUsedGas:  0,
	}
}

func (mp *MockProcessor) SetReturnErr(err error) {
	mp.returnErr = err
}

func (mp *MockProcessor) SetReturnReceipts(receipts types.Receipts) {
	mp.returnReceipts = receipts
}

func (mp *MockProcessor) SetReturnUsedGas(returnUsedGas uint64) {
	mp.returnUsedGas = returnUsedGas
}

func (mp *MockProcessor) Process(block *types.Block, stateDB *state.StateDB) (types.Receipts, []*types.Log, uint64, error) {
	mp.passedBlock = block
	mp.passedStateDB = stateDB
	return mp.returnReceipts, nil, mp.returnUsedGas, mp.returnErr
}

func (mp *MockProcessor) AssertProcessCalledWith(block *types.Block, stateDB *state.StateDB) {
	Expect(mp.passedBlock).To(Equal(block))
	Expect(mp.passedStateDB).To(Equal(stateDB))
}
