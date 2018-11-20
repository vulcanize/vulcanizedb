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
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

type MethodDatastore interface {
	PersistResult(method types.Result, contractAddr, contractName string) error
	CreateMethodTable(contractAddr string, method types.Result) (bool, error)
	CreateContractSchema(contractAddr string) (bool, error)
}

type methodDatastore struct {
	*postgres.DB
}

func NewMethodDatastore(db *postgres.DB) *methodDatastore {

	return &methodDatastore{
		DB: db,
	}
}

func (d *methodDatastore) PersistResult(method types.Result, contractAddr, contractName string) error {
	if len(method.Args) != len(method.Inputs) {
		return errors.New("error: given number of inputs does not match number of method arguments")
	}
	if len(method.Return) != 1 {
		return errors.New("error: given number of outputs does not match number of method return values")
	}

	_, err := d.CreateContractSchema(contractAddr)
	if err != nil {
		return err
	}

	_, err = d.CreateMethodTable(contractAddr, method)
	if err != nil {
		return err
	}

	return d.persistResult(method, contractAddr, contractName)
}

// Creates a custom postgres command to persist logs for the given event
func (d *methodDatastore) persistResult(method types.Result, contractAddr, contractName string) error {
	// Begin postgres string
	pgStr := fmt.Sprintf("INSERT INTO c%s.%s_method ", strings.ToLower(contractAddr), strings.ToLower(method.Name))
	pgStr = pgStr + "(token_name, block"

	// Pack the corresponding variables in a slice
	var data []interface{}
	data = append(data,
		contractName,
		method.Block)

	// Iterate over method args and return value, adding names
	// to the string and pushing values to the slice
	counter := 0 // Keep track of number of inputs
	for i, arg := range method.Args {
		counter += 1
		pgStr = pgStr + fmt.Sprintf(", %s_", strings.ToLower(arg.Name)) // Add underscore after to avoid any collisions with reserved pg words
		data = append(data, method.Inputs[i])
	}

	counter += 1
	pgStr = pgStr + ", returned) VALUES ($1, $2"
	data = append(data, method.Output)

	// For each input entry we created we add its postgres command variable to the string
	for i := 0; i < counter; i++ {
		pgStr = pgStr + fmt.Sprintf(", $%d", i+3)
	}
	pgStr = pgStr + ")"

	_, err := d.DB.Exec(pgStr, data...)
	if err != nil {
		return err
	}

	return nil
}

// Checks for event table and creates it if it does not already exist
func (d *methodDatastore) CreateMethodTable(contractAddr string, method types.Result) (bool, error) {
	tableExists, err := d.checkForTable(contractAddr, method.Name)
	if err != nil {
		return false, err
	}

	if !tableExists {
		err = d.newMethodTable(contractAddr, method)
		if err != nil {
			return false, err
		}
	}

	return !tableExists, nil
}

// Creates a table for the given contract and event
func (d *methodDatastore) newMethodTable(contractAddr string, method types.Result) error {
	// Begin pg string
	pgStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS c%s.%s_method ", strings.ToLower(contractAddr), strings.ToLower(method.Name))
	pgStr = pgStr + "(id SERIAL, token_name CHARACTER VARYING(66) NOT NULL, block INTEGER NOT NULL,"

	// Iterate over method inputs and outputs, using their name and pgType to grow the string
	for _, arg := range method.Args {
		pgStr = pgStr + fmt.Sprintf(" %s_ %s NOT NULL,", strings.ToLower(arg.Name), arg.PgType)
	}

	pgStr = pgStr + fmt.Sprintf(" returned %s NOT NULL)", method.Return[0].PgType)

	_, err := d.DB.Exec(pgStr)

	return err
}

// Checks if a table already exists for the given contract and event
func (d *methodDatastore) checkForTable(contractAddr string, methodName string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'c%s' AND table_name = '%s_method')", strings.ToLower(contractAddr), strings.ToLower(methodName))
	var exists bool
	err := d.DB.Get(&exists, pgStr)

	return exists, err
}

// Checks for contract schema and creates it if it does not already exist
func (d *methodDatastore) CreateContractSchema(contractAddr string) (bool, error) {
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
func (d *methodDatastore) newContractSchema(contractAddr string) error {
	_, err := d.DB.Exec("CREATE SCHEMA IF NOT EXISTS c" + strings.ToLower(contractAddr))

	return err
}

// Checks if a schema already exists for the given contract
func (d *methodDatastore) checkForSchema(contractAddr string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT schema_name FROM information_schema.schemata WHERE schema_name = 'c%s')", strings.ToLower(contractAddr))

	var exists bool
	err := d.DB.Get(&exists, pgStr)

	return exists, err
}
