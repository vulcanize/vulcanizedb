package fakes

import "github.com/ethereum/go-ethereum/core/types"

type MockTransactionSyncer struct {
	SyncTransactionsCalled bool
	SyncTransactionsError  error
}

func (syncer *MockTransactionSyncer) SyncTransactions(headerID int64, logs []types.Log) error {
	syncer.SyncTransactionsCalled = true
	return syncer.SyncTransactionsError
}
