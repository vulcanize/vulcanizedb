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
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

const methodCacheSize = 1000

type MethodRepository interface {
	PersistResults(results []types.Result, methodInfo types.Method, contractAddr, contractName string) error
	CreateMethodTable(contractAddr string, method types.Method) (bool, error)
	CreateContractSchema(contractAddr string) (bool, error)
	CheckSchemaCache(key string) (interface{}, bool)
	CheckTableCache(key string) (interface{}, bool)
}

type methodRepository struct {
	*postgres.DB
	mode    types.Mode
	schemas *lru.Cache // Cache names of recently used schemas to minimize db connections
	tables  *lru.Cache // Cache names of recently used tables to minimize db connections
}

func NewMethodRepository(db *postgres.DB, mode types.Mode) *methodRepository {
	ccs, _ := lru.New(contractCacheSize)
	mcs, _ := lru.New(methodCacheSize)
	return &methodRepository{
		DB:      db,
		mode:    mode,
		schemas: ccs,
		tables:  mcs,
	}
}

// Creates a schema for the contract if needed
// Creates table for the contract method if needed
// Persists method polling data into this custom table
func (r *methodRepository) PersistResults(results []types.Result, methodInfo types.Method, contractAddr, contractName string) error {
	if len(results) == 0 {
		return errors.New("method repository error: passed empty results slice")
	}
	_, err := r.CreateContractSchema(contractAddr)
	if err != nil {
		return err
	}

	_, err = r.CreateMethodTable(contractAddr, methodInfo)
	if err != nil {
		return err
	}

	return r.persistResults(results, methodInfo, contractAddr, contractName)
}

// Creates a custom postgres command to persist logs for the given event
func (r *methodRepository) persistResults(results []types.Result, methodInfo types.Method, contractAddr, contractName string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	for _, result := range results {
		// Begin postgres string
		pgStr := fmt.Sprintf("INSERT INTO %s_%s.%s_method ", r.mode.String(), strings.ToLower(contractAddr), strings.ToLower(result.Name))
		pgStr = pgStr + "(token_name, block"
		ml := len(result.Args)

		// Preallocate slice of needed capacity and proceed to pack variables into it in same order they appear in string
		data := make([]interface{}, 0, 3+ml)
		data = append(data,
			contractName,
			result.Block)

		// Iterate over method args and return value, adding names
		// to the string and pushing values to the slice
		for i, arg := range result.Args {
			pgStr = pgStr + fmt.Sprintf(", %s_", strings.ToLower(arg.Name)) // Add underscore after to avoid any collisions with reserved pg words
			data = append(data, result.Inputs[i])
		}
		pgStr = pgStr + ", returned) VALUES ($1, $2"
		data = append(data, result.Output)

		// For each input entry we created we add its postgres command variable to the string
		for i := 0; i <= ml; i++ {
			pgStr = pgStr + fmt.Sprintf(", $%d", i+3)
		}
		pgStr = pgStr + ")"

		// Add this query to the transaction
		_, err = tx.Exec(pgStr, data...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// Checks for event table and creates it if it does not already exist
func (r *methodRepository) CreateMethodTable(contractAddr string, method types.Method) (bool, error) {
	tableID := fmt.Sprintf("%s_%s.%s_method", r.mode.String(), strings.ToLower(contractAddr), strings.ToLower(method.Name))

	// Check cache before querying pq to see if table exists
	_, ok := r.tables.Get(tableID)
	if ok {
		return false, nil
	}
	tableExists, err := r.checkForTable(contractAddr, method.Name)
	if err != nil {
		return false, err
	}
	if !tableExists {
		err = r.newMethodTable(tableID, method)
		if err != nil {
			return false, err
		}
	}

	// Add schema name to cache
	r.tables.Add(tableID, true)

	return !tableExists, nil
}

// Creates a table for the given contract and event
func (r *methodRepository) newMethodTable(tableID string, method types.Method) error {
	// Begin pg string
	pgStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ", tableID)
	pgStr = pgStr + "(id SERIAL, token_name CHARACTER VARYING(66) NOT NULL, block INTEGER NOT NULL,"

	// Iterate over method inputs and outputs, using their name and pgType to grow the string
	for _, arg := range method.Args {
		pgStr = pgStr + fmt.Sprintf(" %s_ %s NOT NULL,", strings.ToLower(arg.Name), arg.PgType)
	}

	pgStr = pgStr + fmt.Sprintf(" returned %s NOT NULL)", method.Return[0].PgType)

	_, err := r.DB.Exec(pgStr)

	return err
}

// Checks if a table already exists for the given contract and event
func (r *methodRepository) checkForTable(contractAddr string, methodName string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = '%s_%s' AND table_name = '%s_method')", r.mode.String(), strings.ToLower(contractAddr), strings.ToLower(methodName))
	var exists bool
	err := r.DB.Get(&exists, pgStr)

	return exists, err
}

// Checks for contract schema and creates it if it does not already exist
func (r *methodRepository) CreateContractSchema(contractAddr string) (bool, error) {
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
func (r *methodRepository) newContractSchema(contractAddr string) error {
	_, err := r.DB.Exec("CREATE SCHEMA IF NOT EXISTS " + r.mode.String() + "_" + strings.ToLower(contractAddr))

	return err
}

// Checks if a schema already exists for the given contract
func (r *methodRepository) checkForSchema(contractAddr string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT schema_name FROM information_schema.schemata WHERE schema_name = '%s_%s')", r.mode.String(), strings.ToLower(contractAddr))

	var exists bool
	err := r.DB.Get(&exists, pgStr)

	return exists, err
}

func (r *methodRepository) CheckSchemaCache(key string) (interface{}, bool) {
	return r.schemas.Get(key)
}

func (r *methodRepository) CheckTableCache(key string) (interface{}, bool) {
	return r.tables.Get(key)
}
