package repositories_test

import (
	"fmt"
	"strings"

	"github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
	"github.com/8thlight/vulcanizedb/pkg/repositories/testing"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repositories", func() {

	AssertRepositoryBehavior := func(buildRepository func() repositories.Repository) {
		var repository repositories.Repository

		BeforeEach(func() {
			repository = buildRepository()
		})

		Describe("Saving blocks", func() {
			It("starts with no blocks", func() {
				count := repository.BlockCount()
				Expect(count).Should(Equal(0))
			})

			It("increments the block count", func() {
				block := core.Block{Number: 123}

				repository.CreateBlock(block)

				Expect(repository.BlockCount()).To(Equal(1))
			})

			It("saves the attributes of the block", func() {
				blockNumber := int64(123)
				gasLimit := int64(1000000)
				gasUsed := int64(10)
				blockHash := "x123"
				blockParentHash := "x456"
				blockNonce := "0x881db2ca900682e9a9"
				blockTime := int64(1508981640)
				uncleHash := "x789"
				blockSize := int64(1000)
				difficulty := int64(10)
				block := core.Block{
					Difficulty: difficulty,
					GasLimit:   gasLimit,
					GasUsed:    gasUsed,
					Hash:       blockHash,
					Nonce:      blockNonce,
					Number:     blockNumber,
					ParentHash: blockParentHash,
					Size:       blockSize,
					Time:       blockTime,
					UncleHash:  uncleHash,
				}

				repository.CreateBlock(block)

				savedBlock := repository.FindBlockByNumber(blockNumber)
				Expect(savedBlock.Difficulty).To(Equal(difficulty))
				Expect(savedBlock.GasLimit).To(Equal(gasLimit))
				Expect(savedBlock.GasUsed).To(Equal(gasUsed))
				Expect(savedBlock.Hash).To(Equal(blockHash))
				Expect(savedBlock.Nonce).To(Equal(blockNonce))
				Expect(savedBlock.Number).To(Equal(blockNumber))
				Expect(savedBlock.ParentHash).To(Equal(blockParentHash))
				Expect(savedBlock.Size).To(Equal(blockSize))
				Expect(savedBlock.Time).To(Equal(blockTime))
				Expect(savedBlock.UncleHash).To(Equal(uncleHash))
			})

			It("does not find a block when searching for a number that does not exist", func() {
				savedBlock := repository.FindBlockByNumber(111)

				Expect(savedBlock).To(BeNil())
			})

			It("saves one transaction associated to the block", func() {
				block := core.Block{
					Number:       123,
					Transactions: []core.Transaction{{}},
				}

				repository.CreateBlock(block)

				savedBlock := repository.FindBlockByNumber(123)
				Expect(len(savedBlock.Transactions)).To(Equal(1))
			})

			It("saves two transactions associated to the block", func() {
				block := core.Block{
					Number:       123,
					Transactions: []core.Transaction{{}, {}},
				}

				repository.CreateBlock(block)

				savedBlock := repository.FindBlockByNumber(123)
				Expect(len(savedBlock.Transactions)).To(Equal(2))
			})

			It("saves the attributes associated to a transaction", func() {
				gasLimit := int64(5000)
				gasPrice := int64(3)
				nonce := uint64(10000)
				to := "1234567890"
				from := "0987654321"
				value := int64(10)
				transaction := core.Transaction{
					Hash:     "x1234",
					GasPrice: gasPrice,
					GasLimit: gasLimit,
					Nonce:    nonce,
					To:       to,
					From:     from,
					Value:    value,
				}
				block := core.Block{
					Number:       123,
					Transactions: []core.Transaction{transaction},
				}

				repository.CreateBlock(block)

				savedBlock := repository.FindBlockByNumber(123)
				Expect(len(savedBlock.Transactions)).To(Equal(1))
				savedTransaction := savedBlock.Transactions[0]
				Expect(savedTransaction.Hash).To(Equal(transaction.Hash))
				Expect(savedTransaction.To).To(Equal(to))
				Expect(savedTransaction.From).To(Equal(from))
				Expect(savedTransaction.Nonce).To(Equal(nonce))
				Expect(savedTransaction.GasLimit).To(Equal(gasLimit))
				Expect(savedTransaction.GasPrice).To(Equal(gasPrice))
				Expect(savedTransaction.Value).To(Equal(value))
			})

		})

		Describe("The missing block numbers", func() {
			It("is empty the starting block number is the highest known block number", func() {
				repository.CreateBlock(core.Block{Number: 1})

				Expect(len(repository.MissingBlockNumbers(1, 1))).To(Equal(0))
			})

			It("is the only missing block number", func() {
				repository.CreateBlock(core.Block{Number: 2})

				Expect(repository.MissingBlockNumbers(1, 2)).To(Equal([]int64{1}))
			})

			It("is both missing block numbers", func() {
				repository.CreateBlock(core.Block{Number: 3})

				Expect(repository.MissingBlockNumbers(1, 3)).To(Equal([]int64{1, 2}))
			})

			It("goes back to the starting block number", func() {
				repository.CreateBlock(core.Block{Number: 6})

				Expect(repository.MissingBlockNumbers(4, 6)).To(Equal([]int64{4, 5}))
			})

			It("only includes missing block numbers", func() {
				repository.CreateBlock(core.Block{Number: 4})
				repository.CreateBlock(core.Block{Number: 6})

				Expect(repository.MissingBlockNumbers(4, 6)).To(Equal([]int64{5}))
			})

			It("is a list with multiple gaps", func() {
				repository.CreateBlock(core.Block{Number: 4})
				repository.CreateBlock(core.Block{Number: 5})
				repository.CreateBlock(core.Block{Number: 8})
				repository.CreateBlock(core.Block{Number: 10})

				Expect(repository.MissingBlockNumbers(3, 10)).To(Equal([]int64{3, 6, 7, 9}))
			})

			It("returns empty array when lower bound exceeds upper bound", func() {
				Expect(repository.MissingBlockNumbers(10000, 1)).To(Equal([]int64{}))
			})

			It("only returns requested range even when other gaps exist", func() {
				repository.CreateBlock(core.Block{Number: 3})
				repository.CreateBlock(core.Block{Number: 8})

				Expect(repository.MissingBlockNumbers(1, 5)).To(Equal([]int64{1, 2, 4, 5}))
			})
		})

		Describe("The max block numbers", func() {
			It("returns the block number when a single block", func() {
				repository.CreateBlock(core.Block{Number: 1})

				Expect(repository.MaxBlockNumber()).To(Equal(int64(1)))
			})

			It("returns highest known block number when multiple blocks", func() {
				repository.CreateBlock(core.Block{Number: 1})
				repository.CreateBlock(core.Block{Number: 10})

				Expect(repository.MaxBlockNumber()).To(Equal(int64(10)))
			})
		})

		Describe("Creating watched contracts", func() {
			It("returns the watched contract when it exists", func() {
				repository.CreateWatchedContract(core.WatchedContract{Hash: "x123"})

				watchedContract := repository.FindWatchedContract("x123")
				Expect(watchedContract).NotTo(BeNil())
				Expect(watchedContract.Hash).To(Equal("x123"))

				Expect(repository.IsWatchedContract("x123")).To(BeTrue())
				Expect(repository.IsWatchedContract("x456")).To(BeFalse())
			})

			It("returns nil if contract does not exist", func() {
				watchedContract := repository.FindWatchedContract("x123")
				Expect(watchedContract).To(BeNil())
			})

			It("returns empty array when no transactions 'To' a watched contract", func() {
				repository.CreateWatchedContract(core.WatchedContract{Hash: "x123"})
				watchedContract := repository.FindWatchedContract("x123")
				Expect(watchedContract).ToNot(BeNil())
				Expect(watchedContract.Transactions).To(BeEmpty())

			})

			It("returns transactions 'To' a watched contract", func() {
				block := core.Block{
					Number: 123,
					Transactions: []core.Transaction{
						{Hash: "TRANSACTION1", To: "x123"},
						{Hash: "TRANSACTION2", To: "x345"},
						{Hash: "TRANSACTION3", To: "x123"},
					},
				}
				repository.CreateBlock(block)

				repository.CreateWatchedContract(core.WatchedContract{Hash: "x123"})
				watchedContract := repository.FindWatchedContract("x123")
				Expect(watchedContract).ToNot(BeNil())
				Expect(watchedContract.Transactions).To(
					Equal([]core.Transaction{
						{Hash: "TRANSACTION1", To: "x123"},
						{Hash: "TRANSACTION3", To: "x123"},
					}))
			})
		})
	}

	Describe("In memory repository", func() {
		AssertRepositoryBehavior(func() repositories.Repository {
			return repositories.NewInMemory()
		})
	})

	Describe("Postgres repository", func() {
		It("connects to the database", func() {
			cfg, _ := config.NewConfig("private")
			pgConfig := config.DbConnectionString(cfg.Database)
			db, err := sqlx.Connect("postgres", pgConfig)
			Expect(err).Should(BeNil())
			Expect(db).ShouldNot(BeNil())
		})

		It("does not commit block if block is invalid", func() {
			//badNonce violates db Nonce field length
			badNonce := fmt.Sprintf("x %s", strings.Repeat("1", 100))
			badBlock := core.Block{
				Number:       123,
				Nonce:        badNonce,
				Transactions: []core.Transaction{},
			}
			cfg, _ := config.NewConfig("private")
			repository := repositories.NewPostgres(cfg.Database)

			err := repository.CreateBlock(badBlock)
			savedBlock := repository.FindBlockByNumber(123)

			Expect(err).ToNot(BeNil())
			Expect(savedBlock).To(BeNil())
		})

		It("does not commit block or transactions if transaction is invalid", func() {
			//badHash violates db To field length
			badHash := fmt.Sprintf("x %s", strings.Repeat("1", 100))
			badTransaction := core.Transaction{To: badHash}
			block := core.Block{
				Number:       123,
				Transactions: []core.Transaction{badTransaction},
			}
			cfg, _ := config.NewConfig("private")
			repository := repositories.NewPostgres(cfg.Database)

			err := repository.CreateBlock(block)
			savedBlock := repository.FindBlockByNumber(123)

			Expect(err).ToNot(BeNil())
			Expect(savedBlock).To(BeNil())
		})

		AssertRepositoryBehavior(func() repositories.Repository {
			cfg, _ := config.NewConfig("private")
			repository := repositories.NewPostgres(cfg.Database)
			testing.ClearData(repository)
			return repository
		})
	})

})
