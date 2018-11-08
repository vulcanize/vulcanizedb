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

package stability_fee_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Pit file stability fee repository", func() {
	var (
		db                            *postgres.DB
		pitFileStabilityFeeRepository stability_fee.PitFileStabilityFeeRepository
		headerRepository              repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		pitFileStabilityFeeRepository = stability_fee.PitFileStabilityFeeRepository{}
		pitFileStabilityFeeRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.PitFileStabilityFeeModel
		modelWithDifferentLogIdx.LogIndex = modelWithDifferentLogIdx.LogIndex + 1
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.PitFileStabilityFeeChecked,
			LogEventTableName:        "maker.pit_file_stability_fee",
			TestModel:                test_data.PitFileStabilityFeeModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &pitFileStabilityFeeRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a pit file stability fee event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = pitFileStabilityFeeRepository.Create(headerID, []interface{}{test_data.PitFileStabilityFeeModel})

			Expect(err).NotTo(HaveOccurred())
			var dbPitFile stability_fee.PitFileStabilityFeeModel
			err = db.Get(&dbPitFile, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.pit_file_stability_fee WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbPitFile.What).To(Equal(test_data.PitFileStabilityFeeModel.What))
			Expect(dbPitFile.Data).To(Equal(test_data.PitFileStabilityFeeModel.Data))
			Expect(dbPitFile.LogIndex).To(Equal(test_data.PitFileStabilityFeeModel.LogIndex))
			Expect(dbPitFile.TransactionIndex).To(Equal(test_data.PitFileStabilityFeeModel.TransactionIndex))
			Expect(dbPitFile.Raw).To(MatchJSON(test_data.PitFileStabilityFeeModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.PitFileStabilityFeeChecked,
			Repository:              &pitFileStabilityFeeRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &pitFileStabilityFeeRepository,
			RepositoryTwo: &stability_fee.PitFileStabilityFeeRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})
})
