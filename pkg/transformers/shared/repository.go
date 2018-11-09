package shared

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var ErrHeaderMismatch = errors.New("log's block hash does not match associated header's block hash")

func ValidateHeaderConsistency(headerID int64, rawLog []byte, db *postgres.DB) error {
	var log types.Log
	err := json.Unmarshal(rawLog, &log)
	if err != nil {
		return err
	}
	var hash string
	err = db.Get(&hash, `SELECT hash FROM public.headers WHERE id = $1`, headerID)
	if err != nil {
		return err
	}
	if hash != log.BlockHash.String() {
		err = deleteHeader(headerID, db)
		if err != nil {
			return err
		}
		return ErrHeaderMismatch
	}
	return nil
}

func deleteHeader(headerID int64, db *postgres.DB) error {
	_, err := db.Exec(`DELETE FROM public.headers WHERE id = $1`, headerID)
	return err
}

func MarkHeaderChecked(headerID int64, db *postgres.DB, checkedHeadersColumn string) error {
	_, err := db.Exec(`INSERT INTO public.checked_headers (header_id, `+checkedHeadersColumn+`)
		VALUES ($1, $2) 
		ON CONFLICT (header_id) DO
			UPDATE SET `+checkedHeadersColumn+` = $2`, headerID, true)
	return err
}

func MarkHeaderCheckedInTransaction(headerID int64, tx *sql.Tx, checkedHeadersColumn string) error {
	_, err := tx.Exec(`INSERT INTO public.checked_headers (header_id, `+checkedHeadersColumn+`)
		VALUES ($1, $2) 
		ON CONFLICT (header_id) DO
			UPDATE SET `+checkedHeadersColumn+` = $2`, headerID, true)
	return err
}

func MissingHeaders(startingBlockNumber, endingBlockNumber int64, db *postgres.DB, checkedHeadersColumn string) ([]core.Header, error) {
	var result []core.Header
	var query string
	var err error

	if endingBlockNumber == -1 {
		query = `SELECT headers.id, headers.block_number FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
				WHERE (header_id ISNULL OR ` + checkedHeadersColumn + ` IS FALSE)
				AND headers.block_number >= $1
				AND headers.eth_node_fingerprint = $2`
		err = db.Select(&result, query, startingBlockNumber, db.Node.ID)
	} else {
		query = `SELECT headers.id, headers.block_number FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
				WHERE (header_id ISNULL OR ` + checkedHeadersColumn + ` IS FALSE)
				AND headers.block_number >= $1
				AND headers.block_number <= $2
				AND headers.eth_node_fingerprint = $3`
		err = db.Select(&result, query, startingBlockNumber, endingBlockNumber, db.Node.ID)
	}

	return result, err
}
