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

package drip_drip_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_drip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Drip drip repository", func() {
	Describe("Create", func() {
		It("adds a drip drip event", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			dripDripRepository := drip_drip.NewDripDripRepository(db)

			err = dripDripRepository.Create(headerID, test_data.DripDripModel)

			Expect(err).NotTo(HaveOccurred())
			var dbPitFile drip_drip.DripDripModel
			err = db.Get(&dbPitFile, `SELECT ilk, tx_idx, raw_log FROM maker.drip_drip WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbPitFile.Ilk).To(Equal(test_data.DripDripModel.Ilk))
			Expect(dbPitFile.TransactionIndex).To(Equal(test_data.DripDripModel.TransactionIndex))
			Expect(dbPitFile.Raw).To(MatchJSON(test_data.DripDripModel.Raw))
		})

		It("does not duplicate drip drip events", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			dripDripRepository := drip_drip.NewDripDripRepository(db)
			err = dripDripRepository.Create(headerID, test_data.DripDripModel)
			Expect(err).NotTo(HaveOccurred())

			err = dripDripRepository.Create(headerID, test_data.DripDripModel)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes drip drip if corresponding header is deleted", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			dripDripRepository := drip_drip.NewDripDripRepository(db)
			err = dripDripRepository.Create(headerID, test_data.DripDripModel)
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbPitFile drip_drip.DripDripModel
			err = db.Get(&dbPitFile, `SELECT ilk, tx_idx, raw_log FROM maker.drip_drip WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})

	Describe("MissingHeaders", func() {
		It("returns headers with no associated drip drip event", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			startingBlockNumber := int64(1)
			dripDripBlockNumber := int64(2)
			endingBlockNumber := int64(3)
			blockNumbers := []int64{startingBlockNumber, dripDripBlockNumber, endingBlockNumber, endingBlockNumber + 1}
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
			dripDripRepository := drip_drip.NewDripDripRepository(db)
			err := dripDripRepository.Create(headerIDs[1], test_data.DripDripModel)
			Expect(err).NotTo(HaveOccurred())

			headers, err := dripDripRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
		})

		It("only returns headers associated with the current node", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			blockNumbers := []int64{1, 2, 3}
			headerRepository := repositories.NewHeaderRepository(db)
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				Expect(err).NotTo(HaveOccurred())
			}
			dripDripRepository := drip_drip.NewDripDripRepository(db)
			dripDripRepositoryTwo := drip_drip.NewDripDripRepository(dbTwo)
			err := dripDripRepository.Create(headerIDs[0], test_data.DripDripModel)
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := dripDripRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := dripDripRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
