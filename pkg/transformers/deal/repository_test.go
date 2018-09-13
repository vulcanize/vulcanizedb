// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deal_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/deal"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Deal Repository", func() {
	var node core.Node
	var db *postgres.DB
	var dealRepository deal.DealRepository
	var headerRepository repositories.HeaderRepository
	var headerId int64
	var err error

	BeforeEach(func() {
		node = test_config.NewTestNode()
		db = test_config.NewTestDB(node)
		test_config.CleanTestDB(db)
		dealRepository = deal.NewDealRepository(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})
	Describe("Create", func() {
		BeforeEach(func() {
			headerId, err = headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			err := dealRepository.Create(headerId, test_data.DealModel)
			Expect(err).NotTo(HaveOccurred())
		})

		It("persists a deal record", func() {
			var count int
			db.QueryRow(`SELECT count(*) FROM maker.deal`).Scan(&count)
			Expect(count).To(Equal(1))
			var dbResult deal.DealModel
			err = db.Get(&dbResult, `SELECT bid_id, tx_idx, raw_log FROM maker.deal WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.BidId).To(Equal(test_data.DealModel.BidId))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.DealModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.DealModel.Raw))
		})

		It("returns an error if inserting a deal record fails", func() {
			err = dealRepository.Create(headerId, test_data.DealModel)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("deletes the deal record if its corresponding header record is deleted", func() {
			var count int
			err = db.QueryRow(`SELECT count(*) from maker.deal`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
			_, err = db.Exec(`DELETE FROM headers where id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			err = db.QueryRow(`SELECT count(*) from maker.deal`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})

	Describe("MissingHeaders", func() {
		var dealBlockNumber int64
		var startingBlockNumber int64
		var endingBlockNumber int64
		var blockNumbers []int64

		BeforeEach(func() {
			dealBlockNumber = rand.Int63()
			startingBlockNumber = dealBlockNumber - 1
			endingBlockNumber = dealBlockNumber + 1
			outOfRangeBlockNumber := dealBlockNumber + 2
			blockNumbers = []int64{startingBlockNumber, dealBlockNumber, endingBlockNumber, outOfRangeBlockNumber}
			var headerIds []int64
			for _, number := range blockNumbers {
				headerId, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: number})
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)
			}
			dealRepository.Create(headerIds[1], test_data.DealModel)
		})

		It("returns header records that don't have a corresponding deals", func() {
			missingHeaders, err := dealRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(2))
			Expect(missingHeaders[0].BlockNumber).To(Equal(startingBlockNumber))
			Expect(missingHeaders[1].BlockNumber).To(Equal(endingBlockNumber))
		})

		It("only returns missing headers for the given node", func() {
			node2 := core.Node{}
			db2 := test_config.NewTestDB(node2)
			dealRepository2 := deal.NewDealRepository(db2)
			headerRepository2 := repositories.NewHeaderRepository(db2)
			var node2HeaderIds []int64
			for _, number := range blockNumbers {
				id, err := headerRepository2.CreateOrUpdateHeader(core.Header{BlockNumber: number})
				node2HeaderIds = append(node2HeaderIds, id)
				Expect(err).NotTo(HaveOccurred())
			}
			missingHeadersNode1, err := dealRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(missingHeadersNode1)).To(Equal(2))
			Expect(missingHeadersNode1[0].BlockNumber).To(Equal(startingBlockNumber))
			Expect(missingHeadersNode1[1].BlockNumber).To(Equal(endingBlockNumber))
			missingHeadersNode2, err := dealRepository2.MissingHeaders(startingBlockNumber, endingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(missingHeadersNode2)).To(Equal(3))
			Expect(missingHeadersNode2[0].BlockNumber).To(Equal(startingBlockNumber))
			Expect(missingHeadersNode2[1].BlockNumber).To(Equal(dealBlockNumber))
			Expect(missingHeadersNode2[2].BlockNumber).To(Equal(endingBlockNumber))
		})
	})
})
