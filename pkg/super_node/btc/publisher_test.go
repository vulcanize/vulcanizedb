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

package btc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mocks2 "github.com/vulcanize/vulcanizedb/pkg/ipfs/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc/mocks"
)

var (
	mockHeaderDagPutter *mocks2.DagPutter
	mockTrxDagPutter    *mocks2.DagPutter
)

var _ = Describe("Publisher", func() {
	BeforeEach(func() {
		mockHeaderDagPutter = new(mocks2.DagPutter)
		mockTrxDagPutter = new(mocks2.DagPutter)
	})

	Describe("Publish", func() {
		It("Publishes the passed IPLDPayload objects to IPFS and returns a CIDPayload for indexing", func() {
			mockHeaderDagPutter.CIDsToReturn = []string{"mockHeaderCID"}
			mockTrxDagPutter.CIDsToReturn = []string{"mockTrxCID1", "mockTrxCID2", "mockTrxCID3"}
			publisher := btc.IPLDPublisher{
				HeaderPutter:      mockHeaderDagPutter,
				TransactionPutter: mockTrxDagPutter,
			}
			payload, err := publisher.Publish(mocks.MockIPLDPayload)
			Expect(err).ToNot(HaveOccurred())
			cidPayload, ok := payload.(*btc.CIDPayload)
			Expect(ok).To(BeTrue())
			Expect(cidPayload).To(Equal(&mocks.MockCIDPayload))
			Expect(cidPayload.HeaderCID).To(Equal(mocks.MockHeaderMetaData))
			Expect(cidPayload.TransactionCIDs).To(Equal(mocks.MockTxsMetaDataPostPublish))
		})
	})
})
