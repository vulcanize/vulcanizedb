// VulcanizeDB
// Copyright © 2019 Vulcanize

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
	"sort"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Full sync log Repository", func() {
	Describe("Saving logs", func() {
		var db *postgres.DB
		var blockRepository datastore.BlockRepository
		var logsRepository datastore.FullSyncLogRepository
		var receiptRepository datastore.FullSyncReceiptRepository
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
			logsRepository = repositories.FullSyncLogRepository{DB: db}
			receiptRepository = repositories.FullSyncReceiptRepository{DB: db}
		})

		It("returns the log when it exists", func() {
			blockNumber := int64(12345)
			blockId, err := blockRepository.CreateOrUpdateBlock(core.Block{Number: blockNumber})
			Expect(err).NotTo(HaveOccurred())
			tx, _ := db.Beginx()
			receiptId, err := receiptRepository.CreateFullSyncReceiptInTx(blockId, core.Receipt{}, tx)
			tx.Commit()
			Expect(err).NotTo(HaveOccurred())
			err = logsRepository.CreateLogs([]core.FullSyncLog{{
				BlockNumber: blockNumber,
				Index:       0,
				Address:     "x123",
				TxHash:      "x456",
				Topics:      core.Topics{0: "x777", 1: "x888", 2: "x999"},
				Data:        "xabc",
			}}, receiptId)
			Expect(err).NotTo(HaveOccurred())

			log, err := logsRepository.GetLogs("x123", blockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(log).NotTo(BeNil())
			Expect(log[0].BlockNumber).To(Equal(blockNumber))
			Expect(log[0].Address).To(Equal("x123"))
			Expect(log[0].Index).To(Equal(int64(0)))
			Expect(log[0].TxHash).To(Equal("x456"))
			Expect(log[0].Topics[0]).To(Equal("x777"))
			Expect(log[0].Topics[1]).To(Equal("x888"))
			Expect(log[0].Topics[2]).To(Equal("x999"))
			Expect(log[0].Data).To(Equal("xabc"))
		})

		It("returns nil if log does not exist", func() {
			log, err := logsRepository.GetLogs("x123", 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(log).To(BeNil())
		})

		It("filters to the correct block number and address", func() {
			blockNumber := int64(12345)
			blockId, err := blockRepository.CreateOrUpdateBlock(core.Block{Number: blockNumber})
			Expect(err).NotTo(HaveOccurred())
			tx, _ := db.Beginx()
			receiptId, err := receiptRepository.CreateFullSyncReceiptInTx(blockId, core.Receipt{}, tx)
			tx.Commit()
			Expect(err).NotTo(HaveOccurred())

			err = logsRepository.CreateLogs([]core.FullSyncLog{{
				BlockNumber: blockNumber,
				Index:       0,
				Address:     "x123",
				TxHash:      "x456",
				Topics:      core.Topics{0: "x777", 1: "x888", 2: "x999"},
				Data:        "xabc",
			}}, receiptId)
			Expect(err).NotTo(HaveOccurred())

			err = logsRepository.CreateLogs([]core.FullSyncLog{{
				BlockNumber: blockNumber,
				Index:       1,
				Address:     "x123",
				TxHash:      "x789",
				Topics:      core.Topics{0: "x111", 1: "x222", 2: "x333"},
				Data:        "xdef",
			}}, receiptId)
			Expect(err).NotTo(HaveOccurred())

			err = logsRepository.CreateLogs([]core.FullSyncLog{{
				BlockNumber: 2,
				Index:       0,
				Address:     "x123",
				TxHash:      "x456",
				Topics:      core.Topics{0: "x777", 1: "x888", 2: "x999"},
				Data:        "xabc",
			}}, receiptId)
			Expect(err).NotTo(HaveOccurred())

			log, err := logsRepository.GetLogs("x123", blockNumber)
			Expect(err).NotTo(HaveOccurred())

			type logIndex struct {
				blockNumber int64
				Index       int64
			}

			var uniqueBlockNumbers []logIndex
			for _, log := range log {
				uniqueBlockNumbers = append(uniqueBlockNumbers,
					logIndex{log.BlockNumber, log.Index})
			}
			sort.Slice(uniqueBlockNumbers, func(i, j int) bool {
				if uniqueBlockNumbers[i].blockNumber < uniqueBlockNumbers[j].blockNumber {
					return true
				}
				if uniqueBlockNumbers[i].blockNumber > uniqueBlockNumbers[j].blockNumber {
					return false
				}
				return uniqueBlockNumbers[i].Index < uniqueBlockNumbers[j].Index
			})

			Expect(log).NotTo(BeNil())
			Expect(len(log)).To(Equal(2))
			Expect(uniqueBlockNumbers).To(Equal(
				[]logIndex{
					{blockNumber: blockNumber, Index: 0},
					{blockNumber: blockNumber, Index: 1}},
			))
		})

		It("saves the logs attached to a receipt", func() {
			logs := []core.FullSyncLog{{
				Address:     "0x8a4774fe82c63484afef97ca8d89a6ea5e21f973",
				BlockNumber: 4745407,
				Data:        "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000645a68669900000000000000000000000000000000000000000000003397684ab5869b0000000000000000000000000000000000000000000000000000000000005a36053200000000000000000000000099041f808d598b782d5a3e498681c2452a31da08",
				Index:       86,
				Topics: core.Topics{
					0: "0x5a68669900000000000000000000000000000000000000000000000000000000",
					1: "0x000000000000000000000000d0148dad63f73ce6f1b6c607e3413dcf1ff5f030",
					2: "0x00000000000000000000000000000000000000000000003397684ab5869b0000",
					3: "0x000000000000000000000000000000000000000000000000000000005a360532",
				},
				TxHash: "0x002c4799161d809b23f67884eb6598c9df5894929fe1a9ead97ca175d360f547",
			}, {
				Address:     "0x99041f808d598b782d5a3e498681c2452a31da08",
				BlockNumber: 4745407,
				Data:        "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000418178358",
				Index:       87,
				Topics: core.Topics{
					0: "0x1817835800000000000000000000000000000000000000000000000000000000",
					1: "0x0000000000000000000000008a4774fe82c63484afef97ca8d89a6ea5e21f973",
					2: "0x0000000000000000000000000000000000000000000000000000000000000000",
					3: "0x0000000000000000000000000000000000000000000000000000000000000000",
				},
				TxHash: "0x002c4799161d809b23f67884eb6598c9df5894929fe1a9ead97ca175d360f547",
			}, {
				Address:     "0x99041f808d598b782d5a3e498681c2452a31da08",
				BlockNumber: 4745407,
				Data:        "0x00000000000000000000000000000000000000000000003338f64c8423af4000",
				Index:       88,
				Topics: core.Topics{
					0: "0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527",
				},
				TxHash: "0x002c4799161d809b23f67884eb6598c9df5894929fe1a9ead97ca175d360f547",
			},
			}
			receipt := core.Receipt{
				ContractAddress:   "",
				CumulativeGasUsed: 7481414,
				GasUsed:           60711,
				Logs:              logs,
				Bloom:             "0x00000800000000000000001000000000000000400000000080000000000000000000400000010000000000000000000000000000040000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000800004000000000000001000000000000000000000000000002000000480000000000000002000000000000000020000000000000000000000000000000000000000080000000000180000c00000000000000002000002000000040000000000000000000000000000010000000000020000000000000000000002000000000000000000000000400800000000000000000",
				Status:            1,
				TxHash:            "0x002c4799161d809b23f67884eb6598c9df5894929fe1a9ead97ca175d360f547",
			}
			transaction := fakes.GetFakeTransaction(receipt.TxHash, receipt)

			block := core.Block{Transactions: []core.TransactionModel{transaction}}
			_, err := blockRepository.CreateOrUpdateBlock(block)
			Expect(err).To(Not(HaveOccurred()))
			retrievedLogs, err := logsRepository.GetLogs("0x99041f808d598b782d5a3e498681c2452a31da08", 4745407)

			Expect(err).NotTo(HaveOccurred())
			expected := logs[1:]
			Expect(retrievedLogs).To(ConsistOf(expected))
		})
	})
})
