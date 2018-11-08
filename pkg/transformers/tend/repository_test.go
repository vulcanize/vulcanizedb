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

package tend_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("TendRepository", func() {
	var (
		db               *postgres.DB
		tendRepository   tend.TendRepository
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		tendRepository = tend.TendRepository{}
		tendRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.TendModel
		modelWithDifferentLogIdx.LogIndex = modelWithDifferentLogIdx.LogIndex + 1
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.TendChecked,
			LogEventTableName:        "maker.tend",
			TestModel:                test_data.TendModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &tendRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a tend record", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = tendRepository.Create(headerID, []interface{}{test_data.TendModel})

			Expect(err).NotTo(HaveOccurred())
			var count int
			err = db.QueryRow(`SELECT count(*) from maker.tend`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			dbResult := tend.TendModel{}
			err = db.Get(&dbResult, `SELECT bid_id, lot, bid, guy, tic, log_idx, tx_idx, raw_log FROM maker.tend WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbResult.BidId).To(Equal(test_data.TendModel.BidId))
			Expect(dbResult.Lot).To(Equal(test_data.TendModel.Lot))
			Expect(dbResult.Bid).To(Equal(test_data.TendModel.Bid))
			Expect(dbResult.Guy).To(Equal(test_data.TendModel.Guy))
			Expect(dbResult.Tic).To(Equal(test_data.TendModel.Tic))
			Expect(dbResult.LogIndex).To(Equal(test_data.TendModel.LogIndex))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.TendModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.RawLogNoteJson))
		})

	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.TendChecked,
			Repository:              &tendRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &tendRepository,
			RepositoryTwo: &tend.TendRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})
})
