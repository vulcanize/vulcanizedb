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
