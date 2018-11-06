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

package chop_lump_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/chop_lump"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Cat file chop lump repository", func() {
	var (
		catFileRepository chop_lump.CatFileChopLumpRepository
		db                *postgres.DB
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		catFileRepository = chop_lump.CatFileChopLumpRepository{}
		catFileRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.CatFileChopLumpModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  "cat_file_chop_lump_checked",
			LogEventTableName:        "maker.cat_file_chop_lump",
			TestModel:                test_data.CatFileChopLumpModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &catFileRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a cat file chop lump event", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = catFileRepository.Create(headerID, []interface{}{test_data.CatFileChopLumpModel})

			Expect(err).NotTo(HaveOccurred())
			var dbResult chop_lump.CatFileChopLumpModel
			err = db.Get(&dbResult, `SELECT ilk, what, data, tx_idx, log_idx, raw_log FROM maker.cat_file_chop_lump WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(test_data.CatFileChopLumpModel.Ilk))
			Expect(dbResult.What).To(Equal(test_data.CatFileChopLumpModel.What))
			Expect(dbResult.Data).To(Equal(test_data.CatFileChopLumpModel.Data))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.CatFileChopLumpModel.TransactionIndex))
			Expect(dbResult.LogIndex).To(Equal(test_data.CatFileChopLumpModel.LogIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.CatFileChopLumpModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: "cat_file_chop_lump_checked",
			Repository:              &catFileRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &catFileRepository,
			RepositoryTwo: &chop_lump.CatFileChopLumpRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})
})
