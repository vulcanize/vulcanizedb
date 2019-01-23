// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package repositories_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Receipts Repository", func() {
	var blockRepository datastore.BlockRepository
	var logRepository datastore.LogRepository
	var receiptRepository datastore.ReceiptRepository
	var db *postgres.DB
	var node core.Node
	BeforeEach(func() {
		node = core.Node{
			GenesisBlock: "GENESIS",
			NetworkID:    1,
			ID:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
			ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
		}
		db = test_config.NewTestDB(node)
		test_config.CleanTestDB(db)
		blockRepository = repositories.NewBlockRepository(db)
		logRepository = repositories.LogRepository{DB: db}
		receiptRepository = repositories.ReceiptRepository{DB: db}
	})

	Describe("Saving multiple receipts", func() {
		It("persists each receipt and its logs", func() {
			blockNumber := int64(1234567)
			blockId, err := blockRepository.CreateOrUpdateBlock(core.Block{Number: blockNumber})
			Expect(err).NotTo(HaveOccurred())
			txHashOne := "0xTxHashOne"
			txHashTwo := "0xTxHashTwo"
			addressOne := "0xAddressOne"
			addressTwo := "0xAddressTwo"
			logsOne := []core.Log{{
				Address:     addressOne,
				BlockNumber: blockNumber,
				TxHash:      txHashOne,
			}, {
				Address:     addressOne,
				BlockNumber: blockNumber,
				TxHash:      txHashOne,
			}}
			logsTwo := []core.Log{{
				BlockNumber: blockNumber,
				TxHash:      txHashTwo,
				Address:     addressTwo,
			}}
			receiptOne := core.Receipt{
				Logs:   logsOne,
				TxHash: txHashOne,
			}
			receiptTwo := core.Receipt{
				Logs:   logsTwo,
				TxHash: txHashTwo,
			}
			receipts := []core.Receipt{receiptOne, receiptTwo}

			err = receiptRepository.CreateReceiptsAndLogs(blockId, receipts)

			Expect(err).NotTo(HaveOccurred())

			persistedReceiptOne, err := receiptRepository.GetReceipt(txHashOne)
			Expect(err).NotTo(HaveOccurred())
			Expect(persistedReceiptOne).NotTo(BeNil())
			Expect(persistedReceiptOne.TxHash).To(Equal(txHashOne))
			persistedReceiptTwo, err := receiptRepository.GetReceipt(txHashTwo)
			Expect(err).NotTo(HaveOccurred())
			Expect(persistedReceiptTwo).NotTo(BeNil())
			Expect(persistedReceiptTwo.TxHash).To(Equal(txHashTwo))
			persistedAddressOneLogs := logRepository.GetLogs(addressOne, blockNumber)
			Expect(persistedAddressOneLogs).NotTo(BeNil())
			Expect(len(persistedAddressOneLogs)).To(Equal(2))
			persistedAddressTwoLogs := logRepository.GetLogs(addressTwo, blockNumber)
			Expect(persistedAddressTwoLogs).NotTo(BeNil())
			Expect(len(persistedAddressTwoLogs)).To(Equal(1))
		})
	})

	Describe("Saving receipts on a block's transactions", func() {
		It("returns the receipt when it exists", func() {
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

			_, err := blockRepository.CreateOrUpdateBlock(block)

			Expect(err).NotTo(HaveOccurred())
			receipt, err := receiptRepository.GetReceipt("0xe340558980f89d5f86045ac11e5cc34e4bcec20f9f1e2a427aa39d87114e8223")
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
			receipt, err := receiptRepository.GetReceipt("DOES NOT EXIST")
			Expect(err).To(HaveOccurred())
			Expect(receipt).To(BeZero())
		})

		It("still saves receipts without logs", func() {
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

			_, err := blockRepository.CreateOrUpdateBlock(block)

			Expect(err).NotTo(HaveOccurred())
			_, err = receiptRepository.GetReceipt(receipt.TxHash)
			Expect(err).To(Not(HaveOccurred()))
		})
	})
})
