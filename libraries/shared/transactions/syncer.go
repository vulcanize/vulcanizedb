package transactions

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

type ITransactionsSyncer interface {
	SyncTransactions(headerID int64, logs []types.Log) error
}

type TransactionsSyncer struct {
	BlockChain core.BlockChain
	Repository datastore.HeaderRepository
}

func NewTransactionsSyncer(db *postgres.DB, blockChain core.BlockChain) TransactionsSyncer {
	repository := repositories.NewHeaderRepository(db)
	return TransactionsSyncer{
		BlockChain: blockChain,
		Repository: repository,
	}
}

func (syncer TransactionsSyncer) SyncTransactions(headerID int64, logs []types.Log) error {
	transactionHashes := getUniqueTransactionHashes(logs)
	transactions, transactionErr := syncer.BlockChain.GetTransactions(transactionHashes)
	if transactionErr != nil {
		return transactionErr
	}
	writeErr := syncer.Repository.CreateTransactions(headerID, transactions)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func getUniqueTransactionHashes(logs []types.Log) []common.Hash {
	seen := make(map[common.Hash]struct{}, len(logs))
	var result []common.Hash
	for _, log := range logs {
		if _, ok := seen[log.TxHash]; ok {
			continue
		}
		seen[log.TxHash] = struct{}{}
		result = append(result, log.TxHash)
	}
	return result
}
