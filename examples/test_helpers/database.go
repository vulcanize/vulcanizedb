package test_helpers

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/test_config"
)

type TransferDBRow struct {
	ID             int64 `db:"id"`
	VulcanizeLogID int64 `db:"vulcanize_log_id"`
}

func CreateLogRecord(db *postgres.DB, logRepository repositories.LogRepository, log core.Log) {
	blockRepository := repositories.NewBlockRepository(db)
	receiptRepository := repositories.ReceiptRepository{DB: db}

	blockNumber := log.BlockNumber
	blockId, err := blockRepository.CreateOrUpdateBlock(core.Block{Number: blockNumber})
	Expect(err).NotTo(HaveOccurred())

	receiptId, err := receiptRepository.CreateReceipt(blockId, core.Receipt{})
	Expect(err).NotTo(HaveOccurred())

	err = logRepository.CreateLogs([]core.Log{log}, receiptId)
	Expect(err).NotTo(HaveOccurred())
}

func CreateNewDatabase() *postgres.DB {
	var node core.Node
	node = core.Node{
		GenesisBlock: "GENESIS",
		NetworkID:    1,
		ID:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
		ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
	}
	db := test_config.NewTestDB(node)

	_, err := db.Exec(`DELETE FROM logs`)
	Expect(err).NotTo(HaveOccurred())

	return db
}
