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

package ipfs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/mocks"
)

var (
	mockHeaderDagPutter  *mocks.DagPutter
	mockTrxDagPutter     *mocks.DagPutter
	mockRctDagPutter     *mocks.DagPutter
	mockStateDagPutter   *mocks.IncrementingDagPutter
	mockStorageDagPutter *mocks.IncrementingDagPutter
)

var _ = Describe("Publisher", func() {
	BeforeEach(func() {
		mockHeaderDagPutter = new(mocks.DagPutter)
		mockTrxDagPutter = new(mocks.DagPutter)
		mockRctDagPutter = new(mocks.DagPutter)
		mockStateDagPutter = new(mocks.IncrementingDagPutter)
		mockStorageDagPutter = new(mocks.IncrementingDagPutter)
	})

	Describe("Publish", func() {
		It("Publishes the passed IPLDPayload objects to IPFS and returns a CIDPayload for indexing", func() {
			mockHeaderDagPutter.CIDsToReturn = []string{"mockHeaderCID"}
			mockTrxDagPutter.CIDsToReturn = []string{"mockTrxCID1", "mockTrxCID2"}
			mockRctDagPutter.CIDsToReturn = []string{"mockRctCID1", "mockRctCID2"}
			mockStateDagPutter.CIDsToReturn = []string{"mockStateCID1", "mockStateCID2"}
			mockStorageDagPutter.CIDsToReturn = []string{"mockStorageCID"}
			publisher := ipfs.Publisher{
				HeaderPutter:      mockHeaderDagPutter,
				TransactionPutter: mockTrxDagPutter,
				ReceiptPutter:     mockRctDagPutter,
				StatePutter:       mockStateDagPutter,
				StoragePutter:     mockStorageDagPutter,
			}
			cidPayload, err := publisher.Publish(mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			Expect(cidPayload.BlockNumber).To(Equal(mocks.MockCIDPayload.BlockNumber))
			Expect(cidPayload.BlockHash).To(Equal(mocks.MockCIDPayload.BlockHash))
			Expect(cidPayload.UncleCIDs).To(Equal(mocks.MockCIDPayload.UncleCIDs))
			Expect(cidPayload.HeaderCID).To(Equal(mocks.MockCIDPayload.HeaderCID))
			Expect(len(cidPayload.TransactionCIDs)).To(Equal(2))
			Expect(cidPayload.TransactionCIDs[mocks.MockTransactions[0].Hash()]).To(Equal(mocks.MockCIDPayload.TransactionCIDs[mocks.MockTransactions[0].Hash()]))
			Expect(cidPayload.TransactionCIDs[mocks.MockTransactions[1].Hash()]).To(Equal(mocks.MockCIDPayload.TransactionCIDs[mocks.MockTransactions[1].Hash()]))
			Expect(len(cidPayload.ReceiptCIDs)).To(Equal(2))
			Expect(cidPayload.ReceiptCIDs[mocks.MockTransactions[0].Hash()]).To(Equal(mocks.MockCIDPayload.ReceiptCIDs[mocks.MockTransactions[0].Hash()]))
			Expect(cidPayload.ReceiptCIDs[mocks.MockTransactions[1].Hash()]).To(Equal(mocks.MockCIDPayload.ReceiptCIDs[mocks.MockTransactions[1].Hash()]))
			Expect(len(cidPayload.StateNodeCIDs)).To(Equal(2))
			Expect(cidPayload.StateNodeCIDs[mocks.ContractLeafKey]).To(Equal(mocks.MockCIDPayload.StateNodeCIDs[mocks.ContractLeafKey]))
			Expect(cidPayload.StateNodeCIDs[mocks.AnotherContractLeafKey]).To(Equal(mocks.MockCIDPayload.StateNodeCIDs[mocks.AnotherContractLeafKey]))
			Expect(cidPayload.StorageNodeCIDs).To(Equal(mocks.MockCIDPayload.StorageNodeCIDs))
		})
	})
})
