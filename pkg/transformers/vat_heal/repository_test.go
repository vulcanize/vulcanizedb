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

package vat_heal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_heal"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("VatHeal Repository", func() {
	var (
		db         *postgres.DB
		repository vat_heal.VatHealRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = vat_heal.VatHealRepository{}
		repository.SetDB(db)
	})

	type VatHealDBResult struct {
		vat_heal.VatHealModel
		Id       int
		HeaderId int64 `db:"header_id"`
	}

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.VatHealModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatHealChecked,
			LogEventTableName:        "maker.vat_heal",
			TestModel:                test_data.VatHealModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &repository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists vat heal records", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerId, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			anotherVatHeal := test_data.VatHealModel
			anotherVatHeal.LogIndex = test_data.VatHealModel.LogIndex + 1
			err = repository.Create(headerId, []interface{}{test_data.VatHealModel, anotherVatHeal})

			var dbResult []VatHealDBResult
			err = db.Select(&dbResult, `SELECT * from maker.vat_heal where header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbResult)).To(Equal(2))
			Expect(dbResult[0].Urn).To(Equal(test_data.VatHealModel.Urn))
			Expect(dbResult[0].V).To(Equal(test_data.VatHealModel.V))
			Expect(dbResult[0].Rad).To(Equal(test_data.VatHealModel.Rad))
			Expect(dbResult[0].LogIndex).To(Equal(test_data.VatHealModel.LogIndex))
			Expect(dbResult[1].LogIndex).To(Equal(test_data.VatHealModel.LogIndex + 1))
			Expect(dbResult[0].TransactionIndex).To(Equal(test_data.VatHealModel.TransactionIndex))
			Expect(dbResult[0].Raw).To(MatchJSON(test_data.VatHealModel.Raw))
			Expect(dbResult[0].HeaderId).To(Equal(headerId))
		})
	})

	Describe("MarkCheckedHeader", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatHealChecked,
			Repository:              &repository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
