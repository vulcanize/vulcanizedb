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

package eth_test

import (
	"bytes"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

var (
	filterer                  *eth.ResponseFilterer
	expectedRctForStorageRLP1 []byte
	expectedRctForStorageRLP2 []byte
)

var _ = Describe("Filterer", func() {
	Describe("FilterResponse", func() {
		BeforeEach(func() {
			filterer = eth.NewResponseFilterer()
			expectedRctForStorageRLP1 = getReceiptForStorageRLP(mocks.MockReceipts, 0)
			expectedRctForStorageRLP2 = getReceiptForStorageRLP(mocks.MockReceipts, 1)
		})

		It("Transcribes all the data from the IPLDPayload into the StreamPayload if given an open filter", func() {
			payload, err := filterer.Filter(openFilter, mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			superNodePayload, ok := payload.(eth.StreamResponse)
			Expect(ok).To(BeTrue())
			Expect(superNodePayload.BlockNumber.Int64()).To(Equal(mocks.MockSeedNodePayload.BlockNumber.Int64()))
			Expect(superNodePayload.HeadersRlp).To(Equal(mocks.MockSeedNodePayload.HeadersRlp))
			var unclesRlp [][]byte
			Expect(superNodePayload.UnclesRlp).To(Equal(unclesRlp))
			Expect(len(superNodePayload.TransactionsRlp)).To(Equal(2))
			Expect(shared.ListContainsBytes(superNodePayload.TransactionsRlp, mocks.MockTransactions.GetRlp(0))).To(BeTrue())
			Expect(shared.ListContainsBytes(superNodePayload.TransactionsRlp, mocks.MockTransactions.GetRlp(1))).To(BeTrue())
			Expect(len(superNodePayload.ReceiptsRlp)).To(Equal(2))
			Expect(shared.ListContainsBytes(superNodePayload.ReceiptsRlp, expectedRctForStorageRLP1)).To(BeTrue())
			Expect(shared.ListContainsBytes(superNodePayload.ReceiptsRlp, expectedRctForStorageRLP2)).To(BeTrue())
			Expect(len(superNodePayload.StateNodesRlp)).To(Equal(2))
			Expect(superNodePayload.StateNodesRlp[mocks.ContractLeafKey]).To(Equal(mocks.ValueBytes))
			Expect(superNodePayload.StateNodesRlp[mocks.AnotherContractLeafKey]).To(Equal(mocks.AnotherValueBytes))
			Expect(superNodePayload.StorageNodesRlp).To(Equal(mocks.MockSeedNodePayload.StorageNodesRlp))
		})

		It("Applies filters from the provided config.Subscription", func() {
			payload1, err := filterer.Filter(rctContractFilter, mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			superNodePayload1, ok := payload1.(eth.StreamResponse)
			Expect(ok).To(BeTrue())
			Expect(superNodePayload1.BlockNumber.Int64()).To(Equal(mocks.MockSeedNodePayload.BlockNumber.Int64()))
			Expect(len(superNodePayload1.HeadersRlp)).To(Equal(0))
			Expect(len(superNodePayload1.UnclesRlp)).To(Equal(0))
			Expect(len(superNodePayload1.TransactionsRlp)).To(Equal(0))
			Expect(len(superNodePayload1.StorageNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload1.StateNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload1.ReceiptsRlp)).To(Equal(1))
			Expect(superNodePayload1.ReceiptsRlp[0]).To(Equal(expectedRctForStorageRLP2))

			payload2, err := filterer.Filter(rctTopicsFilter, mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			superNodePayload2, ok := payload2.(eth.StreamResponse)
			Expect(ok).To(BeTrue())
			Expect(superNodePayload2.BlockNumber.Int64()).To(Equal(mocks.MockSeedNodePayload.BlockNumber.Int64()))
			Expect(len(superNodePayload2.HeadersRlp)).To(Equal(0))
			Expect(len(superNodePayload2.UnclesRlp)).To(Equal(0))
			Expect(len(superNodePayload2.TransactionsRlp)).To(Equal(0))
			Expect(len(superNodePayload2.StorageNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload2.StateNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload2.ReceiptsRlp)).To(Equal(1))
			Expect(superNodePayload2.ReceiptsRlp[0]).To(Equal(expectedRctForStorageRLP1))

			payload3, err := filterer.Filter(rctTopicsAndContractFilter, mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			superNodePayload3, ok := payload3.(eth.StreamResponse)
			Expect(ok).To(BeTrue())
			Expect(superNodePayload3.BlockNumber.Int64()).To(Equal(mocks.MockSeedNodePayload.BlockNumber.Int64()))
			Expect(len(superNodePayload3.HeadersRlp)).To(Equal(0))
			Expect(len(superNodePayload3.UnclesRlp)).To(Equal(0))
			Expect(len(superNodePayload3.TransactionsRlp)).To(Equal(0))
			Expect(len(superNodePayload3.StorageNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload3.StateNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload3.ReceiptsRlp)).To(Equal(1))
			Expect(superNodePayload3.ReceiptsRlp[0]).To(Equal(expectedRctForStorageRLP1))

			payload4, err := filterer.Filter(rctContractsAndTopicFilter, mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			superNodePayload4, ok := payload4.(eth.StreamResponse)
			Expect(ok).To(BeTrue())
			Expect(superNodePayload4.BlockNumber.Int64()).To(Equal(mocks.MockSeedNodePayload.BlockNumber.Int64()))
			Expect(len(superNodePayload4.HeadersRlp)).To(Equal(0))
			Expect(len(superNodePayload4.UnclesRlp)).To(Equal(0))
			Expect(len(superNodePayload4.TransactionsRlp)).To(Equal(0))
			Expect(len(superNodePayload4.StorageNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload4.StateNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload4.ReceiptsRlp)).To(Equal(1))
			Expect(superNodePayload4.ReceiptsRlp[0]).To(Equal(expectedRctForStorageRLP2))

			payload5, err := filterer.Filter(rctsForAllCollectedTrxs, mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			superNodePayload5, ok := payload5.(eth.StreamResponse)
			Expect(ok).To(BeTrue())
			Expect(superNodePayload5.BlockNumber.Int64()).To(Equal(mocks.MockSeedNodePayload.BlockNumber.Int64()))
			Expect(len(superNodePayload5.HeadersRlp)).To(Equal(0))
			Expect(len(superNodePayload5.UnclesRlp)).To(Equal(0))
			Expect(len(superNodePayload5.TransactionsRlp)).To(Equal(2))
			Expect(shared.ListContainsBytes(superNodePayload5.TransactionsRlp, mocks.MockTransactions.GetRlp(0))).To(BeTrue())
			Expect(shared.ListContainsBytes(superNodePayload5.TransactionsRlp, mocks.MockTransactions.GetRlp(1))).To(BeTrue())
			Expect(len(superNodePayload5.StorageNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload5.StateNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload5.ReceiptsRlp)).To(Equal(2))
			Expect(shared.ListContainsBytes(superNodePayload5.ReceiptsRlp, expectedRctForStorageRLP1)).To(BeTrue())
			Expect(shared.ListContainsBytes(superNodePayload5.ReceiptsRlp, expectedRctForStorageRLP2)).To(BeTrue())

			payload6, err := filterer.Filter(rctsForSelectCollectedTrxs, mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			superNodePayload6, ok := payload6.(eth.StreamResponse)
			Expect(ok).To(BeTrue())
			Expect(superNodePayload6.BlockNumber.Int64()).To(Equal(mocks.MockSeedNodePayload.BlockNumber.Int64()))
			Expect(len(superNodePayload6.HeadersRlp)).To(Equal(0))
			Expect(len(superNodePayload6.UnclesRlp)).To(Equal(0))
			Expect(len(superNodePayload6.TransactionsRlp)).To(Equal(1))
			Expect(shared.ListContainsBytes(superNodePayload5.TransactionsRlp, mocks.MockTransactions.GetRlp(1))).To(BeTrue())
			Expect(len(superNodePayload6.StorageNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload6.StateNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload6.ReceiptsRlp)).To(Equal(1))
			Expect(superNodePayload4.ReceiptsRlp[0]).To(Equal(expectedRctForStorageRLP2))

			payload7, err := filterer.Filter(stateFilter, mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			superNodePayload7, ok := payload7.(eth.StreamResponse)
			Expect(ok).To(BeTrue())
			Expect(superNodePayload7.BlockNumber.Int64()).To(Equal(mocks.MockSeedNodePayload.BlockNumber.Int64()))
			Expect(len(superNodePayload7.HeadersRlp)).To(Equal(0))
			Expect(len(superNodePayload7.UnclesRlp)).To(Equal(0))
			Expect(len(superNodePayload7.TransactionsRlp)).To(Equal(0))
			Expect(len(superNodePayload7.StorageNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload7.ReceiptsRlp)).To(Equal(0))
			Expect(len(superNodePayload7.StateNodesRlp)).To(Equal(1))
			Expect(superNodePayload7.StateNodesRlp[mocks.ContractLeafKey]).To(Equal(mocks.ValueBytes))

			payload8, err := filterer.Filter(rctTopicsAndContractFilterFail, mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			superNodePayload8, ok := payload8.(eth.StreamResponse)
			Expect(ok).To(BeTrue())
			Expect(superNodePayload8.BlockNumber.Int64()).To(Equal(mocks.MockSeedNodePayload.BlockNumber.Int64()))
			Expect(len(superNodePayload8.HeadersRlp)).To(Equal(0))
			Expect(len(superNodePayload8.UnclesRlp)).To(Equal(0))
			Expect(len(superNodePayload8.TransactionsRlp)).To(Equal(0))
			Expect(len(superNodePayload8.StorageNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload8.StateNodesRlp)).To(Equal(0))
			Expect(len(superNodePayload8.ReceiptsRlp)).To(Equal(0))
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
