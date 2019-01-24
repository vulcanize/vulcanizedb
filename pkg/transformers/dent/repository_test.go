// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package dent_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Dent Repository", func() {
	var (
		db               *postgres.DB
		dentRepository   dent.DentRepository
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		dentRepository = dent.DentRepository{}
		dentRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.DentModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.DentChecked,
			LogEventTableName:        "maker.dent",
			TestModel:                test_data.DentModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &dentRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a dent record", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = dentRepository.Create(headerID, []interface{}{test_data.DentModel})
			Expect(err).NotTo(HaveOccurred())

			var count int
			db.QueryRow(`SELECT count(*) FROM maker.dent`).Scan(&count)
			Expect(count).To(Equal(1))

			var dbResult dent.DentModel
			err = db.Get(&dbResult, `SELECT bid_id, lot, bid, guy, log_idx, tx_idx, raw_log FROM maker.dent WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.BidId).To(Equal(test_data.DentModel.BidId))
			Expect(dbResult.Lot).To(Equal(test_data.DentModel.Lot))
			Expect(dbResult.Bid).To(Equal(test_data.DentModel.Bid))
			Expect(dbResult.Guy).To(Equal(test_data.DentModel.Guy))
			Expect(dbResult.LogIndex).To(Equal(test_data.DentModel.LogIndex))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.DentModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.DentModel.Raw))

			var dbTic int64
			err = db.Get(&dbTic, `SELECT tic FROM maker.dent WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbTic).To(Equal(fakes.FakeHeaderTic))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.DentChecked,
			Repository:              &dentRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
