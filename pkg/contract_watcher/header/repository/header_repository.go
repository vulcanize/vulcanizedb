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

package repository

import (
	"fmt"

	"github.com/hashicorp/golang-lru"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const columnCacheSize = 1000

// HeaderRepository interfaces with the header and checked_headers tables
type HeaderRepository interface {
	AddCheckColumn(id string) error
	AddCheckColumns(ids []string) error
	MarkHeaderChecked(headerID int64, eventID string) error
	MarkHeaderCheckedForAll(headerID int64, ids []string) error
	MarkHeadersCheckedForAll(headers []core.Header, ids []string) error
	MissingHeaders(startingBlockNumber int64, endingBlockNumber int64, eventID string) ([]core.Header, error)
	MissingMethodsCheckedEventsIntersection(startingBlockNumber, endingBlockNumber int64, methodIds, eventIds []string) ([]core.Header, error)
	MissingHeadersForAll(startingBlockNumber, endingBlockNumber int64, ids []string) ([]core.Header, error)
	CheckCache(key string) (interface{}, bool)
}

type headerRepository struct {
	db      *postgres.DB
	columns *lru.Cache // Cache created columns to minimize db connections
}

// NewHeaderRepository returns a new HeaderRepository
func NewHeaderRepository(db *postgres.DB) HeaderRepository {
	ccs, _ := lru.New(columnCacheSize)
	return &headerRepository{
		db:      db,
		columns: ccs,
	}
}

// AddCheckColumn adds a checked_header column for the provided column id
func (r *headerRepository) AddCheckColumn(id string) error {
	// Check cache to see if column already exists before querying pg
	_, ok := r.columns.Get(id)
	if ok {
		return nil
	}

	pgStr := "ALTER TABLE public.checked_headers ADD COLUMN IF NOT EXISTS "
	pgStr = pgStr + id + " INTEGER NOT NULL DEFAULT 0"
	_, err := r.db.Exec(pgStr)
	if err != nil {
		return err
	}

	// Add column name to cache
	r.columns.Add(id, true)

	return nil
}

// AddCheckColumns adds a checked_header column for all of the provided column ids
func (r *headerRepository) AddCheckColumns(ids []string) error {
	var err error
	baseQuery := "ALTER TABLE public.checked_headers"
	input := make([]string, 0, len(ids))
	for _, id := range ids {
		_, ok := r.columns.Get(id)
		if !ok {
			baseQuery += " ADD COLUMN IF NOT EXISTS " + id + " INTEGER NOT NULL DEFAULT 0,"
			input = append(input, id)
		}
	}
	if len(input) > 0 {
		_, err = r.db.Exec(baseQuery[:len(baseQuery)-1])
		if err == nil {
			for _, id := range input {
				r.columns.Add(id, true)
			}
		}
	}

	return err
}

// MarkHeaderChecked marks the header checked for the provided column id
func (r *headerRepository) MarkHeaderChecked(headerID int64, id string) error {
	_, err := r.db.Exec(`INSERT INTO public.checked_headers (header_id, `+id+`)
		VALUES ($1, $2) 
		ON CONFLICT (header_id) DO
			UPDATE SET `+id+` = checked_headers.`+id+` + 1`, headerID, 1)
	return err
}

// MarkHeaderCheckedForAll marks the header checked for all of the provided column ids
func (r *headerRepository) MarkHeaderCheckedForAll(headerID int64, ids []string) error {
	pgStr := "INSERT INTO public.checked_headers (header_id, "
	for _, id := range ids {
		pgStr += id + ", "
	}
	pgStr = pgStr[:len(pgStr)-2] + ") VALUES ($1, "
	for i := 0; i < len(ids); i++ {
		pgStr += "1, "
	}
	pgStr = pgStr[:len(pgStr)-2] + ") ON CONFLICT (header_id) DO UPDATE SET "
	for _, id := range ids {
		pgStr += id + `= checked_headers.` + id + ` + 1, `
	}
	pgStr = pgStr[:len(pgStr)-2]
	_, err := r.db.Exec(pgStr, headerID)
	return err
}

// MarkHeadersCheckedForAll marks all of the provided headers checked for each of the provided column ids
func (r *headerRepository) MarkHeadersCheckedForAll(headers []core.Header, ids []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	for _, header := range headers {
		pgStr := "INSERT INTO public.checked_headers (header_id, "
		for _, id := range ids {
			pgStr += id + ", "
		}
		pgStr = pgStr[:len(pgStr)-2] + ") VALUES ($1, "
		for i := 0; i < len(ids); i++ {
			pgStr += "1, "
		}
		pgStr = pgStr[:len(pgStr)-2] + ") ON CONFLICT (header_id) DO UPDATE SET "
		for _, id := range ids {
			pgStr += fmt.Sprintf("%s = checked_headers.%s + 1, ", id, id)
		}
		pgStr = pgStr[:len(pgStr)-2]
		_, err = tx.Exec(pgStr, header.ID)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Warnf("error rolling back transaction: %s", rollbackErr.Error())
			}
			return err
		}
	}
	err = tx.Commit()
	return err
}

