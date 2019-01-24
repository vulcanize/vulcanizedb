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

package tend_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("TendRepository", func() {
	var (
		db               *postgres.DB
		tendRepository   tend.TendRepository
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		tendRepository = tend.TendRepository{}
		tendRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.TendModel
		modelWithDifferentLogIdx.LogIndex = modelWithDifferentLogIdx.LogIndex + 1
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.TendChecked,
			LogEventTableName:        "maker.tend",
			TestModel:                test_data.TendModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &tendRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a tend record", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = tendRepository.Create(headerID, []interface{}{test_data.TendModel})

			Expect(err).NotTo(HaveOccurred())
			var count int
			err = db.QueryRow(`SELECT count(*) from maker.tend`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			dbResult := tend.TendModel{}
			err = db.Get(&dbResult, `SELECT bid_id, lot, bid, guy, log_idx, tx_idx, raw_log FROM maker.tend WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbResult.BidId).To(Equal(test_data.TendModel.BidId))
			Expect(dbResult.Lot).To(Equal(test_data.TendModel.Lot))
			Expect(dbResult.Bid).To(Equal(test_data.TendModel.Bid))
			Expect(dbResult.Guy).To(Equal(test_data.TendModel.Guy))
			Expect(dbResult.LogIndex).To(Equal(test_data.TendModel.LogIndex))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.TendModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.TendModel.Raw))

			var dbTic int64
			err = db.Get(&dbTic, `SELECT tic FROM maker.tend WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbTic).To(Equal(fakes.FakeHeaderTic))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.TendChecked,
			Repository:              &tendRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
