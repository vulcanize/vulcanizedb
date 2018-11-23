// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package fakes

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
)

type MockLevelDatabaseReader struct {
	getBlockCalled               bool
	getBlockNumberCalled         bool
	getBlockNumberPassedHash     common.Hash
	getBlockPassedHash           common.Hash
	getBlockPassedNumber         uint64
	getBlockReceiptsCalled       bool
	getBlockReceiptsPassedHash   common.Hash
	getBlockReceiptsPassedNumber uint64
	getCanonicalHashCalled       bool
	getCanonicalHashPassedNumber uint64
	getCanonicalHashReturnHash   common.Hash
	getHeadBlockHashCalled       bool
	getHeadBlockHashReturnHash   common.Hash
	passedHash                   common.Hash
	returnBlock                  *types.Block
	returnBlockNumber            uint64
	returnReceipts               types.Receipts
}

func NewMockLevelDatabaseReader() *MockLevelDatabaseReader {
	return &MockLevelDatabaseReader{
		getBlockCalled:               false,
		getBlockNumberCalled:         false,
		getBlockNumberPassedHash:     common.Hash{},
		getBlockPassedHash:           common.Hash{},
		getBlockPassedNumber:         0,
		getBlockReceiptsCalled:       false,
		getBlockReceiptsPassedHash:   common.Hash{},
		getBlockReceiptsPassedNumber: 0,
		getCanonicalHashCalled:       false,
		getCanonicalHashPassedNumber: 0,
		getCanonicalHashReturnHash:   common.Hash{},
		getHeadBlockHashCalled:       false,
		getHeadBlockHashReturnHash:   common.Hash{},
		passedHash:                   common.Hash{},
		returnBlock:                  nil,
		returnBlockNumber:            0,
		returnReceipts:               nil,
	}
}

func (mldr *MockLevelDatabaseReader) SetReturnBlock(block *types.Block) {
	mldr.returnBlock = block
}

func (mldr *MockLevelDatabaseReader) SetReturnBlockNumber(n uint64) {
	mldr.returnBlockNumber = n
}

func (mldr *MockLevelDatabaseReader) SetGetCanonicalHashReturnHash(hash common.Hash) {
	mldr.getCanonicalHashReturnHash = hash
}

func (mldr *MockLevelDatabaseReader) SetHeadBlockHashReturnHash(hash common.Hash) {
	mldr.getHeadBlockHashReturnHash = hash
}

func (mldr *MockLevelDatabaseReader) SetReturnReceipts(receipts types.Receipts) {
	mldr.returnReceipts = receipts
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

func (mldr *MockLevelDatabaseReader) GetBlockNumber(hash common.Hash) *uint64 {
	mldr.getBlockNumberCalled = true
	mldr.getBlockNumberPassedHash = hash
	return &mldr.returnBlockNumber
}

func (mldr *MockLevelDatabaseReader) GetCanonicalHash(number uint64) common.Hash {
	mldr.getCanonicalHashCalled = true
	mldr.getCanonicalHashPassedNumber = number
	return mldr.getCanonicalHashReturnHash
}

func (mldr *MockLevelDatabaseReader) GetHeadBlockHash() common.Hash {
	mldr.getHeadBlockHashCalled = true
	return mldr.getHeadBlockHashReturnHash
}

func (mldr *MockLevelDatabaseReader) AssertGetBlockCalledWith(hash common.Hash, number uint64) {
	Expect(mldr.getBlockCalled).To(BeTrue())
	Expect(mldr.getBlockPassedHash).To(Equal(hash))
	Expect(mldr.getBlockPassedNumber).To(Equal(number))
}

func (mldr *MockLevelDatabaseReader) AssertGetBlockNumberCalledWith(hash common.Hash) {
	Expect(mldr.getBlockNumberCalled).To(BeTrue())
	Expect(mldr.getBlockNumberPassedHash).To(Equal(hash))
}

func (mldr *MockLevelDatabaseReader) AssertGetBlockReceiptsCalledWith(hash common.Hash, number uint64) {
	Expect(mldr.getBlockReceiptsCalled).To(BeTrue())
	Expect(mldr.getBlockReceiptsPassedHash).To(Equal(hash))
	Expect(mldr.getBlockReceiptsPassedNumber).To(Equal(number))
}

func (mldr *MockLevelDatabaseReader) AssertGetCanonicalHashCalledWith(number uint64) {
	Expect(mldr.getCanonicalHashCalled).To(BeTrue())
	Expect(mldr.getCanonicalHashPassedNumber).To(Equal(number))
}

func (mldr *MockLevelDatabaseReader) AssertGetHeadBlockHashCalled() {
	Expect(mldr.getHeadBlockHashCalled).To(BeTrue())
}
