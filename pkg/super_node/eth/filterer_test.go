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
			payload, err := filterer.Filter(openFilter, mocks.MockConvertedPayload)
			Expect(err).ToNot(HaveOccurred())
			iplds, ok := payload.(eth.IPLDs)
			Expect(ok).To(BeTrue())
			Expect(iplds.BlockNumber.Int64()).To(Equal(mocks.MockIPLDs.BlockNumber.Int64()))
			Expect(iplds.Headers).To(Equal(mocks.MockIPLDs.Headers))
			var unclesRlp [][]byte
			Expect(iplds.Uncles).To(Equal(unclesRlp))
			Expect(len(iplds.Transactions)).To(Equal(2))
			Expect(shared.ListContainsBytes(iplds.Transactions, mocks.MockTransactions.GetRlp(0))).To(BeTrue())
			Expect(shared.ListContainsBytes(iplds.Transactions, mocks.MockTransactions.GetRlp(1))).To(BeTrue())
			Expect(len(iplds.Receipts)).To(Equal(2))
			Expect(shared.ListContainsBytes(iplds.Receipts, expectedRctForStorageRLP1)).To(BeTrue())
			Expect(shared.ListContainsBytes(iplds.Receipts, expectedRctForStorageRLP2)).To(BeTrue())
			Expect(len(iplds.StateNodes)).To(Equal(2))
			for _, stateNode := range iplds.StateNodes {
				Expect(stateNode.Leaf).To(BeTrue())
				if stateNode.StateTrieKey == mocks.ContractLeafKey {
					Expect(stateNode.IPLD).To(Equal(mocks.ValueBytes))
				}
				if stateNode.StateTrieKey == mocks.AnotherContractLeafKey {
					Expect(stateNode.IPLD).To(Equal(mocks.AnotherValueBytes))
				}
			}
			Expect(iplds.StorageNodes).To(Equal(mocks.MockIPLDs.StorageNodes))
		})

		It("Applies filters from the provided config.Subscription", func() {
			payload1, err := filterer.Filter(rctContractFilter, mocks.MockConvertedPayload)
			Expect(err).ToNot(HaveOccurred())
			iplds1, ok := payload1.(eth.IPLDs)
			Expect(ok).To(BeTrue())
			Expect(iplds1.BlockNumber.Int64()).To(Equal(mocks.MockIPLDs.BlockNumber.Int64()))
			Expect(len(iplds1.Headers)).To(Equal(0))
			Expect(len(iplds1.Uncles)).To(Equal(0))
			Expect(len(iplds1.Transactions)).To(Equal(0))
			Expect(len(iplds1.StorageNodes)).To(Equal(0))
			Expect(len(iplds1.StateNodes)).To(Equal(0))
			Expect(len(iplds1.Receipts)).To(Equal(1))
			Expect(iplds1.Receipts[0]).To(Equal(expectedRctForStorageRLP2))

			payload2, err := filterer.Filter(rctTopicsFilter, mocks.MockConvertedPayload)
			Expect(err).ToNot(HaveOccurred())
			iplds2, ok := payload2.(eth.IPLDs)
			Expect(ok).To(BeTrue())
			Expect(iplds2.BlockNumber.Int64()).To(Equal(mocks.MockIPLDs.BlockNumber.Int64()))
			Expect(len(iplds2.Headers)).To(Equal(0))
			Expect(len(iplds2.Uncles)).To(Equal(0))
			Expect(len(iplds2.Transactions)).To(Equal(0))
			Expect(len(iplds2.StorageNodes)).To(Equal(0))
			Expect(len(iplds2.StateNodes)).To(Equal(0))
			Expect(len(iplds2.Receipts)).To(Equal(1))
			Expect(iplds2.Receipts[0]).To(Equal(expectedRctForStorageRLP1))

			payload3, err := filterer.Filter(rctTopicsAndContractFilter, mocks.MockConvertedPayload)
			Expect(err).ToNot(HaveOccurred())
			iplds3, ok := payload3.(eth.IPLDs)
			Expect(ok).To(BeTrue())
			Expect(iplds3.BlockNumber.Int64()).To(Equal(mocks.MockIPLDs.BlockNumber.Int64()))
			Expect(len(iplds3.Headers)).To(Equal(0))
			Expect(len(iplds3.Uncles)).To(Equal(0))
			Expect(len(iplds3.Transactions)).To(Equal(0))
			Expect(len(iplds3.StorageNodes)).To(Equal(0))
			Expect(len(iplds3.StateNodes)).To(Equal(0))
			Expect(len(iplds3.Receipts)).To(Equal(1))
			Expect(iplds3.Receipts[0]).To(Equal(expectedRctForStorageRLP1))

			payload4, err := filterer.Filter(rctContractsAndTopicFilter, mocks.MockConvertedPayload)
			Expect(err).ToNot(HaveOccurred())
			iplds4, ok := payload4.(eth.IPLDs)
			Expect(ok).To(BeTrue())
			Expect(iplds4.BlockNumber.Int64()).To(Equal(mocks.MockIPLDs.BlockNumber.Int64()))
			Expect(len(iplds4.Headers)).To(Equal(0))
			Expect(len(iplds4.Uncles)).To(Equal(0))
			Expect(len(iplds4.Transactions)).To(Equal(0))
			Expect(len(iplds4.StorageNodes)).To(Equal(0))
			Expect(len(iplds4.StateNodes)).To(Equal(0))
			Expect(len(iplds4.Receipts)).To(Equal(1))
			Expect(iplds4.Receipts[0]).To(Equal(expectedRctForStorageRLP2))

			payload5, err := filterer.Filter(rctsForAllCollectedTrxs, mocks.MockConvertedPayload)
			Expect(err).ToNot(HaveOccurred())
			iplds5, ok := payload5.(eth.IPLDs)
			Expect(ok).To(BeTrue())
			Expect(iplds5.BlockNumber.Int64()).To(Equal(mocks.MockIPLDs.BlockNumber.Int64()))
			Expect(len(iplds5.Headers)).To(Equal(0))
			Expect(len(iplds5.Uncles)).To(Equal(0))
			Expect(len(iplds5.Transactions)).To(Equal(2))
			Expect(shared.ListContainsBytes(iplds5.Transactions, mocks.MockTransactions.GetRlp(0))).To(BeTrue())
			Expect(shared.ListContainsBytes(iplds5.Transactions, mocks.MockTransactions.GetRlp(1))).To(BeTrue())
			Expect(len(iplds5.StorageNodes)).To(Equal(0))
			Expect(len(iplds5.StateNodes)).To(Equal(0))
			Expect(len(iplds5.Receipts)).To(Equal(2))
			Expect(shared.ListContainsBytes(iplds5.Receipts, expectedRctForStorageRLP1)).To(BeTrue())
			Expect(shared.ListContainsBytes(iplds5.Receipts, expectedRctForStorageRLP2)).To(BeTrue())

			payload6, err := filterer.Filter(rctsForSelectCollectedTrxs, mocks.MockConvertedPayload)
			Expect(err).ToNot(HaveOccurred())
			iplds6, ok := payload6.(eth.IPLDs)
			Expect(ok).To(BeTrue())
			Expect(iplds6.BlockNumber.Int64()).To(Equal(mocks.MockIPLDs.BlockNumber.Int64()))
			Expect(len(iplds6.Headers)).To(Equal(0))
			Expect(len(iplds6.Uncles)).To(Equal(0))
			Expect(len(iplds6.Transactions)).To(Equal(1))
			Expect(shared.ListContainsBytes(iplds5.Transactions, mocks.MockTransactions.GetRlp(1))).To(BeTrue())
			Expect(len(iplds6.StorageNodes)).To(Equal(0))
			Expect(len(iplds6.StateNodes)).To(Equal(0))
			Expect(len(iplds6.Receipts)).To(Equal(1))
			Expect(iplds4.Receipts[0]).To(Equal(expectedRctForStorageRLP2))

			payload7, err := filterer.Filter(stateFilter, mocks.MockConvertedPayload)
			Expect(err).ToNot(HaveOccurred())
			iplds7, ok := payload7.(eth.IPLDs)
			Expect(ok).To(BeTrue())
			Expect(iplds7.BlockNumber.Int64()).To(Equal(mocks.MockIPLDs.BlockNumber.Int64()))
			Expect(len(iplds7.Headers)).To(Equal(0))
			Expect(len(iplds7.Uncles)).To(Equal(0))
			Expect(len(iplds7.Transactions)).To(Equal(0))
			Expect(len(iplds7.StorageNodes)).To(Equal(0))
			Expect(len(iplds7.Receipts)).To(Equal(0))
			Expect(len(iplds7.StateNodes)).To(Equal(1))
			Expect(iplds7.StateNodes[0].StateTrieKey).To(Equal(mocks.ContractLeafKey))
			Expect(iplds7.StateNodes[0].IPLD).To(Equal(mocks.ValueBytes))

			payload8, err := filterer.Filter(rctTopicsAndContractFilterFail, mocks.MockConvertedPayload)
			Expect(err).ToNot(HaveOccurred())
			iplds8, ok := payload8.(eth.IPLDs)
			Expect(ok).To(BeTrue())
			Expect(iplds8.BlockNumber.Int64()).To(Equal(mocks.MockIPLDs.BlockNumber.Int64()))
			Expect(len(iplds8.Headers)).To(Equal(0))
			Expect(len(iplds8.Uncles)).To(Equal(0))
			Expect(len(iplds8.Transactions)).To(Equal(0))
			Expect(len(iplds8.StorageNodes)).To(Equal(0))
			Expect(len(iplds8.StateNodes)).To(Equal(0))
			Expect(len(iplds8.Receipts)).To(Equal(0))
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
