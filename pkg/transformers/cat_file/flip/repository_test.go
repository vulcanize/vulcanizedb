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

package flip_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/flip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Cat file flip repository", func() {
	var (
		catFileFlipRepository flip.CatFileFlipRepository
		db                    *postgres.DB
		headerRepository      datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		catFileFlipRepository = flip.CatFileFlipRepository{}
		catFileFlipRepository.SetDB(db)
	})

	Describe("Create", func() {
		catFileFlipWithDifferentIdx := test_data.CatFileFlipModel
		catFileFlipWithDifferentIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.CatFileFlipChecked,
			LogEventTableName:        "maker.cat_file_flip",
			TestModel:                test_data.CatFileFlipModel,
			ModelWithDifferentLogIdx: catFileFlipWithDifferentIdx,
			Repository:               &catFileFlipRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a cat file flip event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = catFileFlipRepository.Create(headerID, []interface{}{test_data.CatFileFlipModel})

			Expect(err).NotTo(HaveOccurred())
			var dbResult flip.CatFileFlipModel
			err = db.Get(&dbResult, `SELECT ilk, what, flip, tx_idx, log_idx, raw_log FROM maker.cat_file_flip WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(test_data.CatFileFlipModel.Ilk))
			Expect(dbResult.What).To(Equal(test_data.CatFileFlipModel.What))
			Expect(dbResult.Flip).To(Equal(test_data.CatFileFlipModel.Flip))
			Expect(dbResult.TransactionIndex).To(Equal(test_data.CatFileFlipModel.TransactionIndex))
			Expect(dbResult.LogIndex).To(Equal(test_data.CatFileFlipModel.LogIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.CatFileFlipModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.CatFileFlipChecked,
			Repository:              &catFileFlipRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
