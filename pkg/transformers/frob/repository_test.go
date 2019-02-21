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

package frob_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"strconv"

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
			ilkID, err := shared.GetOrCreateIlk(test_data.FrobModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbFrob.Ilk).To(Equal(strconv.Itoa(ilkID)))
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
})
