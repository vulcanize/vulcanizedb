// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event_triggered_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/examples/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"math/rand"
	"time"
)

var transferLog = core.Log{
	BlockNumber: 5488076,
	Address:     constants.DaiContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topics: [4]string{
		constants.TransferEventSignature,
		"0x000000000000000000000000000000000000000000000000000000000000af21",
		"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
		"",
	},
	Data: "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var approvalLog = core.Log{
	BlockNumber: 5488076,
	Address:     constants.DaiContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topics: [4]string{
		constants.ApprovalEventSignature,
		"0x000000000000000000000000000000000000000000000000000000000000af21",
		"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
		"",
	},
	Data: "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

//converted transfer to assert against
var logs = []core.Log{
	transferLog,
	approvalLog,
	{
		BlockNumber: 0,
		TxHash:      "",
		Address:     "",
		Topics:      core.Topics{},
		Index:       0,
		Data:        "",
	},
}

var _ = Describe("Integration test with vulcanizedb", func() {
	var db *postgres.DB
	rand.Seed(time.Now().UnixNano())

	BeforeEach(func() {
		var err error
		db, err = postgres.NewDB(config.Database{
			Hostname: "localhost",
			Name:     "vulcanize_private",
			Port:     5432,
		}, core.Node{})
		Expect(err).NotTo(HaveOccurred())

		receiptRepository := repositories.ReceiptRepository{DB: db}
		blockRepository := *repositories.NewBlockRepository(db)

		blockNumber := rand.Int63()
		blockId := test_helpers.CreateBlock(blockNumber, blockRepository)

		receipt := core.Receipt{
			Logs: logs,
		}
		receipts := []core.Receipt{receipt}

		err = receiptRepository.CreateReceiptsAndLogs(blockId, receipts)
		Expect(err).ToNot(HaveOccurred())

		var vulcanizeLogIds []int64
		err = db.Select(&vulcanizeLogIds, `SELECT id FROM logs`)
		Expect(err).ToNot(HaveOccurred())

	})

	AfterEach(func() {
		_, err := db.Exec(`DELETE FROM token_transfers`)
		Expect(err).ToNot(HaveOccurred())
		_, err = db.Exec(`DELETE FROM token_approvals`)
		Expect(err).ToNot(HaveOccurred())
		_, err = db.Exec(`DELETE FROM log_filters`)
		Expect(err).ToNot(HaveOccurred())
		_, err = db.Exec(`DELETE FROM logs`)
		Expect(err).ToNot(HaveOccurred())
	})

	It("creates transfer entry for each Transfer event received", func() {
		transformer := event_triggered.NewTransformer(db, generic.DaiConfig)

		transformer.Execute()

		var count int
		err := db.QueryRow(`SELECT COUNT(*) FROM token_transfers`).Scan(&count)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(1))

		type dbRow struct {
			DBID           uint64 `db:"id"`
			VulcanizeLogID int64  `db:"vulcanize_log_id"`
			event_triggered.TransferModel
		}
		var transfer dbRow
		err = db.Get(&transfer, `SELECT * FROM token_transfers WHERE block=$1`, logs[0].BlockNumber)
		Expect(err).ToNot(HaveOccurred())
		Expect(transfer.TokenName).To(Equal(expectedTransferModel.TokenName))
		Expect(transfer.TokenAddress).To(Equal(expectedTransferModel.TokenAddress))
		Expect(transfer.To).To(Equal(expectedTransferModel.To))
		Expect(transfer.From).To(Equal(expectedTransferModel.From))
		Expect(transfer.Tokens).To(Equal(expectedTransferModel.Tokens))
		Expect(transfer.Block).To(Equal(expectedTransferModel.Block))
		Expect(transfer.TxHash).To(Equal(expectedTransferModel.TxHash))
	})

	It("creates approval entry for each Approval event received", func() {
		transformer := event_triggered.NewTransformer(db, generic.DaiConfig)

		transformer.Execute()

		var count int
		err := db.QueryRow(`SELECT COUNT(*) FROM token_approvals`).Scan(&count)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(1))

		type dbRow struct {
			DBID           uint64 `db:"id"`
			VulcanizeLogID int64  `db:"vulcanize_log_id"`
			event_triggered.ApprovalModel
		}
		var transfer dbRow
		err = db.Get(&transfer, `SELECT * FROM token_approvals WHERE block=$1`, logs[0].BlockNumber)
		Expect(err).ToNot(HaveOccurred())
		Expect(transfer.TokenName).To(Equal(expectedApprovalModel.TokenName))
		Expect(transfer.TokenAddress).To(Equal(expectedApprovalModel.TokenAddress))
		Expect(transfer.Owner).To(Equal(expectedApprovalModel.Owner))
		Expect(transfer.Spender).To(Equal(expectedApprovalModel.Spender))
		Expect(transfer.Tokens).To(Equal(expectedApprovalModel.Tokens))
		Expect(transfer.Block).To(Equal(expectedApprovalModel.Block))
		Expect(transfer.TxHash).To(Equal(expectedApprovalModel.TxHash))
	})

})
