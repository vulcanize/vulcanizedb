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

package repository

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type HeaderRepository interface {
	MarkHeaderChecked(headerID int64, eventID string) error
	MissingHeaders(startingBlockNumber int64, endingBlockNumber int64, eventID string) ([]core.Header, error)
}

type headerRepository struct {
	db *postgres.DB
}

func NewHeaderRepository(db *postgres.DB) *headerRepository {
	return &headerRepository{
		db: db,
	}
}

func (r *headerRepository) MarkHeaderChecked(headerID int64, eventID string) error {
	_, err := r.db.Exec(`INSERT INTO public.checked_headers (header_id, `+eventID+`)
		VALUES ($1, $2) 
		ON CONFLICT (header_id) DO
			UPDATE SET `+eventID+` = $2`, headerID, true)
	return err
}

func (r *headerRepository) MissingHeaders(startingBlockNumber int64, endingBlockNumber int64, eventID string) ([]core.Header, error) {
	var result []core.Header
	var query string
	var err error

	if endingBlockNumber == -1 {
		query = `SELECT headers.id, headers.block_number, headers.hash FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
				WHERE (header_id ISNULL OR ` + eventID + ` IS FALSE)
				AND headers.block_number >= $1
				AND headers.eth_node_fingerprint = $2`
		err = r.db.Select(&result, query, startingBlockNumber, r.db.Node.ID)
	} else {
		query = `SELECT headers.id, headers.block_number, headers.hash FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
				WHERE (header_id ISNULL OR ` + eventID + ` IS FALSE)
				AND headers.block_number >= $1
				AND headers.block_number <= $2
				AND headers.eth_node_fingerprint = $3`
		err = r.db.Select(&result, query, startingBlockNumber, endingBlockNumber, r.db.Node.ID)
	}

	return result, err
}
