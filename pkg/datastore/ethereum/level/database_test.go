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

package level_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/ethereum/level"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Level database", func() {
	Describe("Getting a block", func() {
		It("converts block number to uint64 and hash to common.Hash to fetch block from reader", func() {
			mockReader := fakes.NewMockLevelDatabaseReader()
			ldb := level.NewLevelDatabase(mockReader)
			blockHash := []byte{5, 4, 3, 2, 1}
			blockNumber := int64(12345)

			ldb.GetBlock(blockHash, blockNumber)

			expectedBlockHash := common.BytesToHash(blockHash)
			expectedBlockNumber := uint64(blockNumber)
			mockReader.AssertGetBlockCalledWith(expectedBlockHash, expectedBlockNumber)
		})
	})

	Describe("Getting a block hash", func() {
		It("converts block number to uint64 to fetch hash from reader", func() {
			mockReader := fakes.NewMockLevelDatabaseReader()
			ldb := level.NewLevelDatabase(mockReader)
			blockNumber := int64(12345)

			ldb.GetBlockHash(blockNumber)

			expectedBlockNumber := uint64(blockNumber)
			mockReader.AssertGetCanonicalHashCalledWith(expectedBlockNumber)
		})
	})

	Describe("Getting a block's receipts", func() {
		It("converts block number to uint64 and hash to common.Hash to fetch receipts from reader", func() {
			mockReader := fakes.NewMockLevelDatabaseReader()
			ldb := level.NewLevelDatabase(mockReader)
			blockHash := []byte{5, 4, 3, 2, 1}
			blockNumber := int64(12345)

			ldb.GetBlockReceipts(blockHash, blockNumber)

			expectedBlockHash := common.BytesToHash(blockHash)
			expectedBlockNumber := uint64(blockNumber)
			mockReader.AssertGetBlockReceiptsCalledWith(expectedBlockHash, expectedBlockNumber)
		})
	})

	Describe("Getting the latest block number", func() {
		It("invokes the database reader to get the latest block number by hash and converts result to int64", func() {
			mockReader := fakes.NewMockLevelDatabaseReader()
			fakeHash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
			mockReader.SetHeadBlockHashReturnHash(fakeHash)
			fakeBlockNumber := uint64(123456789)
			mockReader.SetReturnBlockNumber(fakeBlockNumber)
			ldb := level.NewLevelDatabase(mockReader)

			result := ldb.GetHeadBlockNumber()

			mockReader.AssertGetHeadBlockHashCalled()
			mockReader.AssertGetBlockNumberCalledWith(fakeHash)
			Expect(result).To(Equal(int64(fakeBlockNumber)))
		})
	})
})
