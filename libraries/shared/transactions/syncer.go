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

package transactions

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/eth/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
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
	if len(transactionHashes) < 1 {
		return nil
	}
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
