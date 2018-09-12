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

package bite_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"database/sql"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/bite"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Bite repository", func() {
	Describe("Create", func() {
		It("persists a bite record", func() {
			node := core.Node{}
			db := test_config.NewTestDB(node)
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			biteRepository := bite.NewBiteRepository(db)

			err = biteRepository.Create(headerID, test_data.BiteModel)

			Expect(err).NotTo(HaveOccurred())
			var dbBite bite.BiteModel
			err = db.Get(&dbBite, `SELECT id, ilk, urn, ink, art, tab, flip, tx_idx, raw_log FROM maker.bite WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbBite.Id).To(Equal(test_data.BiteModel.Id))
			Expect(dbBite.Ilk).To(Equal(test_data.BiteModel.Ilk))
			Expect(dbBite.Urn).To(Equal(test_data.BiteModel.Urn))
			Expect(dbBite.Art).To(Equal(test_data.BiteModel.Art))
			Expect(dbBite.Tab).To(Equal(test_data.BiteModel.Tab))
			Expect(dbBite.Flip).To(Equal(test_data.BiteModel.Flip))
			Expect(dbBite.TransactionIndex).To(Equal(test_data.BiteModel.TransactionIndex))
			Expect(dbBite.Raw).To(MatchJSON(test_data.BiteModel.Raw))
		})

		It("does not duplicate bite events", func() {
			node := core.Node{}
			db := test_config.NewTestDB(node)
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			biteRepository := bite.NewBiteRepository(db)
			err = biteRepository.Create(headerID, test_data.BiteModel)
			Expect(err).NotTo(HaveOccurred())

			var anotherBiteModel = bite.BiteModel{
				Id:               "11",
				Ilk:              test_data.BiteModel.Ilk,
				Urn:              test_data.BiteModel.Urn,
				Ink:              test_data.BiteModel.Ink,
				Art:              test_data.BiteModel.Art,
				Tab:              test_data.BiteModel.Tab,
				Flip:             test_data.BiteModel.Flip,
				IArt:             test_data.BiteModel.IArt,
				TransactionIndex: test_data.BiteModel.TransactionIndex,
				Raw:              test_data.BiteModel.Raw,
			}

			err = biteRepository.Create(headerID, anotherBiteModel)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes bite if corresponding header is deleted", func() {
			node := core.Node{}
			db := test_config.NewTestDB(node)
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			biteRepository := bite.NewBiteRepository(db)
			err = biteRepository.Create(headerID, test_data.BiteModel)
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbBite bite.BiteModel
			err = db.Get(&dbBite, `SELECT id, ilk, urn, ink, art, tab, flip, tx_idx, raw_log FROM maker.bite WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})

	Describe("MissingHeaders", func() {
		It("returns headers with no associated bite event", func() {
			node := core.Node{}
			db := test_config.NewTestDB(node)
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			startingBlockNumber := int64(1)
			biteBlockNumber := int64(2)
			endingBlockNumber := int64(3)
			blockNumbers := []int64{startingBlockNumber, biteBlockNumber, endingBlockNumber, endingBlockNumber + 1}
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
			biteRepository := bite.NewBiteRepository(db)
			err := biteRepository.Create(headerIDs[1], test_data.BiteModel)
			Expect(err).NotTo(HaveOccurred())

			headers, err := biteRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
		})
	})
})
