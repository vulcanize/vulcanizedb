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

package event

import (
	"database/sql/driver"
	"fmt"
	"github.com/makerdao/vulcanizedb/utils"
	"strings"

	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/sirupsen/logrus"
)

const SetLogTransformedQuery = `UPDATE public.header_sync_logs SET transformed = true WHERE id = $1`

// Repository persists transformed values to the DB
type Repository interface {
	Create(models []InsertionModel) error
	SetDB(db *postgres.DB)
}

// LogFK is the name of log foreign key columns
const LogFK ColumnName = "log_id"

// AddressFK is the name of address foreign key columns
const AddressFK ColumnName = "address_id"

// HeaderFK is the name of header foreign key columns
const HeaderFK ColumnName = "header_id"

// SchemaName is the schema to work with
type SchemaName string

// TableName identifies the table for inserting the data
type TableName string

// ColumnName identifies columns on the given table
type ColumnName string

// ColumnValues maps a column to the value for insertion. This is restricted to []byte, bool, float64, int64, string, time.Time
type ColumnValues map[ColumnName]interface{}

// ErrUnsupportedValue is thrown when a model supplies a type of value the postgres driver cannot handle.
var ErrUnsupportedValue = func(value interface{}) error {
	return fmt.Errorf("unsupported type of value supplied in model: %v (%T)", value, value)
}

// InsertionModel is the generalised data structure a converter returns, and contains everything the repository needs to
// persist the converted data.
type InsertionModel struct {
	SchemaName     SchemaName
	TableName      TableName
	OrderedColumns []ColumnName // Defines the fields to insert, and in which order the table expects them
	ColumnValues   ColumnValues // Associated values for columns, restricted to []byte, bool, float64, int64, string, time.Time
}

// ModelToQuery stores memoised insertion queries to minimise computation
var ModelToQuery = map[string]string{}

// GetMemoizedQuery gets/creates a DB insertion query for the model
func GetMemoizedQuery(model InsertionModel) string {
	// The schema and table name uniquely determines the insertion query, use that for memoization
	queryKey := string(model.SchemaName) + string(model.TableName)
	query, queryMemoized := ModelToQuery[queryKey]
	if !queryMemoized {
		query = GenerateInsertionQuery(model)
		ModelToQuery[queryKey] = query
	}
	return query
}

// GenerateInsertionQuery creates an SQL insertion query from an insertion model.
// Should be called through GetMemoizedQuery, so the query is not generated on each call to Create.
func GenerateInsertionQuery(model InsertionModel) string {
	var valuePlaceholders []string
	var updateOnConflict []string
	for i := 0; i < len(model.OrderedColumns); i++ {
		valuePlaceholder := fmt.Sprintf("$%d", 1+i)
		valuePlaceholders = append(valuePlaceholders, valuePlaceholder)
		updateOnConflict = append(updateOnConflict,
			fmt.Sprintf("%s = %s", model.OrderedColumns[i], valuePlaceholder))
	}

	baseQuery := `INSERT INTO %v.%v (%v) VALUES(%v)
		ON CONFLICT (header_id, log_id) DO UPDATE SET %v;`

	return fmt.Sprintf(baseQuery,
		model.SchemaName,
		model.TableName,
		joinOrderedColumns(model.OrderedColumns),
		strings.Join(valuePlaceholders, ", "),
		strings.Join(updateOnConflict, ", "))
}

/*
Create generates an insertion query and persists to the DB, given a slice of InsertionModels.
ColumnValues are restricted to []byte, bool, float64, int64, string, time.Time.

testModel = shared.InsertionModel{
	SchemaName:     "public"
	TableName:      "testEvent",
	OrderedColumns: []string{"header_id", "log_id", "variable1"},
	ColumnValues: ColumnValues{
		"header_id": 303
		"log_id":   "808",
		"variable1": "value1",
	},
}
*/
func Create(models []InsertionModel, db *postgres.DB) error {
	if len(models) == 0 {
		return fmt.Errorf("repository got empty model slice")
	}

	tx, dbErr := db.Beginx()
	if dbErr != nil {
		return dbErr
	}

	for _, model := range models {
		// Maps can't be iterated over in a reliable manner, so we rely on OrderedColumns to define the order to insert
		// tx.Exec is variadically typed in the args, so if we wrap in []interface{} we can apply them all automatically
		var args []interface{}
		for _, col := range model.OrderedColumns {
			value := model.ColumnValues[col]
			// Check whether or not PG can accept the type of value in the model
			okPgValue := driver.IsValue(value)
			if !okPgValue {
				logrus.WithField("model", model).Errorf("PG cannot handle value of this type: %T", value)
				return ErrUnsupportedValue(value)
			}
			args = append(args, value)
		}

		insertionQuery := GetMemoizedQuery(model)
		_, execErr := tx.Exec(insertionQuery, args...) // couldn't pass varying types in bulk with args :: []string

		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}

		_, logErr := tx.Exec(SetLogTransformedQuery, model.ColumnValues[LogFK])

		if logErr != nil {
			utils.RollbackAndLogFailure(tx, logErr, "header_sync_logs.transformed")
			return logErr
		}
	}

	return tx.Commit()
}

func joinOrderedColumns(columns []ColumnName) string {
	var stringColumns []string
	for _, columnName := range columns {
		stringColumns = append(stringColumns, string(columnName))
	}
	return strings.Join(stringColumns, ", ")
}
