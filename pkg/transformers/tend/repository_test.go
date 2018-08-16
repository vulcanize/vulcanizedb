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

package tend_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("TendRepository", func() {
	var db *postgres.DB
	var tendRepository tend.TendRepository
	var headerRepository repositories.HeaderRepository
	var headerId int64
	var err error

	BeforeEach(func() {
		node := test_config.NewTestNode()
		db = test_config.NewTestDB(node)
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		headerId, err = headerRepository.CreateOrUpdateHeader(core.Header{})
		Expect(err).NotTo(HaveOccurred())

		tendRepository = tend.NewTendRepository(db)
	})

	Describe("Create", func() {
		It("persists a tend record", func() {
			err := tendRepository.Create(headerId, test_data.TendModel)
			Expect(err).NotTo(HaveOccurred())

			var count int
			err = db.QueryRow(`SELECT count(*) from maker.tend`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			dbResult := tend.TendModel{}
			err = db.Get(&dbResult, `SELECT id, lot, bid, guy, tic, era, tx_idx, raw_log FROM maker.tend WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbResult.Id).To(Equal(test_data.TendModel.Id))
			Expect(dbResult.Lot).To(Equal(test_data.TendModel.Lot))
			Expect(dbResult.Bid).To(Equal(test_data.TendModel.Bid))
			Expect(dbResult.Guy).To(Equal(test_data.TendModel.Guy))
			Expect(dbResult.Tic).To(Equal(test_data.TendModel.Tic))
			Expect(dbResult.Era.Equal(test_data.TendModel.Era)).To(BeTrue())
			Expect(dbResult.TransactionIndex).To(Equal(test_data.TendModel.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(test_data.RawJson))
		})

		It("returns an error if inserting a tend record fails", func() {
			err := tendRepository.Create(headerId, test_data.TendModel)
			Expect(err).NotTo(HaveOccurred())

			err = tendRepository.Create(headerId, test_data.TendModel)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("deletes the tend record if its corresponding header record is deleted", func() {
			err := tendRepository.Create(headerId, test_data.TendModel)
			Expect(err).NotTo(HaveOccurred())

			var count int
			err = db.QueryRow(`SELECT count(*) from maker.tend`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			_, err = db.Exec(`DELETE FROM headers where id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())

			err = db.QueryRow(`SELECT count(*) from maker.tend`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})

	Describe("MissingHeaders", func() {
		var tendBlockNumber int64
		var startingBlockNumber int64
		var endingBlockNumber int64
		var outOfRangeBlockNumber int64

		BeforeEach(func() {
			tendBlockNumber = rand.Int63()
			startingBlockNumber = tendBlockNumber - 1
			endingBlockNumber = tendBlockNumber + 1
			outOfRangeBlockNumber = tendBlockNumber + 2
		})

		It("returns headers for which there isn't an associated flip_kick record", func() {
			var headerIds []int64

			for _, number := range []int64{startingBlockNumber, tendBlockNumber, endingBlockNumber, outOfRangeBlockNumber} {
				headerId, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: number})
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)
			}

			err = tendRepository.Create(headerIds[1], test_data.TendModel)
			Expect(err).NotTo(HaveOccurred())

			headers, err := tendRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
		})

		It("only returns missing headers for the current node", func() {
			var headerIds []int64
			node2 := core.Node{}
			db2 := test_config.NewTestDB(node2)
			headerRepository2 := repositories.NewHeaderRepository(db2)
			tendRepository2 := tend.NewTendRepository(db2)

			for _, number := range []int64{startingBlockNumber, tendBlockNumber, endingBlockNumber} {
				headerId, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: number})
				Expect(err).NotTo(HaveOccurred())
				headerIds = append(headerIds, headerId)

				headerRepository2.CreateOrUpdateHeader(core.Header{BlockNumber: number})
				Expect(err).NotTo(HaveOccurred())
			}

			err = tendRepository.Create(headerIds[1], test_data.TendModel)
			Expect(err).NotTo(HaveOccurred())

			node1MissingHeaders, err := tendRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node1MissingHeaders)).To(Equal(2))

			node2MissingHeaders, err := tendRepository2.MissingHeaders(startingBlockNumber, endingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node2MissingHeaders)).To(Equal(3))
		})
	})
})
