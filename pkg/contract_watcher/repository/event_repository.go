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
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/golang-lru"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/sirupsen/logrus"
)

const (
	// Number of contract address and method ids to keep in cache
	contractCacheSize = 100
	eventCacheSize    = 1000
)

// EventRepository is used to persist event data into custom tables
type EventRepository interface {
	PersistLogs(logs []types.Log, eventInfo types.Event, contractAddr, contractName string) error
	CreateEventTable(contractAddr string, event types.Event) (bool, error)
	CreateContractSchema(contractName string) (bool, error)
	CheckSchemaCache(key string) (interface{}, bool)
	CheckTableCache(key string) (interface{}, bool)
}

type eventRepository struct {
	db      *postgres.DB
	schemas *lru.Cache // Cache names of recently used schemas to minimize db connections
	tables  *lru.Cache // Cache names of recently used tables to minimize db connections
}

// NewEventRepository returns a new EventRepository
func NewEventRepository(db *postgres.DB) EventRepository {
	ccs, _ := lru.New(contractCacheSize)
	ecs, _ := lru.New(eventCacheSize)
	return &eventRepository{
		db:      db,
		schemas: ccs,
		tables:  ecs,
	}
}

// PersistLogs creates a schema for the contract if needed
// Creates table for the watched contract event if needed
// Persists converted event log data into this custom table
func (r *eventRepository) PersistLogs(logs []types.Log, eventInfo types.Event, contractAddr, contractName string) error {
	if len(logs) == 0 {
		return errors.New("event repository error: passed empty logs slice")
	}
	_, schemaErr := r.CreateContractSchema(contractAddr)
	if schemaErr != nil {
		return fmt.Errorf("error creating schema for contract %s: %s", contractAddr, schemaErr.Error())
	}

	_, tableErr := r.CreateEventTable(contractAddr, eventInfo)
	if tableErr != nil {
		return fmt.Errorf("error creating table for event %s on contract %s: %s", eventInfo.Name, contractAddr, tableErr.Error())
	}

	return r.persistLogs(logs, eventInfo, contractAddr, contractName)
}

func (r *eventRepository) persistLogs(logs []types.Log, eventInfo types.Event, contractAddr, contractName string) error {
	return r.persistEventLogs(logs, eventInfo, contractAddr, contractName)
}

// Creates a custom postgres command to persist logs for the given event (compatible with header synced vDB)
func (r *eventRepository) persistEventLogs(logs []types.Log, eventInfo types.Event, contractAddr, contractName string) error {
	tx, txErr := r.db.Beginx()
	if txErr != nil {
		return fmt.Errorf("error beginning db transaction: %s", txErr.Error())
	}

	for _, event := range logs {
		// Begin pg query string
		pgStr := fmt.Sprintf("INSERT INTO cw_%s.%s_event ", strings.ToLower(contractAddr), strings.ToLower(eventInfo.Name))
		pgStr = pgStr + "(header_id, token_name, raw_log, log_idx, tx_idx"
		el := len(event.Values)

		// Preallocate slice of needed capacity and proceed to pack variables into it in same order they appear in string
		data := make([]interface{}, 0, 5+el)
		data = append(data,
			event.ID,
			contractName,
			event.Raw,
			event.LogIndex,
			event.TransactionIndex)

		// Iterate over inputs and append name to query string and value to input data
		for inputName, input := range event.Values {
			pgStr = pgStr + fmt.Sprintf(", %s_", strings.ToLower(inputName)) // Add underscore after to avoid any collisions with reserved pg words
			data = append(data, input)
		}

		// For each input entry we created we add its postgres command variable to the string
		pgStr = pgStr + ") VALUES ($1, $2, $3, $4, $5"
		for i := 0; i < el; i++ {
			pgStr = pgStr + fmt.Sprintf(", $%d", i+6)
		}
		pgStr = pgStr + ") ON CONFLICT DO NOTHING"

		logrus.Tracef("query for inserting log: %s", pgStr)
		// Add this query to the transaction
		_, execErr := tx.Exec(pgStr, data...)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Warnf("error rolling back transactions while persisting logs: %s", rollbackErr.Error())
			}
			return fmt.Errorf("error executing query: %s", execErr.Error())
		}
	}

	return tx.Commit()
}

