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

package retriever

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

// Block retriever is used to retrieve the first block for a given contract and the most recent block
// It requires a vDB synced database with blocks, transactions, receipts, and logs
type BlockRetriever interface {
	RetrieveFirstBlock() (int64, error)
	RetrieveMostRecentBlock() (int64, error)
}

type blockRetriever struct {
	db *postgres.DB
}

func NewBlockRetriever(db *postgres.DB) (r *blockRetriever) {
	return &blockRetriever{
		db: db,
	}
}

// Retrieve block number of earliest header in repo
func (r *blockRetriever) RetrieveFirstBlock() (int64, error) {
	var firstBlock int
	err := r.db.Get(
		&firstBlock,
		"SELECT block_number FROM headers ORDER BY block_number LIMIT 1",
	)

	return int64(firstBlock), err
}

// Retrieve block number of latest header in repo
func (r *blockRetriever) RetrieveMostRecentBlock() (int64, error) {
	var lastBlock int
	err := r.db.Get(
		&lastBlock,
		"SELECT block_number FROM headers ORDER BY block_number DESC LIMIT 1",
	)

	return int64(lastBlock), err
}
