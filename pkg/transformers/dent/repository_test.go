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

package dent_test

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Dent Repository", func() {
	var (
		db               *postgres.DB
		dentRepository   dent.DentRepository
		headerRepository repositories.HeaderRepository
		err              error
		rawHeader        []byte
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		dentRepository = dent.NewDentRepository(db)
		headerRepository = repositories.NewHeaderRepository(db)
		rawHeader, err = json.Marshal(types.Header{})
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Create", func() {
		var headerId int64

		BeforeEach(func() {
			headerId, err = headerRepository.CreateOrUpdateHeader(core.Header{Raw: rawHeader})
			Expect(err).NotTo(HaveOccurred())

			err := dentRepository.Create(headerId, []dent.DentModel{test_data.DentModel})
			Expect(err).NotTo(HaveOccurred())
		})

		It("persists a dent record", func() {
			var count int
			db.QueryRow(`SELECT count(*) FROM maker.dent`).Scan(&count)
			Expect(count).To(Equal(1))

			var dbResult dent.DentModel
			err = db.Get(&dbResult, `SELECT bid_id, lot, bid, guy, tic, tx_idx, raw_log FROM maker.dent WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.BidId).To(Equal(test_data.DentModel.BidId))
			Expect(dbResult.Lot).To(Equal(test_data.DentModel.Lot))
			Expect(dbResult.Bid).To(Equal(test_data.DentModel.Bid))
			Expect(dbResult.Guy).To(Equal(test_data.DentModel.Guy))
			Expect(dbResult.Tic).To(Equal(test_data.DentModel.Tic))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.DentModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.DentModel.Raw))
		})

		It("marks header as checked for logs", func() {
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT dent_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("returns an error if inserting a dent record fails", func() {
			err = dentRepository.Create(headerId, []dent.DentModel{test_data.DentModel})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("deletes the tend record if its corresponding header record is deleted", func() {
			var count int
			err = db.QueryRow(`SELECT count(*) from maker.dent`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			_, err = db.Exec(`DELETE FROM headers where id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())

			err = db.QueryRow(`SELECT count(*) from maker.dent`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})

	Describe("MarkHeaderChecked", func() {
		var headerId int64

		BeforeEach(func() {
			headerId, err = headerRepository.CreateOrUpdateHeader(core.Header{Raw: rawHeader})
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a row for a new headerId", func() {
			err = dentRepository.MarkHeaderChecked(headerId)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT dent_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerId already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerId)
			Expect(err).NotTo(HaveOccurred())

			err = dentRepository.MarkHeaderChecked(headerId)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT dent_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		var (
			dentBlock, startingBlock, endingBlock int64
			blockNumbers, headerIds               []int64
		)

		BeforeEach(func() {
			dentBlock = GinkgoRandomSeed()
			startingBlock = dentBlock - 1
			endingBlock = dentBlock + 1
			outOfRangeBlockNumber := dentBlock + 2

			blockNumbers = []int64{startingBlock, dentBlock, endingBlock, outOfRangeBlockNumber}

			headerIds = []int64{}
			for _, number := range blockNumbers {
				headerId, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: number, Raw: rawHeader})
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)
			}
		})

		It("returns header records that haven't been checked", func() {
			dentRepository.MarkHeaderChecked(headerIds[1])
			missingHeaders, err := dentRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(2))
			Expect(missingHeaders[0].BlockNumber).To(Equal(startingBlock))
			Expect(missingHeaders[1].BlockNumber).To(Equal(endingBlock))
		})

		It("only treats headers as checked if dent has been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := dentRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(dentBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(dentBlock)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(dentBlock)))
		})

		It("only returns missing headers for the given node", func() {
			dentRepository.MarkHeaderChecked(headerIds[1])
			node2 := core.Node{}
			db2 := test_config.NewTestDB(node2)
			dentRepository2 := dent.NewDentRepository(db2)
			headerRepository2 := repositories.NewHeaderRepository(db2)
			var node2HeaderIds []int64
			for _, number := range blockNumbers {
				id, err := headerRepository2.CreateOrUpdateHeader(core.Header{BlockNumber: number, Raw: rawHeader})
				node2HeaderIds = append(node2HeaderIds, id)
				Expect(err).NotTo(HaveOccurred())
			}

			missingHeadersNode1, err := dentRepository.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(missingHeadersNode1)).To(Equal(2))
			Expect(missingHeadersNode1[0].BlockNumber).To(Equal(startingBlock))
			Expect(missingHeadersNode1[1].BlockNumber).To(Equal(endingBlock))

			missingHeadersNode2, err := dentRepository2.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(missingHeadersNode2)).To(Equal(3))
			Expect(missingHeadersNode2[0].BlockNumber).To(Or(Equal(startingBlock), Equal(dentBlock), Equal(endingBlock)))
			Expect(missingHeadersNode2[1].BlockNumber).To(Or(Equal(startingBlock), Equal(dentBlock), Equal(endingBlock)))
			Expect(missingHeadersNode2[2].BlockNumber).To(Or(Equal(startingBlock), Equal(dentBlock), Equal(endingBlock)))
		})
	})
})
