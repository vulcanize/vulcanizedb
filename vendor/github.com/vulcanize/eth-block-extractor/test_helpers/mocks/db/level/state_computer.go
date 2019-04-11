package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
)

type MockStateComputer struct {
	computeBlockStateTriePassedCurrentBlock *types.Block
	computeBlockStateTriePassedParentBlock  *types.Block
	computeBlockStateTrieReturnErr          error
	computeBlockStateTrieReturnHash         common.Hash
}

func NewMockStateComputer() *MockStateComputer {
	return &MockStateComputer{
		computeBlockStateTriePassedCurrentBlock: nil,
		computeBlockStateTriePassedParentBlock:  nil,
		computeBlockStateTrieReturnErr:          nil,
		computeBlockStateTrieReturnHash:         common.Hash{},
	}
}

func (msc *MockStateComputer) SetComputeBlockStateTrieReturnErr(err error) {
	msc.computeBlockStateTrieReturnErr = err
}

func (msc *MockStateComputer) SetComputeBlockStateTrieReturnHash(hash common.Hash) {
	msc.computeBlockStateTrieReturnHash = hash
}

func (msc *MockStateComputer) ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) (common.Hash, error) {
	msc.computeBlockStateTriePassedCurrentBlock = currentBlock
	msc.computeBlockStateTriePassedParentBlock = parentBlock
	return msc.computeBlockStateTrieReturnHash, msc.computeBlockStateTrieReturnErr
}

func (msc *MockStateComputer) AssertComputeBlockStateTrieCalledWith(currentBlock *types.Block, parentBlock *types.Block) {
	Expect(msc.computeBlockStateTriePassedCurrentBlock).To(Equal(currentBlock))
	Expect(msc.computeBlockStateTriePassedParentBlock).To(Equal(parentBlock))
}
