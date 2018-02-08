package postgres_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/repositories/postgres"
)

var _ = Describe("Logs Repository", func() {
	var repository repositories.ReceiptRepository
	var node core.Node
	BeforeEach(func() {
		node = core.Node{
			GenesisBlock: "GENESIS",
			NetworkId:    1,
			Id:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
			ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
		}
		repository = postgres.BuildRepository(node)
	})

	Describe("Saving receipts", func() {
		It("returns the receipt when it exists", func() {
			var blockRepository repositories.BlockRepository
			blockRepository = postgres.BuildRepository(node)
			expected := core.Receipt{
				ContractAddress:   "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
				CumulativeGasUsed: 7996119,
				GasUsed:           21000,
				Logs:              []core.Log{},
				StateRoot:         "0x88abf7e73128227370aa7baa3dd4e18d0af70e92ef1f9ef426942fbe2dddb733",
				Status:            1,
				TxHash:            "0xe340558980f89d5f86045ac11e5cc34e4bcec20f9f1e2a427aa39d87114e8223",
			}

			transaction := core.Transaction{
				Hash:    expected.TxHash,
				Receipt: expected,
			}

			block := core.Block{Transactions: []core.Transaction{transaction}}
			blockRepository.CreateOrUpdateBlock(block)
			receipt, err := repository.GetReceipt("0xe340558980f89d5f86045ac11e5cc34e4bcec20f9f1e2a427aa39d87114e8223")

			Expect(err).ToNot(HaveOccurred())
			//Not currently serializing bloom logs
			Expect(receipt.Bloom).To(Equal(core.Receipt{}.Bloom))
			Expect(receipt.TxHash).To(Equal(expected.TxHash))
			Expect(receipt.CumulativeGasUsed).To(Equal(expected.CumulativeGasUsed))
			Expect(receipt.GasUsed).To(Equal(expected.GasUsed))
			Expect(receipt.StateRoot).To(Equal(expected.StateRoot))
			Expect(receipt.Status).To(Equal(expected.Status))
		})

		It("returns ErrReceiptDoesNotExist when receipt does not exist", func() {
			receipt, err := repository.GetReceipt("DOES NOT EXIST")
			Expect(err).To(HaveOccurred())
			Expect(receipt).To(BeZero())
		})

		It("still saves receipts without logs", func() {
			var blockRepository repositories.BlockRepository
			blockRepository = postgres.BuildRepository(node)
			receipt := core.Receipt{
				TxHash: "0x002c4799161d809b23f67884eb6598c9df5894929fe1a9ead97ca175d360f547",
			}
			transaction := core.Transaction{
				Hash:    receipt.TxHash,
				Receipt: receipt,
			}

			block := core.Block{
				Transactions: []core.Transaction{transaction},
			}
			blockRepository.CreateOrUpdateBlock(block)

			_, err := repository.GetReceipt(receipt.TxHash)

			Expect(err).To(Not(HaveOccurred()))
		})
	})
})
