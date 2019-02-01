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

package types

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Event struct {
	Name      string
	Anonymous bool
	Fields    []Field
}

type Field struct {
	abi.Argument        // Name, Type, Indexed
	PgType       string // Holds type used when committing data held in this field to postgres
}

// Struct to hold instance of an event log data
type Log struct {
	Id     int64             // VulcanizeIdLog for full sync and header ID for light sync omni watcher
	Values map[string]string // Map of event input names to their values

	// Used for full sync only
	Block int64
	Tx    string

	// Used for lightSync only
	LogIndex         uint
	TransactionIndex uint
	Raw              []byte // json.Unmarshalled byte array of geth/core/types.Log{}
}

// Unpack abi.Event into our custom Event struct
func NewEvent(e abi.Event) Event {
	fields := make([]Field, len(e.Inputs))
	for i, input := range e.Inputs {
		fields[i] = Field{}
		fields[i].Name = input.Name
		fields[i].Type = input.Type
		fields[i].Indexed = input.Indexed
		// Fill in pg type based on abi type
		switch fields[i].Type.T {
		case abi.HashTy, abi.AddressTy:
			fields[i].PgType = "CHARACTER VARYING(66)"
		case abi.IntTy, abi.UintTy:
			fields[i].PgType = "DECIMAL"
		case abi.BoolTy:
			fields[i].PgType = "BOOLEAN"
		case abi.BytesTy, abi.FixedBytesTy:
			fields[i].PgType = "BYTEA"
		case abi.ArrayTy:
			fields[i].PgType = "TEXT[]"
		case abi.FixedPointTy:
			fields[i].PgType = "MONEY" // use shopspring/decimal for fixed point numbers in go and money type in postgres?
		default:
			fields[i].PgType = "TEXT"
		}
	}

	return Event{
		Name:      e.Name,
		Anonymous: e.Anonymous,
		Fields:    fields,
	}
}

func (e Event) Sig() common.Hash {
	types := make([]string, len(e.Fields))

	for i, input := range e.Fields {
		types[i] = input.Type.String()
	}

	return common.BytesToHash(crypto.Keccak256([]byte(fmt.Sprintf("%v(%v)", e.Name, strings.Join(types, ",")))))
}
