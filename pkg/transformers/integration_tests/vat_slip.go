package integration_tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_slip"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat slip transformer", func() {
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

		err = persistHeader(rpcClient, db, 8953655)
		Expect(err).NotTo(HaveOccurred())
	})

	It("persists vat slip event", func(done Done) {
		config := shared.SingleTransformerConfig{
			ContractAddresses:   []string{"0xcd726790550afcd77e9a7a47e86a3f9010af126b"},
			StartingBlockNumber: 8953655,
			EndingBlockNumber:   8953655,
		}
		initializer := factories.Transformer{
			Config:     config,
			Fetcher:    &shared.Fetcher{},
			Converter:  &vat_slip.VatSlipConverter{},
			Repository: &vat_slip.VatSlipRepository{},
		}
		transformer := initializer.NewTransformer(db, blockChain)

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		var model vat_slip.VatSlipModel
		err = db.Get(&model, `SELECT ilk, guy, rad, tx_idx FROM maker.vat_slip WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.Ilk).To(Equal("ETH"))
		Expect(model.Guy).To(Equal("0xDA15dCE70ab462E66779f23ee14F21d993789eE3"))
		Expect(model.Rad).To(Equal("100000000000000000000000000000000000000000000000"))
		Expect(model.TransactionIndex).To(Equal(uint(0)))
		var headerChecked bool
		err = db.Get(&headerChecked, `SELECT vat_slip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())
		Expect(headerChecked).To(BeTrue())
		close(done)
	}, 15)
})
