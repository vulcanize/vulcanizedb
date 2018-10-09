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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_fold"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("", func() {

	Describe("Create", func() {

		var db *postgres.DB
		var headerID int64
		var vatFoldRepository vat_fold.VatFoldRepository

		BeforeEach(func() {
			db = test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)

			headerRepository := repositories.NewHeaderRepository(db)
			id, err := headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			headerID = id

			vatFoldRepository = vat_fold.NewVatFoldRepository(db)
			err = vatFoldRepository.Create(headerID, []vat_fold.VatFoldModel{test_data.VatFoldModel})
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

		It("does not duplicate vat events", func() {
			err := vatFoldRepository.Create(headerID, []vat_fold.VatFoldModel{test_data.VatFoldModel})

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

	Describe("MissingHeaders", func() {

		It("returns headers with no associated vat event", func() {
			startBlock := int64(1)
			eventBlock := int64(2)
			finalBlock := int64(3)
			blockNumbers := []int64{startBlock, eventBlock, finalBlock, finalBlock + 1}

			repository := initRepository(
				repositoryOptions{
					cleanDB:                true,
					blockNumbers:           blockNumbers,
					storeEvent:             true,
					storedEventBlockNumber: 1,
				},
			)

			headers, err := repository.MissingHeaders(startBlock, finalBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startBlock), Equal(finalBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startBlock), Equal(finalBlock)))
		})

		It("only returns headers associated with the current node", func() {
			blockNumbers := []int64{1, 2, 3}

			repositoryOne := initRepository(
				repositoryOptions{
					blockNumbers: blockNumbers,
					nodeID:       "first",
					cleanDB:      true,
					storeEvent:   true,
				},
			)

			repositoryTwo := initRepository(
				repositoryOptions{
					blockNumbers: blockNumbers,
					nodeID:       "second",
					cleanDB:      false,
					storeEvent:   false,
				},
			)

			nodeOneMissingHeaders, err := repositoryOne.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := repositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})

// -------------------------------------------------------------------------------------------

type repositoryOptions struct {
	blockNumbers           []int64
	nodeID                 string
	cleanDB                bool
	storeEvent             bool
	storedEventBlockNumber int64
}

func initRepository(options repositoryOptions) vat_fold.VatFoldRepository {
	db := test_config.NewTestDB(core.Node{ID: options.nodeID})
	if options.cleanDB {
		test_config.CleanTestDB(db)
	}

	headerRepository := repositories.NewHeaderRepository(db)
	vatfoldRepository := vat_fold.NewVatFoldRepository(db)

	var headerIDs []int64
	for _, n := range options.blockNumbers {
		headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
		headerIDs = append(headerIDs, headerID)
		Expect(err).NotTo(HaveOccurred())
	}

	if options.storeEvent {
		err := vatfoldRepository.Create(
			headerIDs[options.storedEventBlockNumber],
			[]vat_fold.VatFoldModel{test_data.VatFoldModel},
		)
		Expect(err).NotTo(HaveOccurred())
	}

	return vatfoldRepository
}
