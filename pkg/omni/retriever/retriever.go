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

type Retriever interface {
	RetrieveFirstBlock(contractAddr string) (int64, error)
	RetrieveFirstBlockFromLogs(contractAddr string) (int64, error)
	RetrieveFirstBlockFromReceipts(contractAddr string) (int64, error)
}

type retriever struct {
	*postgres.DB
}

func NewRetriever(db *postgres.DB) (r *retriever) {

	return &retriever{
		DB: db,
	}
}

// For some contracts the creation transaction receipt doesn't have the contract address so this doesn't work (e.g. Sai)
func (r *retriever) RetrieveFirstBlockFromReceipts(contractAddr string) (int64, error) {
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

// This servers as a heuristic to find the first block by finding the first contract event log
func (r *retriever) RetrieveFirstBlockFromLogs(contractAddr string) (int64, error) {
	var firstBlock int
	err := r.DB.Get(
		&firstBlock,
		"SELECT block_number FROM logs WHERE address = $1 ORDER BY block_number ASC LIMIT 1",
		contractAddr,
	)

	return int64(firstBlock), err
}

// Try both methods of finding the first block, with the receipt method taking precedence
func (r *retriever) RetrieveFirstBlock(contractAddr string) (int64, error) {
	i, err := r.RetrieveFirstBlockFromReceipts(contractAddr)
	if err != nil {
		i, err = r.RetrieveFirstBlockFromLogs(contractAddr)
	}

	return i, err
}
