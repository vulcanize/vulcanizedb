package rep_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds/rep"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Rep repository", func() {
	It("returns header if matching header does not exist", func() {
		db := test_config.NewTestDB(core.Node{})
		repository := rep.NewRepRepository(db)
		pepToAdd := price_feeds.PriceUpdate{
			BlockNumber: 0,
			HeaderID:    0,
			UsdValue:    "123.456",
		}

		err := repository.CreateRep(pepToAdd)

		Expect(err).To(HaveOccurred())
	})

	It("creates a rep when matching header exists", func() {
		db := test_config.NewTestDB(core.Node{})
		repository := rep.NewRepRepository(db)
		header := core.Header{BlockNumber: 12345}
		headerRepository := repositories.NewHeaderRepository(db)
		headerID, err := headerRepository.CreateOrUpdateHeader(header)
		Expect(err).NotTo(HaveOccurred())
		pepToAdd := price_feeds.PriceUpdate{
			BlockNumber: header.BlockNumber,
			HeaderID:    headerID,
			UsdValue:    "123.456",
		}

		err = repository.CreateRep(pepToAdd)

		Expect(err).NotTo(HaveOccurred())
		var dbRep price_feeds.PriceUpdate
		err = db.Get(&dbRep, `SELECT block_number, header_id, usd_value FROM maker.reps WHERE header_id = $1`, pepToAdd.HeaderID)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbRep).To(Equal(pepToAdd))
	})
})
