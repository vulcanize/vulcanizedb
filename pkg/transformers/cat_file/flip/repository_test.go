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

package flip_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/flip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Cat file flip repository", func() {
	var (
		catFileFlipRepository flip.CatFileFlipRepository
		db                    *postgres.DB
		headerRepository      datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		catFileFlipRepository = flip.CatFileFlipRepository{}
		catFileFlipRepository.SetDB(db)
	})

	Describe("Create", func() {
		catFileFlipWithDifferentIdx := test_data.CatFileFlipModel
		catFileFlipWithDifferentIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.CatFileFlipChecked,
			LogEventTableName:        "maker.cat_file_flip",
			TestModel:                test_data.CatFileFlipModel,
			ModelWithDifferentLogIdx: catFileFlipWithDifferentIdx,
			Repository:               &catFileFlipRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a cat file flip event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = catFileFlipRepository.Create(headerID, []interface{}{test_data.CatFileFlipModel})

			Expect(err).NotTo(HaveOccurred())
			var dbResult flip.CatFileFlipModel
			err = db.Get(&dbResult, `SELECT ilk, what, flip, tx_idx, log_idx, raw_log FROM maker.cat_file_flip WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(test_data.CatFileFlipModel.Ilk))
			Expect(dbResult.What).To(Equal(test_data.CatFileFlipModel.What))
			Expect(dbResult.Flip).To(Equal(test_data.CatFileFlipModel.Flip))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.CatFileFlipModel.TransactionIndex))
			Expect(dbResult.LogIndex).To(Equal(test_data.CatFileFlipModel.LogIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.CatFileFlipModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.CatFileFlipChecked,
			Repository:              &catFileFlipRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
