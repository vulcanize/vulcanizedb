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

package ilk_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Pit file ilk repository", func() {
	var (
		db                   *postgres.DB
		pitFileIlkRepository ilk.PitFileIlkRepository
		headerRepository     repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		pitFileIlkRepository = ilk.PitFileIlkRepository{}
		pitFileIlkRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.PitFileIlkModel
		modelWithDifferentLogIdx.LogIndex = modelWithDifferentLogIdx.LogIndex + 1
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.PitFileIlkChecked,
			LogEventTableName:        "maker.pit_file_ilk",
			TestModel:                test_data.PitFileIlkModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &pitFileIlkRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a pit file ilk event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = pitFileIlkRepository.Create(headerID, []interface{}{test_data.PitFileIlkModel})

			Expect(err).NotTo(HaveOccurred())
			var dbPitFile ilk.PitFileIlkModel
			err = db.Get(&dbPitFile, `SELECT ilk, what, data, log_idx, tx_idx, raw_log FROM maker.pit_file_ilk WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbPitFile.Ilk).To(Equal(test_data.PitFileIlkModel.Ilk))
			Expect(dbPitFile.What).To(Equal(test_data.PitFileIlkModel.What))
			Expect(dbPitFile.Data).To(Equal(test_data.PitFileIlkModel.Data))
			Expect(dbPitFile.LogIndex).To(Equal(test_data.PitFileIlkModel.LogIndex))
			Expect(dbPitFile.TransactionIndex).To(Equal(test_data.PitFileIlkModel.TransactionIndex))
			Expect(dbPitFile.Raw).To(MatchJSON(test_data.PitFileIlkModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.PitFileIlkChecked,
			Repository:              &pitFileIlkRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &pitFileIlkRepository,
			RepositoryTwo: &ilk.PitFileIlkRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})
})
