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
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
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

// Zero out check count for header with the given block number
func (repo CheckedHeadersRepository) MarkSingleHeaderUnchecked(blockNumber int64) error {
	_, err := repo.db.Exec(`UPDATE public.headers SET check_count = 0 WHERE block_number = $1`, blockNumber)
	return err
}

// Return header if check_count  < passed checkCount
func (repo CheckedHeadersRepository) UncheckedHeaders(startingBlockNumber, endingBlockNumber, checkCount int64) ([]core.Header, error) {
	var (
		result                  []core.Header
		query                   string
		err                     error
		recheckOffsetMultiplier = 15
	)

	if endingBlockNumber == -1 {
		query = `SELECT id, block_number, hash
			FROM public.headers
			WHERE (check_count < 1
			           AND block_number >= $1)
			   OR (check_count < $2
			           AND block_number <= ((SELECT MAX(block_number) FROM public.headers) - ($3 * check_count * (check_count + 1) / 2)))`
		err = repo.db.Select(&result, query, startingBlockNumber, checkCount, recheckOffsetMultiplier)
	} else {
		query = `SELECT id, block_number, hash
			FROM public.headers
			WHERE (check_count < 1
			           AND block_number >= $1
			           AND block_number <= $2)
			   OR (check_count < $3
			           AND block_number >= $1
			           AND block_number <= $2
			           AND block_number <= ((SELECT MAX(block_number) FROM public.headers) - ($4 * (check_count * (check_count + 1) / 2))))`
		err = repo.db.Select(&result, query, startingBlockNumber, endingBlockNumber, checkCount, recheckOffsetMultiplier)
	}

	return result, err
}
