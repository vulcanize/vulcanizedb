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

package deal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/deal"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Deal Repository", func() {
	var (
		db               *postgres.DB
		dealRepository   deal.DealRepository
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		dealRepository = deal.DealRepository{}
		dealRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.DealModel
		modelWithDifferentLogIdx.LogIndex = modelWithDifferentLogIdx.LogIndex + 1
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.DealChecked,
			LogEventTableName:        "maker.deal",
			TestModel:                test_data.DealModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &dealRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a deal record", func() {
			headerId, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = dealRepository.Create(headerId, []interface{}{test_data.DealModel})

			Expect(err).NotTo(HaveOccurred())
			var count int
			db.QueryRow(`SELECT count(*) FROM maker.deal`).Scan(&count)
			Expect(count).To(Equal(1))
			var dbResult deal.DealModel
			err = db.Get(&dbResult, `SELECT bid_id, contract_address, log_idx, tx_idx, raw_log FROM maker.deal WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.BidId).To(Equal(test_data.DealModel.BidId))
			Expect(dbResult.ContractAddress).To(Equal(test_data.DealModel.ContractAddress))
			Expect(dbResult.LogIndex).To(Equal(test_data.DealModel.LogIndex))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.DealModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.DealModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.DealChecked,
			Repository:              &dealRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &dealRepository,
			RepositoryTwo: &deal.DealRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})
})
