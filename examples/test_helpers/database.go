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

package test_helpers

import (
	"math/rand"
	"time"

	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/test_config"
)

type TokenSupplyDBRow struct {
	ID           int64  `db:"id"`
	Supply       int64  `db:"supply"`
	BlockID      int64  `db:"block_id"`
	TokenAddress string `db:"token_address"`
}

type TokenBalanceDBRow struct {
	ID                 int64  `db:"id"`
	Balance            int64  `db:"balance"`
	BlockID            int64  `db:"block_id"`
	TokenAddress       string `db:"token_address"`
	TokenHolderAddress string `db:"token_holder_address"`
}

type TokenAllowanceDBRow struct {
	ID                  int64  `db:"id"`
	Allowance           int64  `db:"allowance"`
	BlockID             int64  `db:"block_id"`
	TokenAddress        string `db:"token_address"`
	TokenHolderAddress  string `db:"token_holder_address"`
	TokenSpenderAddress string `db:"token_spender_address"`
}

type TransferDBRow struct {
	ID             int64 `db:"id"`
	VulcanizeLogID int64 `db:"vulcanize_log_id"`
}

func CreateNewDatabase() *postgres.DB {
	var node core.Node
	node = core.Node{
		GenesisBlock: "GENESIS",
		NetworkID:    1,
		ID:           "2ea672a45c4c7b96e3c4b130b21a22af390a552fd0b3cff96420b4bda26568d470dc56e05e453823f64f2556a6e4460ad1d4d00eb2d8b8fc16fcb1be73e86522",
		ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
	}
	db := test_config.NewTestDB(node)

	_, err := db.Exec(`DELETE FROM logs`)
	Expect(err).NotTo(HaveOccurred())

	return db
}

func CreateBlock(blockNumber int64, repository repositories.BlockRepository) (blockId int64) {
	blockId, err := repository.CreateOrUpdateBlock(core.Block{Number: blockNumber})
	Expect(err).NotTo(HaveOccurred())

	return blockId
}

func SetupIntegrationDB(db *postgres.DB, logs []core.Log) *postgres.DB {

	rand.Seed(time.Now().UnixNano())

	db, err := postgres.NewDB(config.Database{
		Hostname: "localhost",
		Name:     "vulcanize_private",
		Port:     5432,
	}, core.Node{})
	Expect(err).NotTo(HaveOccurred())

	receiptRepository := repositories.ReceiptRepository{DB: db}
	blockRepository := *repositories.NewBlockRepository(db)

	blockNumber := rand.Int63()
	blockId := CreateBlock(blockNumber, blockRepository)

	receipt := core.Receipt{
		Logs: logs,
	}
	receipts := []core.Receipt{receipt}

	err = receiptRepository.CreateReceiptsAndLogs(blockId, receipts)
	Expect(err).NotTo(HaveOccurred())

	var vulcanizeLogIds []int64
	err = db.Select(&vulcanizeLogIds, `SELECT id FROM logs`)
	Expect(err).NotTo(HaveOccurred())

	return db
}

func TearDownIntegrationDB(db *postgres.DB) *postgres.DB {

	_, err := db.Exec(`DELETE FROM token_transfers`)
	Expect(err).NotTo(HaveOccurred())

	_, err = db.Exec(`DELETE FROM token_approvals`)
	Expect(err).NotTo(HaveOccurred())

	_, err = db.Exec(`DELETE FROM log_filters`)
	Expect(err).NotTo(HaveOccurred())

	_, err = db.Exec(`DELETE FROM logs`)
	Expect(err).NotTo(HaveOccurred())

	return db
}
