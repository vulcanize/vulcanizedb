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

package flip_kick_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("FlipKick Repository", func() {
	var db *postgres.DB
	var flipKickRepository flip_kick.FlipKickRepository
	var headerId int64
	var blockNumber int64
	var flipKick = test_data.FlipKickModel

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		flipKickRepository = flip_kick.FlipKickRepository{DB: db}
		blockNumber = GinkgoRandomSeed()
		headerId = createHeader(db, blockNumber)

		_, err := db.Exec(`DELETE from maker.flip_kick;`)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Create", func() {
		AfterEach(func() {
			_, err := db.Exec(`DELETE from headers;`)
			Expect(err).NotTo(HaveOccurred())
		})

		It("persists flip_kick records", func() {
			err := flipKickRepository.Create(headerId, []flip_kick.FlipKickModel{flipKick})
			Expect(err).NotTo(HaveOccurred())

			assertDBRecordCount(db, "maker.flip_kick", 1)

			dbResult := test_data.FlipKickDBRow{}
			err = flipKickRepository.DB.QueryRowx(`SELECT * FROM maker.flip_kick`).StructScan(&dbResult)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.HeaderId).To(Equal(headerId))
			Expect(dbResult.BidId).To(Equal(flipKick.BidId))
			Expect(dbResult.Lot).To(Equal(flipKick.Lot))
			Expect(dbResult.Bid).To(Equal(flipKick.Bid))
			Expect(dbResult.Gal).To(Equal(flipKick.Gal))
			Expect(dbResult.End.Equal(flipKick.End)).To(BeTrue())
			Expect(dbResult.Urn).To(Equal(flipKick.Urn))
			Expect(dbResult.Tab).To(Equal(flipKick.Tab))
			Expect(dbResult.TransactionIndex).To(Equal(flipKick.TransactionIndex))
			Expect(dbResult.Raw).To(MatchJSON(flipKick.Raw))
		})

		It("marks header checked", func() {
			err := flipKickRepository.Create(headerId, []flip_kick.FlipKickModel{flipKick})
			Expect(err).NotTo(HaveOccurred())

			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT flip_kick_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("returns an error if inserting the flip_kick record fails", func() {
			err := flipKickRepository.Create(headerId, []flip_kick.FlipKickModel{flipKick})
			Expect(err).NotTo(HaveOccurred())

			err = flipKickRepository.Create(headerId, []flip_kick.FlipKickModel{flipKick})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("deletes the flip_kick records if its corresponding header record is deleted", func() {
			err := flipKickRepository.Create(headerId, []flip_kick.FlipKickModel{flipKick})
			Expect(err).NotTo(HaveOccurred())
			assertDBRecordCount(db, "maker.flip_kick", 1)
			assertDBRecordCount(db, "headers", 1)

			_, err = db.Exec(`DELETE FROM headers where id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())

			assertDBRecordCount(db, "headers", 0)
			assertDBRecordCount(db, "maker.flip_kick", 0)
		})
	})

	Describe("MarkHeaderChecked", func() {
		It("creates a row for a new headerID", func() {
			err := flipKickRepository.MarkHeaderChecked(headerId)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT flip_kick_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerId)
			Expect(err).NotTo(HaveOccurred())

			err = flipKickRepository.MarkHeaderChecked(headerId)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT flip_kick_checked FROM public.checked_headers WHERE header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("When there are multiple nodes", func() {
		var db2 *postgres.DB
		var flipKickRepository2 flip_kick.FlipKickRepository

		BeforeEach(func() {
			//create database for the second node
			node2 := core.Node{
				GenesisBlock: "GENESIS",
				NetworkID:    1,
				ID:           "node2",
				ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
			}
			db2 = test_config.NewTestDB(node2)
			flipKickRepository2 = flip_kick.FlipKickRepository{DB: db2}
			createHeader(db2, blockNumber)

			_, err := db2.Exec(`DELETE from maker.flip_kick;`)
			Expect(err).NotTo(HaveOccurred())
		})

		It("only includes missing headers for the current node", func() {
			node1missingHeaders, err := flipKickRepository.MissingHeaders(blockNumber, blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node1missingHeaders)).To(Equal(1))

			node2MissingHeaders, err := flipKickRepository2.MissingHeaders(blockNumber, blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node2MissingHeaders)).To(Equal(1))
		})
	})

	Describe("MissingHeaders", func() {
		It("returns headers that haven't been marked as checked", func() {
			startingBlock := blockNumber - 3
			endingBlock := blockNumber + 3
			err := flipKickRepository.MarkHeaderChecked(headerId)
			Expect(err).NotTo(HaveOccurred())

			newBlockNumber := blockNumber + 3
			newHeaderId := createHeader(db, newBlockNumber)
			createHeader(db, blockNumber+10) //this one is out of the block range and shouldn't be included
			headers, err := flipKickRepository.MissingHeaders(startingBlock, endingBlock)
			Expect(len(headers)).To(Equal(1))
			Expect(headers[0].Id).To(Equal(newHeaderId))
			Expect(headers[0].BlockNumber).To(Equal(newBlockNumber))
		})
	})
})

func assertDBRecordCount(db *postgres.DB, dbTable string, expectedCount int) {
	var count int
	query := `SELECT count(*) FROM ` + dbTable
	err := db.QueryRow(query).Scan(&count)
	Expect(err).NotTo(HaveOccurred())
	Expect(count).To(Equal(expectedCount))
}

func createHeader(db *postgres.DB, blockNumber int64) (headerId int64) {
	headerRepository := repositories.NewHeaderRepository(db)
	rawHeader, err := json.Marshal(types.Header{})
	Expect(err).NotTo(HaveOccurred())
	header := core.Header{
		BlockNumber: blockNumber,
		Hash:        common.BytesToHash([]byte{1, 2, 3, 4, 5}).Hex(),
		Raw:         rawHeader,
	}
	_, err = headerRepository.CreateOrUpdateHeader(header)
	Expect(err).NotTo(HaveOccurred())

	var dbHeader core.Header
	err = db.Get(&dbHeader, `SELECT id, block_number, hash, raw FROM public.headers WHERE block_number = $1 AND eth_node_id = $2`, header.BlockNumber, db.NodeID)
	Expect(err).NotTo(HaveOccurred())
	return dbHeader.Id
}
