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

package vat_flux_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"strconv"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_flux"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("VatFlux Repository", func() {
	var (
		db         *postgres.DB
		repository vat_flux.VatFluxRepository
	)

	BeforeEach(func() {
		node := test_config.NewTestNode()
		db = test_config.NewTestDB(node)
		test_config.CleanTestDB(db)
		repository = vat_flux.VatFluxRepository{}
		repository.SetDB(db)
	})

	type VatFluxDBResult struct {
		vat_flux.VatFluxModel
		Id       int
		HeaderId int64 `db:"header_id"`
	}

	Describe("Create", func() {
		vatFluxWithDifferentLogIdx := test_data.VatFluxModel
		vatFluxWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatFluxChecked,
			LogEventTableName:        "maker.vat_flux",
			TestModel:                test_data.VatFluxModel,
			ModelWithDifferentLogIdx: vatFluxWithDifferentLogIdx,
			Repository:               &repository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists vat flux records", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerId, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			anotherVatFlux := test_data.VatFluxModel
			anotherVatFlux.TransactionIndex = test_data.VatFluxModel.TransactionIndex + 1
			err = repository.Create(headerId, []interface{}{test_data.VatFluxModel, anotherVatFlux})

			var dbResult []VatFluxDBResult
			err = db.Select(&dbResult, `SELECT * from maker.vat_flux where header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbResult)).To(Equal(2))
			ilkID, err := shared.GetOrCreateIlk(test_data.VatFluxModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
			Expect(dbResult[0].Dst).To(Equal(test_data.VatFluxModel.Dst))
			Expect(dbResult[0].Src).To(Equal(test_data.VatFluxModel.Src))
			Expect(dbResult[0].Rad).To(Equal(test_data.VatFluxModel.Rad))
			Expect(dbResult[0].TransactionIndex).To(Equal(test_data.VatFluxModel.TransactionIndex))
			Expect(dbResult[1].TransactionIndex).To(Equal(test_data.VatFluxModel.TransactionIndex + 1))
			Expect(dbResult[0].LogIndex).To(Equal(test_data.VatFluxModel.LogIndex))
			Expect(dbResult[0].Raw).To(MatchJSON(test_data.VatFluxModel.Raw))
			Expect(dbResult[0].HeaderId).To(Equal(headerId))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatFluxChecked,
			Repository:              &repository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
