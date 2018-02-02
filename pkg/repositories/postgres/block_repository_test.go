package postgres_test

import (
	"math/big"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/repositories/postgres"
)

var _ = Describe("Saving blocks", func() {
	var repository repositories.BlockRepository
	BeforeEach(func() {
		node := core.Node{
			GenesisBlock: "GENESIS",
			NetworkId:    1,
			Id:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
			ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
		}
		repository = postgres.BuildRepository(node)
	})

	It("associates blocks to a node", func() {
		block := core.Block{
			Number: 123,
		}
		repository.CreateOrUpdateBlock(block)
		nodeTwo := core.Node{
			GenesisBlock: "0x456",
			NetworkId:    1,
			Id:           "x123456",
			ClientName:   "Geth",
		}
		repositoryTwo := postgres.BuildRepository(nodeTwo)

		_, err := repositoryTwo.FindBlockByNumber(123)
		Expect(err).To(HaveOccurred())
	})

	It("saves the attributes of the block", func() {
		blockNumber := int64(123)
		gasLimit := int64(1000000)
		gasUsed := int64(10)
		blockHash := "x123"
		blockParentHash := "x456"
		blockNonce := "0x881db2ca900682e9a9"
		miner := "x123"
		extraData := "xextraData"
		blockTime := int64(1508981640)
		uncleHash := "x789"
		blockSize := int64(1000)
		difficulty := int64(10)
		blockReward := float64(5.132)
		unclesReward := float64(3.580)
		block := core.Block{
			Reward:       blockReward,
			Difficulty:   difficulty,
			GasLimit:     gasLimit,
			GasUsed:      gasUsed,
			Hash:         blockHash,
			ExtraData:    extraData,
			Nonce:        blockNonce,
			Miner:        miner,
			Number:       blockNumber,
			ParentHash:   blockParentHash,
			Size:         blockSize,
			Time:         blockTime,
			UncleHash:    uncleHash,
			UnclesReward: unclesReward,
		}

		repository.CreateOrUpdateBlock(block)

		savedBlock, err := repository.FindBlockByNumber(blockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(savedBlock.Reward).To(Equal(blockReward))
		Expect(savedBlock.Difficulty).To(Equal(difficulty))
		Expect(savedBlock.GasLimit).To(Equal(gasLimit))
		Expect(savedBlock.GasUsed).To(Equal(gasUsed))
		Expect(savedBlock.Hash).To(Equal(blockHash))
		Expect(savedBlock.Nonce).To(Equal(blockNonce))
		Expect(savedBlock.Miner).To(Equal(miner))
		Expect(savedBlock.ExtraData).To(Equal(extraData))
		Expect(savedBlock.Number).To(Equal(blockNumber))
		Expect(savedBlock.ParentHash).To(Equal(blockParentHash))
		Expect(savedBlock.Size).To(Equal(blockSize))
		Expect(savedBlock.Time).To(Equal(blockTime))
		Expect(savedBlock.UncleHash).To(Equal(uncleHash))
		Expect(savedBlock.UnclesReward).To(Equal(unclesReward))
	})

	It("does not find a block when searching for a number that does not exist", func() {
		_, err := repository.FindBlockByNumber(111)

		Expect(err).To(HaveOccurred())
	})

	It("saves one transaction associated to the block", func() {
		block := core.Block{
			Number:       123,
			Transactions: []core.Transaction{{}},
		}

		repository.CreateOrUpdateBlock(block)

		savedBlock, _ := repository.FindBlockByNumber(123)
		Expect(len(savedBlock.Transactions)).To(Equal(1))
	})

	It("saves two transactions associated to the block", func() {
		block := core.Block{
			Number:       123,
			Transactions: []core.Transaction{{}, {}},
		}

		repository.CreateOrUpdateBlock(block)

		savedBlock, _ := repository.FindBlockByNumber(123)
		Expect(len(savedBlock.Transactions)).To(Equal(2))
	})

	It(`replaces blocks and transactions associated to the block
			when a more new block is in conflict (same block number + nodeid)`, func() {
		blockOne := core.Block{
			Number:       123,
			Hash:         "xabc",
			Transactions: []core.Transaction{{Hash: "x123"}, {Hash: "x345"}},
		}
		blockTwo := core.Block{
			Number:       123,
			Hash:         "xdef",
			Transactions: []core.Transaction{{Hash: "x678"}, {Hash: "x9ab"}},
		}

		repository.CreateOrUpdateBlock(blockOne)
		repository.CreateOrUpdateBlock(blockTwo)

		savedBlock, _ := repository.FindBlockByNumber(123)
		Expect(len(savedBlock.Transactions)).To(Equal(2))
		Expect(savedBlock.Transactions[0].Hash).To(Equal("x678"))
		Expect(savedBlock.Transactions[1].Hash).To(Equal("x9ab"))
	})

	It(`does not replace blocks when block number is not unique
			     but block number + node id is`, func() {
		blockOne := core.Block{
			Number:       123,
			Transactions: []core.Transaction{{Hash: "x123"}, {Hash: "x345"}},
		}
		blockTwo := core.Block{
			Number:       123,
			Transactions: []core.Transaction{{Hash: "x678"}, {Hash: "x9ab"}},
		}
		repository.CreateOrUpdateBlock(blockOne)
		nodeTwo := core.Node{
			GenesisBlock: "0x456",
			NetworkId:    1,
		}
		repositoryTwo := postgres.BuildRepository(nodeTwo)

		repository.CreateOrUpdateBlock(blockOne)
		repositoryTwo.CreateOrUpdateBlock(blockTwo)
		retrievedBlockOne, _ := repository.FindBlockByNumber(123)
		retrievedBlockTwo, _ := repositoryTwo.FindBlockByNumber(123)

		Expect(retrievedBlockOne.Transactions[0].Hash).To(Equal("x123"))
		Expect(retrievedBlockTwo.Transactions[0].Hash).To(Equal("x678"))
	})

	It("saves the attributes associated to a transaction", func() {
		gasLimit := int64(5000)
		gasPrice := int64(3)
		nonce := uint64(10000)
		to := "1234567890"
		from := "0987654321"
		var value = new(big.Int)
		value.SetString("34940183920000000000", 10)
		inputData := "0xf7d8c8830000000000000000000000000000000000000000000000000000000000037788000000000000000000000000000000000000000000000000000000000003bd14"
		transaction := core.Transaction{
			Hash:     "x1234",
			GasPrice: gasPrice,
			GasLimit: gasLimit,
			Nonce:    nonce,
			To:       to,
			From:     from,
			Value:    value.String(),
			Data:     inputData,
		}
		block := core.Block{
			Number:       123,
			Transactions: []core.Transaction{transaction},
		}

		repository.CreateOrUpdateBlock(block)

		savedBlock, _ := repository.FindBlockByNumber(123)
		Expect(len(savedBlock.Transactions)).To(Equal(1))
		savedTransaction := savedBlock.Transactions[0]
		Expect(savedTransaction.Data).To(Equal(transaction.Data))
		Expect(savedTransaction.Hash).To(Equal(transaction.Hash))
		Expect(savedTransaction.To).To(Equal(to))
		Expect(savedTransaction.From).To(Equal(from))
		Expect(savedTransaction.Nonce).To(Equal(nonce))
		Expect(savedTransaction.GasLimit).To(Equal(gasLimit))
		Expect(savedTransaction.GasPrice).To(Equal(gasPrice))
		Expect(savedTransaction.Value).To(Equal(value.String()))
	})

	Describe("The missing block numbers", func() {
		It("is empty the starting block number is the highest known block number", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 1})

			Expect(len(repository.MissingBlockNumbers(1, 1))).To(Equal(0))
		})

		It("is the only missing block number", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 2})

			Expect(repository.MissingBlockNumbers(1, 2)).To(Equal([]int64{1}))
		})

		It("is both missing block numbers", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 3})

			Expect(repository.MissingBlockNumbers(1, 3)).To(Equal([]int64{1, 2}))
		})

		It("goes back to the starting block number", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 6})

			Expect(repository.MissingBlockNumbers(4, 6)).To(Equal([]int64{4, 5}))
		})

		It("only includes missing block numbers", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 4})
			repository.CreateOrUpdateBlock(core.Block{Number: 6})

			Expect(repository.MissingBlockNumbers(4, 6)).To(Equal([]int64{5}))
		})

		It("is a list with multiple gaps", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 4})
			repository.CreateOrUpdateBlock(core.Block{Number: 5})
			repository.CreateOrUpdateBlock(core.Block{Number: 8})
			repository.CreateOrUpdateBlock(core.Block{Number: 10})

			Expect(repository.MissingBlockNumbers(3, 10)).To(Equal([]int64{3, 6, 7, 9}))
		})

		It("returns empty array when lower bound exceeds upper bound", func() {
			Expect(repository.MissingBlockNumbers(10000, 1)).To(Equal([]int64{}))
		})

		It("only returns requested range even when other gaps exist", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 3})
			repository.CreateOrUpdateBlock(core.Block{Number: 8})

			Expect(repository.MissingBlockNumbers(1, 5)).To(Equal([]int64{1, 2, 4, 5}))
		})
	})

	Describe("The block status", func() {
		It("sets the status of blocks within n-20 of chain HEAD as final", func() {
			blockNumberOfChainHead := 25
			for i := 0; i < blockNumberOfChainHead; i++ {
				repository.CreateOrUpdateBlock(core.Block{Number: int64(i), Hash: strconv.Itoa(i)})
			}

			repository.SetBlocksStatus(int64(blockNumberOfChainHead))

			blockOne, err := repository.FindBlockByNumber(1)
			Expect(err).ToNot(HaveOccurred())
			Expect(blockOne.IsFinal).To(Equal(true))
			blockTwo, err := repository.FindBlockByNumber(24)
			Expect(err).ToNot(HaveOccurred())
			Expect(blockTwo.IsFinal).To(BeFalse())
		})

	})
})
