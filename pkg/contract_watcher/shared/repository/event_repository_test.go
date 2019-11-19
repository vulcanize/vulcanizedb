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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fc "github.com/makerdao/vulcanizedb/pkg/contract_watcher/full/converter"
	lc "github.com/makerdao/vulcanizedb/pkg/contract_watcher/header/converter"
	lr "github.com/makerdao/vulcanizedb/pkg/contract_watcher/header/repository"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/shared/contract"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/shared/helpers/test_helpers"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/shared/helpers/test_helpers/mocks"
	sr "github.com/makerdao/vulcanizedb/pkg/contract_watcher/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/shared/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
)

var _ = Describe("Repository", func() {
	var db *postgres.DB
	var dataStore sr.EventRepository
	var err error
	var log *types.Log
	var logs []types.Log
	var con *contract.Contract
	var vulcanizeLogId int64
	var wantedEvents = []string{"Transfer"}
	var wantedMethods = []string{"balanceOf"}
	var event types.Event
	var headerID int64
	var mockEvent = mocks.MockTranferEvent
	var mockLog1 = mocks.MockTransferLog1
	var mockLog2 = mocks.MockTransferLog2

	BeforeEach(func() {
		db, con = test_helpers.SetupTusdRepo(&vulcanizeLogId, wantedEvents, wantedMethods)
		mockEvent.LogID = vulcanizeLogId

		event = con.Events["Transfer"]
		err = con.GenerateFilters()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	Describe("Full sync mode", func() {
		BeforeEach(func() {
			dataStore = sr.NewEventRepository(db, types.FullSync)
		})

		Describe("CreateContractSchema", func() {
			It("Creates schema if it doesn't exist", func() {
				created, err := dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				created, err = dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(false))
			})

			It("Caches schema it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
				_, ok := dataStore.CheckSchemaCache(con.Address)
				Expect(ok).To(Equal(false))

				created, err := dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				v, ok := dataStore.CheckSchemaCache(con.Address)
				Expect(ok).To(Equal(true))
				Expect(v).To(Equal(true))
			})
		})

		Describe("CreateEventTable", func() {
			It("Creates table if it doesn't exist", func() {
				created, err := dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				created, err = dataStore.CreateEventTable(con.Address, event)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				created, err = dataStore.CreateEventTable(con.Address, event)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(false))
			})

			It("Caches table it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
				created, err := dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				tableID := fmt.Sprintf("%s_%s.%s_event", types.FullSync, strings.ToLower(con.Address), strings.ToLower(event.Name))
				_, ok := dataStore.CheckTableCache(tableID)
				Expect(ok).To(Equal(false))

				created, err = dataStore.CreateEventTable(con.Address, event)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				v, ok := dataStore.CheckTableCache(tableID)
				Expect(ok).To(Equal(true))
				Expect(v).To(Equal(true))
			})
		})

		Describe("PersistLogs", func() {
			BeforeEach(func() {
				c := fc.Converter{}
				c.Update(con)
				log, err = c.Convert(mockEvent, event)
				Expect(err).ToNot(HaveOccurred())
			})

			It("Persists contract event log values into custom tables", func() {
				err = dataStore.PersistLogs([]types.Log{*log}, event, con.Address, con.Name)
				Expect(err).ToNot(HaveOccurred())

				b, ok := con.EmittedAddrs[common.HexToAddress("0x000000000000000000000000000000000000Af21")]
				Expect(ok).To(Equal(true))
				Expect(b).To(Equal(true))

				b, ok = con.EmittedAddrs[common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391")]
				Expect(ok).To(Equal(true))
				Expect(b).To(Equal(true))

				scanLog := test_helpers.TransferLog{}

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.transfer_event", constants.TusdContractAddress)).StructScan(&scanLog)
				Expect(err).ToNot(HaveOccurred())
				expectedLog := test_helpers.TransferLog{
					Id:             1,
					VulcanizeLogId: vulcanizeLogId,
					TokenName:      "TrueUSD",
					Block:          5488076,
					Tx:             "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
					From:           "0x000000000000000000000000000000000000Af21",
					To:             "0x09BbBBE21a5975cAc061D82f7b843bCE061BA391",
					Value:          "1097077688018008265106216665536940668749033598146",
				}
				Expect(scanLog).To(Equal(expectedLog))
			})

			It("Doesn't persist duplicate event logs", func() {
				// Try to persist the same log twice in a single call
				err = dataStore.PersistLogs([]types.Log{*log, *log}, event, con.Address, con.Name)
				Expect(err).ToNot(HaveOccurred())

				scanLog := test_helpers.TransferLog{}

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.transfer_event", constants.TusdContractAddress)).StructScan(&scanLog)
				Expect(err).ToNot(HaveOccurred())
				expectedLog := test_helpers.TransferLog{
					Id:             1,
					VulcanizeLogId: vulcanizeLogId,
					TokenName:      "TrueUSD",
					Block:          5488076,
					Tx:             "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
					From:           "0x000000000000000000000000000000000000Af21",
					To:             "0x09BbBBE21a5975cAc061D82f7b843bCE061BA391",
					Value:          "1097077688018008265106216665536940668749033598146",
				}
				Expect(scanLog).To(Equal(expectedLog))

				// Attempt to persist the same log again in separate call
				err = dataStore.PersistLogs([]types.Log{*log}, event, con.Address, con.Name)
				Expect(err).ToNot(HaveOccurred())

				// Show that no new logs were entered
				var count int
				err = db.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM full_%s.transfer_event", constants.TusdContractAddress))
				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("Fails with empty log", func() {
				err = dataStore.PersistLogs([]types.Log{}, event, con.Address, con.Name)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Header sync mode", func() {
		BeforeEach(func() {
			dataStore = sr.NewEventRepository(db, types.HeaderSync)
		})

		Describe("CreateContractSchema", func() {
			It("Creates schema if it doesn't exist", func() {
				created, err := dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				created, err = dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(false))
			})

			It("Caches schema it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
				_, ok := dataStore.CheckSchemaCache(con.Address)
				Expect(ok).To(Equal(false))

				created, err := dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				v, ok := dataStore.CheckSchemaCache(con.Address)
				Expect(ok).To(Equal(true))
				Expect(v).To(Equal(true))
			})

			It("Caches table it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
				created, err := dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				tableID := fmt.Sprintf("%s_%s.%s_event", types.HeaderSync, strings.ToLower(con.Address), strings.ToLower(event.Name))
				_, ok := dataStore.CheckTableCache(tableID)
				Expect(ok).To(Equal(false))

				created, err = dataStore.CreateEventTable(con.Address, event)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				v, ok := dataStore.CheckTableCache(tableID)
				Expect(ok).To(Equal(true))
				Expect(v).To(Equal(true))
			})
		})

		Describe("CreateEventTable", func() {
			It("Creates table if it doesn't exist", func() {
				created, err := dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				created, err = dataStore.CreateEventTable(con.Address, event)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				created, err = dataStore.CreateEventTable(con.Address, event)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(false))
			})
		})

		Describe("PersistLogs", func() {
			BeforeEach(func() {
				headerRepository := repositories.NewHeaderRepository(db)
				headerID, err = headerRepository.CreateOrUpdateHeader(mocks.MockHeader1)
				Expect(err).ToNot(HaveOccurred())
				c := lc.Converter{}
				c.Update(con)
				logs, err = c.Convert([]geth.Log{mockLog1, mockLog2}, event, headerID)
				Expect(err).ToNot(HaveOccurred())
			})

			It("Persists contract event log values into custom tables", func() {
				hr := lr.NewHeaderRepository(db)
				err = hr.AddCheckColumn(event.Name + "_" + con.Address)
				Expect(err).ToNot(HaveOccurred())

				err = dataStore.PersistLogs(logs, event, con.Address, con.Name)
				Expect(err).ToNot(HaveOccurred())

				var count int
				err = db.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM header_%s.transfer_event", constants.TusdContractAddress))
				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(2))

				scanLog := test_helpers.HeaderSyncTransferLog{}
				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM header_%s.transfer_event LIMIT 1", constants.TusdContractAddress)).StructScan(&scanLog)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanLog.HeaderID).To(Equal(headerID))
				Expect(scanLog.TokenName).To(Equal("TrueUSD"))
				Expect(scanLog.TxIndex).To(Equal(int64(110)))
				Expect(scanLog.LogIndex).To(Equal(int64(1)))
				Expect(scanLog.From).To(Equal("0x000000000000000000000000000000000000Af21"))
				Expect(scanLog.To).To(Equal("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"))
				Expect(scanLog.Value).To(Equal("1097077688018008265106216665536940668749033598146"))

				var expectedRawLog, rawLog geth.Log
				err = json.Unmarshal(logs[0].Raw, &expectedRawLog)
				Expect(err).ToNot(HaveOccurred())
				err = json.Unmarshal(scanLog.RawLog, &rawLog)
				Expect(err).ToNot(HaveOccurred())
				Expect(rawLog).To(Equal(expectedRawLog))
			})

			It("Doesn't persist duplicate event logs", func() {
				hr := lr.NewHeaderRepository(db)
				err = hr.AddCheckColumn(event.Name + "_" + con.Address)
				Expect(err).ToNot(HaveOccurred())

				// Successfully persist the two unique logs
				err = dataStore.PersistLogs(logs, event, con.Address, con.Name)
				Expect(err).ToNot(HaveOccurred())

				// Try to insert the same logs again
				err = dataStore.PersistLogs(logs, event, con.Address, con.Name)
				Expect(err).ToNot(HaveOccurred())

				// Show that no new logs were entered
				var count int
				err = db.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM header_%s.transfer_event", constants.TusdContractAddress))
				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(2))
			})

			It("inserts additional log if only some are duplicate", func() {
				hr := lr.NewHeaderRepository(db)
				err = hr.AddCheckColumn(event.Name + "_" + con.Address)
				Expect(err).ToNot(HaveOccurred())

				// Successfully persist first log
				err = dataStore.PersistLogs([]types.Log{logs[0]}, event, con.Address, con.Name)
				Expect(err).ToNot(HaveOccurred())

				// Successfully persist second log even though first already persisted
				err = dataStore.PersistLogs(logs, event, con.Address, con.Name)
				Expect(err).ToNot(HaveOccurred())

				// Show that both logs were entered
				var count int
				err = db.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM header_%s.transfer_event", constants.TusdContractAddress))
				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(2))
			})

			It("Fails with empty log", func() {
				err = dataStore.PersistLogs([]types.Log{}, event, con.Address, con.Name)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
