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

package retriever

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

// Block retriever is used to retrieve the first block for a given contract and the most recent block
// It requires a vDB synced database with blocks, transactions, receipts, and logs
type BlockRetriever interface {
	RetrieveFirstBlock(contractAddr string) (int64, error)
	RetrieveMostRecentBlock() (int64, error)
}

type blockRetriever struct {
	*postgres.DB
}

func NewBlockRetriever(db *postgres.DB) (r *blockRetriever) {

	return &blockRetriever{
		DB: db,
	}
}

// Try both methods of finding the first block, with the receipt method taking precedence
func (r *blockRetriever) RetrieveFirstBlock(contractAddr string) (int64, error) {
	i, err := r.retrieveFirstBlockFromReceipts(contractAddr)
	if err != nil {
		i, err = r.retrieveFirstBlockFromLogs(contractAddr)
	}

	return i, err
}

// For some contracts the contract creation transaction receipt doesn't have the contract address so this doesn't work (e.g. Sai)
func (r *blockRetriever) retrieveFirstBlockFromReceipts(contractAddr string) (int64, error) {
	var firstBlock int
	err := r.DB.Get(
		&firstBlock,
		`SELECT number FROM blocks
		       WHERE id = (SELECT block_id FROM receipts
               	           WHERE contract_address = $1
		                   ORDER BY block_id ASC
					       LIMIT 1)`,
		contractAddr,
	)

	return int64(firstBlock), err
}

// In which case this servers as a heuristic to find the first block by finding the first contract event log
func (r *blockRetriever) retrieveFirstBlockFromLogs(contractAddr string) (int64, error) {
	var firstBlock int
	err := r.DB.Get(
		&firstBlock,
		"SELECT block_number FROM logs WHERE address = $1 ORDER BY block_number ASC LIMIT 1",
		contractAddr,
	)

	return int64(firstBlock), err
}

// Method to retrieve the most recent block in vDB
func (r *blockRetriever) RetrieveMostRecentBlock() (int64, error) {
	var lastBlock int64
	err := r.DB.Get(
		&lastBlock,
		"SELECT number FROM blocks ORDER BY number DESC LIMIT 1",
	)

	return lastBlock, err
}
