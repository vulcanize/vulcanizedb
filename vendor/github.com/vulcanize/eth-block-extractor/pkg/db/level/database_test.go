package level_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/eth-block-extractor/pkg/db/level"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	level_wrapper "github.com/vulcanize/eth-block-extractor/test_helpers/mocks/db/level"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/core/rawdb"
)

var _ = Describe("Database", func() {
	Describe("Computing state trie nodes", func() {
		It("invokes state computer to build historical state", func() {
			mockStateComputer := level_wrapper.NewMockStateComputer()
			db := level.NewLevelDatabase(rawdb.NewMockAccessorsChain(), mockStateComputer, level_wrapper.NewMockStateTrieReader())
			currentBlock := &types.Block{}
			parentBlock := &types.Block{}

			_, err := db.ComputeBlockStateTrie(currentBlock, parentBlock)

			Expect(err).NotTo(HaveOccurred())
			mockStateComputer.AssertComputeBlockStateTrieCalledWith(currentBlock, parentBlock)
		})

		It("returns err if state computer returns err", func() {
			mockStateComputer := level_wrapper.NewMockStateComputer()
			mockStateComputer.SetComputeBlockStateTrieReturnErr(test_helpers.FakeError)
			db := level.NewLevelDatabase(rawdb.NewMockAccessorsChain(), mockStateComputer, level_wrapper.NewMockStateTrieReader())

			_, err := db.ComputeBlockStateTrie(&types.Block{}, &types.Block{})

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(test_helpers.FakeError))
		})
	})

	Describe("Getting block body data", func() {
		It("invokes the chain accessor to query for block hash by block number", func() {
			mockAccessorsChain := rawdb.NewMockAccessorsChain()
			db := level.NewLevelDatabase(mockAccessorsChain, level_wrapper.NewMockStateComputer(), level_wrapper.NewMockStateTrieReader())
			num := int64(123456)

			db.GetBlockBodyByBlockNumber(num)

			mockAccessorsChain.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the chain accessor to query for block body data", func() {
			mockAccessorsChain := rawdb.NewMockAccessorsChain()
			mockAccessorsChain.SetGetCanonicalHashReturnHash(test_helpers.FakeHash)
			db := level.NewLevelDatabase(mockAccessorsChain, level_wrapper.NewMockStateComputer(), level_wrapper.NewMockStateTrieReader())
			num := int64(123456)

			db.GetBlockBodyByBlockNumber(num)

			mockAccessorsChain.AssertGetBodyRLPCalledWith(test_helpers.FakeHash, uint64(num))
		})
	})

	Describe("Getting block", func() {
		It("invokes the chain accessor to query for block hash by block number", func() {
			mockAccessorsChain := rawdb.NewMockAccessorsChain()
			db := level.NewLevelDatabase(mockAccessorsChain, level_wrapper.NewMockStateComputer(), level_wrapper.NewMockStateTrieReader())
			num := int64(123456)

			db.GetBlockByBlockNumber(num)

			mockAccessorsChain.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the chain accessor to query for block", func() {
			mockAccessorsChain := rawdb.NewMockAccessorsChain()
			mockAccessorsChain.SetGetCanonicalHashReturnHash(test_helpers.FakeHash)
			db := level.NewLevelDatabase(mockAccessorsChain, level_wrapper.NewMockStateComputer(), level_wrapper.NewMockStateTrieReader())
			num := int64(123456)

			db.GetBlockByBlockNumber(num)

			mockAccessorsChain.AssertGetBlockCalledWith(test_helpers.FakeHash, uint64(num))
		})
	})

	Describe("Getting block header", func() {
		It("invokes the chain accessor to query for block hash by block number", func() {
			mockAccessorsChain := rawdb.NewMockAccessorsChain()
			db := level.NewLevelDatabase(mockAccessorsChain, level_wrapper.NewMockStateComputer(), level_wrapper.NewMockStateTrieReader())
			num := int64(123456)

			db.GetBlockHeaderByBlockNumber(num)

			mockAccessorsChain.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the chain accessor to query for block header", func() {
			mockAccessorsChain := rawdb.NewMockAccessorsChain()
			mockAccessorsChain.SetGetCanonicalHashReturnHash(test_helpers.FakeHash)
			db := level.NewLevelDatabase(mockAccessorsChain, level_wrapper.NewMockStateComputer(), level_wrapper.NewMockStateTrieReader())
			num := int64(123456)

			db.GetBlockHeaderByBlockNumber(num)

			mockAccessorsChain.AssertGetHeaderCalledWith(test_helpers.FakeHash, uint64(num))
		})
	})

	Describe("Getting raw block header data", func() {
		It("invokes the chain accessor to query for block hash by block number", func() {
			mockAccessorsChain := rawdb.NewMockAccessorsChain()
			db := level.NewLevelDatabase(mockAccessorsChain, level_wrapper.NewMockStateComputer(), level_wrapper.NewMockStateTrieReader())
			num := int64(123456)

			db.GetRawBlockHeaderByBlockNumber(num)

			mockAccessorsChain.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the chain accessor to query for block header data", func() {
			mockAccessorsChain := rawdb.NewMockAccessorsChain()
			mockAccessorsChain.SetGetCanonicalHashReturnHash(test_helpers.FakeHash)
			db := level.NewLevelDatabase(mockAccessorsChain, level_wrapper.NewMockStateComputer(), level_wrapper.NewMockStateTrieReader())
			num := int64(123456)

			db.GetRawBlockHeaderByBlockNumber(num)

			mockAccessorsChain.AssertGetHeaderRLPCalledWith(test_helpers.FakeHash, uint64(num))
		})
	})

	Describe("Getting block receipts", func() {
		It("invokes the chain accessor to query for block hash by block number", func() {
			mockAccessorsChain := rawdb.NewMockAccessorsChain()
			db := level.NewLevelDatabase(mockAccessorsChain, level_wrapper.NewMockStateComputer(), level_wrapper.NewMockStateTrieReader())
			num := int64(123456)

			db.GetBlockReceipts(num)

			mockAccessorsChain.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the chain accessor to query for block receipts", func() {
			mockAccessorsChain := rawdb.NewMockAccessorsChain()
			mockAccessorsChain.SetGetCanonicalHashReturnHash(test_helpers.FakeHash)
			db := level.NewLevelDatabase(mockAccessorsChain, level_wrapper.NewMockStateComputer(), level_wrapper.NewMockStateTrieReader())
			num := int64(123456)

			db.GetBlockReceipts(num)

			mockAccessorsChain.AssertGetBlockReceiptsCalledWith(test_helpers.FakeHash, uint64(num))
		})
	})

	Describe("Getting state trie nodes", func() {
		It("invokes the chain accessor to query for state trie data", func() {
			mockStateTrieReader := level_wrapper.NewMockStateTrieReader()
			db := level.NewLevelDatabase(rawdb.NewMockAccessorsChain(), level_wrapper.NewMockStateComputer(), mockStateTrieReader)
			root := common.HexToHash("abcde")

			_, _, err := db.GetStateAndStorageTrieNodes(root)

			Expect(err).NotTo(HaveOccurred())
			mockStateTrieReader.AssertGetStateAndStorageTrieNodesCalledWith(root)
		})
	})
})
