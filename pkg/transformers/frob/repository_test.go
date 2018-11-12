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

package frob_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Frob repository", func() {
	var (
		db             *postgres.DB
		frobRepository frob.FrobRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		frobRepository = frob.FrobRepository{}
		frobRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.FrobModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.FrobChecked,
			LogEventTableName:        "maker.frob",
			TestModel:                test_data.FrobModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &frobRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a frob", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = frobRepository.Create(headerID, []interface{}{test_data.FrobModel})
			Expect(err).NotTo(HaveOccurred())
			var dbFrob frob.FrobModel
			err = db.Get(&dbFrob, `SELECT art, dart, dink, iart, ilk, ink, urn, log_idx, tx_idx, raw_log FROM maker.frob WHERE header_id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			Expect(dbFrob.Ilk).To(Equal(test_data.FrobModel.Ilk))
			Expect(dbFrob.Urn).To(Equal(test_data.FrobModel.Urn))
			Expect(dbFrob.Ink).To(Equal(test_data.FrobModel.Ink))
			Expect(dbFrob.Art).To(Equal(test_data.FrobModel.Art))
			Expect(dbFrob.Dink).To(Equal(test_data.FrobModel.Dink))
			Expect(dbFrob.Dart).To(Equal(test_data.FrobModel.Dart))
			Expect(dbFrob.IArt).To(Equal(test_data.FrobModel.IArt))
			Expect(dbFrob.LogIndex).To(Equal(test_data.FrobModel.LogIndex))
			Expect(dbFrob.TransactionIndex).To(Equal(test_data.FrobModel.TransactionIndex))
			Expect(dbFrob.Raw).To(MatchJSON(test_data.FrobModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.FrobChecked,
			Repository:              &frobRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &frobRepository,
			RepositoryTwo: &frob.FrobRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})
})
