// VulcanizeDB
// Copyright © 2019 Vulcanize

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
	var (
		blockChain *fakes.MockBlockChain
		syncer     transactions.TransactionsSyncer
	)

	BeforeEach(func() {
		db := test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		blockChain = fakes.NewMockBlockChain()
		syncer = transactions.NewTransactionsSyncer(db, blockChain)
	})

	It("fetches transactions for logs", func() {
		err := syncer.SyncTransactions(0, []types.Log{{TxHash: fakes.FakeHash}})

		Expect(err).NotTo(HaveOccurred())
		Expect(blockChain.GetTransactionsCalled).To(BeTrue())
	})

	It("does not fetch transactions if no logs", func() {
		err := syncer.SyncTransactions(0, []types.Log{})

		Expect(err).NotTo(HaveOccurred())
		Expect(blockChain.GetTransactionsCalled).To(BeFalse())
	})

	It("only fetches transactions with unique hashes", func() {
		err := syncer.SyncTransactions(0, []types.Log{{
			TxHash: fakes.FakeHash,
		}, {
			TxHash: fakes.FakeHash,
		}})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(blockChain.GetTransactionsPassedHashes)).To(Equal(1))
	})

	It("returns error if fetching transactions fails", func() {
		blockChain.GetTransactionsError = fakes.FakeError

		err := syncer.SyncTransactions(0, []types.Log{{TxHash: fakes.FakeHash}})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("passes transactions to repository for persistence", func() {
		blockChain.Transactions = []core.TransactionModel{{}}
		mockHeaderRepository := fakes.NewMockHeaderRepository()
		syncer.Repository = mockHeaderRepository

		err := syncer.SyncTransactions(0, []types.Log{{TxHash: fakes.FakeHash}})

		Expect(err).NotTo(HaveOccurred())
		Expect(mockHeaderRepository.CreateTransactionsCalled).To(BeTrue())
	})

	It("returns error if persisting transactions fails", func() {
		blockChain.Transactions = []core.TransactionModel{{}}
		mockHeaderRepository := fakes.NewMockHeaderRepository()
		mockHeaderRepository.CreateTransactionsError = fakes.FakeError
		syncer.Repository = mockHeaderRepository

		err := syncer.SyncTransactions(0, []types.Log{{TxHash: fakes.FakeHash}})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
