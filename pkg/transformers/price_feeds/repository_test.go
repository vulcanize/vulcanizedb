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

package price_feeds_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Price feeds repository", func() {
	var (
		db                  *postgres.DB
		priceFeedRepository price_feeds.PriceFeedRepository
		headerRepository    repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		priceFeedRepository = price_feeds.PriceFeedRepository{}
		priceFeedRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.PriceFeedModel
		modelWithDifferentLogIdx.LogIndex = modelWithDifferentLogIdx.LogIndex + 1
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.PriceFeedsChecked,
			LogEventTableName:        "maker.price_feeds",
			TestModel:                test_data.PriceFeedModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &priceFeedRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists a price feed update", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = priceFeedRepository.Create(headerID, []interface{}{test_data.PriceFeedModel})

			Expect(err).NotTo(HaveOccurred())
			var dbPriceFeedUpdate price_feeds.PriceFeedModel
			err = db.Get(&dbPriceFeedUpdate, `SELECT block_number, medianizer_address, usd_value, log_idx, tx_idx, raw_log FROM maker.price_feeds WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbPriceFeedUpdate.BlockNumber).To(Equal(test_data.PriceFeedModel.BlockNumber))
			Expect(dbPriceFeedUpdate.MedianizerAddress).To(Equal(test_data.PriceFeedModel.MedianizerAddress))
			Expect(dbPriceFeedUpdate.UsdValue).To(Equal(test_data.PriceFeedModel.UsdValue))
			Expect(dbPriceFeedUpdate.LogIndex).To(Equal(test_data.PriceFeedModel.LogIndex))
			Expect(dbPriceFeedUpdate.TransactionIndex).To(Equal(test_data.PriceFeedModel.TransactionIndex))
			Expect(dbPriceFeedUpdate.Raw).To(MatchJSON(test_data.PriceFeedModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.PriceFeedsChecked,
			Repository:              &priceFeedRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
