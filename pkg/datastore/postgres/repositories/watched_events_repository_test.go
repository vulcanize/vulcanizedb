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

package repositories_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Watched Events Repository", func() {
	var db *postgres.DB
	var blocksRepository datastore.BlockRepository
	var filterRepository datastore.FilterRepository
	var logRepository datastore.LogRepository
	var receiptRepository datastore.ReceiptRepository
	var watchedEventRepository datastore.WatchedEventRepository

	BeforeEach(func() {
		db = test_config.NewTestDB(core.Node{})
		test_config.CleanTestDB(db)
		blocksRepository = repositories.NewBlockRepository(db)
		filterRepository = repositories.FilterRepository{DB: db}
		logRepository = repositories.LogRepository{DB: db}
		receiptRepository = repositories.ReceiptRepository{DB: db}
		watchedEventRepository = repositories.WatchedEventRepository{DB: db}
	})

	It("retrieves watched event logs that match the event filter", func() {
		filter := filters.LogFilter{
			Name:      "Filter1",
			FromBlock: 0,
			ToBlock:   10,
			Address:   "0x123",
			Topics:    core.Topics{0: "event1=10", 2: "event3=hello"},
		}
		logs := []core.Log{
			{
				BlockNumber: 0,
				TxHash:      "0x1",
				Address:     "0x123",
				Topics:      core.Topics{0: "event1=10", 2: "event3=hello"},
				Index:       0,
				Data:        "",
			},
		}
		expectedWatchedEventLog := []*core.WatchedEvent{
			{
				Name:        "Filter1",
				BlockNumber: 0,
				TxHash:      "0x1",
				Address:     "0x123",
				Topic0:      "event1=10",
				Topic2:      "event3=hello",
				Index:       0,
				Data:        "",
			},
		}
		err := filterRepository.CreateFilter(filter)
		Expect(err).ToNot(HaveOccurred())
		blockId, err := blocksRepository.CreateOrUpdateBlock(core.Block{})
		Expect(err).NotTo(HaveOccurred())
		receiptId, err := receiptRepository.CreateReceipt(blockId, core.Receipt{})
		Expect(err).NotTo(HaveOccurred())
		err = logRepository.CreateLogs(logs, receiptId)
		Expect(err).ToNot(HaveOccurred())
		matchingLogs, err := watchedEventRepository.GetWatchedEvents("Filter1")
		Expect(err).ToNot(HaveOccurred())
		Expect(len(matchingLogs)).To(Equal(1))
		Expect(matchingLogs[0].Name).To(Equal(expectedWatchedEventLog[0].Name))
		Expect(matchingLogs[0].BlockNumber).To(Equal(expectedWatchedEventLog[0].BlockNumber))
		Expect(matchingLogs[0].TxHash).To(Equal(expectedWatchedEventLog[0].TxHash))
		Expect(matchingLogs[0].Address).To(Equal(expectedWatchedEventLog[0].Address))
		Expect(matchingLogs[0].Topic0).To(Equal(expectedWatchedEventLog[0].Topic0))
		Expect(matchingLogs[0].Topic1).To(Equal(expectedWatchedEventLog[0].Topic1))
		Expect(matchingLogs[0].Topic2).To(Equal(expectedWatchedEventLog[0].Topic2))
		Expect(matchingLogs[0].Data).To(Equal(expectedWatchedEventLog[0].Data))
	})

	It("retrieves a watched event log by name", func() {
		filter := filters.LogFilter{
			Name:      "Filter1",
			FromBlock: 0,
			ToBlock:   10,
			Address:   "0x123",
			Topics:    core.Topics{0: "event1=10", 2: "event3=hello"},
		}
		logs := []core.Log{
			{
				BlockNumber: 0,
				TxHash:      "0x1",
				Address:     "0x123",
				Topics:      core.Topics{0: "event1=10", 2: "event3=hello"},
				Index:       0,
				Data:        "",
			},
			{
				BlockNumber: 100,
				TxHash:      "",
				Address:     "",
				Topics:      core.Topics{},
				Index:       0,
				Data:        "",
			},
		}
		expectedWatchedEventLog := []*core.WatchedEvent{{
			Name:        "Filter1",
			BlockNumber: 0,
			TxHash:      "0x1",
			Address:     "0x123",
			Topic0:      "event1=10",
			Topic2:      "event3=hello",
			Index:       0,
			Data:        "",
		}}
		err := filterRepository.CreateFilter(filter)
		Expect(err).ToNot(HaveOccurred())
		blockId, err := blocksRepository.CreateOrUpdateBlock(core.Block{Hash: "Ox123"})
		Expect(err).NotTo(HaveOccurred())
		receiptId, err := receiptRepository.CreateReceipt(blockId, core.Receipt{TxHash: "0x123"})
		Expect(err).NotTo(HaveOccurred())
		err = logRepository.CreateLogs(logs, receiptId)
		Expect(err).ToNot(HaveOccurred())
		matchingLogs, err := watchedEventRepository.GetWatchedEvents("Filter1")
		Expect(err).ToNot(HaveOccurred())
		Expect(len(matchingLogs)).To(Equal(1))
		Expect(matchingLogs[0].Name).To(Equal(expectedWatchedEventLog[0].Name))
		Expect(matchingLogs[0].BlockNumber).To(Equal(expectedWatchedEventLog[0].BlockNumber))
		Expect(matchingLogs[0].TxHash).To(Equal(expectedWatchedEventLog[0].TxHash))
		Expect(matchingLogs[0].Address).To(Equal(expectedWatchedEventLog[0].Address))
		Expect(matchingLogs[0].Topic0).To(Equal(expectedWatchedEventLog[0].Topic0))
		Expect(matchingLogs[0].Topic1).To(Equal(expectedWatchedEventLog[0].Topic1))
		Expect(matchingLogs[0].Topic2).To(Equal(expectedWatchedEventLog[0].Topic2))
		Expect(matchingLogs[0].Data).To(Equal(expectedWatchedEventLog[0].Data))
	})
})
