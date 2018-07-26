package pip_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds/pip"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Pip repository", func() {
	It("does not create a pip if no matching header", func() {
		db := test_config.NewTestDB(core.Node{})
		repository := pip.NewPipRepository(db)
		priceUpdate := price_feeds.PriceUpdate{
			BlockNumber: 0,
			HeaderID:    0,
			UsdValue:    "123",
		}

		err := repository.CreatePip(priceUpdate)

		Expect(err).To(HaveOccurred())
	})

	It("creates a pip when header exists", func() {
		db := test_config.NewTestDB(core.Node{})
		test_config.CleanTestDB(db)
		repository := pip.NewPipRepository(db)
		headerRepository := repositories.NewHeaderRepository(db)
		header := core.Header{BlockNumber: 12345}
		headerID, err := headerRepository.CreateOrUpdateHeader(header)
		Expect(err).NotTo(HaveOccurred())
		priceUpdate := price_feeds.PriceUpdate{
			BlockNumber: header.BlockNumber,
			HeaderID:    headerID,
			UsdValue:    "777.777",
		}

		err = repository.CreatePip(priceUpdate)

		Expect(err).NotTo(HaveOccurred())
		var dbPip price_feeds.PriceUpdate
		err = db.Get(&dbPip, `SELECT block_number, header_id, usd_value FROM maker.pips WHERE block_number = $1`, header.BlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPip).To(Equal(priceUpdate))
	})
})
