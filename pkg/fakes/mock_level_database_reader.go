package fakes

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
)

type MockLevelDatabaseReader struct {
	getBlockCalled               bool
	getBlockReceiptsCalled       bool
	getCanonicalHashCalled       bool
	passedHash                   common.Hash
	getCanonicalHashPassedNumber uint64
	getBlockPassedHash           common.Hash
	getBlockPassedNumber         uint64
	getBlockReceiptsPassedHash   common.Hash
	getBlockReceiptsPassedNumber uint64
	returnBlock                  *types.Block
	returnHash                   common.Hash
	returnReceipts               types.Receipts
}

func NewMockLevelDatabaseReader() *MockLevelDatabaseReader {
	return &MockLevelDatabaseReader{
		getBlockCalled:               false,
		getBlockReceiptsCalled:       false,
		getCanonicalHashCalled:       false,
		passedHash:                   common.Hash{},
		getCanonicalHashPassedNumber: 0,
		getBlockPassedHash:           common.Hash{},
		getBlockPassedNumber:         0,
		getBlockReceiptsPassedHash:   common.Hash{},
		getBlockReceiptsPassedNumber: 0,
		returnBlock:                  nil,
		returnHash:                   common.Hash{},
		returnReceipts:               nil,
	}
}

func (mldr *MockLevelDatabaseReader) SetReturnBlock(block *types.Block) {
	mldr.returnBlock = block
}

func (mldr *MockLevelDatabaseReader) SetReturnHash(hash common.Hash) {
	mldr.returnHash = hash
}

func (mldr *MockLevelDatabaseReader) SetReturnReceipts(receipts types.Receipts) {
	mldr.returnReceipts = receipts
}

func (mldr *MockLevelDatabaseReader) GetCanonicalHash(number uint64) common.Hash {
	mldr.getCanonicalHashCalled = true
	mldr.getCanonicalHashPassedNumber = number
	return mldr.returnHash
}

func (mldr *MockLevelDatabaseReader) GetBlock(hash common.Hash, number uint64) *types.Block {
	mldr.getBlockCalled = true
	mldr.getBlockPassedHash = hash
	mldr.getBlockPassedNumber = number
	return mldr.returnBlock
}

func (mldr *MockLevelDatabaseReader) GetBlockReceipts(hash common.Hash, number uint64) types.Receipts {
	mldr.getBlockReceiptsCalled = true
	mldr.getBlockReceiptsPassedHash = hash
	mldr.getBlockReceiptsPassedNumber = number
	return mldr.returnReceipts
}

func (mldr *MockLevelDatabaseReader) AssertGetCanonicalHashCalledWith(number uint64) {
	Expect(mldr.getCanonicalHashCalled).To(BeTrue())
	Expect(mldr.getCanonicalHashPassedNumber).To(Equal(number))
}

func (mldr *MockLevelDatabaseReader) AssertGetBlockCalledWith(hash common.Hash, number uint64) {
	Expect(mldr.getBlockCalled).To(BeTrue())
	Expect(mldr.getBlockPassedHash).To(Equal(hash))
	Expect(mldr.getBlockPassedNumber).To(Equal(number))
}

func (mldr *MockLevelDatabaseReader) AssertGetBlockReceiptsCalledWith(hash common.Hash, number uint64) {
	Expect(mldr.getBlockReceiptsCalled).To(BeTrue())
	Expect(mldr.getBlockReceiptsPassedHash).To(Equal(hash))
	Expect(mldr.getBlockReceiptsPassedNumber).To(Equal(number))
}
