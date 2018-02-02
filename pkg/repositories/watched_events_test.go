package repositories_test

import (
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/repositories/testing"

	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

var _ = Describe("Watched Events Repository", func() {
	var repository repositories.Postgres

	BeforeEach(func() {
		cfg, err := config.NewConfig("private")
		if err != nil {
			log.Fatal(err)
		}
		repository, err = repositories.NewPostgres(cfg.Database, core.Node{})
		if err != nil {
			log.Fatal(err)
		}
		testing.ClearData(repository)
	})

	It("retrieves watched logs that match the event filter", func() {
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
		expectedWatchedEventLog := []*repositories.WatchedEventLog{
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
		err := repository.AddFilter(filter)
		Expect(err).ToNot(HaveOccurred())
		err = repository.CreateLogs(logs)
		Expect(err).ToNot(HaveOccurred())
		matchingLogs, err := repository.AllWatchedEventLogs()
		Expect(err).ToNot(HaveOccurred())
		Expect(matchingLogs).To(Equal(expectedWatchedEventLog))

	})
})
