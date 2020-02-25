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

package wasm

import (
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
)

// Instantiator is used to instantiate WASM functions in Postgres
type Instantiator struct {
	db        *postgres.DB
	instances []WasmFunction // WASM file paths and namespaces
}

type WasmFunction struct {
	BinaryPath string
	Namespace  string
}

// NewWASMInstantiator returns a pointer to a new Instantiator
func NewWASMInstantiator(db *postgres.DB, instances []WasmFunction) *Instantiator {
	return &Instantiator{
		db:        db,
		instances: instances,
	}
}

// Instantiate is used to load the WASM functions into Postgres
func (i *Instantiator) Instantiate() error {
	tx, err := i.db.Beginx()
	if err != nil {
		return err
	}
	for _, pn := range i.instances {
		_, err := tx.Exec(`SELECT wasm_new_instance('$1', '$2')`, pn.BinaryPath, pn.Namespace)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
