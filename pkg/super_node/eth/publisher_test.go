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
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth/mocks"
)

var (
	mockHeaderDagPutter  *mocks.DagPutter
	mockTrxDagPutter     *mocks.DagPutter
	mockRctDagPutter     *mocks.DagPutter
	mockStateDagPutter   *mocks.MappedDagPutter
	mockStorageDagPutter *mocks.DagPutter
)

var _ = Describe("Publisher", func() {
	BeforeEach(func() {
		mockHeaderDagPutter = new(mocks.DagPutter)
		mockTrxDagPutter = new(mocks.DagPutter)
		mockRctDagPutter = new(mocks.DagPutter)
		mockStateDagPutter = new(mocks.MappedDagPutter)
		mockStorageDagPutter = new(mocks.DagPutter)
	})

	Describe("Publish", func() {
		It("Publishes the passed IPLDPayload objects to IPFS and returns a CIDPayload for indexing", func() {
			mockHeaderDagPutter.CIDsToReturn = []string{"mockHeaderCID"}
			mockTrxDagPutter.CIDsToReturn = []string{"mockTrxCID1", "mockTrxCID2"}
			mockRctDagPutter.CIDsToReturn = []string{"mockRctCID1", "mockRctCID2"}
			val1 := common.BytesToHash(mocks.MockIPLDPayload.StateNodes[0].Value)
			val2 := common.BytesToHash(mocks.MockIPLDPayload.StateNodes[1].Value)
			mockStateDagPutter.CIDsToReturn = map[common.Hash][]string{
				val1: {"mockStateCID1"},
				val2: {"mockStateCID2"},
			}
			mockStorageDagPutter.CIDsToReturn = []string{"mockStorageCID"}
			publisher := eth.IPLDPublisher{
				HeaderPutter:      mockHeaderDagPutter,
				TransactionPutter: mockTrxDagPutter,
				ReceiptPutter:     mockRctDagPutter,
				StatePutter:       mockStateDagPutter,
				StoragePutter:     mockStorageDagPutter,
			}
			payload, err := publisher.Publish(mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			cidPayload, ok := payload.(*eth.CIDPayload)
			Expect(ok).To(BeTrue())
			Expect(cidPayload.HeaderCID.TotalDifficulty).To(Equal(mocks.MockIPLDPayload.TotalDifficulty.String()))
			Expect(cidPayload.HeaderCID.BlockNumber).To(Equal(mocks.MockCIDPayload.HeaderCID.BlockNumber))
			Expect(cidPayload.HeaderCID.BlockHash).To(Equal(mocks.MockCIDPayload.HeaderCID.BlockHash))
			Expect(cidPayload.UncleCIDs).To(Equal(mocks.MockCIDPayload.UncleCIDs))
			Expect(cidPayload.HeaderCID).To(Equal(mocks.MockCIDPayload.HeaderCID))
			Expect(len(cidPayload.TransactionCIDs)).To(Equal(2))
			Expect(cidPayload.TransactionCIDs[0]).To(Equal(mocks.MockCIDPayload.TransactionCIDs[0]))
			Expect(cidPayload.TransactionCIDs[1]).To(Equal(mocks.MockCIDPayload.TransactionCIDs[1]))
			Expect(len(cidPayload.ReceiptCIDs)).To(Equal(2))
			Expect(cidPayload.ReceiptCIDs[mocks.MockTransactions[0].Hash()]).To(Equal(mocks.MockCIDPayload.ReceiptCIDs[mocks.MockTransactions[0].Hash()]))
			Expect(cidPayload.ReceiptCIDs[mocks.MockTransactions[1].Hash()]).To(Equal(mocks.MockCIDPayload.ReceiptCIDs[mocks.MockTransactions[1].Hash()]))
			Expect(len(cidPayload.StateNodeCIDs)).To(Equal(2))
			Expect(cidPayload.StateNodeCIDs[0]).To(Equal(mocks.MockCIDPayload.StateNodeCIDs[0]))
			Expect(cidPayload.StateNodeCIDs[1]).To(Equal(mocks.MockCIDPayload.StateNodeCIDs[1]))
			Expect(cidPayload.StorageNodeCIDs).To(Equal(mocks.MockCIDPayload.StorageNodeCIDs))
		})
	})
})
