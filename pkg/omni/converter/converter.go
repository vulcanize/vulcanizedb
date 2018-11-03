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

package converter

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/examples/generic/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

// Converter is used to convert watched event logs to
// custom logs containing event input name => value maps
type Converter interface {
	Convert(watchedEvent core.WatchedEvent, event *types.Event) error
	Update(info types.ContractInfo)
}

type converter struct {
	contractInfo types.ContractInfo
}

func NewConverter(info types.ContractInfo) *converter {

	return &converter{
		contractInfo: info,
	}
}

func (c *converter) Update(info types.ContractInfo) {
	c.contractInfo = info
}

// Convert the given watched event log into a types.Log for the given event
func (c *converter) Convert(watchedEvent core.WatchedEvent, event *types.Event) error {
	contract := bind.NewBoundContract(common.HexToAddress(c.contractInfo.Address), c.contractInfo.ParsedAbi, nil, nil, nil)
	values := make(map[string]interface{})

	for _, field := range event.Fields {
		var i interface{}
		values[field.Name] = i

		switch field.Type.T {
		case abi.StringTy, abi.HashTy, abi.AddressTy:
			field.PgType = "CHARACTER VARYING(66)"
		case abi.IntTy, abi.UintTy:
			field.PgType = "DECIMAL"
		case abi.BoolTy:
			field.PgType = "BOOLEAN"
		case abi.BytesTy, abi.FixedBytesTy:
			field.PgType = "BYTEA"
		case abi.ArrayTy:
			field.PgType = "TEXT[]"
		case abi.FixedPointTy:
			field.PgType = "MONEY" // use shopspring/decimal for fixed point numbers in go and money type in postgres?
		case abi.FunctionTy:
			field.PgType = "TEXT"
		default:
			field.PgType = "TEXT"
		}

	}

	log := helpers.ConvertToLog(watchedEvent)
	err := contract.UnpackLogIntoMap(values, event.Name, log)
	if err != nil {
		return err
	}

	eventLog := types.Log{
		Values: values,
		Block:  watchedEvent.BlockNumber,
		Tx:     watchedEvent.TxHash,
	}

	event.Logs[watchedEvent.LogID] = eventLog

	return nil
}
