package history_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"math/big"
)

var _ = Describe("Populating blocks", func() {
	var blockRepository *fakes.MockBlockRepository

	BeforeEach(func() {
		blockRepository = fakes.NewMockBlockRepository()
	})

	It("fills in the only missing block (BlockNumber 1)", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(2))
		blockRepository.SetMissingBlockNumbersReturnArray([]int64{2})

		blocksAdded := history.PopulateMissingBlocks(blockChain, blockRepository, 1)
		_, err := blockRepository.GetBlock(1)

		Expect(blocksAdded).To(Equal(1))
		Expect(err).ToNot(HaveOccurred())
	})

	It("fills in the three missing blocks (Numbers: 5,8,10)", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(13))
		blockRepository.SetMissingBlockNumbersReturnArray([]int64{5, 8, 10})

		blocksAdded := history.PopulateMissingBlocks(blockChain, blockRepository, 5)

		Expect(blocksAdded).To(Equal(3))
		blockRepository.AssertCreateOrUpdateBlocksCallCountAndBlockNumbersEquals(3, []int64{5, 8, 10})
	})

	It("returns the number of blocks created", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(6))
		blockRepository.SetMissingBlockNumbersReturnArray([]int64{4, 5})

		numberOfBlocksCreated := history.PopulateMissingBlocks(blockChain, blockRepository, 3)

		Expect(numberOfBlocksCreated).To(Equal(2))
	})

	It("updates the repository with a range of blocks w/in the range ", func() {
		blockChain := fakes.NewMockBlockChain()

		history.RetrieveAndUpdateBlocks(blockChain, blockRepository, history.MakeRange(2, 5))

		blockRepository.AssertCreateOrUpdateBlocksCallCountAndBlockNumbersEquals(4, []int64{2, 3, 4, 5})
	})

	It("does not call repository create block when there is an error", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetGetBlockByNumberErr(fakes.FakeError)
		blocks := history.MakeRange(1, 10)

		history.RetrieveAndUpdateBlocks(blockChain, blockRepository, blocks)

		blockRepository.AssertCreateOrUpdateBlockCallCountEquals(0)
	})
})
