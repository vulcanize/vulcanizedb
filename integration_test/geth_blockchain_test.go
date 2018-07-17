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

	var blockchain *geth.BlockChain
	var inMemory *inmemory.InMemory

	BeforeEach(func() {
		blockchain = geth.NewBlockChain(test_config.InfuraClient.IPCPath)
		inMemory = inmemory.NewInMemory()
	})

	It("reads two blocks", func(done Done) {
		blocks := &inmemory.BlockRepository{InMemory: inMemory}
		lastBlock := blockchain.LastBlock()
		queriedBlocks := []int64{lastBlock.Int64() - 5, lastBlock.Int64() - 6}
		history.RetrieveAndUpdateBlocks(blockchain, blocks, queriedBlocks)
		Expect(blocks.BlockCount()).To(Equal(2))
		close(done)
	}, 30)

	It("retrieves the genesis block and first block", func(done Done) {
		genesisBlock, err := blockchain.GetBlockByNumber(int64(0))
		Expect(err).ToNot(HaveOccurred())
		firstBlock, err := blockchain.GetBlockByNumber(int64(1))
		Expect(err).ToNot(HaveOccurred())
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

	//Benchmarking test: remove skip to test performance of block retrieval
	XMeasure("retrieving n blocks", func(b Benchmarker) {
		b.Time("runtime", func() {
			var blocks []core.Block
			n := 10
			for i := 5327459; i > 5327459-n; i-- {
				block, err := blockchain.GetBlockByNumber(int64(i))
				Expect(err).ToNot(HaveOccurred())
				blocks = append(blocks, block)
			}
			Expect(len(blocks)).To(Equal(n))
		})
	}, 10)
})
