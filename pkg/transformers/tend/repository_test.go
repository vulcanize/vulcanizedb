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

package tend_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("TendRepository", func() {
	var (
		db               *postgres.DB
		tendRepository   tend.TendRepository
		headerRepository repositories.HeaderRepository
		headerId         int64
		err              error
	)

	BeforeEach(func() {
		node := test_config.NewTestNode()
		db = test_config.NewTestDB(node)
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		tendRepository = tend.TendRepository{}
		tendRepository.SetDB(db)
	})

	Describe("Create", func() {
		BeforeEach(func() {
			headerId, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = tendRepository.Create(headerId, []interface{}{test_data.TendModel})
			Expect(err).NotTo(HaveOccurred())
		})

		It("persists a tend record", func() {
			var count int
			err = db.QueryRow(`SELECT count(*) from maker.tend`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			dbResult := tend.TendModel{}
			err = db.Get(&dbResult, `SELECT bid_id, lot, bid, guy, tic, tx_idx, raw_log FROM maker.tend WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbResult.BidId).To(Equal(test_data.TendModel.BidId))
			Expect(dbResult.Lot).To(Equal(test_data.TendModel.Lot))
			Expect(dbResult.Bid).To(Equal(test_data.TendModel.Bid))
			Expect(dbResult.Guy).To(Equal(test_data.TendModel.Guy))
			Expect(dbResult.Tic).To(Equal(test_data.TendModel.Tic))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.TendModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.RawLogNoteJson))
		})

		It("marks header as checked", func() {
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT tend_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("returns an error if inserting a tend record fails", func() {
			err = tendRepository.Create(headerId, []interface{}{test_data.TendModel})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("deletes the tend record if its corresponding header record is deleted", func() {
			var count int
			err = db.QueryRow(`SELECT count(*) from maker.tend`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			_, err = db.Exec(`DELETE FROM headers where id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())

			err = db.QueryRow(`SELECT count(*) from maker.tend`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})

		It("Returns an error if model is of wrong type", func() {
			err = tendRepository.Create(headerId, []interface{}{test_data.WrongModel{}})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("model of type"))
		})
	})

	Describe("MarkHeaderChecked", func() {
		BeforeEach(func() {
			headerId, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = tendRepository.Create(headerId, []interface{}{test_data.TendModel})
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a row for a new headerID", func() {
			err = tendRepository.MarkHeaderChecked(headerId)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT tend_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerId)

			err = tendRepository.MarkHeaderChecked(headerId)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT tend_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		var (
			tendBlock, startingBlock, endingBlock, outOfRangeBlock int64
			headerIds                                              []int64
		)

		BeforeEach(func() {
			tendBlock = GinkgoRandomSeed()
			startingBlock = tendBlock - 1
			endingBlock = tendBlock + 1
			outOfRangeBlock = tendBlock + 2

			headerIds = []int64{}
			for _, number := range []int64{startingBlock, tendBlock, endingBlock, outOfRangeBlock} {
				headerId, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(number))
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)
			}
		})

		It("returns headers for which there isn't an associated tend record", func() {
			err = tendRepository.MarkHeaderChecked(headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := tendRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
		})

		It("only treats headers as checked if deal have been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := tendRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(tendBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(tendBlock)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(tendBlock)))
		})

		It("only returns missing headers for the current node", func() {
			node2 := core.Node{}
			db2 := test_config.NewTestDB(node2)
			headerRepository2 := repositories.NewHeaderRepository(db2)
			tendRepository2 := tend.TendRepository{}
			tendRepository2.SetDB(db2)

			for _, number := range []int64{startingBlock, tendBlock, endingBlock} {
				headerRepository2.CreateOrUpdateHeader(fakes.GetFakeHeader(number))
				Expect(err).NotTo(HaveOccurred())
			}

			err = tendRepository.MarkHeaderChecked(headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			node1MissingHeaders, err := tendRepository.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node1MissingHeaders)).To(Equal(2))

			node2MissingHeaders, err := tendRepository2.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node2MissingHeaders)).To(Equal(3))
		})
	})
})
