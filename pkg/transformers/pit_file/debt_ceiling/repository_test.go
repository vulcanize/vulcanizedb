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

package debt_ceiling_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/debt_ceiling"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Pit file debt ceiling repository", func() {
	var (
		db                           *postgres.DB
		pitFileDebtCeilingRepository debt_ceiling.PitFileDebtCeilingRepository
		headerRepository             repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		pitFileDebtCeilingRepository = debt_ceiling.PitFileDebtCeilingRepository{}
		pitFileDebtCeilingRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.PitFileDebtCeilingModel
		modelWithDifferentLogIdx.LogIndex = modelWithDifferentLogIdx.LogIndex + 1
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.PitFileDebtCeilingChecked,
			LogEventTableName:        "maker.pit_file_debt_ceiling",
			TestModel:                test_data.PitFileDebtCeilingModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &pitFileDebtCeilingRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a pit file debt ceiling event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
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
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.PitFileDebtCeilingChecked,
			Repository:              &pitFileDebtCeilingRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
