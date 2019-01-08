package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
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
	})

	It("persists vat slip event", func() {
		blockNumber := int64(8953655)
		config := vat_slip.VatSlipConfig
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		fetcher := shared.NewFetcher(blockChain)
		logs, err := fetcher.FetchLogs(
			shared.HexStringsToAddresses(config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &vat_slip.VatSlipConverter{},
			Repository: &vat_slip.VatSlipRepository{},
		}.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header)

		Expect(err).NotTo(HaveOccurred())
		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
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
	})
})
