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
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/golang-lru"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/light/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

const (
	// Number of contract address and method ids to keep in cache
	contractCacheSize = 100
	eventChacheSize   = 1000
)

// Event repository is used to persist event data into custom tables
type EventRepository interface {
	PersistLogs(logs []types.Log, eventInfo types.Event, contractAddr, contractName string) error
	CreateEventTable(contractAddr string, event types.Event) (bool, error)
	CreateContractSchema(contractName string) (bool, error)
	CheckSchemaCache(key string) (interface{}, bool)
	CheckTableCache(key string) (interface{}, bool)
}

type eventRepository struct {
	db      *postgres.DB
	mode    types.Mode
	schemas *lru.Cache // Cache names of recently used schemas to minimize db connections
	tables  *lru.Cache // Cache names of recently used tables to minimize db connections
}

func NewEventRepository(db *postgres.DB, mode types.Mode) *eventRepository {
	ccs, _ := lru.New(contractCacheSize)
	ecs, _ := lru.New(eventChacheSize)
	return &eventRepository{
		db:      db,
		mode:    mode,
		schemas: ccs,
		tables:  ecs,
	}
}

// Creates a schema for the contract if needed
// Creates table for the watched contract event if needed
// Persists converted event log data into this custom table
func (r *eventRepository) PersistLogs(logs []types.Log, eventInfo types.Event, contractAddr, contractName string) error {
	if len(logs) == 0 {
		return errors.New("event repository error: passed empty logs slice")
	}
	_, err := r.CreateContractSchema(contractAddr)
	if err != nil {
		return err
	}

	_, err = r.CreateEventTable(contractAddr, eventInfo)
	if err != nil {
		return err
	}

	return r.persistLogs(logs, eventInfo, contractAddr, contractName)
}

func (r *eventRepository) persistLogs(logs []types.Log, eventInfo types.Event, contractAddr, contractName string) error {
	var err error
	switch r.mode {
	case types.LightSync:
		err = r.persistLightSyncLogs(logs, eventInfo, contractAddr, contractName)
	case types.FullSync:
		err = r.persistFullSyncLogs(logs, eventInfo, contractAddr, contractName)
	default:
		return errors.New("event repository error: unhandled mode")
	}

	return err
}

