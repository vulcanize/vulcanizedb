package level_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/ethereum/level"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Level database", func() {
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
})
