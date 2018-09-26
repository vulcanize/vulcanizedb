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

package flop_kick_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("FlopRepository", func() {
	var db *postgres.DB
	var repository flop_kick.FlopKickRepository
	var headerRepository repositories.HeaderRepository
	var headerId int64
	var err error
	var dbResult test_data.FlopKickDBResult

	BeforeEach(func() {
		node := test_config.NewTestNode()
		db = test_config.NewTestDB(node)
		test_config.CleanTestDB(db)
		repository = flop_kick.NewFlopKickRepository(db)
		headerRepository = repositories.NewHeaderRepository(db)
		headerId, err = headerRepository.CreateOrUpdateHeader(core.Header{})
		Expect(err).NotTo(HaveOccurred())
		dbResult = test_data.FlopKickDBResult{}
	})

	Describe("Create", func() {
		It("creates FlopKick records", func() {
			err := repository.Create(headerId, []flop_kick.Model{test_data.FlopKickModel})
			Expect(err).NotTo(HaveOccurred())

			err = db.QueryRowx(`SELECT * FROM maker.flop_kick WHERE header_id = $1`, headerId).StructScan(&dbResult)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.HeaderId).To(Equal(headerId))
			Expect(dbResult.BidId).To(Equal(test_data.FlopKickModel.BidId))
			Expect(dbResult.Lot).To(Equal(test_data.FlopKickModel.Lot))
			Expect(dbResult.Bid).To(Equal(test_data.FlopKickModel.Bid))
			Expect(dbResult.Gal).To(Equal(test_data.FlopKickModel.Gal))
			Expect(dbResult.End.Equal(test_data.FlopKickModel.End)).To(BeTrue())
			Expect(dbResult.TransactionIndex).To(Equal(test_data.FlopKickModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.FlopKickModel.Raw))
		})

		It("marks headerId as checked for flop kick logs", func() {
			err = repository.Create(headerId, []flop_kick.Model{test_data.FlopKickModel})
			Expect(err).NotTo(HaveOccurred())

			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT flop_kick_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("returns an error if inserting the flop_kick record fails", func() {
			err = repository.Create(headerId, []flop_kick.Model{test_data.FlopKickModel})
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(headerId, []flop_kick.Model{test_data.FlopKickModel})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("deletes the flop_kick records if its corresponding header record is deleted", func() {
			err = repository.Create(headerId, []flop_kick.Model{test_data.FlopKickModel})
			Expect(err).NotTo(HaveOccurred())

			var flopKickCount int
			err = db.QueryRow(`SELECT count(*) FROM maker.flop_kick`).Scan(&flopKickCount)
			Expect(err).NotTo(HaveOccurred())
			Expect(flopKickCount).To(Equal(1))

			_, err = db.Exec(`DELETE FROM headers where id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())

			err = db.QueryRow(`SELECT count(*) FROM maker.flop_kick`).Scan(&flopKickCount)
			Expect(err).NotTo(HaveOccurred())
			Expect(flopKickCount).To(Equal(0))
		})
	})

	Describe("MarkedHeadersChecked", func() {
		It("creates a row for a new headerId", func() {
			err := repository.MarkHeaderChecked(headerId)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT flop_kick_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerId already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerId)
			err = repository.MarkHeaderChecked(headerId)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT flop_kick_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		var flopKickBlockNumber = rand.Int63()
		var startingBlockNumber = flopKickBlockNumber - 1
		var endingBlockNumber = flopKickBlockNumber + 1
		var outOfRangeBlockNumber = flopKickBlockNumber + 2

		It("returns headers haven't been checked", func() {
			var headerIds []int64

			for _, number := range []int64{startingBlockNumber, flopKickBlockNumber, endingBlockNumber, outOfRangeBlockNumber} {
				headerId, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: number})
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)
			}

			err = repository.MarkHeaderChecked(headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := repository.MissingHeaders(startingBlockNumber, endingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
		})

		It("only treats headers as checked if flop_kicks have been checked", func() {
			var headerIds []int64
			for _, number := range []int64{startingBlockNumber, flopKickBlockNumber, endingBlockNumber, outOfRangeBlockNumber} {
				headerId, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: number})
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)
			}

			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := repository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(flopKickBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(flopKickBlockNumber)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(flopKickBlockNumber)))
		})

		It("only returns missing headers for the current node", func() {
			var headerIds []int64
			node2 := core.Node{}
			db2 := test_config.NewTestDB(node2)
			headerRepository2 := repositories.NewHeaderRepository(db2)
			flopKickRepository2 := flop_kick.NewFlopKickRepository(db2)

			for _, number := range []int64{startingBlockNumber, flopKickBlockNumber, endingBlockNumber} {
				headerId, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: number})
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)

				headerRepository2.CreateOrUpdateHeader(core.Header{BlockNumber: number})
				Expect(err).NotTo(HaveOccurred())
			}

			err = repository.MarkHeaderChecked(headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			node1MissingHeaders, err := repository.MissingHeaders(startingBlockNumber, endingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node1MissingHeaders)).To(Equal(2))

			node2MissingHeaders, err := flopKickRepository2.MissingHeaders(startingBlockNumber, endingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node2MissingHeaders)).To(Equal(3))
		})
	})
})
