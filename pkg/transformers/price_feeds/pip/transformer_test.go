package pip_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds/pip"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Pip transformer", func() {
	It("returns nil if no logs found", func() {
		chain := fakes.NewMockBlockChain()
		db := test_config.NewTestDB(core.Node{})
		transformer := pip.NewPipTransformer(chain, db, "pip-contract-address")

		err := transformer.Execute(core.Header{}, 123)

		Expect(err).NotTo(HaveOccurred())
	})

	It("creates pip row for found log", func() {
		chain := fakes.NewMockBlockChain()
		chain.SetGetEthLogsWithCustomQueryReturnLogs([]types.Log{{Data: []byte{1, 2, 3, 4, 5}}})
		db := test_config.NewTestDB(core.Node{})
		test_config.CleanTestDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		header := core.Header{BlockNumber: 12345}
		headerID, err := headerRepository.CreateOrUpdateHeader(header)
		Expect(err).NotTo(HaveOccurred())
		transformer := pip.NewPipTransformer(chain, db, "pip-contract-address")

		err = transformer.Execute(header, headerID)

		Expect(err).NotTo(HaveOccurred())
		var dbPip price_feeds.PriceUpdate
		err = db.Get(&dbPip, `SELECT block_number, header_id, usd_value FROM maker.pips WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPip.BlockNumber).To(Equal(header.BlockNumber))
	})
})
