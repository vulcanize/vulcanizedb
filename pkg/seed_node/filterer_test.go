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

package seed_node_test

import (
	"bytes"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/seed_node"
)

var (
	filterer                  seed_node.ResponseFilterer
	expectedRctForStorageRLP1 []byte
	expectedRctForStorageRLP2 []byte
)

var _ = Describe("Filterer", func() {
	Describe("FilterResponse", func() {
		BeforeEach(func() {
			filterer = seed_node.NewResponseFilterer()
			expectedRctForStorageRLP1 = getReceiptForStorageRLP(mocks.MockReceipts, 0)
			expectedRctForStorageRLP2 = getReceiptForStorageRLP(mocks.MockReceipts, 1)
		})

		It("Transcribes all the data from the IPLDPayload into the SeedNodePayload if given an open filter", func() {
			seedNodePayload, err := filterer.FilterResponse(openFilter, *mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			Expect(seedNodePayload.BlockNumber.Int64()).To(Equal(mocks.MockSeeNodePayload.BlockNumber.Int64()))
			Expect(seedNodePayload.HeadersRlp).To(Equal(mocks.MockSeeNodePayload.HeadersRlp))
			Expect(seedNodePayload.UnclesRlp).To(Equal(mocks.MockSeeNodePayload.UnclesRlp))
			Expect(len(seedNodePayload.TransactionsRlp)).To(Equal(2))
			Expect(seed_node.ListContainsBytes(seedNodePayload.TransactionsRlp, mocks.MockTransactions.GetRlp(0))).To(BeTrue())
			Expect(seed_node.ListContainsBytes(seedNodePayload.TransactionsRlp, mocks.MockTransactions.GetRlp(1))).To(BeTrue())
			Expect(len(seedNodePayload.ReceiptsRlp)).To(Equal(2))
			Expect(seed_node.ListContainsBytes(seedNodePayload.ReceiptsRlp, expectedRctForStorageRLP1)).To(BeTrue())
			Expect(seed_node.ListContainsBytes(seedNodePayload.ReceiptsRlp, expectedRctForStorageRLP2)).To(BeTrue())
			Expect(len(seedNodePayload.StateNodesRlp)).To(Equal(2))
			Expect(seedNodePayload.StateNodesRlp[mocks.ContractLeafKey]).To(Equal(mocks.ValueBytes))
			Expect(seedNodePayload.StateNodesRlp[mocks.AnotherContractLeafKey]).To(Equal(mocks.AnotherValueBytes))
			Expect(seedNodePayload.StorageNodesRlp).To(Equal(mocks.MockSeeNodePayload.StorageNodesRlp))
		})

		It("Applies filters from the provided config.Subscription", func() {
			seedNodePayload1, err := filterer.FilterResponse(rctContractFilter, *mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			Expect(seedNodePayload1.BlockNumber.Int64()).To(Equal(mocks.MockSeeNodePayload.BlockNumber.Int64()))
			Expect(len(seedNodePayload1.HeadersRlp)).To(Equal(0))
			Expect(len(seedNodePayload1.UnclesRlp)).To(Equal(0))
			Expect(len(seedNodePayload1.TransactionsRlp)).To(Equal(0))
			Expect(len(seedNodePayload1.StorageNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload1.StateNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload1.ReceiptsRlp)).To(Equal(1))
			Expect(seedNodePayload1.ReceiptsRlp[0]).To(Equal(expectedRctForStorageRLP2))

			seedNodePayload2, err := filterer.FilterResponse(rctTopicsFilter, *mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			Expect(seedNodePayload2.BlockNumber.Int64()).To(Equal(mocks.MockSeeNodePayload.BlockNumber.Int64()))
			Expect(len(seedNodePayload2.HeadersRlp)).To(Equal(0))
			Expect(len(seedNodePayload2.UnclesRlp)).To(Equal(0))
			Expect(len(seedNodePayload2.TransactionsRlp)).To(Equal(0))
			Expect(len(seedNodePayload2.StorageNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload2.StateNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload2.ReceiptsRlp)).To(Equal(1))
			Expect(seedNodePayload2.ReceiptsRlp[0]).To(Equal(expectedRctForStorageRLP1))

			seedNodePayload3, err := filterer.FilterResponse(rctTopicsAndContractFilter, *mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			Expect(seedNodePayload3.BlockNumber.Int64()).To(Equal(mocks.MockSeeNodePayload.BlockNumber.Int64()))
			Expect(len(seedNodePayload3.HeadersRlp)).To(Equal(0))
			Expect(len(seedNodePayload3.UnclesRlp)).To(Equal(0))
			Expect(len(seedNodePayload3.TransactionsRlp)).To(Equal(0))
			Expect(len(seedNodePayload3.StorageNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload3.StateNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload3.ReceiptsRlp)).To(Equal(1))
			Expect(seedNodePayload3.ReceiptsRlp[0]).To(Equal(expectedRctForStorageRLP1))

			seedNodePayload4, err := filterer.FilterResponse(rctContractsAndTopicFilter, *mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			Expect(seedNodePayload4.BlockNumber.Int64()).To(Equal(mocks.MockSeeNodePayload.BlockNumber.Int64()))
			Expect(len(seedNodePayload4.HeadersRlp)).To(Equal(0))
			Expect(len(seedNodePayload4.UnclesRlp)).To(Equal(0))
			Expect(len(seedNodePayload4.TransactionsRlp)).To(Equal(0))
			Expect(len(seedNodePayload4.StorageNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload4.StateNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload4.ReceiptsRlp)).To(Equal(1))
			Expect(seedNodePayload4.ReceiptsRlp[0]).To(Equal(expectedRctForStorageRLP2))

			seedNodePayload5, err := filterer.FilterResponse(rctsForAllCollectedTrxs, *mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			Expect(seedNodePayload5.BlockNumber.Int64()).To(Equal(mocks.MockSeeNodePayload.BlockNumber.Int64()))
			Expect(len(seedNodePayload5.HeadersRlp)).To(Equal(0))
			Expect(len(seedNodePayload5.UnclesRlp)).To(Equal(0))
			Expect(len(seedNodePayload5.TransactionsRlp)).To(Equal(2))
			Expect(seed_node.ListContainsBytes(seedNodePayload5.TransactionsRlp, mocks.MockTransactions.GetRlp(0))).To(BeTrue())
			Expect(seed_node.ListContainsBytes(seedNodePayload5.TransactionsRlp, mocks.MockTransactions.GetRlp(1))).To(BeTrue())
			Expect(len(seedNodePayload5.StorageNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload5.StateNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload5.ReceiptsRlp)).To(Equal(2))
			Expect(seed_node.ListContainsBytes(seedNodePayload5.ReceiptsRlp, expectedRctForStorageRLP1)).To(BeTrue())
			Expect(seed_node.ListContainsBytes(seedNodePayload5.ReceiptsRlp, expectedRctForStorageRLP2)).To(BeTrue())

			seedNodePayload6, err := filterer.FilterResponse(rctsForSelectCollectedTrxs, *mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			Expect(seedNodePayload6.BlockNumber.Int64()).To(Equal(mocks.MockSeeNodePayload.BlockNumber.Int64()))
			Expect(len(seedNodePayload6.HeadersRlp)).To(Equal(0))
			Expect(len(seedNodePayload6.UnclesRlp)).To(Equal(0))
			Expect(len(seedNodePayload6.TransactionsRlp)).To(Equal(1))
			Expect(seed_node.ListContainsBytes(seedNodePayload5.TransactionsRlp, mocks.MockTransactions.GetRlp(1))).To(BeTrue())
			Expect(len(seedNodePayload6.StorageNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload6.StateNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload6.ReceiptsRlp)).To(Equal(1))
			Expect(seedNodePayload4.ReceiptsRlp[0]).To(Equal(expectedRctForStorageRLP2))

			seedNodePayload7, err := filterer.FilterResponse(stateFilter, *mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			Expect(seedNodePayload7.BlockNumber.Int64()).To(Equal(mocks.MockSeeNodePayload.BlockNumber.Int64()))
			Expect(len(seedNodePayload7.HeadersRlp)).To(Equal(0))
			Expect(len(seedNodePayload7.UnclesRlp)).To(Equal(0))
			Expect(len(seedNodePayload7.TransactionsRlp)).To(Equal(0))
			Expect(len(seedNodePayload7.StorageNodesRlp)).To(Equal(0))
			Expect(len(seedNodePayload7.ReceiptsRlp)).To(Equal(0))
			Expect(len(seedNodePayload7.StateNodesRlp)).To(Equal(1))
			Expect(seedNodePayload7.StateNodesRlp[mocks.ContractLeafKey]).To(Equal(mocks.ValueBytes))
		})
	})
})

func getReceiptForStorageRLP(receipts types.Receipts, i int) []byte {
	receiptForStorage := (*types.ReceiptForStorage)(receipts[i])
	receiptBuffer := new(bytes.Buffer)
	err := receiptForStorage.EncodeRLP(receiptBuffer)
	Expect(err).ToNot(HaveOccurred())
	return receiptBuffer.Bytes()
}
