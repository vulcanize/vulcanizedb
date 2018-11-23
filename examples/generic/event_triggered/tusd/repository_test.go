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

package tusd_test

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/generic/event_triggered"
	"github.com/vulcanize/vulcanizedb/examples/generic/event_triggered/tusd"
	"github.com/vulcanize/vulcanizedb/examples/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers"
)

var burnEntity = &tusd.BurnEntity{
	TokenName:    "Tusd",
	TokenAddress: common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"),
	Burner:       common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"),
	Value:        helpers.BigFromString("1097077688018008265106216665536940668749033598146"),
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var mintEntity = &tusd.MintEntity{
	TokenName:    "Tusd",
	TokenAddress: common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"),
	To:           common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"),
	Amount:       helpers.BigFromString("1097077688018008265106216665536940668749033598146"),
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var _ = Describe("Approval and Transfer Repository Tests", func() {
	var db *postgres.DB
	var converter tusd.GenericConverter
	var repository event_triggered.GenericEventRepository
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

		repository = event_triggered.GenericEventRepository{DB: db}
		converter = tusd.GenericConverter{}
	})

	AfterEach(func() {
		db.Query(`DELETE FROM logs`)
		db.Query(`DELETE FROM log_filters`)
		db.Query(`DELETE FROM token_burns`)
		db.Query(`DELETE FROM token_mints`)

	})

	It("Creates a new Burn record", func() {
		model := converter.ToBurnModel(burnEntity)
		err := repository.CreateBurn(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())
		type DBRow struct {
			DBID           uint64 `db:"id"`
			VulcanizeLogID int64  `db:"vulcanize_log_id"`
			event_triggered.BurnModel
		}
		dbResult := DBRow{}

		err = repository.QueryRowx(`SELECT * FROM token_burns`).StructScan(&dbResult)
		Expect(err).ToNot(HaveOccurred())

		Expect(dbResult.VulcanizeLogID).To(Equal(vulcanizeLogId))
		Expect(dbResult.TokenName).To(Equal(model.TokenName))
		Expect(dbResult.TokenAddress).To(Equal(model.TokenAddress))
		Expect(dbResult.Burner).To(Equal(model.Burner))
		Expect(dbResult.Tokens).To(Equal(model.Tokens))
		Expect(dbResult.Block).To(Equal(model.Block))
		Expect(dbResult.TxHash).To(Equal(model.TxHash))
	})

	It("does not duplicate token_transfers that have already been seen", func() {
		model := converter.ToBurnModel(burnEntity)

		err := repository.CreateBurn(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())
		err = repository.CreateBurn(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())

		var count int
		err = repository.DB.QueryRowx(`SELECT count(*) FROM token_burns`).Scan(&count)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("Removes a Burn record when the corresponding log is removed", func() {
		var exists bool

		model := converter.ToBurnModel(burnEntity)
		err := repository.CreateBurn(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())

		err = repository.DB.QueryRow(`SELECT exists (SELECT * FROM token_burns WHERE vulcanize_log_id = $1)`, vulcanizeLogId).Scan(&exists)
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
			`SELECT count(*) FROM token_burns WHERE vulcanize_log_id = $1`, vulcanizeLogId).Scan(&LogKillCount)
		Expect(err).ToNot(HaveOccurred())
		Expect(LogKillCount).To(BeZero())
	})

	It("Creates a new Mint record", func() {
		model := converter.ToMintModel(mintEntity)
		err := repository.CreateMint(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())
		type DBRow struct {
			DBID           uint64 `db:"id"`
			VulcanizeLogID int64  `db:"vulcanize_log_id"`
			event_triggered.MintModel
		}
		dbResult := DBRow{}

		err = repository.QueryRowx(`SELECT * FROM token_mints`).StructScan(&dbResult)
		Expect(err).ToNot(HaveOccurred())

		Expect(dbResult.VulcanizeLogID).To(Equal(vulcanizeLogId))
		Expect(dbResult.TokenName).To(Equal(model.TokenName))
		Expect(dbResult.TokenAddress).To(Equal(model.TokenAddress))
		Expect(dbResult.Mintee).To(Equal(model.Mintee))
		Expect(dbResult.Minter).To(Equal(model.Minter))
		Expect(dbResult.Tokens).To(Equal(model.Tokens))
		Expect(dbResult.Block).To(Equal(model.Block))
		Expect(dbResult.TxHash).To(Equal(model.TxHash))
	})

	It("does not duplicate token_mints that have already been seen", func() {
		model := converter.ToMintModel(mintEntity)

		err := repository.CreateMint(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())
		err = repository.CreateMint(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())

		var count int
		err = repository.DB.QueryRowx(`SELECT count(*) FROM token_mints`).Scan(&count)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("Removes a Mint record when the corresponding log is removed", func() {
		var exists bool

		model := converter.ToMintModel(mintEntity)
		err := repository.CreateMint(model, vulcanizeLogId)
		Expect(err).ToNot(HaveOccurred())

		err = repository.DB.QueryRow(`SELECT exists (SELECT * FROM token_mints WHERE vulcanize_log_id = $1)`, vulcanizeLogId).Scan(&exists)
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
			`SELECT count(*) FROM token_mints WHERE vulcanize_log_id = $1`, vulcanizeLogId).Scan(&LogKillCount)
		Expect(err).ToNot(HaveOccurred())
		Expect(LogKillCount).To(BeZero())
	})
})
