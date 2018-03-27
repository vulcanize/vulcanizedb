package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/inmemory"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Reading from the Geth blockchain", func() {

	var blockchain *geth.Blockchain
	var inMemory *inmemory.InMemory

	BeforeEach(func() {
		blockchain = geth.NewBlockchain(test_config.InfuraClient.IPCPath)
		inMemory = inmemory.NewInMemory()
	})

	It("reads two blocks", func(done Done) {
		blocks := &inmemory.BlockRepository{InMemory: inMemory}
		validator := history.NewBlockValidator(blockchain, blocks, 2)
		validator.ValidateBlocks()
		Expect(blocks.BlockCount()).To(Equal(2))
		close(done)
	}, 30)

	It("retrieves the genesis block and first block", func(done Done) {
		genesisBlock := blockchain.GetBlockByNumber(int64(0))
		firstBlock := blockchain.GetBlockByNumber(int64(1))
		lastBlockNumber := blockchain.LastBlock()

		Expect(genesisBlock.Number).To(Equal(int64(0)))
		Expect(firstBlock.Number).To(Equal(int64(1)))
		Expect(lastBlockNumber.Int64()).To(BeNumerically(">", 0))
		close(done)
	}, 15)

	It("retrieves the node info", func(done Done) {
		node := blockchain.Node()
		mainnetID := float64(1)

		Expect(node.GenesisBlock).ToNot(BeNil())
		Expect(node.NetworkID).To(Equal(mainnetID))
		Expect(len(node.ID)).ToNot(BeZero())
		Expect(node.ClientName).ToNot(BeZero())

		close(done)
	}, 15)

	XMeasure("retrieving n blocks", func(b Benchmarker) {
		b.Time("runtime", func() {
			var blocks []core.Block
			n := 10
			for i := 5327459; i > 5327459-n; i-- {
				block := blockchain.GetBlockByNumber(int64(i))
				blocks = append(blocks, block)
			}
			Expect(len(blocks)).To(Equal(n))
		})
	}, 10)
})
