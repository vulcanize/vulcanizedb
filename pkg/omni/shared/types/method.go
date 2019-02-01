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

type Method struct {
	Name   string
	Const  bool
	Args   []Field
	Return []Field
}

// Struct to hold instance of result from method call with given inputs and block
type Result struct {
	Method
	Inputs []interface{} // Will only use addresses
	Output interface{}
	PgType string // Holds output pg type
	Block  int64
}

// Unpack abi.Method into our custom Method struct
func NewMethod(m abi.Method) Method {
	inputs := make([]Field, len(m.Inputs))
	for i, input := range m.Inputs {
		inputs[i] = Field{}
		inputs[i].Name = input.Name
		inputs[i].Type = input.Type
		inputs[i].Indexed = input.Indexed
		switch inputs[i].Type.T {
		case abi.HashTy, abi.AddressTy:
			inputs[i].PgType = "CHARACTER VARYING(66)"
		case abi.IntTy, abi.UintTy:
			inputs[i].PgType = "DECIMAL"
		case abi.BoolTy:
			inputs[i].PgType = "BOOLEAN"
		case abi.BytesTy, abi.FixedBytesTy:
			inputs[i].PgType = "BYTEA"
		case abi.ArrayTy:
			inputs[i].PgType = "TEXT[]"
		case abi.FixedPointTy:
			inputs[i].PgType = "MONEY" // use shopspring/decimal for fixed point numbers in go and money type in postgres?
		default:
			inputs[i].PgType = "TEXT"
		}
	}

	outputs := make([]Field, len(m.Outputs))
	for i, output := range m.Outputs {
		outputs[i] = Field{}
		outputs[i].Name = output.Name
		outputs[i].Type = output.Type
		outputs[i].Indexed = output.Indexed
		switch outputs[i].Type.T {
		case abi.HashTy, abi.AddressTy:
			outputs[i].PgType = "CHARACTER VARYING(66)"
		case abi.IntTy, abi.UintTy:
			outputs[i].PgType = "DECIMAL"
		case abi.BoolTy:
			outputs[i].PgType = "BOOLEAN"
		case abi.BytesTy, abi.FixedBytesTy:
			outputs[i].PgType = "BYTEA"
		case abi.ArrayTy:
			outputs[i].PgType = "TEXT[]"
		case abi.FixedPointTy:
			outputs[i].PgType = "MONEY" // use shopspring/decimal for fixed point numbers in go and money type in postgres?
		default:
			outputs[i].PgType = "TEXT"
		}
	}

	return Method{
		Name:   m.Name,
		Const:  m.Const,
		Args:   inputs,
		Return: outputs,
	}
}

func (m Method) Sig() common.Hash {
	types := make([]string, len(m.Args))
	i := 0
	for _, arg := range m.Args {
		types[i] = arg.Type.String()
		i++
	}

	return common.BytesToHash(crypto.Keccak256([]byte(fmt.Sprintf("%v(%v)", m.Name, strings.Join(types, ",")))))
}
