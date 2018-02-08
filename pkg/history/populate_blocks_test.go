package history_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"github.com/vulcanize/vulcanizedb/pkg/repositories/inmemory"
)

var _ = Describe("Populating blocks", func() {
	var inmemoryDB *inmemory.InMemory
	var blocksRepo *inmemory.Blocks

	BeforeEach(func() {
		inmemoryDB = inmemory.NewInMemory()
		blocksRepo = &inmemory.Blocks{InMemory: inmemoryDB}
	})

	It("fills in the only missing block (Number 1)", func() {
		blocks := []core.Block{
			{Number: 1},
			{Number: 2},
		}
		blockchain := fakes.NewBlockchainWithBlocks(blocks)

		blocksRepo.CreateOrUpdateBlock(core.Block{Number: 2})

		blocksAdded := history.PopulateMissingBlocks(blockchain, blocksRepo, 1)
		_, err := blocksRepo.GetBlock(1)

		Expect(blocksAdded).To(Equal(1))
		Expect(err).ToNot(HaveOccurred())
	})

	It("fills in the three missing blocks (Numbers: 5,8,10)", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 4},
			{Number: 5},
			{Number: 6},
			{Number: 7},
			{Number: 8},
			{Number: 9},
			{Number: 10},
			{Number: 11},
			{Number: 12},
			{Number: 13},
		})
		blocksRepo.CreateOrUpdateBlock(core.Block{Number: 1})
		blocksRepo.CreateOrUpdateBlock(core.Block{Number: 2})
		blocksRepo.CreateOrUpdateBlock(core.Block{Number: 3})
		blocksRepo.CreateOrUpdateBlock(core.Block{Number: 6})
		blocksRepo.CreateOrUpdateBlock(core.Block{Number: 7})
		blocksRepo.CreateOrUpdateBlock(core.Block{Number: 9})
		blocksRepo.CreateOrUpdateBlock(core.Block{Number: 11})
		blocksRepo.CreateOrUpdateBlock(core.Block{Number: 12})

		blocksAdded := history.PopulateMissingBlocks(blockchain, blocksRepo, 5)

		Expect(blocksAdded).To(Equal(3))
		Expect(blocksRepo.BlockCount()).To(Equal(11))
		_, err := blocksRepo.GetBlock(4)
		Expect(err).To(HaveOccurred())
		_, err = blocksRepo.GetBlock(5)
		Expect(err).ToNot(HaveOccurred())
		_, err = blocksRepo.GetBlock(8)
		Expect(err).ToNot(HaveOccurred())
		_, err = blocksRepo.GetBlock(10)
		Expect(err).ToNot(HaveOccurred())
		_, err = blocksRepo.GetBlock(13)
		Expect(err).To(HaveOccurred())
	})

	It("returns the number of blocks created", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 4},
			{Number: 5},
			{Number: 6},
		})
		blocksRepo.CreateOrUpdateBlock(core.Block{Number: 3})
		blocksRepo.CreateOrUpdateBlock(core.Block{Number: 6})

		numberOfBlocksCreated := history.PopulateMissingBlocks(blockchain, blocksRepo, 3)

		Expect(numberOfBlocksCreated).To(Equal(2))
	})

	It("updates the repository with a range of blocks w/in the range ", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 1},
			{Number: 2},
			{Number: 3},
			{Number: 4},
			{Number: 5},
		})

		history.RetrieveAndUpdateBlocks(blockchain, blocksRepo, history.MakeRange(2, 5))
		Expect(blocksRepo.BlockCount()).To(Equal(3))
		Expect(blocksRepo.CreateOrUpdateBlockCallCount).To(Equal(3))
	})

})
