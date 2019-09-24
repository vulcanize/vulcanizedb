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

package repositories_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	repository2 "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/libraries/shared/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Header sync log repository", func() {
	var (
		db         *postgres.DB
		headerID   int64
		repository datastore.HeaderSyncLogRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var headerErr error
		headerID, headerErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(headerErr).NotTo(HaveOccurred())
		repository = repositories.NewHeaderSyncLogRepository(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("CreateHeaderSyncLogs", func() {
		type headerSyncLog struct {
			ID          int64
			HeaderID    int64 `db:"header_id"`
			Address     int64
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

		It("writes a log to the db", func() {
			log := test_data.GenericTestLog()

			err := repository.CreateHeaderSyncLogs(headerID, []types.Log{log})

			Expect(err).NotTo(HaveOccurred())
			var dbLog headerSyncLog
			lookupErr := db.Get(&dbLog, `SELECT * FROM header_sync_logs`)
			Expect(lookupErr).NotTo(HaveOccurred())
			Expect(dbLog.ID).NotTo(BeZero())
			Expect(dbLog.HeaderID).To(Equal(headerID))
			actualAddress, addressErr := repository2.GetAddressById(db, dbLog.Address)
			Expect(addressErr).NotTo(HaveOccurred())
			Expect(actualAddress).To(Equal(log.Address.Hex()))
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

			err := repository.CreateHeaderSyncLogs(headerID, logs)

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

			err := repository.CreateHeaderSyncLogs(headerID, []types.Log{log})

			Expect(err).NotTo(HaveOccurred())
			var dbLog headerSyncLog
			lookupErr := db.Get(&dbLog, `SELECT * FROM header_sync_logs`)
			Expect(lookupErr).NotTo(HaveOccurred())

			var logTopics []common.Hash
			for _, topic := range dbLog.Topics {
				logTopics = append(logTopics, common.BytesToHash(topic))
			}

			actualAddress, addressErr := repository2.GetAddressById(db, dbLog.Address)
			Expect(addressErr).NotTo(HaveOccurred())
			reconstructedLog := types.Log{
				Address:     common.HexToAddress(actualAddress),
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

			err := repository.CreateHeaderSyncLogs(headerID, []types.Log{log, log})

			Expect(err).NotTo(HaveOccurred())
			var count int
			lookupErr := db.Get(&count, `SELECT COUNT(*) FROM header_sync_logs`)
			Expect(lookupErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})

	Describe("GetUntransformedHeaderSyncLogs", func() {
		Describe("when there are no logs", func() {
			It("returns empty collection", func() {
				result, err := repository.GetUntransformedHeaderSyncLogs()

				Expect(err).NotTo(HaveOccurred())
				Expect(len(result)).To(BeZero())
			})
		})

		Describe("when there are logs", func() {
			var log1, log2 types.Log

			BeforeEach(func() {
				log1 = test_data.GenericTestLog()
				log2 = test_data.GenericTestLog()
				logs := []types.Log{log1, log2}
				logsErr := repository.CreateHeaderSyncLogs(headerID, logs)
				Expect(logsErr).NotTo(HaveOccurred())
			})

			It("returns persisted logs", func() {
				result, err := repository.GetUntransformedHeaderSyncLogs()

				Expect(err).NotTo(HaveOccurred())
				Expect(len(result)).To(Equal(2))
				Expect(result[0].Log).To(Or(Equal(log1), Equal(log2)))
				Expect(result[1].Log).To(Or(Equal(log1), Equal(log2)))
				Expect(result[0].Log).NotTo(Equal(result[1].Log))
			})

			It("excludes logs that have been transformed", func() {
				_, insertErr := db.Exec(`UPDATE public.header_sync_logs SET transformed = true WHERE tx_hash = $1`, log1.TxHash.Hex())
				Expect(insertErr).NotTo(HaveOccurred())

				result, err := repository.GetUntransformedHeaderSyncLogs()

				Expect(err).NotTo(HaveOccurred())
				Expect(len(result)).To(Equal(1))
				Expect(result[0].Log).To(Equal(log2))
			})

			It("returns empty collection if all logs transformed", func() {
				_, insertErr := db.Exec(`UPDATE public.header_sync_logs SET transformed = true WHERE header_id = $1`, headerID)
				Expect(insertErr).NotTo(HaveOccurred())

				result, err := repository.GetUntransformedHeaderSyncLogs()

				Expect(err).NotTo(HaveOccurred())
				Expect(len(result)).To(BeZero())
			})
		})
	})
})
