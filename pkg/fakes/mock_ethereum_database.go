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
	getHeadBlockNumberCalled       bool
	getHeadBlockNumberReturnVal    int64
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
		getHeadBlockNumberCalled:       false,
		getHeadBlockNumberReturnVal:    0,
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

func (med *MockEthereumDatabase) GetHeadBlockNumber() int64 {
	med.getHeadBlockNumberCalled = true
	return med.getHeadBlockNumberReturnVal
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
