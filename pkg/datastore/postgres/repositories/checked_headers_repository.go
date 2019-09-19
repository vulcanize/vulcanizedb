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

package repositories

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertCheckedHeaderQuery = `UPDATE public.headers SET check_count = (SELECT check_count WHERE id = $1) + 1 WHERE id = $1`
)

type CheckedHeadersRepository struct {
	db *postgres.DB
}

func NewCheckedHeadersRepository(db *postgres.DB) CheckedHeadersRepository {
	return CheckedHeadersRepository{db: db}
}

// Increment check_count for header
func (repo CheckedHeadersRepository) MarkHeaderChecked(headerID int64) error {
	_, err := repo.db.Exec(insertCheckedHeaderQuery, headerID)
	return err
}

// Zero out check count for headers with block number >= startingBlockNumber
func (repo CheckedHeadersRepository) MarkHeadersUnchecked(startingBlockNumber int64) error {
	_, err := repo.db.Exec(`UPDATE public.headers SET check_count = 0 WHERE block_number >= $1`, startingBlockNumber)
	return err
}

// Return header if check_count  < passed checkCount
func (repo CheckedHeadersRepository) UncheckedHeaders(startingBlockNumber, endingBlockNumber, checkCount int64) ([]core.Header, error) {
	var result []core.Header
	var query string
	var err error

	if endingBlockNumber == -1 {
		query = `SELECT id, block_number, hash
				FROM headers
				WHERE check_count < $2
				AND block_number >= $1
				AND eth_node_fingerprint = $3`
		err = repo.db.Select(&result, query, startingBlockNumber, checkCount, repo.db.Node.ID)
	} else {
		query = `SELECT id, block_number, hash
				FROM headers
				WHERE check_count < $3
				AND block_number >= $1
				AND block_number <= $2
				AND eth_node_fingerprint = $4`
		err = repo.db.Select(&result, query, startingBlockNumber, endingBlockNumber, checkCount, repo.db.Node.ID)
	}

	return result, err
}
