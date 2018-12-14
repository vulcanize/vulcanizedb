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
	"database/sql"
	"fmt"

	"github.com/hashicorp/golang-lru"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const columnCacheSize = 1000

type HeaderRepository interface {
	AddCheckColumn(eventID string) error
	MarkHeaderChecked(headerID int64, eventID string) error
	MarkHeadersChecked(headers []core.Header, ids []string) error
	MissingHeaders(startingBlockNumber int64, endingBlockNumber int64, eventID string) ([]core.Header, error)
	MissingMethodsCheckedEventsIntersection(startingBlockNumber, endingBlockNumber int64, methodIds, eventIds []string) ([]core.Header, error)
	CheckCache(key string) (interface{}, bool)
}

type headerRepository struct {
	db      *postgres.DB
	columns *lru.Cache // Cache created columns to minimize db connections
}

func NewHeaderRepository(db *postgres.DB) *headerRepository {
	ccs, _ := lru.New(columnCacheSize)
	return &headerRepository{
		db:      db,
		columns: ccs,
	}
}

func (r *headerRepository) AddCheckColumn(id string) error {
	// Check cache to see if column already exists before querying pg
	_, ok := r.columns.Get(id)
	if ok {
		return nil
	}

	pgStr := "ALTER TABLE public.checked_headers ADD COLUMN IF NOT EXISTS "
	pgStr = pgStr + id + " BOOLEAN NOT NULL DEFAULT FALSE"
	_, err := r.db.Exec(pgStr)
	if err != nil {
		return err
	}

	// Add column name to cache
	r.columns.Add(id, true)

	return nil
}

func (r *headerRepository) MarkHeaderChecked(headerID int64, id string) error {
	_, err := r.db.Exec(`INSERT INTO public.checked_headers (header_id, `+id+`)
		VALUES ($1, $2) 
		ON CONFLICT (header_id) DO
			UPDATE SET `+id+` = $2`, headerID, true)

	return err
}

func (r *headerRepository) MarkHeadersChecked(headers []core.Header, ids []string) error {
	tx, err := r.db.Begin()
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
			pgStr += "true, "
		}
		pgStr = pgStr[:len(pgStr)-2] + ") ON CONFLICT (header_id) DO UPDATE SET "
		for _, id := range ids {
			pgStr += fmt.Sprintf("%s = true, ", id)
		}
		pgStr = pgStr[:len(pgStr)-2]
		_, err = tx.Exec(pgStr, header.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *headerRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64, id string) ([]core.Header, error) {
	var result []core.Header
	var query string
	var err error

	if endingBlockNumber == -1 {
		query = `SELECT headers.id, headers.block_number, headers.hash FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
				WHERE (header_id ISNULL OR ` + id + ` IS FALSE)
				AND headers.block_number >= $1
				AND headers.eth_node_fingerprint = $2
				ORDER BY headers.block_number`
		err = r.db.Select(&result, query, startingBlockNumber, r.db.Node.ID)
	} else {
		query = `SELECT headers.id, headers.block_number, headers.hash FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
				WHERE (header_id ISNULL OR ` + id + ` IS FALSE)
				AND headers.block_number >= $1
				AND headers.block_number <= $2
				AND headers.eth_node_fingerprint = $3
				ORDER BY headers.block_number`
		err = r.db.Select(&result, query, startingBlockNumber, endingBlockNumber, r.db.Node.ID)
	}

	return result, err
}

func (r *headerRepository) MissingMethodsCheckedEventsIntersection(startingBlockNumber, endingBlockNumber int64, methodIds, eventIds []string) ([]core.Header, error) {
	var result []core.Header
	var query string
	var err error
	baseQuery := `SELECT headers.id, headers.block_number, headers.hash FROM headers
				  LEFT JOIN checked_headers on headers.id = header_id
				  WHERE (header_id IS NOT NULL`
	for _, id := range eventIds {
		baseQuery += ` AND ` + id + ` IS TRUE`
	}
	baseQuery += `) AND (`
	for _, id := range methodIds {
		baseQuery += id + ` IS FALSE AND `
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

	return result, err
}

func (r *headerRepository) CheckCache(key string) (interface{}, bool) {
	return r.columns.Get(key)
}

func MarkHeaderCheckedInTransaction(headerID int64, tx *sql.Tx, eventID string) error {
	_, err := tx.Exec(`INSERT INTO public.checked_headers (header_id, `+eventID+`)
		VALUES ($1, $2) 
		ON CONFLICT (header_id) DO
			UPDATE SET `+eventID+` = $2`, headerID, true)
	return err
}
