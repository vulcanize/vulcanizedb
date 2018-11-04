// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repository

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

// Repository is used to
type DataStore interface {
	PersistEvents(info *contract.Contract) error
}

type dataStore struct {
	*postgres.DB
}

func NewDataStore(db *postgres.DB) *dataStore {

	return &dataStore{
		DB: db,
	}
}

// Creates a schema for the contract
// Creates tables for the watched contract events
// Persists converted event log data into these custom tables
func (d *dataStore) PersistEvents(con *contract.Contract) error {

	schemaExists, err := d.CheckForSchema(con.Name)
	if err != nil {
		return err
	}

	if !schemaExists {
		err = d.CreateContractSchema(con.Name)
		if err != nil {
			return err
		}
	}

	for eventName := range con.Filters {
		event := con.Events[eventName]
		if len(event.Logs) == 0 {
			break
		}

		tableExists, err := d.CheckForTable(con.Name, eventName)
		if err != nil {
			return err
		}

		if !tableExists {
			err = d.CreateEventTable(con.Name, event)
			if err != nil {
				return err
			}
		}

		for id, log := range event.Logs {
			// Create postgres command to persist any given event
			pgStr := fmt.Sprintf("INSERT INTO %s.%s ", strings.ToLower(con.Name), strings.ToLower(eventName))
			pgStr = pgStr + "(vulcanize_log_id, token_name, token_address, event_name, block, tx"
			var data []interface{}
			data = append(data,
				id,
				strings.ToLower(con.Name),
				strings.ToLower(con.Address),
				strings.ToLower(eventName),
				log.Block,
				log.Tx)

			counter := 0
			for inputName, input := range log.Values {
				// postgres cannot handle custom types, resolve to strings
				switch input.(type) {
				case big.Int:
					var b big.Int
					b = input.(big.Int)
					input = b.String()
				case *big.Int:
					var b *big.Int
					b = input.(*big.Int)
					input = b.String()
				case common.Address:
					var a common.Address
					a = input.(common.Address)
					input = a.String()         // this also gives us a chance to add any event emitted address
					con.AddAddress(a.String()) // to a list of token-related addresses, growing it as we go
				case common.Hash:
					var h common.Hash
					h = input.(common.Hash)
					input = h.String()
				}

				counter += 1
				pgStr = pgStr + fmt.Sprintf(", _%s", strings.ToLower(inputName))
				data = append(data, input)
			}

			appendStr := "VALUES ($1, $2, $3, $4, $5, $6"

			for i := 0; i < counter; i++ {
				appendStr = appendStr + fmt.Sprintf(", $%d", i+7)
			}

			appendStr = appendStr + ") "
			appendStr = appendStr + "ON CONFLICT (vulcanize_log_id) DO NOTHING"
			pgStr = pgStr + fmt.Sprintf(") %s", appendStr)

			_, err := d.DB.Exec(pgStr, data...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Creates a table for the given contract and event
func (d *dataStore) CreateEventTable(contractName string, event *types.Event) error {
	// Create postgres command to create table for any given event
	pgStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s ", strings.ToLower(contractName), strings.ToLower(event.Name))
	pgStr = pgStr + "(id SERIAL, vulcanize_log_id INTEGER NOT NULL UNIQUE, token_name CHARACTER VARYING(66) NOT NULL, token_address CHARACTER VARYING(66) NOT NULL, event_name CHARACTER VARYING(66) NOT NULL, block INTEGER NOT NULL, tx CHARACTER VARYING(66) NOT NULL,"
	for _, field := range event.Fields {
		pgStr = pgStr + fmt.Sprintf(" _%s %s NOT NULL,", field.Name, field.PgType)
	}
	pgStr = pgStr + " CONSTRAINT log_index_fk FOREIGN KEY (vulcanize_log_id) REFERENCES logs (id) ON DELETE CASCADE)"
	_, err := d.DB.Exec(pgStr)

	return err
}

// Checks if a table already exists for the given contract and event
func (d *dataStore) CheckForTable(contractName string, eventName string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = '%s' AND table_name = '%s')", contractName, eventName)
	var exists bool
	err := d.DB.Get(&exists, pgStr)

	return exists, err
}

// Creates a schema for the given contract
func (d *dataStore) CreateContractSchema(contractName string) error {
	_, err := d.DB.Exec("CREATE SCHEMA IF NOT EXISTS " + contractName)

	return err
}

// Checks if a schema already exists for the given contract
func (d *dataStore) CheckForSchema(contractName string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT schema_name FROM information_schema.schemata WHERE schema_name = '%s')", contractName)

	var exists bool
	err := d.DB.Get(&exists, pgStr)

	return exists, err
}
