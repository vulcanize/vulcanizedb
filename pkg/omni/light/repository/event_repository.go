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

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

type EventRepository interface {
	PersistLog(event types.Log, contractAddr, contractName string) error
	CreateEventTable(contractName string, event types.Log) (bool, error)
	CreateContractSchema(contractName string) (bool, error)
}

type eventRepository struct {
	db *postgres.DB
}

func NewEventRepository(db *postgres.DB) *eventRepository {
	return &eventRepository{
		db: db,
	}
}

// Creates a schema for the contract if needed
// Creates table for the watched contract event if needed
// Persists converted event log data into this custom table
func (r *eventRepository) PersistLog(event types.Log, contractAddr, contractName string) error {
	_, err := r.CreateContractSchema(contractAddr)
	if err != nil {
		return err
	}

	_, err = r.CreateEventTable(contractAddr, event)
	if err != nil {
		return err
	}

	return r.persistLog(event, contractAddr, contractName)
}

// Creates a custom postgres command to persist logs for the given event
func (r *eventRepository) persistLog(event types.Log, contractAddr, contractName string) error {
	// Begin postgres string
	pgStr := fmt.Sprintf("INSERT INTO l%s.%s_event ", strings.ToLower(contractAddr), strings.ToLower(event.Name))
	pgStr = pgStr + "(header_id, token_name, raw_log, log_idx, tx_idx"

	// Pack the corresponding variables in a slice
	var data []interface{}
	data = append(data,
		event.Id,
		contractName,
		event.Raw,
		event.LogIndex,
		event.TransactionIndex)

	// Iterate over name-value pairs in the log adding
	// names to the string and pushing values to the slice
	counter := 0 // Keep track of number of inputs
	for inputName, input := range event.Values {
		counter += 1
		pgStr = pgStr + fmt.Sprintf(", %s_", strings.ToLower(inputName)) // Add underscore after to avoid any collisions with reserved pg words
		data = append(data, input)
	}

	// Finish off the string and execute the command using the packed data
	// For each input entry we created we add its postgres command variable to the string
	pgStr = pgStr + ") VALUES ($1, $2, $3, $4, $5"
	for i := 0; i < counter; i++ {
		pgStr = pgStr + fmt.Sprintf(", $%d", i+6)
	}
	pgStr = pgStr + ")"

	_, err := r.db.Exec(pgStr, data...)
	if err != nil {
		return err
	}

	return nil
}

// Checks for event table and creates it if it does not already exist
func (r *eventRepository) CreateEventTable(contractAddr string, event types.Log) (bool, error) {
	tableExists, err := r.checkForTable(contractAddr, event.Name)
	if err != nil {
		return false, err
	}

	if !tableExists {
		err = r.newEventTable(contractAddr, event)
		if err != nil {
			return false, err
		}
	}

	return !tableExists, nil
}

// Creates a table for the given contract and event
func (r *eventRepository) newEventTable(contractAddr string, event types.Log) error {
	// Begin pg string
	pgStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS l%s.%s_event ", strings.ToLower(contractAddr), strings.ToLower(event.Name))
	pgStr = pgStr + "(id SERIAL, header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE, token_name CHARACTER VARYING(66) NOT NULL, raw_log JSONB, log_idx INTEGER NOT NULL, tx_idx INTEGER NOT NULL,"

	// Iterate over event fields, using their name and pgType to grow the string
	for _, field := range event.Fields {
		pgStr = pgStr + fmt.Sprintf(" %s_ %s NOT NULL,", strings.ToLower(field.Name), field.PgType)
	}

	pgStr = pgStr + " UNIQUE (header_id, tx_idx, log_idx))"
	_, err := r.db.Exec(pgStr)

	return err
}

// Checks if a table already exists for the given contract and event
func (r *eventRepository) checkForTable(contractAddr string, eventName string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'l%s' AND table_name = '%s_event')", strings.ToLower(contractAddr), strings.ToLower(eventName))

	var exists bool
	err := r.db.Get(&exists, pgStr)

	return exists, err
}

// Checks for contract schema and creates it if it does not already exist
func (r *eventRepository) CreateContractSchema(contractAddr string) (bool, error) {
	if contractAddr == "" {
		return false, errors.New("error: no contract address specified")
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

	return !schemaExists, nil
}

// Creates a schema for the given contract
func (r *eventRepository) newContractSchema(contractAddr string) error {
	_, err := r.db.Exec("CREATE SCHEMA IF NOT EXISTS l" + strings.ToLower(contractAddr))

	return err
}

// Checks if a schema already exists for the given contract
func (r *eventRepository) checkForSchema(contractAddr string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT schema_name FROM information_schema.schemata WHERE schema_name = 'l%s')", strings.ToLower(contractAddr))

	var exists bool
	err := r.db.QueryRow(pgStr).Scan(&exists)

	return exists, err
}
