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

package debt_ceiling_test

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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/debt_ceiling"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Pit file debt ceiling repository", func() {
	var (
		db                           *postgres.DB
		pitFileDebtCeilingRepository debt_ceiling.PitFileDebtCeilingRepository
		err                          error
		headerRepository             datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(core.Node{})
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		pitFileDebtCeilingRepository = debt_ceiling.PitFileDebtCeilingRepository{}
		pitFileDebtCeilingRepository.SetDB(db)
	})

	Describe("Create", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("adds a pit file debt ceiling event", func() {
			err = pitFileDebtCeilingRepository.Create(headerID, []interface{}{test_data.PitFileDebtCeilingModel})

			Expect(err).NotTo(HaveOccurred())
			var dbPitFile debt_ceiling.PitFileDebtCeilingModel
			err = db.Get(&dbPitFile, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.pit_file_debt_ceiling WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbPitFile.What).To(Equal(test_data.PitFileDebtCeilingModel.What))
			Expect(dbPitFile.Data).To(Equal(test_data.PitFileDebtCeilingModel.Data))
			Expect(dbPitFile.LogIndex).To(Equal(test_data.PitFileDebtCeilingModel.LogIndex))
			Expect(dbPitFile.TransactionIndex).To(Equal(test_data.PitFileDebtCeilingModel.TransactionIndex))
			Expect(dbPitFile.Raw).To(MatchJSON(test_data.PitFileDebtCeilingModel.Raw))
		})

		It("marks header as checked for logs", func() {
			err = pitFileDebtCeilingRepository.Create(headerID, []interface{}{test_data.PitFileDebtCeilingModel})

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT pit_file_debt_ceiling_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates the header to checked if checked headers row already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)
			Expect(err).NotTo(HaveOccurred())

			err = pitFileDebtCeilingRepository.Create(headerID, []interface{}{test_data.PitFileDebtCeilingModel})

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT pit_file_debt_ceiling_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("does not duplicate pit file events", func() {
			err = pitFileDebtCeilingRepository.Create(headerID, []interface{}{test_data.PitFileDebtCeilingModel})
			Expect(err).NotTo(HaveOccurred())

			err = pitFileDebtCeilingRepository.Create(headerID, []interface{}{test_data.PitFileDebtCeilingModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes pit file if corresponding header is deleted", func() {
			err = pitFileDebtCeilingRepository.Create(headerID, []interface{}{test_data.PitFileDebtCeilingModel})
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbPitFile debt_ceiling.PitFileDebtCeilingModel
			err = db.Get(&dbPitFile, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.pit_file_debt_ceiling WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})

		It("returns an error if model is of wrong type", func() {
			err = pitFileDebtCeilingRepository.Create(headerID, []interface{}{test_data.WrongModel{}})

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
			err = pitFileDebtCeilingRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT pit_file_debt_ceiling_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)

			err = pitFileDebtCeilingRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT pit_file_debt_ceiling_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
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
			startingBlock = rand.Int63()
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

			err := pitFileDebtCeilingRepository.MarkHeaderChecked(headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := pitFileDebtCeilingRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
		})

		It("only treats headers as checked if pit file debt ceiling logs have been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := pitFileDebtCeilingRepository.MissingHeaders(startingBlock, endingBlock)

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
			pitFileDebtCeilingRepositoryTwo := debt_ceiling.PitFileDebtCeilingRepository{}
			pitFileDebtCeilingRepositoryTwo.SetDB(dbTwo)
			err := pitFileDebtCeilingRepository.MarkHeaderChecked(headerIDs[0])
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := pitFileDebtCeilingRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := pitFileDebtCeilingRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
