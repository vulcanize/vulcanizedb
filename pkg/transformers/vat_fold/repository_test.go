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

package vat_fold_test

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

			ilkID, err := shared.GetOrCreateIlk(test_data.VatFoldModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbVatFold.Ilk).To(Equal(strconv.Itoa(ilkID)))
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
})
