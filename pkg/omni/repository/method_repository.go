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
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
	"strings"
)

type MethodDatastore interface {
	PersistContractMethods(con *contract.Contract) error
	PersistResults(methodName string, con *contract.Contract) error
	CreateMethodTable(contractName string, method *types.Method) error
	CreateContractSchema(contractName string) error
}

type methodDatastore struct {
	*postgres.DB
}

func NewMethodDatastore(db *postgres.DB) *methodDatastore {

	return &methodDatastore{
		DB: db,
	}
}

func (d *methodDatastore) PersistContractMethods(con *contract.Contract) error {
	err := d.CreateContractSchema(con.Name)
	if err != nil {
		return err
	}

	for _, method := range con.Methods {
		err = d.CreateMethodTable(con.Name, method)
		if err != nil {
			return err
		}

		//TODO: Persist method data

	}

	return nil
}

// Creates a custom postgres command to persist logs for the given event
func (d *methodDatastore) PersistResults(methodName string, con *contract.Contract) error {
	for _, result := range con.Methods[methodName].Results {
		println(result)
		//TODO: Persist result data
	}

	return nil
}

// Checks for event table and creates it if it does not already exist
func (d *methodDatastore) CreateMethodTable(contractAddr string, method *types.Method) error {
	tableExists, err := d.checkForTable(contractAddr, method.Name)
	if err != nil {
		return err
	}

	if !tableExists {
		err = d.newMethodTable(contractAddr, method)
		if err != nil {
			return err
		}
	}

	return nil
}

// Creates a table for the given contract and event
func (d *methodDatastore) newMethodTable(contractAddr string, method *types.Method) error {
	// Begin pg string
	pgStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS _%s.%s ", strings.ToLower(contractAddr), strings.ToLower(method.Name))
	pgStr = pgStr + "(id SERIAL, token_name CHARACTER VARYING(66) NOT NULL, block INTEGER NOT NULL,"

	// Iterate over method inputs and outputs, using their name and pgType to grow the string
	for _, input := range method.Inputs {
		pgStr = pgStr + fmt.Sprintf("%s_ %s NOT NULL,", strings.ToLower(input.Name), input.PgType)
	}

	for _, output := range method.Outputs {
		pgStr = pgStr + fmt.Sprintf(" %s_ %s NOT NULL,", strings.ToLower(output.Name), output.PgType)
	}

	pgStr = pgStr[:len(pgStr)-1] + ")"
	_, err := d.DB.Exec(pgStr)

	return err
}

// Checks if a table already exists for the given contract and event
func (d *methodDatastore) checkForTable(contractAddr string, methodName string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = '_%s' AND table_name = '%s')", strings.ToLower(contractAddr), strings.ToLower(methodName))
	var exists bool
	err := d.DB.Get(&exists, pgStr)

	return exists, err
}

// Checks for contract schema and creates it if it does not already exist
func (d *methodDatastore) CreateContractSchema(contractName string) error {
	schemaExists, err := d.checkForSchema(contractName)
	if err != nil {
		return err
	}

	if !schemaExists {
		err = d.newContractSchema(contractName)
		if err != nil {
			return err
		}
	}

	return nil
}

// Creates a schema for the given contract
func (d *methodDatastore) newContractSchema(contractAddr string) error {
	_, err := d.DB.Exec("CREATE SCHEMA IF NOT EXISTS _" + strings.ToLower(contractAddr))

	return err
}

// Checks if a schema already exists for the given contract
func (d *methodDatastore) checkForSchema(contractAddr string) (bool, error) {
	pgStr := fmt.Sprintf("SELECT EXISTS (SELECT schema_name FROM information_schema.schemata WHERE schema_name = '_%s')", strings.ToLower(contractAddr))

	var exists bool
	err := d.DB.Get(&exists, pgStr)

	return exists, err
}