// MissingHeaders returns missing headers for the provided checked_headers column id
func (r *headerRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64, id string) ([]core.Header, error) {
	var result []core.Header
	var query string
	var err error
	if endingBlockNumber == -1 {
		query = `SELECT headers.id, headers.block_number, headers.hash FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
				WHERE (header_id ISNULL OR checked_headers.` + id + `=0)
				AND headers.block_number >= $1
				AND headers.eth_node_fingerprint = $2
				ORDER BY headers.block_number`
		err = r.db.Select(&result, query, startingBlockNumber, r.db.Node.ID)
	} else {
		query = `SELECT headers.id, headers.block_number, headers.hash FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
				WHERE (header_id ISNULL OR checked_headers.` + id + `=0)
				AND headers.block_number >= $1
				AND headers.block_number <= $2
				AND headers.eth_node_fingerprint = $3
				ORDER BY headers.block_number`
		err = r.db.Select(&result, query, startingBlockNumber, endingBlockNumber, r.db.Node.ID)
	}
	return continuousHeaders(result), err
}

// MissingHeadersForAll returns missing headers for all of the provided checked_headers column ids
func (r *headerRepository) MissingHeadersForAll(startingBlockNumber, endingBlockNumber int64, ids []string) ([]core.Header, error) {
	var result []core.Header
	var query string
	var err error
	baseQuery := `SELECT headers.id, headers.block_number, headers.hash FROM headers
				  LEFT JOIN checked_headers on headers.id = header_id
				  WHERE (header_id ISNULL`
	for _, id := range ids {
		baseQuery += ` OR checked_headers.` + id + `= 0`
	}
	if endingBlockNumber == -1 {
		endStr := `) AND headers.block_number >= $1
				  AND headers.eth_node_fingerprint = $2
				  ORDER BY headers.block_number`
		query = baseQuery + endStr
		err = r.db.Select(&result, query, startingBlockNumber, r.db.Node.ID)
	} else {
		endStr := `) AND headers.block_number >= $1
				  AND headers.block_number <= $2
				  AND headers.eth_node_fingerprint = $3
				  ORDER BY headers.block_number`
		query = baseQuery + endStr
		err = r.db.Select(&result, query, startingBlockNumber, endingBlockNumber, r.db.Node.ID)
	}
	return continuousHeaders(result), err
}

// MissingMethodsCheckedEventsIntersection returns headers that have been checked for all of the provided event ids but not for the provided method ids
func (r *headerRepository) MissingMethodsCheckedEventsIntersection(startingBlockNumber, endingBlockNumber int64, methodIds, eventIds []string) ([]core.Header, error) {
	var result []core.Header
	var query string
	var err error
	baseQuery := `SELECT headers.id, headers.block_number, headers.hash FROM headers
				  LEFT JOIN checked_headers on headers.id = header_id
				  WHERE (header_id IS NOT NULL`
	for _, id := range eventIds {
		baseQuery += ` AND ` + id + `!=0`
	}
	baseQuery += `) AND (`
	for _, id := range methodIds {
		baseQuery += id + ` =0 AND `
	}
	baseQuery = baseQuery[:len(baseQuery)-5] + `) `
	if endingBlockNumber == -1 {
		endStr := `AND headers.block_number >= $1
				  AND headers.eth_node_fingerprint = $2
				  ORDER BY headers.block_number`
		query = baseQuery + endStr
		err = r.db.Select(&result, query, startingBlockNumber, r.db.Node.ID)
	} else {
		endStr := `AND headers.block_number >= $1
				  AND headers.block_number <= $2
				  AND headers.eth_node_fingerprint = $3
				  ORDER BY headers.block_number`
		query = baseQuery + endStr
		err = r.db.Select(&result, query, startingBlockNumber, endingBlockNumber, r.db.Node.ID)
	}
	return continuousHeaders(result), err
}

// Returns a continuous set of headers
func continuousHeaders(headers []core.Header) []core.Header {
	if len(headers) < 1 {
		logrus.Trace("no headers to arrange continuously")
		return headers
	}
	previousHeader := headers[0].BlockNumber
	for i := 1; i < len(headers); i++ {
		previousHeader++
		if headers[i].BlockNumber != previousHeader {
			return headers[:i]
		}
	}

	return headers
}

// CheckCache checks the repositories column id cache for a value
func (r *headerRepository) CheckCache(key string) (interface{}, bool) {
	return r.columns.Get(key)
}
