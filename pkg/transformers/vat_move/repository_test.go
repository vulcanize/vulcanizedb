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

package vat_move_test

import (
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_move"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat Move", func() {
	var db *postgres.DB
	var headerRepository repositories.HeaderRepository
	var vatMoveRepository vat_move.VatMoveRepository

	BeforeEach(func() {
		db = test_config.NewTestDB(core.Node{})
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		vatMoveRepository = vat_move.VatMoveRepository{DB: db}
	})

	Describe("Create", func() {
		var headerID int64
		var err error

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = vatMoveRepository.Create(headerID, []interface{}{test_data.VatMoveModel})
			Expect(err).NotTo(HaveOccurred())
		})

		It("adds a vat move event", func() {
			var dbVatMove vat_move.VatMoveModel
			err = db.Get(&dbVatMove, `SELECT src, dst, rad, tx_idx, raw_log FROM maker.vat_move WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbVatMove.Src).To(Equal(test_data.VatMoveModel.Src))
			Expect(dbVatMove.Dst).To(Equal(test_data.VatMoveModel.Dst))
			Expect(dbVatMove.Rad).To(Equal(test_data.VatMoveModel.Rad))
			Expect(dbVatMove.TransactionIndex).To(Equal(test_data.VatMoveModel.TransactionIndex))
			Expect(dbVatMove.Raw).To(MatchJSON(test_data.VatMoveModel.Raw))
		})

		It("marks header id as checked for vat move logs", func() {
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_move_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("returns an error if insertion fails", func() {
			err = vatMoveRepository.Create(headerID, []interface{}{test_data.VatMoveModel})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes vat move event if corresponding header is deleted", func() {
			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())

			var dbVatMove vat_move.VatMoveModel
			err = db.Get(&dbVatMove, `SELECT src, dst, rad, tx_idx, raw_log FROM maker.vat_move WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})

		It("Returns an error if model is of wrong type", func() {
			err = vatMoveRepository.Create(headerID, []interface{}{test_data.WrongModel{}})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("model of type"))
		})
	})

	Describe("MissingHeaders", func() {
		var eventBlockNumber = GinkgoRandomSeed()
		var startingBlockNumber = eventBlockNumber - 1
		var endingBlockNumber = eventBlockNumber + 1
		var outOfRangeBlockNumber = eventBlockNumber + 2

		It("returns headers haven't been checked", func() {
			var headerIds []int64

			for _, number := range []int64{startingBlockNumber, eventBlockNumber, endingBlockNumber, outOfRangeBlockNumber} {
				headerId, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(number))
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)
			}

			err := vatMoveRepository.MarkHeaderChecked(headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := vatMoveRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
		})

		It("only treats headers as checked if vat_move has been checked", func() {
			var headerIds []int64
			for _, number := range []int64{startingBlockNumber, eventBlockNumber, endingBlockNumber, outOfRangeBlockNumber} {
				headerId, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(number))
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)
			}

			// Just creates row, doesn't set this header as checked
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := vatMoveRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventBlockNumber)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventBlockNumber)))
		})

		It("only returns headers associated with the current node", func() {
			blockNumbers := []int64{1, 2, 3}
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
			}
			vatMoveRepositoryTwo := vat_move.VatMoveRepository{DB: dbTwo}
			err := vatMoveRepository.Create(headerIDs[0], []interface{}{test_data.VatMoveModel})
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := vatMoveRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := vatMoveRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})

	Describe("MarkHeaderChecked", func() {
		var headerID int64
		var err error

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a new row for a new headerID", func() {
			err = vatMoveRepository.MarkHeaderChecked(headerID)
			Expect(err).NotTo(HaveOccurred())

			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_move_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerId already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)
			Expect(err).NotTo(HaveOccurred())
			err = vatMoveRepository.MarkHeaderChecked(headerID)
			Expect(err).NotTo(HaveOccurred())

			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_move_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("SetDB", func() {
		It("sets the repository db", func() {
			vatMoveRepository.DB = nil
			Expect(vatMoveRepository.DB).To(BeNil())
			vatMoveRepository.SetDB(db)
			Expect(vatMoveRepository.DB).To(Equal(db))
		})
	})
})
