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
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const insertHeaderSyncLogQuery = `INSERT INTO header_sync_logs
	(header_id, address, topics, data, block_number, block_hash, tx_index, tx_hash, log_index, raw)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT DO NOTHING RETURNING id`

func MarkHeaderChecked(headerID int64, db *postgres.DB, checkedHeadersColumn string) error {
	_, err := db.Exec(`INSERT INTO public.checked_headers (header_id, `+checkedHeadersColumn+`)
		VALUES ($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET `+checkedHeadersColumn+` = checked_headers.`+checkedHeadersColumn+` + 1`, headerID, 1)
	return err
}

func MarkHeaderCheckedInTransaction(headerID int64, tx *sqlx.Tx, checkedHeadersColumn string) error {
	_, err := tx.Exec(`INSERT INTO public.checked_headers (header_id, `+checkedHeadersColumn+`)
		VALUES ($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET `+checkedHeadersColumn+` = checked_headers.`+checkedHeadersColumn+` + 1`, headerID, 1)
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
			AND column_name <> 'id'
			AND column_name <> 'header_id';`

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

// Builds a SQL string that checks if any column should be checked/rechecked.
// Defaults to FALSE when no columns are provided.
// Ex: ["columnA", "columnB"] => "NOT (columnA!=0 AND columnB!=0)"
//     [] => "FALSE"
func CreateHeaderCheckedPredicateSQL(boolColumns []string, recheckHeaders constants.TransformerExecution) string {
	if len(boolColumns) == 0 {
		return "FALSE"
	}

	if recheckHeaders {
		return createHeaderCheckedPredicateSQLForRecheckedHeaders(boolColumns)
	} else {
		return createHeaderCheckedPredicateSQLForMissingHeaders(boolColumns)
	}
}

func CreateLogs(headerID int64, logs []types.Log, db *postgres.DB) ([]core.HeaderSyncLog, error) {
	tx, txErr := db.Beginx()
	if txErr != nil {
		return nil, txErr
	}
	var results []core.HeaderSyncLog
	for _, log := range logs {
		logID, err := insertLog(headerID, log, tx)
		if err != nil {
			if logWasADuplicate(err) {
				continue
			}
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Errorf("failed to rollback header sync log insert: %s", rollbackErr.Error())
			}
			return nil, err
		}
		results = append(results, buildLog(logID, headerID, log))
	}
	return results, tx.Commit()
}

func logWasADuplicate(err error) bool {
	return err == sql.ErrNoRows
}

func insertLog(headerID int64, log types.Log, tx *sqlx.Tx) (int64, error) {
	topics := buildTopics(log)
	raw, jsonErr := log.MarshalJSON()
	if jsonErr != nil {
		return 0, jsonErr
	}
	var logID int64
	err := tx.QueryRowx(insertHeaderSyncLogQuery, headerID, log.Address.Hex(), topics, log.Data, log.BlockNumber,
		log.BlockHash.Hex(), log.TxIndex, log.TxHash.Hex(), log.Index, raw).Scan(&logID)
	return logID, err
}

func buildLog(logID int64, headerID int64, log types.Log) core.HeaderSyncLog {
	return core.HeaderSyncLog{
		ID:          logID,
		HeaderID:    headerID,
		Log:         log,
		Transformed: false,
	}
}

func buildTopics(log types.Log) pq.ByteaArray {
	var topics pq.ByteaArray
	for _, topic := range log.Topics {
		topics = append(topics, topic.Bytes())
	}
	return topics
}

func createHeaderCheckedPredicateSQLForMissingHeaders(boolColumns []string) string {
	var result bytes.Buffer
	result.WriteString(" (")

	// Loop excluding last column name
	for _, column := range boolColumns[:len(boolColumns)-1] {
		result.WriteString(fmt.Sprintf("%v=0 OR ", column))
	}

	result.WriteString(fmt.Sprintf("%v=0)", boolColumns[len(boolColumns)-1]))

	return result.String()
}

func createHeaderCheckedPredicateSQLForRecheckedHeaders(boolColumns []string) string {
	var result bytes.Buffer
	result.WriteString(" (")

	// Loop excluding last column name
	for _, column := range boolColumns[:len(boolColumns)-1] {
		result.WriteString(fmt.Sprintf("%v<%s OR ", column, constants.RecheckHeaderCap))
	}

	// No trailing "OR" for the last column name
	result.WriteString(fmt.Sprintf("%v<%s)", boolColumns[len(boolColumns)-1], constants.RecheckHeaderCap))

	return result.String()
}
