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

package flap_kick_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flap_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Flap Kick Repository", func() {
	var (
		db                 *postgres.DB
		flapKickRepository flap_kick.FlapKickRepository
		headerRepository   repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		flapKickRepository = flap_kick.FlapKickRepository{}
		flapKickRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.FlapKickModel
		modelWithDifferentLogIdx.LogIndex = modelWithDifferentLogIdx.LogIndex + 1
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  "flap_kick_checked",
			LogEventTableName:        "maker.flap_kick",
			TestModel:                test_data.FlapKickModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &flapKickRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a flap kick record", func() {
			headerId, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = flapKickRepository.Create(headerId, []interface{}{test_data.FlapKickModel})

			Expect(err).NotTo(HaveOccurred())
			var count int
			db.QueryRow(`SELECT count(*) FROM maker.flap_kick`).Scan(&count)
			Expect(count).To(Equal(1))
			var dbResult flap_kick.FlapKickModel
			err = db.Get(&dbResult, `SELECT bid, bid_id, "end", gal, lot, log_idx, tx_idx, raw_log FROM maker.flap_kick WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Bid).To(Equal(test_data.FlapKickModel.Bid))
			Expect(dbResult.BidId).To(Equal(test_data.FlapKickModel.BidId))
			Expect(dbResult.End.Equal(test_data.FlapKickModel.End)).To(BeTrue())
			Expect(dbResult.Gal).To(Equal(test_data.FlapKickModel.Gal))
			Expect(dbResult.Lot).To(Equal(test_data.FlapKickModel.Lot))
			Expect(dbResult.LogIndex).To(Equal(test_data.FlapKickModel.LogIndex))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.FlapKickModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.FlapKickModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: "flap_kick_checked",
			Repository:              &flapKickRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &flapKickRepository,
			RepositoryTwo: &flap_kick.FlapKickRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})
})
