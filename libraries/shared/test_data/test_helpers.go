package test_data

import (
	"crypto/sha256"
	"fmt"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/watcher"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/gomega"
)

// Create a header sync log to reference in an event, returning inserted header sync log
func CreateTestLog(headerID int64, db *postgres.DB) core.EventLog {
	txHash := getRandomHash()
	log := types.Log{
		Address:     common.Address{},
		Topics:      nil,
		Data:        nil,
		BlockNumber: 0,
		TxHash:      txHash,
		TxIndex:     uint(rand.Int31()),
		BlockHash:   common.Hash{},
		Index:       0,
		Removed:     false,
	}

	tx := getFakeTransactionFromHash(txHash)
	headerRepository := repositories.NewHeaderRepository(db)
	txErr := headerRepository.CreateTransactions(headerID, []core.TransactionModel{tx})
	Expect(txErr).NotTo(HaveOccurred())

	eventLogRepository := repositories.NewEventLogRepository(db)
	insertLogsErr := eventLogRepository.CreateEventLogs(headerID, []types.Log{log})
	Expect(insertLogsErr).NotTo(HaveOccurred())

	// TODO: consider better calibrated limit depending on testing needs
	eventLogs, getLogsErr := eventLogRepository.GetUntransformedEventLogs(0, watcher.ResultsLimit)
	Expect(getLogsErr).NotTo(HaveOccurred())
	for _, eventLog := range eventLogs {
		if eventLog.Log.TxIndex == log.TxIndex {
			return eventLog
		}
	}
	panic("couldn't find inserted test log")
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
