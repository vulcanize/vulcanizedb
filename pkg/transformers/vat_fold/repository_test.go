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

package vat_fold_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_fold"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat.fold repository", func() {
	var (
		db         *postgres.DB
		repository vat_fold.VatFoldRepository
	)

	BeforeEach(func() {
		node := test_config.NewTestNode()
		db = test_config.NewTestDB(node)
		test_config.CleanTestDB(db)
		repository = vat_fold.VatFoldRepository{}
		repository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.VatFoldModel
		modelWithDifferentLogIdx.LogIndex++

		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatFoldChecked,
			LogEventTableName:        "maker.vat_fold",
			TestModel:                test_data.VatFoldModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &repository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a vat fold event", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = repository.Create(headerID, []interface{}{test_data.VatFoldModel})

			Expect(err).NotTo(HaveOccurred())
			var dbVatFold vat_fold.VatFoldModel
			err = db.Get(&dbVatFold, `SELECT ilk, urn, rate, log_idx, tx_idx, raw_log FROM maker.vat_fold WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbVatFold.Ilk).To(Equal(test_data.VatFoldModel.Ilk))
			Expect(dbVatFold.Urn).To(Equal(test_data.VatFoldModel.Urn))
			Expect(dbVatFold.Rate).To(Equal(test_data.VatFoldModel.Rate))
			Expect(dbVatFold.LogIndex).To(Equal(test_data.VatFoldModel.LogIndex))
			Expect(dbVatFold.TransactionIndex).To(Equal(test_data.VatFoldModel.TransactionIndex))
			Expect(dbVatFold.Raw).To(MatchJSON(test_data.VatFoldModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatFoldChecked,
			Repository:              &repository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &repository,
			RepositoryTwo: &vat_fold.VatFoldRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})
})
