package rep_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds/rep"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Rep transformer", func() {
	It("returns nil if no logs found", func() {
		chain := fakes.NewMockBlockChain()
		db := test_config.NewTestDB(core.Node{})
		transformer := rep.NewRepTransformer(chain, db)

		err := transformer.Execute(core.Header{}, 123)

		Expect(err).NotTo(HaveOccurred())
	})

	It("creates rep row for found log", func() {
		chain := fakes.NewMockBlockChain()
		chain.SetGetLogsReturnLogs([]core.Log{{Data: common.ToHex([]byte{1, 2, 3, 4, 5})}})
		db := test_config.NewTestDB(core.Node{})
		headerRepository := repositories.NewHeaderRepository(db)
		header := core.Header{BlockNumber: 12345}
		headerID, err := headerRepository.CreateOrUpdateHeader(header)
		Expect(err).NotTo(HaveOccurred())
		transformer := rep.NewRepTransformer(chain, db)

		err = transformer.Execute(header, headerID)

		Expect(err).NotTo(HaveOccurred())
		var dbRep price_feeds.PriceUpdate
		err = db.Get(&dbRep, `SELECT block_number, header_id, usd_value FROM maker.reps WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbRep.BlockNumber).To(Equal(header.BlockNumber))
	})
})
