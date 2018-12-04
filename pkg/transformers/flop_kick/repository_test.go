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

package flop_kick_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("FlopRepository", func() {
	var (
		db         *postgres.DB
		repository flop_kick.FlopKickRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = flop_kick.FlopKickRepository{}
		repository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.FlopKickModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.FlopKickChecked,
			LogEventTableName:        "maker.flop_kick",
			TestModel:                test_data.FlopKickModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &repository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("creates FlopKick records", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerId, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(headerId, []interface{}{test_data.FlopKickModel})
			Expect(err).NotTo(HaveOccurred())

			dbResult := test_data.FlopKickDBResult{}
			err = db.QueryRowx(`SELECT * FROM maker.flop_kick WHERE header_id = $1`, headerId).StructScan(&dbResult)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.HeaderId).To(Equal(headerId))
			Expect(dbResult.BidId).To(Equal(test_data.FlopKickModel.BidId))
			Expect(dbResult.Lot).To(Equal(test_data.FlopKickModel.Lot))
			Expect(dbResult.Bid).To(Equal(test_data.FlopKickModel.Bid))
			Expect(dbResult.Gal).To(Equal(test_data.FlopKickModel.Gal))
			Expect(dbResult.End.Equal(test_data.FlopKickModel.End)).To(BeTrue())
			Expect(dbResult.TransactionIndex).To(Equal(test_data.FlopKickModel.TransactionIndex))
			Expect(dbResult.LogIndex).To(Equal(test_data.FlopKickModel.LogIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.FlopKickModel.Raw))
		})
	})

	Describe("MarkedHeadersChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.FlopKickChecked,
			Repository:              &repository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
