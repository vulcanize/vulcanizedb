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

package pit_vow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/pit_vow"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Cat file pit vow repository", func() {
	var (
		catFilePitVowRepository pit_vow.CatFilePitVowRepository
		db                      *postgres.DB
		headerRepository        datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		catFilePitVowRepository = pit_vow.CatFilePitVowRepository{}
		catFilePitVowRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.CatFilePitVowModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.CatFilePitVowChecked,
			LogEventTableName:        "maker.cat_file_pit_vow",
			TestModel:                test_data.CatFilePitVowModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &catFilePitVowRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a cat file pit vow event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = catFilePitVowRepository.Create(headerID, []interface{}{test_data.CatFilePitVowModel})

			Expect(err).NotTo(HaveOccurred())
			var dbResult pit_vow.CatFilePitVowModel
			err = db.Get(&dbResult, `SELECT what, data, tx_idx, log_idx, raw_log FROM maker.cat_file_pit_vow WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.What).To(Equal(test_data.CatFilePitVowModel.What))
			Expect(dbResult.Data).To(Equal(test_data.CatFilePitVowModel.Data))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.CatFilePitVowModel.TransactionIndex))
			Expect(dbResult.LogIndex).To(Equal(test_data.CatFilePitVowModel.LogIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.CatFilePitVowModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.CatFilePitVowChecked,
			Repository:              &catFilePitVowRepository,
		}
		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
