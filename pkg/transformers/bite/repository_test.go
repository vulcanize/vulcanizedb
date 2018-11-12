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

package bite_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/bite"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Bite repository", func() {
	var (
		biteRepository bite.BiteRepository
		db             *postgres.DB
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		biteRepository = bite.BiteRepository{}
		biteRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.BiteModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.BiteChecked,
			LogEventTableName:        "maker.bite",
			TestModel:                test_data.BiteModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &biteRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a bite record", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = biteRepository.Create(headerID, []interface{}{test_data.BiteModel})

			Expect(err).NotTo(HaveOccurred())
			var dbBite bite.BiteModel
			err = db.Get(&dbBite, `SELECT ilk, urn, ink, art, tab, nflip, iart, log_idx, tx_idx, raw_log FROM maker.bite WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbBite.Ilk).To(Equal(test_data.BiteModel.Ilk))
			Expect(dbBite.Urn).To(Equal(test_data.BiteModel.Urn))
			Expect(dbBite.Ink).To(Equal(test_data.BiteModel.Ink))
			Expect(dbBite.Art).To(Equal(test_data.BiteModel.Art))
			Expect(dbBite.Tab).To(Equal(test_data.BiteModel.Tab))
			Expect(dbBite.NFlip).To(Equal(test_data.BiteModel.NFlip))
			Expect(dbBite.IArt).To(Equal(test_data.BiteModel.IArt))
			Expect(dbBite.LogIndex).To(Equal(test_data.BiteModel.LogIndex))
			Expect(dbBite.TransactionIndex).To(Equal(test_data.BiteModel.TransactionIndex))
			Expect(dbBite.Raw).To(MatchJSON(test_data.BiteModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.BiteChecked,
			Repository:              &biteRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})

	Describe("MissingHeaders", func() {

		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &biteRepository,
			RepositoryTwo: &bite.BiteRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})
})
