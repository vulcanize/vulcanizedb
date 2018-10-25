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

package flog_test

import (
	"database/sql"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flog"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Flog repository", func() {
	var (
		db               *postgres.DB
		flogRepository   flog.FlogRepository
		err              error
		headerRepository datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(core.Node{})
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		flogRepository = flog.FlogRepository{}
		flogRepository.SetDB(db)
	})

	Describe("Create", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("adds a flog event", func() {
			err = flogRepository.Create(headerID, []interface{}{test_data.FlogModel})
			Expect(err).NotTo(HaveOccurred())

			var dbFlog flog.FlogModel
			err = db.Get(&dbFlog, `SELECT era, log_idx, tx_idx, raw_log FROM maker.vow_flog WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbFlog.Era).To(Equal(test_data.FlogModel.Era))
			Expect(dbFlog.LogIndex).To(Equal(test_data.FlogModel.LogIndex))
			Expect(dbFlog.TransactionIndex).To(Equal(test_data.FlogModel.TransactionIndex))
			Expect(dbFlog.Raw).To(MatchJSON(test_data.FlogModel.Raw))
		})

		It("marks header as checked for logs", func() {
			err = flogRepository.Create(headerID, []interface{}{test_data.FlogModel})
			Expect(err).NotTo(HaveOccurred())

			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vow_flog_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates the header to checked if checked headers row already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)
			Expect(err).NotTo(HaveOccurred())

			err = flogRepository.Create(headerID, []interface{}{test_data.FlogModel})

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vow_flog_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("does not duplicate flog events", func() {
			err = flogRepository.Create(headerID, []interface{}{test_data.FlogModel})
			Expect(err).NotTo(HaveOccurred())

			err = flogRepository.Create(headerID, []interface{}{test_data.FlogModel})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes flog events if corresponding header is deleted", func() {
			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbFlog flog.FlogModel
			err = db.Get(&dbFlog, `SELECT era, log_idx, tx_idx, raw_log FROM maker.vow_flog WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})

		It("returns an error if model is of wrong type", func() {
			err = flogRepository.Create(headerID, []interface{}{test_data.WrongModel{}})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("model of type"))
		})
	})

	Describe("MarkHeaderChecked", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a row for a new headerID", func() {
			err = flogRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vow_flog_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)

			err = flogRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vow_flog_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		var (
			startingBlockNumber, flogBlockNumber, endingBlockNumber int64
			blockNumbers, headerIDs                                 []int64
		)

		BeforeEach(func() {
			startingBlockNumber = rand.Int63()
			flogBlockNumber = startingBlockNumber + 1
			endingBlockNumber = startingBlockNumber + 2

			blockNumbers = []int64{startingBlockNumber, flogBlockNumber, endingBlockNumber, endingBlockNumber + 1}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
			}
		})

		It("returns headers that haven't been checked", func() {
			err := flogRepository.MarkHeaderChecked(headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := flogRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
		})

		It("only treats headers as checked if flog logs have been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := flogRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(flogBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(flogBlockNumber)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(flogBlockNumber)))
		})

		It("only returns headers associated with the current node", func() {
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			for _, n := range blockNumbers {
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
			}
			flogRepositoryTwo := flog.FlogRepository{}
			flogRepositoryTwo.SetDB(dbTwo)
			err := flogRepository.MarkHeaderChecked(headerIDs[0])
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := flogRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := flogRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
