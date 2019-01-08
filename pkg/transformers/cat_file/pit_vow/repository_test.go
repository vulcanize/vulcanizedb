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
