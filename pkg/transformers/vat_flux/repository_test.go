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

package vat_flux_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

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
			Expect(dbResult[0].Ilk).To(Equal(test_data.VatFluxModel.Ilk))
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

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &repository,
			RepositoryTwo: &vat_flux.VatFluxRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatFluxChecked,
			Repository:              &repository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
