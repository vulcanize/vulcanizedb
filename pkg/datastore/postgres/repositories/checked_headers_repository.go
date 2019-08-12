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
	insertCheckedHeaderQuery = `INSERT INTO public.checked_headers (header_id) VALUES ($1)
		ON CONFLICT (header_id) DO UPDATE
		SET check_count = (SELECT check_count FROM public.checked_headers WHERE header_id = $1) + 1`
)

type CheckedHeadersRepository struct {
	db *postgres.DB
}

func NewCheckedHeadersRepository(db *postgres.DB) CheckedHeadersRepository {
	return CheckedHeadersRepository{db: db}
}

// Adds header_id to the checked_headers table, or increment check_count if header_id already present
func (repo CheckedHeadersRepository) MarkHeaderChecked(headerID int64) error {
	_, err := repo.db.Exec(insertCheckedHeaderQuery, headerID)
	return err
}

// Return header_id if not present in checked_headers or its check_count is < passed checkCount
func (repo CheckedHeadersRepository) MissingHeaders(startingBlockNumber, endingBlockNumber, checkCount int64) ([]core.Header, error) {
	var result []core.Header
	var query string
	var err error

	if endingBlockNumber == -1 {
		query = `SELECT headers.id, headers.block_number, headers.hash FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
				WHERE (header_id ISNULL OR check_count < $2)
				AND headers.block_number >= $1
				AND headers.eth_node_fingerprint = $3`
		err = repo.db.Select(&result, query, startingBlockNumber, checkCount, repo.db.Node.ID)
	} else {
		query = `SELECT headers.id, headers.block_number, headers.hash FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
				WHERE (header_id ISNULL OR check_count < $3)
				AND headers.block_number >= $1
				AND headers.block_number <= $2
				AND headers.eth_node_fingerprint = $4`
		err = repo.db.Select(&result, query, startingBlockNumber, endingBlockNumber, checkCount, repo.db.Node.ID)
	}

	return result, err
}
