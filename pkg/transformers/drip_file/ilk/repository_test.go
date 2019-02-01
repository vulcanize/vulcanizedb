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

package ilk_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Drip file ilk repository", func() {
	var (
		db                    *postgres.DB
		dripFileIlkRepository ilk.DripFileIlkRepository
		headerRepository      datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		dripFileIlkRepository = ilk.DripFileIlkRepository{}
		dripFileIlkRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.DripFileIlkModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.DripFileIlkChecked,
			LogEventTableName:        "maker.drip_file_ilk",
			TestModel:                test_data.DripFileIlkModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &dripFileIlkRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a drip file ilk event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = dripFileIlkRepository.Create(headerID, []interface{}{test_data.DripFileIlkModel})

			Expect(err).NotTo(HaveOccurred())
			var dbDripFileIlk ilk.DripFileIlkModel
			err = db.Get(&dbDripFileIlk, `SELECT ilk, vow, tax, log_idx, tx_idx, raw_log FROM maker.drip_file_ilk WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbDripFileIlk.Ilk).To(Equal(test_data.DripFileIlkModel.Ilk))
			Expect(dbDripFileIlk.Vow).To(Equal(test_data.DripFileIlkModel.Vow))
			Expect(dbDripFileIlk.Tax).To(Equal(test_data.DripFileIlkModel.Tax))
			Expect(dbDripFileIlk.LogIndex).To(Equal(test_data.DripFileIlkModel.LogIndex))
			Expect(dbDripFileIlk.TransactionIndex).To(Equal(test_data.DripFileIlkModel.TransactionIndex))
			Expect(dbDripFileIlk.Raw).To(MatchJSON(test_data.DripFileIlkModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.DripFileIlkChecked,
			Repository:              &dripFileIlkRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
