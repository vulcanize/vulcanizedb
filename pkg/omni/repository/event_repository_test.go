// VulcanizeDB
// Copyright © 2018 Vulcanize

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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/converter"
	"github.com/vulcanize/vulcanizedb/pkg/omni/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

var mockEvent = core.WatchedEvent{
	Name:        constants.TransferEvent.String(),
	BlockNumber: 5488076,
	Address:     constants.TusdContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topic0:      constants.TransferEvent.Signature(),
	Topic1:      "0x000000000000000000000000000000000000000000000000000000000000af21",
	Topic2:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
	Topic3:      "",
	Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var _ = Describe("Repository", func() {
	var db *postgres.DB
	var dataStore repository.EventDatastore
	var err error
	var log *types.Log
	var con *contract.Contract
	var vulcanizeLogId int64
	var wantedEvents = []string{"Transfer"}
	var event types.Event

	BeforeEach(func() {
		db, con = test_helpers.SetupTusdRepo(&vulcanizeLogId, wantedEvents, []string{})
		mockEvent.LogID = vulcanizeLogId

		event = con.Events["Transfer"]
		err = con.GenerateFilters()
		Expect(err).ToNot(HaveOccurred())

		c := converter.NewConverter(con)
		log, err = c.Convert(mockEvent, event)
		Expect(err).ToNot(HaveOccurred())

		dataStore = repository.NewEventDataStore(db)
	})

	AfterEach(func() {
		test_helpers.TearDown(db)
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
	})

	Describe("CreateEventTable", func() {
		It("Creates table if it doesn't exist", func() {
			created, err := dataStore.CreateContractSchema(con.Address)
			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(Equal(true))

			created, err = dataStore.CreateEventTable(con.Address, *log)
			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(Equal(true))

			created, err = dataStore.CreateEventTable(con.Address, *log)
			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(Equal(false))
		})
	})

	Describe("PersistLog", func() {
		It("Persists contract event log values into custom tables, adding any addresses to a growing list of contract associated addresses", func() {
			err = dataStore.PersistLog(*log, con.Address, con.Name)
			Expect(err).ToNot(HaveOccurred())

			b, ok := con.TknHolderAddrs["0x000000000000000000000000000000000000Af21"]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = con.TknHolderAddrs["0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			scanLog := test_helpers.TransferLog{}

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM c%s.transfer_event", constants.TusdContractAddress)).StructScan(&scanLog)
			expectedLog := test_helpers.TransferLog{
				Id:             1,
				VulvanizeLogId: vulcanizeLogId,
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
			// Perist once
			err = dataStore.PersistLog(*log, con.Address, con.Name)
			Expect(err).ToNot(HaveOccurred())

			scanLog := test_helpers.TransferLog{}

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM c%s.transfer_event", constants.TusdContractAddress)).StructScan(&scanLog)
			expectedLog := test_helpers.TransferLog{
				Id:             1,
				VulvanizeLogId: vulcanizeLogId,
				TokenName:      "TrueUSD",
				Block:          5488076,
				Tx:             "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
				From:           "0x000000000000000000000000000000000000Af21",
				To:             "0x09BbBBE21a5975cAc061D82f7b843bCE061BA391",
				Value:          "1097077688018008265106216665536940668749033598146",
			}

			Expect(scanLog).To(Equal(expectedLog))

			// Attempt to persist the same log again
			err = dataStore.PersistLog(*log, con.Address, con.Name)
			Expect(err).ToNot(HaveOccurred())

			// Show that no new logs were entered
			var count int
			err = db.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM c%s.transfer_event", constants.TusdContractAddress))
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("Fails with empty log", func() {
			err = dataStore.PersistLog(types.Log{}, con.Address, con.Name)
			Expect(err).To(HaveOccurred())
		})
	})
})
