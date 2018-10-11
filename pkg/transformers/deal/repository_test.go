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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/deal"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Deal Repository", func() {
	var (
		db               *postgres.DB
		dealRepository   deal.DealRepository
		headerRepository repositories.HeaderRepository
		err              error
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		dealRepository = deal.NewDealRepository(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		var headerId int64

		BeforeEach(func() {
			headerId, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err := dealRepository.Create(headerId, []deal.DealModel{test_data.DealModel})
			Expect(err).NotTo(HaveOccurred())
		})

		It("persists a deal record", func() {
			var count int
			db.QueryRow(`SELECT count(*) FROM maker.deal`).Scan(&count)
			Expect(count).To(Equal(1))
			var dbResult deal.DealModel
			err = db.Get(&dbResult, `SELECT bid_id, contract_address, tx_idx, raw_log FROM maker.deal WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.BidId).To(Equal(test_data.DealModel.BidId))
			Expect(dbResult.ContractAddress).To(Equal(test_data.DealModel.ContractAddress))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.DealModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.DealModel.Raw))
		})

		It("marks header as checked for logs", func() {
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT deal_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("returns an error if inserting a deal record fails", func() {
			err = dealRepository.Create(headerId, []deal.DealModel{test_data.DealModel})
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

	Describe("MarkHeaderChecked", func() {
		var headerId int64

		BeforeEach(func() {
			headerId, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a row for a new headerID", func() {
			err = dealRepository.MarkHeaderChecked(headerId)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT deal_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerId)

			err = dealRepository.MarkHeaderChecked(headerId)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT deal_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		var (
			dealBlock, startingBlock, endingBlock int64
			blockNumbers, headerIds               []int64
		)

		BeforeEach(func() {
			dealBlock = GinkgoRandomSeed()
			startingBlock = dealBlock - 1
			endingBlock = dealBlock + 1
			outOfRangeBlockNumber := dealBlock + 2

			blockNumbers = []int64{startingBlock, dealBlock, endingBlock, outOfRangeBlockNumber}

			headerIds = []int64{}
			for _, number := range blockNumbers {
				headerId, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(number))
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)
			}
		})

		It("returns header records that don't have a corresponding deals", func() {
			err = dealRepository.MarkHeaderChecked(headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			missingHeaders, err := dealRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(2))
			Expect(missingHeaders[0].BlockNumber).To(Equal(startingBlock))
			Expect(missingHeaders[1].BlockNumber).To(Equal(endingBlock))
		})

		It("only treats headers as checked if deal have been checked", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := dealRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(dealBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(dealBlock)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(dealBlock)))
		})

		It("only returns missing headers for the given node", func() {
			err = dealRepository.MarkHeaderChecked(headerIds[1])
			Expect(err).NotTo(HaveOccurred())
			node2 := core.Node{}
			db2 := test_config.NewTestDB(node2)
			dealRepository2 := deal.NewDealRepository(db2)
			headerRepository2 := repositories.NewHeaderRepository(db2)
			var node2HeaderIds []int64
			for _, number := range blockNumbers {
				id, err := headerRepository2.CreateOrUpdateHeader(fakes.GetFakeHeader(number))
				node2HeaderIds = append(node2HeaderIds, id)
				Expect(err).NotTo(HaveOccurred())
			}

			missingHeadersNode1, err := dealRepository.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(missingHeadersNode1)).To(Equal(2))
			Expect(missingHeadersNode1[0].BlockNumber).To(Equal(startingBlock))
			Expect(missingHeadersNode1[1].BlockNumber).To(Equal(endingBlock))

			missingHeadersNode2, err := dealRepository2.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(missingHeadersNode2)).To(Equal(3))
			Expect(missingHeadersNode2[0].BlockNumber).To(Or(Equal(startingBlock), Equal(dealBlock), Equal(endingBlock)))
			Expect(missingHeadersNode2[1].BlockNumber).To(Or(Equal(startingBlock), Equal(dealBlock), Equal(endingBlock)))
			Expect(missingHeadersNode2[2].BlockNumber).To(Or(Equal(startingBlock), Equal(dealBlock), Equal(endingBlock)))
		})
	})
})
