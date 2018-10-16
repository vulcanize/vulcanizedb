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

package ilk_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Pit file ilk repository", func() {
	var (
		db                *postgres.DB
		pitFileRepository ilk.PitFileIlkRepository
		err               error
		headerRepository  datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(core.Node{})
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		pitFileRepository = ilk.PitFileIlkRepository{DB: db}
	})

	Describe("Create", func() {
		var (
			db                   *postgres.DB
			pitFileIlkRepository ilk.PitFileIlkRepository
			err                  error
			headerID             int64
		)

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = pitFileRepository.Create(headerID, []interface{}{test_data.PitFileIlkModel})
			Expect(err).NotTo(HaveOccurred())
			pitFileIlkRepository = ilk.PitFileIlkRepository{DB: db}
		})

		It("adds a pit file ilk event", func() {
			var dbPitFile ilk.PitFileIlkModel
			err = db.Get(&dbPitFile, `SELECT ilk, what, data, tx_idx, raw_log FROM maker.pit_file_ilk WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbPitFile.Ilk).To(Equal(test_data.PitFileIlkModel.Ilk))
			Expect(dbPitFile.What).To(Equal(test_data.PitFileIlkModel.What))
			Expect(dbPitFile.Data).To(Equal(test_data.PitFileIlkModel.Data))
			Expect(dbPitFile.TransactionIndex).To(Equal(test_data.PitFileIlkModel.TransactionIndex))
			Expect(dbPitFile.Raw).To(MatchJSON(test_data.PitFileIlkModel.Raw))
		})

		It("marks header as checked for logs", func() {
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT pit_file_ilk_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("does not duplicate pit file ilk events", func() {
			err = pitFileRepository.Create(headerID, []interface{}{test_data.PitFileIlkModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes pit file ilk if corresponding header is deleted", func() {
			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbPitFile ilk.PitFileIlkModel
			err = db.Get(&dbPitFile, `SELECT ilk, what, data, tx_idx, raw_log FROM maker.pit_file_ilk WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})

	Describe("MarkHeaderChecked", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a row for a new headerID", func() {
			err = pitFileRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT pit_file_ilk_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)

			err = pitFileRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT pit_file_ilk_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		var (
			startingBlock, pitFileBlock, endingBlock int64
			blockNumbers, headerIDs                  []int64
		)

		BeforeEach(func() {
			startingBlock = GinkgoRandomSeed()
			pitFileBlock = startingBlock + 1
			endingBlock = startingBlock + 2

			blockNumbers = []int64{startingBlock, pitFileBlock, endingBlock, endingBlock + 1}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
			}
		})

		It("returns headers that haven't been checked", func() {
			err := pitFileRepository.MarkHeaderChecked(headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := pitFileRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
		})

		It("only treats headers as checked if pit file ilk logs have been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := pitFileRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(pitFileBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(pitFileBlock)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(pitFileBlock)))
		})

		It("only returns headers associated with the current node", func() {
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			for _, n := range blockNumbers {
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
			}

			pitFileRepositoryTwo := ilk.PitFileIlkRepository{DB: dbTwo}
			err := pitFileRepository.MarkHeaderChecked(headerIDs[0])
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := pitFileRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := pitFileRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
