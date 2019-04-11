package eth_block_receipts_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_receipts"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/ipfs"
)

var _ = Describe("Eth block receipts dag putter", func() {
	It("adds a node for each receipt", func() {
		mockAdder := ipfs.NewMockAdder()
		dagPutter := eth_block_receipts.NewEthBlockReceiptDagPutter(mockAdder)
		fakeReceipts := types.Receipts{
			&types.Receipt{},
			&types.Receipt{},
		}

		_, err := dagPutter.DagPut(fakeReceipts)

		Expect(err).NotTo(HaveOccurred())
		mockAdder.AssertAddCalled(2, &eth_block_receipts.EthReceiptNode{})
	})

	It("returns error if adding node fails", func() {
		mockAdder := ipfs.NewMockAdder()
		mockAdder.SetError(test_helpers.FakeError)
		dagPutter := eth_block_receipts.NewEthBlockReceiptDagPutter(mockAdder)
		fakeReceipts := types.Receipts{
			&types.Receipt{},
			&types.Receipt{},
		}

		_, err := dagPutter.DagPut(fakeReceipts)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})
})
