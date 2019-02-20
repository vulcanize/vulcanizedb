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

package drip_drip_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_drip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
	"strconv"
)

var _ = Describe("Drip drip repository", func() {
	var (
		db                 *postgres.DB
		dripDripRepository drip_drip.DripDripRepository
		headerRepository   datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		dripDripRepository = drip_drip.DripDripRepository{}
		dripDripRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.DripDripModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.DripDripChecked,
			LogEventTableName:        "maker.drip_drip",
			TestModel:                test_data.DripDripModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &dripDripRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a drip drip event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = dripDripRepository.Create(headerID, []interface{}{test_data.DripDripModel})

			Expect(err).NotTo(HaveOccurred())
			var dbDripDrip drip_drip.DripDripModel
			err = db.Get(&dbDripDrip, `SELECT ilk, log_idx, tx_idx, raw_log FROM maker.drip_drip WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_data.DripDripModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbDripDrip.Ilk).To(Equal(strconv.Itoa(ilkID)))
			Expect(dbDripDrip.LogIndex).To(Equal(test_data.DripDripModel.LogIndex))
			Expect(dbDripDrip.TransactionIndex).To(Equal(test_data.DripDripModel.TransactionIndex))
			Expect(dbDripDrip.Raw).To(MatchJSON(test_data.DripDripModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.DripDripChecked,
			Repository:              &dripDripRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
