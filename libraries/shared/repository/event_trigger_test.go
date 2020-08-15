package repository_test

import (
	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("event updated trigger", func() {
	var (
		db       = test_config.NewTestDB(test_config.NewTestNode())
		headerID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())

	})

	type dbEvent struct {
		Created string
		Updated string
	}

	It("indicates when a record was created or updated", func() {
		var eventLogRes dbEvent
		test_data.CreateTestLog(headerID, db)
		initialEventLogErr := db.Get(&eventLogRes, `SELECT created, updated FROM public.event_logs`)
		Expect(initialEventLogErr).NotTo(HaveOccurred())
		Expect(eventLogRes.Created).To(Equal(eventLogRes.Updated))

		_, updateErr := db.Exec(`UPDATE public.event_logs SET block_hash = '{"new_block_hash"}' where header_id = $1`, headerID)
		Expect(updateErr).NotTo(HaveOccurred())
		updatedEventErr := db.Get(&eventLogRes, `SELECT created, updated FROM public.event_logs`)
		Expect(updatedEventErr).NotTo(HaveOccurred())
		Expect(eventLogRes.Created).NotTo(Equal(eventLogRes.Updated))
	})
})
