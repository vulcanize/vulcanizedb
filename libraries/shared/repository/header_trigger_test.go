package repository_test

import (
	"math/rand"

	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("header updated trigger", func() {
	var db = test_config.NewTestDB(test_config.NewTestNode())

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	type dbHeader struct {
		Created string
		Updated string
	}

	It("updates time updated when record is changed", func() {
		blockNumber := rand.Int63()
		headerRepo := repositories.NewHeaderRepository(db)
		header := fakes.GetFakeHeader(blockNumber)
		headerID, insertErr := headerRepo.CreateOrUpdateHeader(header)
		Expect(insertErr).NotTo(HaveOccurred())

		var headerRes dbHeader
		initialHeaderErr := db.Get(&headerRes, `SELECT created, updated FROM public.headers`)
		Expect(initialHeaderErr).NotTo(HaveOccurred())
		Expect(headerRes.Created).To(Equal(headerRes.Updated))

		_, updateErr := db.Exec(`UPDATE public.headers SET hash = '{"new_hash"}' WHERE id = $1`, headerID)
		Expect(updateErr).NotTo(HaveOccurred())
		updatedHeaderErr := db.Get(&headerRes, `SELECT created, updated FROM public.headers`)
		Expect(updatedHeaderErr).NotTo(HaveOccurred())
		Expect(headerRes.Created).NotTo(Equal(headerRes.Updated))
	})
})