// Creates a custom postgres command to persist logs for the given event (compatible with light synced vDB)
func (r *eventRepository) persistLightSyncLogs(logs []types.Log, eventInfo types.Event, contractAddr, contractName string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, event := range logs {
		// Begin pg query string
		pgStr := fmt.Sprintf("INSERT INTO %s_%s.%s_event ", r.mode.String(), strings.ToLower(contractAddr), strings.ToLower(eventInfo.Name))
		pgStr = pgStr + "(header_id, token_name, raw_log, log_idx, tx_idx"
		el := len(event.Values)

		// Preallocate slice of needed capacity and proceed to pack variables into it in same order they appear in string
		data := make([]interface{}, 0, 5+el)
		data = append(data,
			event.Id,
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
		pgStr = pgStr + ")"

		// Add this query to the transaction
		_, err = tx.Exec(pgStr, data...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Mark header as checked for this eventId
	eventId := strings.ToLower(eventInfo.Name + "_" + contractAddr)
	err = repository.MarkHeaderCheckedInTransaction(logs[0].Id, tx, eventId) // This assumes all logs are from same block
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Creates a custom postgres command to persist logs for the given event (compatible with fully synced vDB)
func (r *eventRepository) persistFullSyncLogs(logs []types.Log, eventInfo types.Event, contractAddr, contractName string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, event := range logs {
		pgStr := fmt.Sprintf("INSERT INTO %s_%s.%s_event ", r.mode.String(), strings.ToLower(contractAddr), strings.ToLower(eventInfo.Name))
		pgStr = pgStr + "(vulcanize_log_id, token_name, block, tx"
		el := len(event.Values)

		data := make([]interface{}, 0, 4+el)
		data = append(data,
			event.Id,
			contractName,
			event.Block,
			event.Tx)

		for inputName, input := range event.Values {
			pgStr = pgStr + fmt.Sprintf(", %s_", strings.ToLower(inputName))
			data = append(data, input)
		}

		pgStr = pgStr + ") VALUES ($1, $2, $3, $4"
		for i := 0; i < el; i++ {
			pgStr = pgStr + fmt.Sprintf(", $%d", i+5)
		}
		pgStr = pgStr + ") ON CONFLICT (vulcanize_log_id) DO NOTHING"

		_, err = tx.Exec(pgStr, data...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// Checks for event table and creates it if it does not already exist
// Returns true if it created a new table; returns false if table already existed
func (r *eventRepository) CreateEventTable(contractAddr string, event types.Event) (bool, error) {
	tableID := fmt.Sprintf("%s_%s.%s_event", r.mode.String(), strings.ToLower(contractAddr), strings.ToLower(event.Name))
	// Check cache before querying pq to see if table exists
	_, ok := r.tables.Get(tableID)
	if ok {
		return false, nil
	}
	tableExists, err := r.checkForTable(contractAddr, event.Name)
	if err != nil {
		return false, err
	}

	if !tableExists {
		err = r.newEventTable(tableID, event)
		if err != nil {
			return false, err
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
	var err error

	// Handle different modes
	switch r.mode {
	case types.FullSync:
		pgStr = pgStr + "(id SERIAL, vulcanize_log_id INTEGER NOT NULL UNIQUE, token_name CHARACTER VARYING(66) NOT NULL, block INTEGER NOT NULL, tx CHARACTER VARYING(66) NOT NULL,"

		// Iterate over event fields, using their name and pgType to grow the string
		for _, field := range event.Fields {
			pgStr = pgStr + fmt.Sprintf(" %s_ %s NOT NULL,", strings.ToLower(field.Name), field.PgType)
		}
		pgStr = pgStr + " CONSTRAINT log_index_fk FOREIGN KEY (vulcanize_log_id) REFERENCES logs (id) ON DELETE CASCADE)"
	case types.LightSync:
		pgStr = pgStr + "(id SERIAL, header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE, token_name CHARACTER VARYING(66) NOT NULL, raw_log JSONB, log_idx INTEGER NOT NULL, tx_idx INTEGER NOT NULL,"

		for _, field := range event.Fields {
			pgStr = pgStr + fmt.Sprintf(" %s_ %s NOT NULL,", strings.ToLower(field.Name), field.PgType)
		}
		pgStr = pgStr + " UNIQUE (header_id, tx_idx, log_idx))"
	default:
		return errors.New("unhandled repository mode")
	}

	_, err = r.db.Exec(pgStr)

	return err
}

// Checks if a table already exists for the given contract and event
func (r *eventRepository) checkForTable(contractAddr string, eventName string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = '%s_%s' AND table_name = '%s_event')", r.mode.String(), strings.ToLower(contractAddr), strings.ToLower(eventName))

	var exists bool
	err := r.db.Get(&exists, pgStr)

	return exists, err
}

// Checks for contract schema and creates it if it does not already exist
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
	schemaExists, err := r.checkForSchema(contractAddr)
	if err != nil {
		return false, err
	}
	if !schemaExists {
		err = r.newContractSchema(contractAddr)
		if err != nil {
			return false, err
		}
	}

	// Add schema name to cache
	r.schemas.Add(contractAddr, true)

	return !schemaExists, nil
}

// Creates a schema for the given contract
func (r *eventRepository) newContractSchema(contractAddr string) error {
	_, err := r.db.Exec("CREATE SCHEMA IF NOT EXISTS " + r.mode.String() + "_" + strings.ToLower(contractAddr))

	return err
}

// Checks if a schema already exists for the given contract
func (r *eventRepository) checkForSchema(contractAddr string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT schema_name FROM information_schema.schemata WHERE schema_name = '%s_%s')", r.mode.String(), strings.ToLower(contractAddr))

	var exists bool
	err := r.db.QueryRow(pgStr).Scan(&exists)

	return exists, err
}

func (r *eventRepository) CheckSchemaCache(key string) (interface{}, bool) {
	return r.schemas.Get(key)
}

func (r *eventRepository) CheckTableCache(key string) (interface{}, bool) {
	return r.tables.Get(key)
}
