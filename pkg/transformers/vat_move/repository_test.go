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

package vat_move_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_move"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat Move", func() {
	var db *postgres.DB
	var repository vat_move.VatMoveRepository

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = vat_move.VatMoveRepository{}
		repository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.VatMoveModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatMoveChecked,
			LogEventTableName:        "maker.vat_move",
			TestModel:                test_data.VatMoveModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &repository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists vat move records", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(headerID, []interface{}{test_data.VatMoveModel})

			Expect(err).NotTo(HaveOccurred())
			var dbVatMove vat_move.VatMoveModel
			err = db.Get(&dbVatMove, `SELECT src, dst, rad, log_idx, tx_idx, raw_log FROM maker.vat_move WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbVatMove.Src).To(Equal(test_data.VatMoveModel.Src))
			Expect(dbVatMove.Dst).To(Equal(test_data.VatMoveModel.Dst))
			Expect(dbVatMove.Rad).To(Equal(test_data.VatMoveModel.Rad))
			Expect(dbVatMove.LogIndex).To(Equal(test_data.VatMoveModel.LogIndex))
			Expect(dbVatMove.TransactionIndex).To(Equal(test_data.VatMoveModel.TransactionIndex))
			Expect(dbVatMove.Raw).To(MatchJSON(test_data.VatMoveModel.Raw))
		})
	})

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &repository,
			RepositoryTwo: &vat_move.VatMoveRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatMoveChecked,
			Repository:              &repository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
