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
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

// Event datastore is used to persist event data into custom tables
type EventDatastore interface {
	PersistContractEvents(con *contract.Contract) error
	CreateEventTable(contractName string, event *types.Event) (bool, error)
	CreateContractSchema(contractName string) (bool, error)
}

type eventDatastore struct {
	*postgres.DB
}

func NewEventDataStore(db *postgres.DB) *eventDatastore {

	return &eventDatastore{
		DB: db,
	}
}

// Creates a schema for the contract
// Creates tables for the watched contract events
// Persists converted event log data into these custom tables
// Edit this method to accept a single event
func (d *eventDatastore) PersistContractEvents(con *contract.Contract) error {
	_, err := d.CreateContractSchema(con.Address)
	if err != nil {
		return err
	}

	for eventName := range con.Filters {
		event := con.Events[eventName]

		// Move to next event if there are no logs to process
		if len(event.Logs) == 0 {
			break
		}

		// Create a pg table for the event if one does not already exist
		_, err = d.CreateEventTable(con.Address, event)
		if err != nil {
			return err
		}

		// Persist event log data in this table
		err = d.persistLogs(eventName, con)
		if err != nil {
			return err
		}

		// Clear logs when we are done using them
		event.Logs = map[int64]types.Log{}
	}

	return nil
}

// Creates a custom postgres command to persist logs for the given event
func (d *eventDatastore) persistLogs(eventName string, con *contract.Contract) error {
	for id, log := range con.Events[eventName].Logs {
		// Begin postgres string
		pgStr := fmt.Sprintf("INSERT INTO c%s.%s ", strings.ToLower(con.Address), strings.ToLower(eventName))
		pgStr = pgStr + "(vulcanize_log_id, token_name, block, tx"

		// Pack the corresponding variables in a slice
		var data []interface{}
		data = append(data,
			id,
			con.Name,
			log.Block,
			log.Tx)

		// Iterate over name-value pairs in the log adding
		// names to the string and pushing values to the slice
		counter := 0 // Keep track of number of inputs
		for inputName, input := range log.Values {
			counter += 1
			pgStr = pgStr + fmt.Sprintf(", %s_", strings.ToLower(inputName)) // Add underscore after to avoid any collisions with reserved pg words
			data = append(data, input)
		}

		// Finish off the string and execute the command using the packed data
		// For each input entry we created we add its postgres command variable to the string
		pgStr = pgStr + ") VALUES ($1, $2, $3, $4"
		for i := 0; i < counter; i++ {
			pgStr = pgStr + fmt.Sprintf(", $%d", i+5)
		}
		pgStr = pgStr + ") ON CONFLICT (vulcanize_log_id) DO NOTHING"

		_, err := d.DB.Exec(pgStr, data...)
		if err != nil {
			return err
		}
	}

	return nil
}

// Checks for event table and creates it if it does not already exist
func (d *eventDatastore) CreateEventTable(contractName string, event *types.Event) (bool, error) {
	tableExists, err := d.checkForTable(contractName, event.Name)
	if err != nil {
		return false, err
	}

	if !tableExists {
		err = d.newEventTable(contractName, event)
		if err != nil {
			return false, err
		}
	}

	return !tableExists, nil
}

// Creates a table for the given contract and event
func (d *eventDatastore) newEventTable(contractAddr string, event *types.Event) error {
	// Begin pg string
	pgStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS c%s.%s ", strings.ToLower(contractAddr), strings.ToLower(event.Name))
	pgStr = pgStr + "(id SERIAL, vulcanize_log_id INTEGER NOT NULL UNIQUE, token_name CHARACTER VARYING(66) NOT NULL, block INTEGER NOT NULL, tx CHARACTER VARYING(66) NOT NULL,"

	// Iterate over event fields, using their name and pgType to grow the string
	for _, field := range event.Fields {
		pgStr = pgStr + fmt.Sprintf(" %s_ %s NOT NULL,", strings.ToLower(field.Name), field.PgType)
	}

	pgStr = pgStr + " CONSTRAINT log_index_fk FOREIGN KEY (vulcanize_log_id) REFERENCES logs (id) ON DELETE CASCADE)"
	_, err := d.DB.Exec(pgStr)

	return err
}

// Checks if a table already exists for the given contract and event
func (d *eventDatastore) checkForTable(contractAddr string, eventName string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'c%s' AND table_name = '%s')", strings.ToLower(contractAddr), strings.ToLower(eventName))

	var exists bool
	err := d.DB.Get(&exists, pgStr)

	return exists, err
}

// Checks for contract schema and creates it if it does not already exist
func (d *eventDatastore) CreateContractSchema(contractAddr string) (bool, error) {
	if contractAddr == "" {
		return false, errors.New("error: no contract address specified")
	}
	schemaExists, err := d.checkForSchema(contractAddr)
	if err != nil {
		return false, err
	}

	if !schemaExists {
		err = d.newContractSchema(contractAddr)
		if err != nil {
			return false, err
		}
	}

	return !schemaExists, nil
}

// Creates a schema for the given contract
func (d *eventDatastore) newContractSchema(contractAddr string) error {
	_, err := d.DB.Exec("CREATE SCHEMA IF NOT EXISTS c" + strings.ToLower(contractAddr))

	return err
}

// Checks if a schema already exists for the given contract
func (d *eventDatastore) checkForSchema(contractAddr string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT schema_name FROM information_schema.schemata WHERE schema_name = 'c%s')", strings.ToLower(contractAddr))

	var exists bool
	err := d.DB.QueryRow(pgStr).Scan(&exists)

	return exists, err
}
