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
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"
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
	var (
		db               *postgres.DB
		repository       flop_kick.FlopKickRepository
		headerRepository repositories.HeaderRepository
		err              error
		dbResult         test_data.FlopKickDBResult
		rawHeader        []byte
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = flop_kick.NewFlopKickRepository(db)
		headerRepository = repositories.NewHeaderRepository(db)
		dbResult = test_data.FlopKickDBResult{}
		rawHeader, err = json.Marshal(types.Header{})
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Create", func() {
		var headerId int64

		BeforeEach(func() {
			headerId, err = headerRepository.CreateOrUpdateHeader(core.Header{Raw: rawHeader})
			Expect(err).NotTo(HaveOccurred())

			err := repository.Create(headerId, []flop_kick.Model{test_data.FlopKickModel})
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates FlopKick records", func() {
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
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT flop_kick_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("returns an error if inserting the flop_kick record fails", func() {
			err = repository.Create(headerId, []flop_kick.Model{test_data.FlopKickModel})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("deletes the flop_kick records if its corresponding header record is deleted", func() {
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
		var headerId int64

		BeforeEach(func() {
			headerId, err = headerRepository.CreateOrUpdateHeader(core.Header{Raw: rawHeader})
			Expect(err).NotTo(HaveOccurred())
		})

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
		var (
			flopKickBlock, startingBlock, endingBlock, outOfRangeBlock int64
			headerIds                                                  []int64
		)

		BeforeEach(func() {
			flopKickBlock = GinkgoRandomSeed()
			startingBlock = flopKickBlock - 1
			endingBlock = flopKickBlock + 1
			outOfRangeBlock = flopKickBlock + 2

			headerIds = []int64{}
			for _, number := range []int64{startingBlock, flopKickBlock, endingBlock, outOfRangeBlock} {
				headerId, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: number, Raw: rawHeader})
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)
			}
		})

		It("returns headers haven't been checked", func() {
			err = repository.MarkHeaderChecked(headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := repository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
		})

		It("only treats headers as checked if flop_kicks have been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := repository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(flopKickBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(flopKickBlock)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(flopKickBlock)))
		})

		It("only returns missing headers for the current node", func() {
			node2 := core.Node{}
			db2 := test_config.NewTestDB(node2)
			headerRepository2 := repositories.NewHeaderRepository(db2)
			flopKickRepository2 := flop_kick.NewFlopKickRepository(db2)

			for _, number := range []int64{startingBlock, flopKickBlock, endingBlock} {
				headerRepository2.CreateOrUpdateHeader(core.Header{BlockNumber: number, Raw: rawHeader})
				Expect(err).NotTo(HaveOccurred())
			}

			err = repository.MarkHeaderChecked(headerIds[1])
			Expect(err).NotTo(HaveOccurred())

			node1MissingHeaders, err := repository.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node1MissingHeaders)).To(Equal(2))

			node2MissingHeaders, err := flopKickRepository2.MissingHeaders(startingBlock, endingBlock)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node2MissingHeaders)).To(Equal(3))
		})
	})
})
