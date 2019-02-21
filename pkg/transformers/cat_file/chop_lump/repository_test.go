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

package chop_lump_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/chop_lump"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
	"strconv"
)

var _ = Describe("Cat file chop lump repository", func() {
	var (
		catFileRepository chop_lump.CatFileChopLumpRepository
		db                *postgres.DB
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		catFileRepository = chop_lump.CatFileChopLumpRepository{}
		catFileRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.CatFileChopModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.CatFileChopLumpChecked,
			LogEventTableName:        "maker.cat_file_chop_lump",
			TestModel:                test_data.CatFileChopModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &catFileRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a cat file chop event", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = catFileRepository.Create(headerID, []interface{}{test_data.CatFileChopModel})

			Expect(err).NotTo(HaveOccurred())
			var dbResult chop_lump.CatFileChopLumpModel
			err = db.Get(&dbResult, `SELECT ilk, what, data, tx_idx, log_idx, raw_log FROM maker.cat_file_chop_lump WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_data.CatFileChopModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(strconv.Itoa(ilkID)))
			Expect(dbResult.What).To(Equal(test_data.CatFileChopModel.What))
			Expect(dbResult.Data).To(Equal(test_data.CatFileChopModel.Data))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.CatFileChopModel.TransactionIndex))
			Expect(dbResult.LogIndex).To(Equal(test_data.CatFileChopModel.LogIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.CatFileChopModel.Raw))
		})

		It("adds a cat file lump event", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = catFileRepository.Create(headerID, []interface{}{test_data.CatFileLumpModel})

			Expect(err).NotTo(HaveOccurred())
			var dbResult chop_lump.CatFileChopLumpModel
			err = db.Get(&dbResult, `SELECT ilk, what, data, tx_idx, log_idx, raw_log FROM maker.cat_file_chop_lump WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_data.CatFileLumpModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(strconv.Itoa(ilkID)))
			Expect(dbResult.What).To(Equal(test_data.CatFileLumpModel.What))
			Expect(dbResult.Data).To(Equal(test_data.CatFileLumpModel.Data))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.CatFileLumpModel.TransactionIndex))
			Expect(dbResult.LogIndex).To(Equal(test_data.CatFileLumpModel.LogIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.CatFileLumpModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.CatFileChopLumpChecked,
			Repository:              &catFileRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
