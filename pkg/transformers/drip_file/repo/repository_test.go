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

package repo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/repo"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Drip file repo repository", func() {
	var (
		db                     *postgres.DB
		dripFileRepoRepository repo.DripFileRepoRepository
		headerRepository       datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		dripFileRepoRepository = repo.DripFileRepoRepository{}
		dripFileRepoRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.DripFileRepoModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.DripFileRepoChecked,
			LogEventTableName:        "maker.drip_file_repo",
			TestModel:                test_data.DripFileRepoModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &dripFileRepoRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a drip file repo event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = dripFileRepoRepository.Create(headerID, []interface{}{test_data.DripFileRepoModel})

			Expect(err).NotTo(HaveOccurred())
			var dbDripFileRepo repo.DripFileRepoModel
			err = db.Get(&dbDripFileRepo, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.drip_file_repo WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbDripFileRepo.What).To(Equal(test_data.DripFileRepoModel.What))
			Expect(dbDripFileRepo.Data).To(Equal(test_data.DripFileRepoModel.Data))
			Expect(dbDripFileRepo.LogIndex).To(Equal(test_data.DripFileRepoModel.LogIndex))
			Expect(dbDripFileRepo.TransactionIndex).To(Equal(test_data.DripFileRepoModel.TransactionIndex))
			Expect(dbDripFileRepo.Raw).To(MatchJSON(test_data.DripFileRepoModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.DripFileRepoChecked,
			Repository:              &dripFileRepoRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
