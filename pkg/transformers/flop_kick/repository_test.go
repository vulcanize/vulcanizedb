// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
