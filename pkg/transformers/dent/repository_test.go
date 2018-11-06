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

package dent_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Dent Repository", func() {
	var (
		db               *postgres.DB
		dentRepository   dent.DentRepository
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		dentRepository = dent.DentRepository{}
		dentRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.DentModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  "dent_checked",
			LogEventTableName:        "maker.dent",
			TestModel:                test_data.DentModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &dentRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a dent record", func() {
			headerId, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = dentRepository.Create(headerId, []interface{}{test_data.DentModel})
			Expect(err).NotTo(HaveOccurred())

			var count int
			db.QueryRow(`SELECT count(*) FROM maker.dent`).Scan(&count)
			Expect(count).To(Equal(1))

			var dbResult dent.DentModel
			err = db.Get(&dbResult, `SELECT bid_id, lot, bid, guy, tic, log_idx, tx_idx, raw_log FROM maker.dent WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.BidId).To(Equal(test_data.DentModel.BidId))
			Expect(dbResult.Lot).To(Equal(test_data.DentModel.Lot))
			Expect(dbResult.Bid).To(Equal(test_data.DentModel.Bid))
			Expect(dbResult.Guy).To(Equal(test_data.DentModel.Guy))
			Expect(dbResult.Tic).To(Equal(test_data.DentModel.Tic))
			Expect(dbResult.LogIndex).To(Equal(test_data.DentModel.LogIndex))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.DentModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.DentModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: "dent_checked",
			Repository:              &dentRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &dentRepository,
			RepositoryTwo: &dent.DentRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})
})