// CreateEventTable checks for event table and creates it if it does not already exist
// Returns true if it created a new table; returns false if table already existed
func (r *eventRepository) CreateEventTable(contractAddr string, event types.Event) (bool, error) {
	tableID := fmt.Sprintf("cw_%s.%s_event", strings.ToLower(contractAddr), strings.ToLower(event.Name))
	// Check cache before querying pq to see if table exists
	_, ok := r.tables.Get(tableID)
	if ok {
		return false, nil
	}
	tableExists, checkTableErr := r.checkForTable(contractAddr, event.Name)
	if checkTableErr != nil {
		return false, fmt.Errorf("error checking for table: %s", checkTableErr)
	}

	if !tableExists {
		createTableErr := r.newEventTable(tableID, event)
		if createTableErr != nil {
			return false, fmt.Errorf("error creating table: %s", createTableErr.Error())
		}
	}

	// Add table id to cache
	r.tables.Add(tableID, true)

	return !tableExists, nil
}

// Creates a table for the given contract and event
func (r *eventRepository) newEventTable(tableID string, event types.Event) error {
	// Begin pg string
	var pgStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ", tableID)

	pgStr = pgStr + "(id SERIAL, header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE, token_name CHARACTER VARYING(66) NOT NULL, raw_log JSONB, log_idx INTEGER NOT NULL, tx_idx INTEGER NOT NULL,"

	for _, field := range event.Fields {
		pgStr = pgStr + fmt.Sprintf(" %s_ %s NOT NULL,", strings.ToLower(field.Name), field.PgType)
	}
	pgStr = pgStr + " UNIQUE (header_id, tx_idx, log_idx))"

	_, err := r.db.Exec(pgStr)
	return err
}

// Checks if a table already exists for the given contract and event
func (r *eventRepository) checkForTable(contractAddr string, eventName string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'cw_%s' AND table_name = '%s_event')", strings.ToLower(contractAddr), strings.ToLower(eventName))

	var exists bool
	err := r.db.Get(&exists, pgStr)

	return exists, err
}

// CreateContractSchema checks for contract schema and creates it if it does not already exist
// Returns true if it created a new schema; returns false if schema already existed
func (r *eventRepository) CreateContractSchema(contractAddr string) (bool, error) {
	if contractAddr == "" {
		return false, errors.New("error: no contract address specified")
	}

	// Check cache before querying pq to see if schema exists
	_, ok := r.schemas.Get(contractAddr)
	if ok {
		return false, nil
	}
	schemaExists, checkSchemaErr := r.checkForSchema(contractAddr)
	if checkSchemaErr != nil {
		return false, fmt.Errorf("error checking for schema: %s", checkSchemaErr.Error())
	}
	if !schemaExists {
		createSchemaErr := r.newContractSchema(contractAddr)
		if createSchemaErr != nil {
			return false, fmt.Errorf("error creating schema: %s", createSchemaErr.Error())
		}
	}

	// Add schema name to cache
	r.schemas.Add(contractAddr, true)

	return !schemaExists, nil
}

// Creates a schema for the given contract
func (r *eventRepository) newContractSchema(contractAddr string) error {
	_, err := r.db.Exec("CREATE SCHEMA IF NOT EXISTS cw_" + strings.ToLower(contractAddr))

	return err
}

// Checks if a schema already exists for the given contract
func (r *eventRepository) checkForSchema(contractAddr string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT schema_name FROM information_schema.schemata WHERE schema_name = 'cw_%s')", strings.ToLower(contractAddr))

	var exists bool
	err := r.db.QueryRow(pgStr).Scan(&exists)

	return exists, err
}

// CheckSchemaCache is used to query the schema name cache
func (r *eventRepository) CheckSchemaCache(key string) (interface{}, bool) {
	return r.schemas.Get(key)
}

// CheckTableCache is used to query the table name cache
func (r *eventRepository) CheckTableCache(key string) (interface{}, bool) {
	return r.tables.Get(key)
}
