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

package vat_fold_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_fold"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat.fold repository", func() {
	var (
		db               *postgres.DB
		repository       vat_fold.VatFoldRepository
		headerRepository repositories.HeaderRepository
		err              error
	)

	BeforeEach(func() {
		node := test_config.NewTestNode()
		db = test_config.NewTestDB(node)
		test_config.CleanTestDB(db)
		repository = vat_fold.NewVatFoldRepository(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(headerID, []vat_fold.VatFoldModel{test_data.VatFoldModel})
			Expect(err).NotTo(HaveOccurred())
		})

		It("adds a vat event", func() {
			var dbVatFold vat_fold.VatFoldModel
			err := db.Get(&dbVatFold, `SELECT ilk, urn, rate, tx_idx, raw_log FROM maker.vat_fold WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbVatFold.Ilk).To(Equal(test_data.VatFoldModel.Ilk))
			Expect(dbVatFold.Urn).To(Equal(test_data.VatFoldModel.Urn))
			Expect(dbVatFold.Rate).To(Equal(test_data.VatFoldModel.Rate))
			Expect(dbVatFold.TransactionIndex).To(Equal(test_data.VatFoldModel.TransactionIndex))
			Expect(dbVatFold.Raw).To(MatchJSON(test_data.VatFoldModel.Raw))
		})

		It("marks header as checked for logs", func() {
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_fold_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("does not duplicate vat events", func() {
			err := repository.Create(headerID, []vat_fold.VatFoldModel{test_data.VatFoldModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes vat if corresponding header is deleted", func() {
			_, err := db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())

			var dbVatFold vat_fold.VatFoldModel
			err = db.Get(&dbVatFold, `SELECT ilk, tx_idx, raw_log FROM maker.vat_fold WHERE header_id = $1`, headerID)
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

		type CheckedHeaderResult struct {
			VatFoldChecked bool `db:"vat_fold_checked"`
		}

		It("creates a new checked header record", func() {
			err := repository.MarkHeaderChecked(headerID)
			Expect(err).NotTo(HaveOccurred())

			var checkedHeaderResult = CheckedHeaderResult{}
			err = db.Get(&checkedHeaderResult, `SELECT vat_fold_checked FROM checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(checkedHeaderResult.VatFoldChecked).To(BeTrue())
		})

		It("updates an existing checked header", func() {
			_, err := repository.DB.Exec(`INSERT INTO checked_headers (header_id) VALUES($1)`, headerID)
			Expect(err).NotTo(HaveOccurred())

			var checkedHeaderResult CheckedHeaderResult
			err = db.Get(&checkedHeaderResult, `SELECT vat_fold_checked FROM checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(checkedHeaderResult.VatFoldChecked).To(BeFalse())

			err = repository.MarkHeaderChecked(headerID)
			Expect(err).NotTo(HaveOccurred())

			err = db.Get(&checkedHeaderResult, `SELECT vat_fold_checked FROM checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(checkedHeaderResult.VatFoldChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		var (
			startingBlock, vatFoldBlock, endingBlock int64
			blockNumbers, headerIDs                  []int64
		)

		BeforeEach(func() {
			startingBlock = GinkgoRandomSeed()
			vatFoldBlock = startingBlock + 1
			endingBlock = startingBlock + 2

			blockNumbers = []int64{startingBlock, vatFoldBlock, endingBlock, endingBlock + 1}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
			}
		})

		It("returns headers that haven't been checked", func() {
			err := repository.MarkHeaderChecked(headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := repository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
		})

		It("only treats headers as checked if vat fold logs have been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := repository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(vatFoldBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(vatFoldBlock)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(vatFoldBlock)))
		})

		It("only returns headers associated with the current node", func() {
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)

			for _, n := range blockNumbers {
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
			}

			repositoryTwo := vat_fold.NewVatFoldRepository(dbTwo)

			err := repository.MarkHeaderChecked(headerIDs[0])
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := repository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := repositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
