package fakes

import (
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/core/types"
)

type MockEthereumDatabase struct {
	getBlockCalled                 bool
	getBlockPassedHash             []byte
	getBlockPassedNumber           int64
	getBlockReturnBlock            *types.Block
	getBlockHashCalled             bool
	getBlockHashPassedNumber       int64
	getBlockHashReturnHash         []byte
	getBlockReceiptsCalled         bool
	getBlockReceiptsPassedHash     []byte
	getBlockReceiptsPassedNumber   int64
	getBlockReceiptsReturnReceipts types.Receipts
}

func NewMockEthereumDatabase() *MockEthereumDatabase {
	return &MockEthereumDatabase{
		getBlockCalled:                 false,
		getBlockPassedHash:             nil,
		getBlockPassedNumber:           0,
		getBlockReturnBlock:            nil,
		getBlockHashCalled:             false,
		getBlockHashPassedNumber:       0,
		getBlockHashReturnHash:         nil,
		getBlockReceiptsCalled:         false,
		getBlockReceiptsPassedHash:     nil,
		getBlockReceiptsPassedNumber:   0,
		getBlockReceiptsReturnReceipts: nil,
	}
}

func (med *MockEthereumDatabase) SetReturnBlock(block *types.Block) {
	med.getBlockReturnBlock = block
}

func (med *MockEthereumDatabase) SetReturnHash(hash []byte) {
	med.getBlockHashReturnHash = hash
}

func (med *MockEthereumDatabase) SetReturnReceipts(receipts types.Receipts) {
	med.getBlockReceiptsReturnReceipts = receipts
}

func (med *MockEthereumDatabase) GetBlock(hash []byte, blockNumber int64) *types.Block {
	med.getBlockCalled = true
	med.getBlockPassedHash = hash
	med.getBlockPassedNumber = blockNumber
	return med.getBlockReturnBlock
}

func (med *MockEthereumDatabase) GetBlockHash(blockNumber int64) []byte {
	med.getBlockHashCalled = true
	med.getBlockHashPassedNumber = blockNumber
	return med.getBlockHashReturnHash
}

func (med *MockEthereumDatabase) GetBlockReceipts(blockHash []byte, blockNumber int64) types.Receipts {
	med.getBlockReceiptsCalled = true
	med.getBlockReceiptsPassedHash = blockHash
	med.getBlockReceiptsPassedNumber = blockNumber
	return med.getBlockReceiptsReturnReceipts
}

func (med *MockEthereumDatabase) AssertGetBlockCalledWith(hash []byte, blockNumber int64) {
	Expect(med.getBlockCalled).To(BeTrue())
	Expect(med.getBlockPassedHash).To(Equal(hash))
	Expect(med.getBlockPassedNumber).To(Equal(blockNumber))
}

func (med *MockEthereumDatabase) AssertGetBlockHashCalledWith(blockNumber int64) {
	Expect(med.getBlockHashCalled).To(BeTrue())
	Expect(med.getBlockHashPassedNumber).To(Equal(blockNumber))
}

func (med *MockEthereumDatabase) AssertGetBlockReceiptsCalledWith(blockHash []byte, blockNumber int64) {
	Expect(med.getBlockReceiptsCalled).To(BeTrue())
	Expect(med.getBlockReceiptsPassedHash).To(Equal(blockHash))
	Expect(med.getBlockReceiptsPassedNumber).To(Equal(blockNumber))
}
