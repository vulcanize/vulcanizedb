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
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"strings"
)

type Repository interface {
	Create(models []InsertionModel) error
	SetDB(db *postgres.DB)
}

const LogFK ColumnName = "log_id"
const AddressFK ColumnName = "address_id"
const HeaderFK ColumnName = "header_id"

// Type aliases to reduce fat-finger bugs, forcing plugins to use declared variables for schema, table, and column names
type SchemaName string
type TableName string
type ColumnName string
type ColumnValues map[ColumnName]interface{}

type InsertionModel struct {
	SchemaName     SchemaName
	TableName      TableName
	OrderedColumns []ColumnName // Defines the fields to insert, and in which order the table expects them
	// ColumnValues needs to be typed interface{}, since  values to insert can be of mixed types
	ColumnValues ColumnValues // Associated values for columns.
}

// Stores memoised insertion queries to minimise computation
var ModelToQuery = map[string]string{}

// Get or create a DB insertion query for the model
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

// Creates an insertion query from an insertion model. Should be called through GetMemoizedQuery, so the query is not
// generated on each call to Create.
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

/* Given an instance of InsertionModel, example below, generates an insertion query and persists to the DB.

testModel = shared.InsertionModel{
	SchemaName:     "maker"
	TableName:      "testEvent",
	OrderedColumns: []string{"header_id", "log_id", constants.IlkFK, constants.UrnFK, "variable1"},
	ColumnValues: ColumnValues{
		"header_id": 303
		"log_id":   "1",
		"variable1": "value1",
		constants.IlkFK: 808,
		constants.UrnFK: 909,
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
			args = append(args, model.ColumnValues[col])
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

		_, logErr := tx.Exec(`UPDATE public.header_sync_logs SET transformed = true WHERE id = $1`, model.ColumnValues[LogFK])

		if logErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
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
