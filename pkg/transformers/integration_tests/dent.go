package integration_tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Dent transformer", func() {
	var (
		db         *postgres.DB
		blockChain core.BlockChain
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)
	})

	It("persists a flop dent log event", func() {
		blockNumber := int64(8955613)
		err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		config := dent.DentConfig
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &dent.DentConverter{},
			Repository: &dent.DentRepository{},
			Fetcher:    &shared.Fetcher{},
		}
		transformer := initializer.NewLogNoteTransformer(db, blockChain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResult []dent.DentModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, guy, lot FROM maker.dent`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("10000000000000000000000"))
		Expect(dbResult[0].BidId).To(Equal("2"))
		Expect(dbResult[0].Guy).To(Equal("0x0000d8b4147eDa80Fec7122AE16DA2479Cbd7ffB"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000000000000"))
	})

	It("persists a flip dent log event", func() {
		//TODO: There are currently no Flip.dent events on Kovan
	})
})
