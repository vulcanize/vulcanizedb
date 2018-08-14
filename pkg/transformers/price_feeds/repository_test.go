// Copyright © 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package price_feeds_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Price feeds repository", func() {
	Describe("Create", func() {
		It("persists a price feed update", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			blockNumber := uint64(12345)
			header := core.Header{BlockNumber: int64(blockNumber)}
			headerID, err := headerRepository.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			priceFeedUpdate := price_feeds.PriceFeedModel{
				BlockNumber:       blockNumber,
				HeaderID:          headerID,
				MedianizerAddress: []byte{1, 2, 3, 4, 5},
				UsdValue:          "123.45",
				TransactionIndex:  1,
			}
			priceFeedRepository := price_feeds.NewPriceFeedRepository(db)

			err = priceFeedRepository.Create(priceFeedUpdate)

			Expect(err).NotTo(HaveOccurred())
			var dbPriceFeedUpdate price_feeds.PriceFeedModel
			err = db.Get(&dbPriceFeedUpdate, `SELECT block_number, header_id, medianizer_address, usd_value, tx_idx FROM maker.price_feeds WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbPriceFeedUpdate).To(Equal(priceFeedUpdate))
		})

		It("does not duplicate price feed updates", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			blockNumber := uint64(12345)
			header := core.Header{BlockNumber: int64(blockNumber)}
			headerID, err := headerRepository.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			priceFeedUpdate := price_feeds.PriceFeedModel{
				BlockNumber:       blockNumber,
				HeaderID:          headerID,
				MedianizerAddress: []byte{1, 2, 3, 4, 5},
				UsdValue:          "123.45",
				TransactionIndex:  1,
			}
			priceFeedRepository := price_feeds.NewPriceFeedRepository(db)
			err = priceFeedRepository.Create(priceFeedUpdate)
			Expect(err).NotTo(HaveOccurred())

			err = priceFeedRepository.Create(priceFeedUpdate)

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("MissingHeaders", func() {
		It("returns headers with no associated price feed event", func() {
			node := core.Node{}
			db := test_config.NewTestDB(node)
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			startingBlockNumber := int64(1)
			priceFeedBlockNumber := int64(2)
			endingBlockNumber := int64(3)
			blockNumbers := []int64{startingBlockNumber, priceFeedBlockNumber, endingBlockNumber, endingBlockNumber + 1}
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
			priceFeedRepository := price_feeds.NewPriceFeedRepository(db)
			priceFeedUpdate := price_feeds.PriceFeedModel{
				BlockNumber: uint64(blockNumbers[1]),
				HeaderID:    headerIDs[1],
				UsdValue:    "123.45",
			}
			err := priceFeedRepository.Create(priceFeedUpdate)
			Expect(err).NotTo(HaveOccurred())

			headers, err := priceFeedRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
		})

		It("only returns headers associated with the current node", func() {
			nodeOne := core.Node{}
			db := test_config.NewTestDB(nodeOne)
			test_config.CleanTestDB(db)
			blockNumbers := []int64{1, 2, 3}
			headerRepository := repositories.NewHeaderRepository(db)
			nodeTwo := core.Node{ID: "second"}
			dbTwo := test_config.NewTestDB(nodeTwo)
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				Expect(err).NotTo(HaveOccurred())
			}
			priceFeedRepository := price_feeds.NewPriceFeedRepository(db)
			priceFeedRepositoryTwo := price_feeds.NewPriceFeedRepository(dbTwo)
			err := priceFeedRepository.Create(price_feeds.PriceFeedModel{
				HeaderID: headerIDs[0],
				UsdValue: "123.45",
			})
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := priceFeedRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := priceFeedRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
