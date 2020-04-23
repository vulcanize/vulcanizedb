package test_data

import (
	"crypto/sha256"
	"fmt"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/gomega"
)

// Create an event log to reference in an event, returning inserted event log
func CreateTestLog(headerID int64, db *postgres.DB) core.EventLog {
	txHash := getRandomHash()
	log := types.Log{
		Address:     common.Address{},
		BlockNumber: 0,
		TxHash:      txHash,
		TxIndex:     uint(rand.Int31()),
		BlockHash:   common.Hash{},
		Index:       0,
	}

	tx := getFakeTransactionFromHash(txHash)
	headerRepository := repositories.NewHeaderRepository(db)
	txErr := headerRepository.CreateTransactions(headerID, []core.TransactionModel{tx})
	Expect(txErr).NotTo(HaveOccurred())

	eventLogRepository := repositories.NewEventLogRepository(db)
	insertLogsErr := eventLogRepository.CreateEventLogs(headerID, []types.Log{log})
	Expect(insertLogsErr).NotTo(HaveOccurred())

	type persistedEventLog struct {
		ID          int64
		HeaderID    int64 `db:"header_id"`
		Transformed bool
	}
	var eventLog persistedEventLog
	getLogErr := db.Get(&eventLog, `SELECT id, header_id, transformed FROM public.event_logs WHERE tx_hash = $1`, log.TxHash.Hex())
	Expect(getLogErr).NotTo(HaveOccurred())
	result := core.EventLog{
		ID:          eventLog.ID,
		HeaderID:    eventLog.HeaderID,
		Log:         log,
		Transformed: eventLog.Transformed,
	}
	return result
}

func getFakeTransactionFromHash(txHash common.Hash) core.TransactionModel {
	return core.TransactionModel{
		Data:     nil,
		From:     getRandomAddress(),
		GasLimit: 0,
		GasPrice: 0,
		Hash:     hashToPrefixedString(txHash),
		Nonce:    0,
		Raw:      nil,
		Receipt:  core.Receipt{},
		To:       getRandomAddress(),
		TxIndex:  0,
		Value:    "0",
	}
}

func CreateMatchingTx(log types.Log, headerID int64, headerRepo repositories.HeaderRepository) {
	fakeHashTx := getFakeTransactionFromHash(log.TxHash)
	txErr := headerRepo.CreateTransactions(headerID, []core.TransactionModel{fakeHashTx})
	Expect(txErr).NotTo(HaveOccurred())
}

func getRandomHash() common.Hash {
	seed := randomString(10)
	return sha256.Sum256([]byte(seed))
}

func hashToPrefixedString(hash common.Hash) string {
	return fmt.Sprintf("0x%x", hash)
}

func getRandomAddress() string {
	hash := getRandomHash()
	stringHash := hashToPrefixedString(hash)
	return stringHash[:42]
}
