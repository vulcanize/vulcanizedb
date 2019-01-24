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

package vow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/vow"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Drip file vow repository", func() {
	var (
		db                    *postgres.DB
		dripFileVowRepository vow.DripFileVowRepository
		headerRepository      datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		dripFileVowRepository = vow.DripFileVowRepository{}
		dripFileVowRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.DripFileVowModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.DripFileVowChecked,
			LogEventTableName:        "maker.drip_file_vow",
			TestModel:                test_data.DripFileVowModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &dripFileVowRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a drip file vow event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = dripFileVowRepository.Create(headerID, []interface{}{test_data.DripFileVowModel})

			Expect(err).NotTo(HaveOccurred())
			var dbDripFileVow vow.DripFileVowModel
			err = db.Get(&dbDripFileVow, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.drip_file_vow WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbDripFileVow.What).To(Equal(test_data.DripFileVowModel.What))
			Expect(dbDripFileVow.Data).To(Equal(test_data.DripFileVowModel.Data))
			Expect(dbDripFileVow.LogIndex).To(Equal(test_data.DripFileVowModel.LogIndex))
			Expect(dbDripFileVow.TransactionIndex).To(Equal(test_data.DripFileVowModel.TransactionIndex))
			Expect(dbDripFileVow.Raw).To(MatchJSON(test_data.DripFileVowModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.DripFileVowChecked,
			Repository:              &dripFileVowRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
