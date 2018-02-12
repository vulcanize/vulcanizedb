package postgres_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/pkg/repositories/postgres"
)

var _ = Describe("Watched Events Repository", func() {
	var db *postgres.DB
	var logRepository postgres.LogRepository
	var filterRepository postgres.FilterRepository
	var watchedEventRepository postgres.WatchedEventRepository

	BeforeEach(func() {
		db = postgres.NewTestDB(core.Node{})
		logRepository = postgres.LogRepository{DB: db}
		filterRepository = postgres.FilterRepository{DB: db}
		watchedEventRepository = postgres.WatchedEventRepository{DB: db}
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
		err = logRepository.CreateLogs(logs)
		Expect(err).ToNot(HaveOccurred())
		matchingLogs, err := watchedEventRepository.GetWatchedEvents("Filter1")
		Expect(err).ToNot(HaveOccurred())
		Expect(matchingLogs).To(Equal(expectedWatchedEventLog))

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
		err = logRepository.CreateLogs(logs)
		Expect(err).ToNot(HaveOccurred())
		matchingLogs, err := watchedEventRepository.GetWatchedEvents("Filter1")
		Expect(err).ToNot(HaveOccurred())
		Expect(matchingLogs).To(Equal(expectedWatchedEventLog))

	})
})
