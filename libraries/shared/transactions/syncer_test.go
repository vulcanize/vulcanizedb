package transactions_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transactions"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Transaction syncer", func() {
	It("fetches transactions for logs", func() {
		db := test_config.NewTestDB(test_config.NewTestNode())
		blockChain := fakes.NewMockBlockChain()
		syncer := transactions.NewTransactionsSyncer(db, blockChain)

		err := syncer.SyncTransactions(0, []types.Log{})

		Expect(err).NotTo(HaveOccurred())
		Expect(blockChain.GetTransactionsCalled).To(BeTrue())
	})

	It("only fetches transactions with unique hashes", func() {
		db := test_config.NewTestDB(test_config.NewTestNode())
		blockChain := fakes.NewMockBlockChain()
		syncer := transactions.NewTransactionsSyncer(db, blockChain)

		err := syncer.SyncTransactions(0, []types.Log{{
			TxHash: fakes.FakeHash,
		}, {
			TxHash: fakes.FakeHash,
		}})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(blockChain.GetTransactionsPassedHashes)).To(Equal(1))
	})

	It("returns error if fetching transactions fails", func() {
		db := test_config.NewTestDB(test_config.NewTestNode())
		blockChain := fakes.NewMockBlockChain()
		blockChain.GetTransactionsError = fakes.FakeError
		syncer := transactions.NewTransactionsSyncer(db, blockChain)

		err := syncer.SyncTransactions(0, []types.Log{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("passes transactions to repository for persistence", func() {
		db := test_config.NewTestDB(test_config.NewTestNode())
		blockChain := fakes.NewMockBlockChain()
		blockChain.Transactions = []core.TransactionModel{{}}
		syncer := transactions.NewTransactionsSyncer(db, blockChain)
		mockHeaderRepository := fakes.NewMockHeaderRepository()
		syncer.Repository = mockHeaderRepository

		err := syncer.SyncTransactions(0, []types.Log{})

		Expect(err).NotTo(HaveOccurred())
		Expect(mockHeaderRepository.CreateTransactionCalled).To(BeTrue())
	})

	It("returns error if persisting transactions fails", func() {
		db := test_config.NewTestDB(test_config.NewTestNode())
		blockChain := fakes.NewMockBlockChain()
		blockChain.Transactions = []core.TransactionModel{{}}
		syncer := transactions.NewTransactionsSyncer(db, blockChain)
		mockHeaderRepository := fakes.NewMockHeaderRepository()
		mockHeaderRepository.CreateTransactionError = fakes.FakeError
		syncer.Repository = mockHeaderRepository

		err := syncer.SyncTransactions(0, []types.Log{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
