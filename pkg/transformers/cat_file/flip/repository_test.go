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

package flip_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/flip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Cat file flip repository", func() {
	Describe("Create", func() {
		var catFileRepository flip.CatFileFlipRepository
		var db *postgres.DB
		var err error
		var headerID int64

		BeforeEach(func() {
			db = test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err = headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			catFileRepository = flip.NewCatFileFlipRepository(db)
		})

		It("adds a cat file flip event", func() {
			err = catFileRepository.Create(headerID, []flip.CatFileFlipModel{test_data.CatFileFlipModel})

			Expect(err).NotTo(HaveOccurred())
			var dbPitFile flip.CatFileFlipModel
			err = db.Get(&dbPitFile, `SELECT ilk, what, flip, tx_idx, raw_log FROM maker.cat_file_flip WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbPitFile.Ilk).To(Equal(test_data.CatFileFlipModel.Ilk))
			Expect(dbPitFile.What).To(Equal(test_data.CatFileFlipModel.What))
			Expect(dbPitFile.Flip).To(Equal(test_data.CatFileFlipModel.Flip))
			Expect(dbPitFile.TransactionIndex).To(Equal(test_data.CatFileFlipModel.TransactionIndex))
			Expect(dbPitFile.Raw).To(MatchJSON(test_data.CatFileFlipModel.Raw))
		})

		It("marks header as checked for logs", func() {
			err = catFileRepository.Create(headerID, []flip.CatFileFlipModel{test_data.CatFileFlipModel})

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT cat_file_flip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("does not duplicate cat file flip events", func() {
			err = catFileRepository.Create(headerID, []flip.CatFileFlipModel{test_data.CatFileFlipModel})
			Expect(err).NotTo(HaveOccurred())

			err = catFileRepository.Create(headerID, []flip.CatFileFlipModel{test_data.CatFileFlipModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes cat file flip if corresponding header is deleted", func() {
			err = catFileRepository.Create(headerID, []flip.CatFileFlipModel{test_data.CatFileFlipModel})
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbPitFile flip.CatFileFlipModel
			err = db.Get(&dbPitFile, `SELECT ilk, what, flip, tx_idx, raw_log FROM maker.cat_file_flip WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})

	Describe("MarkHeaderChecked", func() {
		It("creates a row for a new headerID", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			catFileRepository := flip.NewCatFileFlipRepository(db)

			err = catFileRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT cat_file_flip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			catFileRepository := flip.NewCatFileFlipRepository(db)
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)

			err = catFileRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT cat_file_flip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		It("returns headers that haven't been checked", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			startingBlockNumber := int64(1)
			catFileBlockNumber := int64(2)
			endingBlockNumber := int64(3)
			blockNumbers := []int64{startingBlockNumber, catFileBlockNumber, endingBlockNumber, endingBlockNumber + 1}
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
			catFileRepository := flip.NewCatFileFlipRepository(db)
			err := catFileRepository.MarkHeaderChecked(headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := catFileRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
		})

		It("only treats headers as checked if cat file flip logs have been checked", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			startingBlockNumber := int64(1)
			catFiledBlockNumber := int64(2)
			endingBlockNumber := int64(3)
			blockNumbers := []int64{startingBlockNumber, catFiledBlockNumber, endingBlockNumber, endingBlockNumber + 1}
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
			catFiledRepository := flip.NewCatFileFlipRepository(db)
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := catFiledRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(catFiledBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(catFiledBlockNumber)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(catFiledBlockNumber)))
		})

		It("only returns headers associated with the current node", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			blockNumbers := []int64{1, 2, 3}
			headerRepository := repositories.NewHeaderRepository(db)
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				Expect(err).NotTo(HaveOccurred())
			}
			catFileRepository := flip.NewCatFileFlipRepository(db)
			catFileRepositoryTwo := flip.NewCatFileFlipRepository(dbTwo)
			err := catFileRepository.MarkHeaderChecked(headerIDs[0])
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := catFileRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := catFileRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
