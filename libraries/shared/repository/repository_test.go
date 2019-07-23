// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package repository_test

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lib/pq"
	"github.com/vulcanize/vulcanizedb/libraries/shared/test_data"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	shared "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	r2 "github.com/vulcanize/vulcanizedb/pkg/contract_watcher/header/repository"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Repository", func() {
	var (
		checkedHeadersColumn string
		db                   *postgres.DB
	)

	Describe("MarkHeaderChecked", func() {
		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)

			checkedHeadersColumn = "test_column_checked"
			_, migrateErr := db.Exec(`ALTER TABLE public.checked_headers
				ADD COLUMN ` + checkedHeadersColumn + ` integer`)
			Expect(migrateErr).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			_, cleanupMigrateErr := db.Exec(`ALTER TABLE public.checked_headers DROP COLUMN ` + checkedHeadersColumn)
			Expect(cleanupMigrateErr).NotTo(HaveOccurred())
		})

		It("marks passed column as checked for passed header", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, headerErr := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(headerErr).NotTo(HaveOccurred())

			err := shared.MarkHeaderChecked(headerID, db, checkedHeadersColumn)

			Expect(err).NotTo(HaveOccurred())
			var checkedCount int
			fetchErr := db.Get(&checkedCount, `SELECT `+checkedHeadersColumn+` FROM public.checked_headers LIMIT 1`)
			Expect(fetchErr).NotTo(HaveOccurred())
			Expect(checkedCount).To(Equal(1))
		})
	})

	Describe("MarkHeaderCheckedInTransaction", func() {
		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)

			checkedHeadersColumn = "test_column_checked"
			_, migrateErr := db.Exec(`ALTER TABLE public.checked_headers
				ADD COLUMN ` + checkedHeadersColumn + ` integer`)
			Expect(migrateErr).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			_, cleanupMigrateErr := db.Exec(`ALTER TABLE public.checked_headers DROP COLUMN ` + checkedHeadersColumn)
			Expect(cleanupMigrateErr).NotTo(HaveOccurred())
		})

		It("marks passed column as checked for passed header within a passed transaction", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, headerErr := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(headerErr).NotTo(HaveOccurred())
			tx, txErr := db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())

			err := shared.MarkHeaderCheckedInTransaction(headerID, tx, checkedHeadersColumn)

			Expect(err).NotTo(HaveOccurred())
			commitErr := tx.Commit()
			Expect(commitErr).NotTo(HaveOccurred())
			var checkedCount int
			fetchErr := db.Get(&checkedCount, `SELECT `+checkedHeadersColumn+` FROM public.checked_headers LIMIT 1`)
			Expect(fetchErr).NotTo(HaveOccurred())
			Expect(checkedCount).To(Equal(1))
		})
	})

	Describe("MissingHeaders", func() {
		var (
			headerRepository         datastore.HeaderRepository
			startingBlockNumber      int64
			endingBlockNumber        int64
			eventSpecificBlockNumber int64
			outOfRangeBlockNumber    int64
			blockNumbers             []int64
			headerIDs                []int64
			notCheckedSQL            string
			err                      error
			hr                       r2.HeaderRepository
			columnNames              []string
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			headerRepository = repositories.NewHeaderRepository(db)
			hr = r2.NewHeaderRepository(db)
			hr.AddCheckColumns(getExpectedColumnNames())

			columnNames, err = shared.GetCheckedColumnNames(db)
			Expect(err).NotTo(HaveOccurred())
			notCheckedSQL = shared.CreateHeaderCheckedPredicateSQL(columnNames, constants.HeaderMissing)

			startingBlockNumber = rand.Int63()
			eventSpecificBlockNumber = startingBlockNumber + 1
			endingBlockNumber = startingBlockNumber + 2
			outOfRangeBlockNumber = endingBlockNumber + 1

			blockNumbers = []int64{startingBlockNumber, eventSpecificBlockNumber, endingBlockNumber, outOfRangeBlockNumber}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		AfterEach(func() {
			test_config.CleanCheckedHeadersTable(db, getExpectedColumnNames())
		})

		It("only treats headers as checked if the event specific logs have been checked", func() {
			//add a checked_header record, but don't mark it check for any of the columns
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := shared.MissingHeaders(startingBlockNumber, endingBlockNumber, db, notCheckedSQL)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber)))
		})

		It("only returns headers associated with the current node", func() {
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			for _, n := range blockNumbers {
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n + 10))
				Expect(err).NotTo(HaveOccurred())
			}

			Expect(err).NotTo(HaveOccurred())
			nodeOneMissingHeaders, err := shared.MissingHeaders(startingBlockNumber, endingBlockNumber, db, notCheckedSQL)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(3))
			Expect(nodeOneMissingHeaders[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))
			Expect(nodeOneMissingHeaders[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))
			Expect(nodeOneMissingHeaders[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(startingBlockNumber), Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))

			nodeTwoMissingHeaders, err := shared.MissingHeaders(startingBlockNumber, endingBlockNumber+10, dbTwo, notCheckedSQL)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(3))
			Expect(nodeTwoMissingHeaders[0].BlockNumber).To(Or(Equal(startingBlockNumber+10), Equal(eventSpecificBlockNumber+10), Equal(endingBlockNumber+10)))
			Expect(nodeTwoMissingHeaders[1].BlockNumber).To(Or(Equal(startingBlockNumber+10), Equal(eventSpecificBlockNumber+10), Equal(endingBlockNumber+10)))
		})

		It("handles an ending block of -1 ", func() {
			endingBlock := int64(-1)
			headers, err := shared.MissingHeaders(startingBlockNumber, endingBlock, db, notCheckedSQL)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(4))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber), Equal(outOfRangeBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber), Equal(outOfRangeBlockNumber)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber), Equal(outOfRangeBlockNumber)))
			Expect(headers[3].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber), Equal(outOfRangeBlockNumber)))

		})

		It("when a the `notCheckedSQL` argument allows for rechecks it returns headers where the checked count is less than the maximum", func() {
			columnName := columnNames[0]
			recheckedSQL := shared.CreateHeaderCheckedPredicateSQL([]string{columnName}, constants.HeaderRecheck)
			// mark every header checked at least once
			// header 4 is marked the maximum number of times, it it is not longer checked

			maxCheckCount, intConversionErr := strconv.Atoi(constants.RecheckHeaderCap)
			Expect(intConversionErr).NotTo(HaveOccurred())

			markHeaderOneErr := shared.MarkHeaderChecked(headerIDs[0], db, columnName)
			Expect(markHeaderOneErr).NotTo(HaveOccurred())
			markHeaderTwoErr := shared.MarkHeaderChecked(headerIDs[1], db, columnName)
			Expect(markHeaderTwoErr).NotTo(HaveOccurred())
			markHeaderThreeErr := shared.MarkHeaderChecked(headerIDs[2], db, columnName)
			Expect(markHeaderThreeErr).NotTo(HaveOccurred())
			for i := 0; i <= maxCheckCount; i++ {
				markHeaderFourErr := shared.MarkHeaderChecked(headerIDs[3], db, columnName)
				Expect(markHeaderFourErr).NotTo(HaveOccurred())
			}

			headers, err := shared.MissingHeaders(1, -1, db, recheckedSQL)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].Id).To(Or(Equal(headerIDs[0]), Equal(headerIDs[1]), Equal(headerIDs[2])))
			Expect(headers[1].Id).To(Or(Equal(headerIDs[0]), Equal(headerIDs[1]), Equal(headerIDs[2])))
			Expect(headers[2].Id).To(Or(Equal(headerIDs[0]), Equal(headerIDs[1]), Equal(headerIDs[2])))
		})
	})

	Describe("GetCheckedColumnNames", func() {
		It("gets the column names from checked_headers", func() {
			db := test_config.NewTestDB(test_config.NewTestNode())
			hr := r2.NewHeaderRepository(db)
			hr.AddCheckColumns(getExpectedColumnNames())
			test_config.CleanTestDB(db)
			expectedColumnNames := getExpectedColumnNames()
			actualColumnNames, err := shared.GetCheckedColumnNames(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualColumnNames).To(Equal(expectedColumnNames))
			test_config.CleanCheckedHeadersTable(db, getExpectedColumnNames())
		})
	})

	Describe("CreateHeaderCheckedPredicateSQL", func() {
		Describe("for headers that haven't been checked for logs", func() {
			It("generates a correct SQL string for one column", func() {
				columns := []string{"columnA"}
				expected := " (columnA=0)"
				actual := shared.CreateHeaderCheckedPredicateSQL(columns, constants.HeaderMissing)
				Expect(actual).To(Equal(expected))
			})

			It("generates a correct SQL string for several columns", func() {
				columns := []string{"columnA", "columnB"}
				expected := " (columnA=0 OR columnB=0)"
				actual := shared.CreateHeaderCheckedPredicateSQL(columns, constants.HeaderMissing)
				Expect(actual).To(Equal(expected))
			})

			It("defaults to FALSE when there are no columns", func() {
				expected := "FALSE"
				actual := shared.CreateHeaderCheckedPredicateSQL([]string{}, constants.HeaderMissing)
				Expect(actual).To(Equal(expected))
			})
		})

		Describe("for headers that are being rechecked for logs", func() {
			It("generates a correct SQL string for rechecking headers for one column", func() {
				columns := []string{"columnA"}
				expected := fmt.Sprintf(" (columnA<%s)", constants.RecheckHeaderCap)
				actual := shared.CreateHeaderCheckedPredicateSQL(columns, constants.HeaderRecheck)
				Expect(actual).To(Equal(expected))
			})

			It("generates a correct SQL string for rechecking headers for several columns", func() {
				columns := []string{"columnA", "columnB"}
				expected := fmt.Sprintf(" (columnA<%s OR columnB<%s)", constants.RecheckHeaderCap, constants.RecheckHeaderCap)
				actual := shared.CreateHeaderCheckedPredicateSQL(columns, constants.HeaderRecheck)
				Expect(actual).To(Equal(expected))
			})

			It("defaults to FALSE when there are no columns", func() {
				expected := "FALSE"
				actual := shared.CreateHeaderCheckedPredicateSQL([]string{}, constants.HeaderRecheck)
				Expect(actual).To(Equal(expected))
			})
		})
	})

	Describe("CreateHeaderSyncLogs", func() {
		var headerID int64

		type HeaderSyncLog struct {
			ID          int64
			HeaderID    int64 `db:"header_id"`
			Address     string
			Topics      pq.ByteaArray
			Data        []byte
			BlockNumber uint64 `db:"block_number"`
			BlockHash   string `db:"block_hash"`
			TxHash      string `db:"tx_hash"`
			TxIndex     uint   `db:"tx_index"`
			LogIndex    uint   `db:"log_index"`
			Transformed bool
			Raw         []byte
		}

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			var headerErr error
			headerID, headerErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(headerErr).NotTo(HaveOccurred())
		})

		It("writes a log to the db", func() {
			log := test_data.GenericTestLog()

			_, err := shared.CreateLogs(headerID, []types.Log{log}, db)

			Expect(err).NotTo(HaveOccurred())
			var dbLog HeaderSyncLog
			lookupErr := db.Get(&dbLog, `SELECT * FROM header_sync_logs`)
			Expect(lookupErr).NotTo(HaveOccurred())
			Expect(dbLog.ID).NotTo(BeZero())
			Expect(dbLog.HeaderID).To(Equal(headerID))
			Expect(dbLog.Address).To(Equal(log.Address.Hex()))
			Expect(dbLog.Topics[0]).To(Equal(log.Topics[0].Bytes()))
			Expect(dbLog.Topics[1]).To(Equal(log.Topics[1].Bytes()))
			Expect(dbLog.Data).To(Equal(log.Data))
			Expect(dbLog.BlockNumber).To(Equal(log.BlockNumber))
			Expect(dbLog.BlockHash).To(Equal(log.BlockHash.Hex()))
			Expect(dbLog.TxIndex).To(Equal(log.TxIndex))
			Expect(dbLog.TxHash).To(Equal(log.TxHash.Hex()))
			Expect(dbLog.LogIndex).To(Equal(log.Index))
			expectedRaw, jsonErr := log.MarshalJSON()
			Expect(jsonErr).NotTo(HaveOccurred())
			Expect(dbLog.Raw).To(MatchJSON(expectedRaw))
			Expect(dbLog.Transformed).To(BeFalse())
		})

		It("writes several logs to the db", func() {
			log1 := test_data.GenericTestLog()
			log2 := test_data.GenericTestLog()
			logs := []types.Log{log1, log2}

			_, err := shared.CreateLogs(headerID, logs, db)

			Expect(err).NotTo(HaveOccurred())
			var count int
			lookupErr := db.Get(&count, `SELECT COUNT(*) FROM header_sync_logs`)
			Expect(lookupErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(len(logs)))
		})

		It("persists record that can be unpacked into types.Log", func() {
			// important if we want to decouple log persistence from transforming and still make use of
			// tools on types.Log like abi.Unpack

			log := test_data.GenericTestLog()

			_, err := shared.CreateLogs(headerID, []types.Log{log}, db)

			Expect(err).NotTo(HaveOccurred())
			var dbLog HeaderSyncLog
			lookupErr := db.Get(&dbLog, `SELECT * FROM header_sync_logs`)
			Expect(lookupErr).NotTo(HaveOccurred())

			var logTopics []common.Hash
			for _, topic := range dbLog.Topics {
				logTopics = append(logTopics, common.BytesToHash(topic))
			}

			reconstructedLog := types.Log{
				Address:     common.HexToAddress(dbLog.Address),
				Topics:      logTopics,
				Data:        dbLog.Data,
				BlockNumber: dbLog.BlockNumber,
				TxHash:      common.HexToHash(dbLog.TxHash),
				TxIndex:     dbLog.TxIndex,
				BlockHash:   common.HexToHash(dbLog.BlockHash),
				Index:       dbLog.LogIndex,
				Removed:     false,
			}
			Expect(reconstructedLog).To(Equal(log))
		})

		It("does not duplicate logs", func() {
			log := test_data.GenericTestLog()

			results, err := shared.CreateLogs(headerID, []types.Log{log, log}, db)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(results)).To(Equal(1))
			var count int
			lookupErr := db.Get(&count, `SELECT COUNT(*) FROM header_sync_logs`)
			Expect(lookupErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns results with log id and header id for persisted logs", func() {
			log1 := test_data.GenericTestLog()
			log2 := test_data.GenericTestLog()
			logs := []types.Log{log1, log2}

			results, err := shared.CreateLogs(headerID, logs, db)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(results)).To(Equal(len(logs)))
			var log1ID, log2ID int64
			lookupErr := db.Get(&log1ID, `SELECT id FROM header_sync_logs WHERE data = $1`, log1.Data)
			Expect(lookupErr).NotTo(HaveOccurred())
			lookup2Err := db.Get(&log2ID, `SELECT id FROM header_sync_logs WHERE data = $1`, log2.Data)
			Expect(lookup2Err).NotTo(HaveOccurred())
			Expect(results[0].ID).To(Or(Equal(log1ID), Equal(log2ID)))
			Expect(results[1].ID).To(Or(Equal(log1ID), Equal(log2ID)))
			Expect(results[0].HeaderID).To(Equal(headerID))
			Expect(results[1].HeaderID).To(Equal(headerID))
		})

		It("returns results with properties for persisted logs", func() {
			log1 := test_data.GenericTestLog()
			log2 := test_data.GenericTestLog()
			logs := []types.Log{log1, log2}

			results, err := shared.CreateLogs(headerID, logs, db)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(results)).To(Equal(len(logs)))
			Expect(results[0].Log).To(Or(Equal(log1), Equal(log2)))
			Expect(results[1].Log).To(Or(Equal(log1), Equal(log2)))
			Expect(results[0].Transformed).To(BeFalse())
			Expect(results[1].Transformed).To(BeFalse())
		})
	})
})

func getExpectedColumnNames() []string {
	return []string{
		"column_1",
		"column_2",
		"column_3",
		"column_4",
	}
}
