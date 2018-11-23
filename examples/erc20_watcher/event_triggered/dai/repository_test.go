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

package dai_test

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered/dai"
	"github.com/vulcanize/vulcanizedb/examples/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers"
)

var transferEntity = &dai.TransferEntity{
	TokenName:    "Dai",
	TokenAddress: common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"),
	Src:          common.HexToAddress("0x000000000000000000000000000000000000Af21"),
	Dst:          common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"),
	Wad:          helpers.BigFromString("1097077688018008265106216665536940668749033598146"),
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var approvalEntity = &dai.ApprovalEntity{
	TokenName:    "Dai",
	TokenAddress: common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"),
	Src:          common.HexToAddress("0x000000000000000000000000000000000000Af21"),
	Guy:          common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"),
	Wad:          helpers.BigFromString("1097077688018008265106216665536940668749033598146"),
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var _ = Describe("Approval and Transfer Repository Tests", func() {
	var db *postgres.DB
	var converter dai.ERC20Converter
	var repository event_triggered.ERC20EventRepository
	var logRepository repositories.LogRepository
	var blockRepository repositories.BlockRepository
	var receiptRepository repositories.ReceiptRepository
	var blockNumber int64
	var blockId int64
	var vulcanizeLogId int64
	rand.Seed(time.Now().UnixNano())

	BeforeEach(func() {
		var err error
		db, err = postgres.NewDB(config.Database{
			Hostname: "localhost",
			Name:     "vulcanize_private",
			Port:     5432,
		}, core.Node{})
		Expect(err).NotTo(HaveOccurred())

		receiptRepository = repositories.ReceiptRepository{DB: db}
		logRepository = repositories.LogRepository{DB: db}
		blockRepository = *repositories.NewBlockRepository(db)

		blockNumber = rand.Int63()
		blockId = test_helpers.CreateBlock(blockNumber, blockRepository)

		log := core.Log{}
		logs := []core.Log{log}
		receipt := core.Receipt{
			Logs: logs,
		}
		receipts := []core.Receipt{receipt}

		err = receiptRepository.CreateReceiptsAndLogs(blockId, receipts)
		Expect(err).ToNot(HaveOccurred())

		err = logRepository.Get(&vulcanizeLogId, `SELECT id FROM logs`)
		Expect(err).ToNot(HaveOccurred())

		repository = event_triggered.ERC20EventRepository{DB: db}
		converter = dai.ERC20Converter{}

	})

	AfterEach(func() {
		db.Query(`DELETE FROM logs`)
		db.Query(`DELETE FROM log_filters`)
		db.Query(`DELETE FROM token_transfers`)
		db.Query(`DELETE FROM token_approvals`)

	})

	It("Creates a new Transfer record", func() {
		model := converter.ToTransferModel(transferEntity)
		err := repository.CreateTransfer(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())
		type DBRow struct {
			DBID           uint64 `db:"id"`
			VulcanizeLogID int64  `db:"vulcanize_log_id"`
			event_triggered.TransferModel
		}
		dbResult := DBRow{}

		err = repository.QueryRowx(`SELECT * FROM token_transfers`).StructScan(&dbResult)
		Expect(err).ToNot(HaveOccurred())

		Expect(dbResult.VulcanizeLogID).To(Equal(vulcanizeLogId))
		Expect(dbResult.TokenName).To(Equal(model.TokenName))
		Expect(dbResult.TokenAddress).To(Equal(model.TokenAddress))
		Expect(dbResult.To).To(Equal(model.To))
		Expect(dbResult.From).To(Equal(model.From))
		Expect(dbResult.Tokens).To(Equal(model.Tokens))
		Expect(dbResult.Block).To(Equal(model.Block))
		Expect(dbResult.TxHash).To(Equal(model.TxHash))
	})

	It("does not duplicate token_transfers that have already been seen", func() {
		model := converter.ToTransferModel(transferEntity)

		err := repository.CreateTransfer(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())
		err = repository.CreateTransfer(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())

		var count int
		err = repository.DB.QueryRowx(`SELECT count(*) FROM token_transfers`).Scan(&count)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("Removes a Transfer record when the corresponding log is removed", func() {
		var exists bool

		model := converter.ToTransferModel(transferEntity)
		err := repository.CreateTransfer(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())

		err = repository.DB.QueryRow(`SELECT exists (SELECT * FROM token_transfers WHERE vulcanize_log_id = $1)`, vulcanizeLogId).Scan(&exists)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(BeTrue())

		var logCount int
		_, err = logRepository.DB.Exec(`DELETE FROM logs WHERE id = $1`, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())
		err = logRepository.Get(&logCount, `SELECT count(*) FROM logs WHERE id = $1`, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())
		Expect(logCount).To(BeZero())

		var LogKillCount int
		err = repository.DB.QueryRowx(
			`SELECT count(*) FROM token_transfers WHERE vulcanize_log_id = $1`, vulcanizeLogId).Scan(&LogKillCount)
		Expect(err).ToNot(HaveOccurred())
		Expect(LogKillCount).To(BeZero())
	})

	It("Creates a new Approval record", func() {
		model := converter.ToApprovalModel(approvalEntity)
		err := repository.CreateApproval(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())
		type DBRow struct {
			DBID           uint64 `db:"id"`
			VulcanizeLogID int64  `db:"vulcanize_log_id"`
			event_triggered.ApprovalModel
		}
		dbResult := DBRow{}

		err = repository.QueryRowx(`SELECT * FROM token_approvals`).StructScan(&dbResult)
		Expect(err).ToNot(HaveOccurred())

		Expect(dbResult.VulcanizeLogID).To(Equal(vulcanizeLogId))
		Expect(dbResult.TokenName).To(Equal(model.TokenName))
		Expect(dbResult.TokenAddress).To(Equal(model.TokenAddress))
		Expect(dbResult.Owner).To(Equal(model.Owner))
		Expect(dbResult.Spender).To(Equal(model.Spender))
		Expect(dbResult.Tokens).To(Equal(model.Tokens))
		Expect(dbResult.Block).To(Equal(model.Block))
		Expect(dbResult.TxHash).To(Equal(model.TxHash))
	})

	It("does not duplicate token_approvals that have already been seen", func() {
		model := converter.ToApprovalModel(approvalEntity)

		err := repository.CreateApproval(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())
		err = repository.CreateApproval(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())

		var count int
		err = repository.DB.QueryRowx(`SELECT count(*) FROM token_approvals`).Scan(&count)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("Removes a Approval record when the corresponding log is removed", func() {
		var exists bool

		model := converter.ToApprovalModel(approvalEntity)
		err := repository.CreateApproval(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())

		err = repository.DB.QueryRow(`SELECT exists (SELECT * FROM token_approvals WHERE vulcanize_log_id = $1)`, vulcanizeLogId).Scan(&exists)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(BeTrue())

		var logCount int
		_, err = logRepository.DB.Exec(`DELETE FROM logs WHERE id = $1`, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())
		err = logRepository.Get(&logCount, `SELECT count(*) FROM logs WHERE id = $1`, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())
		Expect(logCount).To(BeZero())

		var LogKillCount int
		err = repository.DB.QueryRowx(
			`SELECT count(*) FROM token_approvals WHERE vulcanize_log_id = $1`, vulcanizeLogId).Scan(&LogKillCount)
		Expect(err).ToNot(HaveOccurred())
		Expect(LogKillCount).To(BeZero())
	})
})
