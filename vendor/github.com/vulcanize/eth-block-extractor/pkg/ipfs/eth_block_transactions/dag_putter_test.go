package eth_block_transactions_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_transactions"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/ipfs"
)

var _ = Describe("Eth block transactions dag putter", func() {
	It("adds a node for each transaction on the block", func() {
		mockAdder := ipfs.NewMockAdder()
		fakeTransactionOne := &types.Transaction{}
		fakeTransactionTwo := &types.Transaction{}
		fakeBlockBody := &types.Body{
			Transactions: types.Transactions{fakeTransactionOne, fakeTransactionTwo},
			Uncles:       nil,
		}
		dagPutter := eth_block_transactions.NewBlockTransactionsDagPutter(mockAdder)

		_, err := dagPutter.DagPut(fakeBlockBody)

		Expect(err).NotTo(HaveOccurred())
		mockAdder.AssertAddCalled(2, &eth_block_transactions.EthTransactionNode{})
	})

	It("returns error if adding node fails", func() {
		mockAdder := ipfs.NewMockAdder()
		mockAdder.SetError(test_helpers.FakeError)
		dagPutter := eth_block_transactions.NewBlockTransactionsDagPutter(mockAdder)

		_, err := dagPutter.DagPut(&types.Body{Transactions: types.Transactions{{}}})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})
})
