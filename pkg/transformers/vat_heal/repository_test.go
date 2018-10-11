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

package vat_heal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_heal"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("VatHeal Repository", func() {
	var (
		db               *postgres.DB
		repository       vat_heal.VatHealRepository
		headerRepository repositories.HeaderRepository
		err              error
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = vat_heal.VatHealRepository{DB: db}
		headerRepository = repositories.NewHeaderRepository(db)
	})

	type VatHealDBResult struct {
		vat_heal.VatHealModel
		Id       int
		HeaderId int64 `db:"header_id"`
	}

	type CheckedHeaderResult struct {
		VatHealChecked bool `db:"vat_heal_checked"`
	}

	Describe("Create", func() {
		var headerId int64

		BeforeEach(func() {
			headerId, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("persists vat heal records", func() {
			anotherVatHeal := test_data.VatHealModel
			anotherVatHeal.TransactionIndex = test_data.VatHealModel.TransactionIndex + 1
			err = repository.Create(headerId, []vat_heal.VatHealModel{test_data.VatHealModel, anotherVatHeal})

			var dbResult []VatHealDBResult
			err = db.Select(&dbResult, `SELECT * from maker.vat_heal where header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbResult)).To(Equal(2))
			Expect(dbResult[0].Urn).To(Equal(test_data.VatHealModel.Urn))
			Expect(dbResult[0].V).To(Equal(test_data.VatHealModel.V))
			Expect(dbResult[0].Rad).To(Equal(test_data.VatHealModel.Rad))
			Expect(dbResult[0].TransactionIndex).To(Equal(test_data.VatHealModel.TransactionIndex))
			Expect(dbResult[1].TransactionIndex).To(Equal(test_data.VatHealModel.TransactionIndex + 1))
			Expect(dbResult[0].Raw).To(MatchJSON(test_data.VatHealModel.Raw))
			Expect(dbResult[0].HeaderId).To(Equal(headerId))
		})

		It("returns an error if the insertion fails", func() {
			err = repository.Create(headerId, []vat_heal.VatHealModel{test_data.VatHealModel})
			Expect(err).NotTo(HaveOccurred())
			err = repository.Create(headerId, []vat_heal.VatHealModel{test_data.VatHealModel})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("marks the header as checked for vat heal logs", func() {
			err = repository.Create(headerId, []vat_heal.VatHealModel{test_data.VatHealModel})
			Expect(err).NotTo(HaveOccurred())

			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_heal_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("removes vat heal if corresponding header is deleted", func() {
			err = repository.Create(headerId, []vat_heal.VatHealModel{test_data.VatHealModel})
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerId)

			Expect(err).NotTo(HaveOccurred())
			var count int
			err = db.QueryRow(`SELECT count(*) from maker.vat_heal`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})

		It("wraps create in a transaction", func() {
			err = repository.Create(headerId, []vat_heal.VatHealModel{test_data.VatHealModel, test_data.VatHealModel})
			Expect(err).To(HaveOccurred())
			var count int
			err = repository.DB.QueryRowx(`SELECT count(*) FROM maker.vat_heal`).Scan(&count)
			Expect(count).To(Equal(0))
		})
	})

	Describe("MissingHeaders", func() {
		var (
			startingBlock, vatHealBlock, endingBlock, outsideRangeBlock int64
			blockNumbers, headerIds                                     []int64
		)

		BeforeEach(func() {
			startingBlock = GinkgoRandomSeed()
			vatHealBlock = startingBlock + 1
			endingBlock = startingBlock + 2
			outsideRangeBlock = startingBlock + 3

			headerIds = []int64{}
			blockNumbers = []int64{startingBlock, vatHealBlock, endingBlock, outsideRangeBlock}
			for _, n := range blockNumbers {
				headerId, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)
			}
		})

		It("returns headers that haven't been checked", func() {
			err = repository.MarkCheckedHeader(headerIds[0])
			Expect(err).NotTo(HaveOccurred())

			headers, err := repository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(headers[0].Id).To(Or(Equal(headerIds[1]), Equal(headerIds[2])))
			Expect(headers[1].Id).To(Or(Equal(headerIds[1]), Equal(headerIds[2])))
			Expect(len(headers)).To(Equal(2))
		})

		It("returns header ids when checked_headers.vat_heal is false", func() {
			err = repository.MarkCheckedHeader(headerIds[0])
			_, err = repository.DB.Exec(`INSERT INTO checked_headers (header_id) VALUES ($1)`, headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := repository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(headers[0].Id).To(Or(Equal(headerIds[1]), Equal(headerIds[2])))
			Expect(headers[1].Id).To(Or(Equal(headerIds[1]), Equal(headerIds[2])))
			Expect(len(headers)).To(Equal(2))
		})

		It("only returns header ids for the current node", func() {
			db2 := test_config.NewTestDB(core.Node{ID: "second node"})
			headerRepository2 := repositories.NewHeaderRepository(db2)
			repository2 := vat_heal.NewVatHealRepository(db2)

			for _, n := range blockNumbers {
				_, err = headerRepository2.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
			}

			err = repository.MarkCheckedHeader(headerIds[0])
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := repository.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(2))

			nodeTwoMissingHeaders, err := repository2.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(3))
		})
	})

	Describe("MarkCheckedHeader", func() {
		var headerId int64

		BeforeEach(func() {
			headerId, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a new checked_header record", func() {
			err := repository.MarkCheckedHeader(headerId)
			Expect(err).NotTo(HaveOccurred())

			var checkedHeaderResult = CheckedHeaderResult{}
			err = db.Get(&checkedHeaderResult, `SELECT vat_heal_checked FROM checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(checkedHeaderResult.VatHealChecked).To(BeTrue())
		})

		It("updates an existing checked header", func() {
			_, err := repository.DB.Exec(`INSERT INTO checked_headers (header_id) VALUES($1)`, headerId)
			Expect(err).NotTo(HaveOccurred())

			var checkedHeaderResult CheckedHeaderResult
			err = db.Get(&checkedHeaderResult, `SELECT vat_heal_checked FROM checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(checkedHeaderResult.VatHealChecked).To(BeFalse())

			err = repository.MarkCheckedHeader(headerId)
			Expect(err).NotTo(HaveOccurred())

			err = db.Get(&checkedHeaderResult, `SELECT vat_heal_checked FROM checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(checkedHeaderResult.VatHealChecked).To(BeTrue())
		})
	})
})
