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
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_flux"
	"github.com/vulcanize/vulcanizedb/test_config"
	"math/rand"
)

var _ = Describe("VatFlux Repository", func() {
	var (
		db               *postgres.DB
		repository       vat_flux.VatFluxRepository
		headerRepository repositories.HeaderRepository
		headerId         int64
		err              error
	)

	BeforeEach(func() {
		node := test_config.NewTestNode()
		db = test_config.NewTestDB(node)
		test_config.CleanTestDB(db)
		repository = vat_flux.VatFluxRepository{}
		repository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		headerId, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(err).NotTo(HaveOccurred())
	})

	type VatFluxDBResult struct {
		vat_flux.VatFluxModel
		Id       int
		HeaderId int64 `db:"header_id"`
	}

	type CheckedHeaderResult struct {
		VatFluxChecked bool `db:"vat_flux_checked"`
	}

	Describe("Create", func() {
		It("persists vat flux records", func() {
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

		It("returns an error if the insertion fails", func() {
			err = repository.Create(headerId, []interface{}{test_data.VatFluxModel})
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(headerId, []interface{}{test_data.VatFluxModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("allows for multiple vat flux events in one transaction if they have different log indexes", func() {
			err = repository.Create(headerId, []interface{}{test_data.VatFluxModel})
			Expect(err).NotTo(HaveOccurred())

			anotherVatFlux := test_data.VatFluxModel
			anotherVatFlux.LogIndex = anotherVatFlux.LogIndex + 1
			err = repository.Create(headerId, []interface{}{anotherVatFlux})

			Expect(err).NotTo(HaveOccurred())
		})

		It("marks the header as checked for vat flux logs", func() {
			err = repository.Create(headerId, []interface{}{test_data.VatFluxModel})

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_flux_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates the header to checked if checked headers row already exists", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerId)
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(headerId, []interface{}{test_data.VatFluxModel})

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_flux_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("removes vat flux if corresponding header is deleted", func() {
			err = repository.Create(headerId, []interface{}{test_data.VatFluxModel})
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerId)

			Expect(err).NotTo(HaveOccurred())
			var count int
			err = db.QueryRow(`SELECT count(*) from maker.vat_flux`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})

		It("wraps create in a transaction", func() {
			err = repository.Create(headerId, []interface{}{test_data.VatFluxModel, test_data.VatFluxModel})

			Expect(err).To(HaveOccurred())
			var count int
			err = db.QueryRowx(`SELECT count(*) FROM maker.vat_flux`).Scan(&count)
			Expect(count).To(Equal(0))
		})

		It("returns an error if model is of wrong type", func() {
			err = repository.Create(headerId, []interface{}{test_data.WrongModel{}})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("model of type"))
		})
	})

	Describe("MissingHeaders", func() {
		var (
			startingBlock, vatFluxBlock, endingBlock, outsideRangeBlock int64
			headerIds, blockNumbers                                     []int64
		)

		BeforeEach(func() {
			startingBlock = rand.Int63()
			vatFluxBlock = startingBlock + 1
			endingBlock = startingBlock + 2
			outsideRangeBlock = startingBlock + 3

			blockNumbers = []int64{startingBlock, vatFluxBlock, endingBlock, outsideRangeBlock}
			headerIds = []int64{}
			for _, n := range blockNumbers {
				headerId, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)
			}
		})

		It("returns headers that haven't been checked", func() {
			err = repository.MarkHeaderChecked(headerIds[0])
			Expect(err).NotTo(HaveOccurred())

			headers, err := repository.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(headers[0].Id).To(Or(Equal(headerIds[1]), Equal(headerIds[2])))
			Expect(headers[1].Id).To(Or(Equal(headerIds[1]), Equal(headerIds[2])))
			Expect(len(headers)).To(Equal(2))
		})

		It("returns header ids when checked_headers.vat_flux is false", func() {
			err = repository.MarkHeaderChecked(headerIds[0])
			_, err = db.Exec(`INSERT INTO checked_headers (header_id) VALUES ($1)`, headerIds[1])
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
			repository2 := vat_flux.VatFluxRepository{}
			repository2.SetDB(db2)

			for _, n := range blockNumbers {
				_, err = headerRepository2.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
			}

			err = repository.MarkHeaderChecked(headerIds[0])
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := repository.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(2))

			nodeTwoMissingHeaders, err := repository2.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(3))
		})
	})

	Describe("MarkHeaderChecked", func() {
		It("creates a new checked_header record", func() {
			err := repository.MarkHeaderChecked(headerId)
			Expect(err).NotTo(HaveOccurred())

			var checkedHeaderResult = CheckedHeaderResult{}
			err = db.Get(&checkedHeaderResult, `SELECT vat_flux_checked FROM checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(checkedHeaderResult.VatFluxChecked).To(BeTrue())
		})

		It("updates an existing checked header", func() {
			_, err := db.Exec(`INSERT INTO checked_headers (header_id) VALUES($1)`, headerId)
			Expect(err).NotTo(HaveOccurred())

			var checkedHeaderResult CheckedHeaderResult
			err = db.Get(&checkedHeaderResult, `SELECT vat_flux_checked FROM checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(checkedHeaderResult.VatFluxChecked).To(BeFalse())

			err = repository.MarkHeaderChecked(headerId)
			Expect(err).NotTo(HaveOccurred())

			err = db.Get(&checkedHeaderResult, `SELECT vat_flux_checked FROM checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(checkedHeaderResult.VatFluxChecked).To(BeTrue())
		})
	})
})
