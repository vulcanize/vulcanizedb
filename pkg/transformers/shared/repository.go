package shared

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

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

// Treats a header as missing if it's not in the headers table, or not checked for some log type
func MissingHeaders(startingBlockNumber, endingBlockNumber int64, db *postgres.DB, notCheckedSQL string) ([]core.Header, error) {
	var result []core.Header
	var query string
	var err error

	if endingBlockNumber == -1 {
		query = `SELECT headers.id, headers.block_number, headers.hash FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
				WHERE (header_id ISNULL OR ` + notCheckedSQL + `)
				AND headers.block_number >= $1
				AND headers.eth_node_fingerprint = $2`
		err = db.Select(&result, query, startingBlockNumber, db.Node.ID)
	} else {
		query = `SELECT headers.id, headers.block_number, headers.hash FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
				WHERE (header_id ISNULL OR ` + notCheckedSQL + `)
				AND headers.block_number >= $1
				AND headers.block_number <= $2
				AND headers.eth_node_fingerprint = $3`
		err = db.Select(&result, query, startingBlockNumber, endingBlockNumber, db.Node.ID)
	}

	return result, err
}

func GetCheckedColumnNames(db *postgres.DB) ([]string, error) {
	// Query returns `[]driver.Value`, nullable polymorphic interface
	var queryResult []driver.Value
	columnNamesQuery :=
		`SELECT column_name FROM information_schema.columns
		WHERE table_schema = 'public'
			AND table_name = 'checked_headers'
			AND column_name ~ '_checked';`

	err := db.Select(&queryResult, columnNamesQuery)
	if err != nil {
		return []string{}, err
	}

	// Transform column names from `driver.Value` to strings
	var columnNames []string
	for _, result := range queryResult {
		if columnName, ok := result.(string); ok {
			columnNames = append(columnNames, columnName)
		} else {
			return []string{}, fmt.Errorf("incorrect value for checked_headers column name")
		}
	}

	return columnNames, nil
}

// Builds a SQL string that checks if any column value is FALSE, given the column names.
// Defaults to FALSE when no columns are provided.
// Ex: ["columnA", "columnB"] => "NOT (columnA AND columnB)"
//     [] => "FALSE"
func CreateNotCheckedSQL(boolColumns []string) string {
	var result bytes.Buffer

	if len(boolColumns) == 0 {
		return "FALSE"
	}

	result.WriteString("NOT (")
	// Loop excluding last column name
	for _, column := range boolColumns[:len(boolColumns)-1] {
		result.WriteString(fmt.Sprintf("%v AND ", column))
	}

	// No trailing "OR" for the last column name
	result.WriteString(fmt.Sprintf("%v)", boolColumns[len(boolColumns)-1]))

	return result.String()
}

func GetTicInTx(headerID int64, tx *sql.Tx) (int64, error) {
	var blockTimestamp int64
	err := tx.QueryRow(`SELECT block_timestamp FROM public.headers WHERE id = $1;`, headerID).Scan(&blockTimestamp)
	if err != nil {
		return 0, err
	}

	tic := blockTimestamp + constants.TTL
	return tic, nil
}
