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

package chop_lump_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/chop_lump"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Cat file chop lump repository", func() {
	var (
		catFileRepository chop_lump.CatFileChopLumpRepository
		db                *postgres.DB
		err               error
		headerRepository  datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		catFileRepository = chop_lump.NewCatFileChopLumpRepository(db)
	})

	Describe("Create", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = catFileRepository.Create(headerID, []chop_lump.CatFileChopLumpModel{test_data.CatFileChopLumpModel})
			Expect(err).NotTo(HaveOccurred())
		})

		It("adds a cat file chop lump event", func() {
			var dbResult chop_lump.CatFileChopLumpModel
			err = db.Get(&dbResult, `SELECT ilk, what, data, tx_idx, raw_log FROM maker.cat_file_chop_lump WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(test_data.CatFileChopLumpModel.Ilk))
			Expect(dbResult.What).To(Equal(test_data.CatFileChopLumpModel.What))
			Expect(dbResult.Data).To(Equal(test_data.CatFileChopLumpModel.Data))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.CatFileChopLumpModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.CatFileChopLumpModel.Raw))
		})

		It("marks header as checked for logs", func() {
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT cat_file_chop_lump_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("does not duplicate cat file chop lump events", func() {
			err = catFileRepository.Create(headerID, []chop_lump.CatFileChopLumpModel{test_data.CatFileChopLumpModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes cat file chop lump if corresponding header is deleted", func() {
			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbResult chop_lump.CatFileChopLumpModel
			err = db.Get(&dbResult, `SELECT ilk, what, data, tx_idx, raw_log FROM maker.cat_file_chop_lump WHERE header_id = $1`, headerID)
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
			err = catFileRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT cat_file_chop_lump_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)

			err = catFileRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT cat_file_chop_lump_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		var (
			startingBlock, endingBlock, catFileBlock int64
			blockNumbers, headerIDs                  []int64
		)

		BeforeEach(func() {
			startingBlock = GinkgoRandomSeed()
			catFileBlock = startingBlock + 1
			endingBlock = startingBlock + 2

			blockNumbers = []int64{startingBlock, catFileBlock, endingBlock, endingBlock + 1}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("returns headers that haven't been checked", func() {
			err := catFileRepository.MarkHeaderChecked(headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := catFileRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
		})

		It("only treats headers as checked if cat file chop lump logs have been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := catFileRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(catFileBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(catFileBlock)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(catFileBlock)))
		})

		It("only returns headers associated with the current node", func() {
			err := catFileRepository.MarkHeaderChecked(headerIDs[0])
			Expect(err).NotTo(HaveOccurred())
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			for _, n := range blockNumbers {
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
			}
			catFileRepositoryTwo := chop_lump.NewCatFileChopLumpRepository(dbTwo)

			nodeOneMissingHeaders, err := catFileRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := catFileRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
