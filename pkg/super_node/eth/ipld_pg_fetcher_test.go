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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

var (
	db            *postgres.DB
	pubAndIndexer *eth.IPLDPublisherAndIndexer
	fetcher       *eth.IPLDPGFetcher
)

var _ = Describe("IPLDPGFetcher", func() {
	Describe("Fetch", func() {
		BeforeEach(func() {
			var err error
			db, err = shared.SetupDB()
			Expect(err).ToNot(HaveOccurred())
			pubAndIndexer = eth.NewIPLDPublisherAndIndexer(db)
			_, err = pubAndIndexer.Publish(mocks.MockConvertedPayload)
			Expect(err).ToNot(HaveOccurred())
			fetcher = eth.NewIPLDPGFetcher(db)
		})
		AfterEach(func() {
			eth.TearDownDB(db)
		})

		It("Fetches and returns IPLDs for the CIDs provided in the CIDWrapper", func() {
			i, err := fetcher.Fetch(mocks.MockCIDWrapper)
			Expect(err).ToNot(HaveOccurred())
			iplds, ok := i.(eth.IPLDs)
			Expect(ok).To(BeTrue())
			Expect(iplds.TotalDifficulty).To(Equal(mocks.MockConvertedPayload.TotalDifficulty))
			Expect(iplds.BlockNumber).To(Equal(mocks.MockConvertedPayload.Block.Number()))
			Expect(iplds.Header).To(Equal(mocks.MockIPLDs.Header))
			Expect(len(iplds.Uncles)).To(Equal(0))
			Expect(iplds.Transactions).To(Equal(mocks.MockIPLDs.Transactions))
			Expect(iplds.Receipts).To(Equal(mocks.MockIPLDs.Receipts))
			Expect(iplds.StateNodes).To(Equal(mocks.MockIPLDs.StateNodes))
			Expect(iplds.StorageNodes).To(Equal(mocks.MockIPLDs.StorageNodes))
		})
	})
})
