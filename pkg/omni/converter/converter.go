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

package converter

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

// Converter is used to convert watched event logs to
// custom logs containing event input name => value maps
type Converter interface {
	Convert(watchedEvent core.WatchedEvent, event *types.Event) error
	Update(info *contract.Contract)
}

type converter struct {
	contractInfo *contract.Contract
}

func NewConverter(info *contract.Contract) *converter {

	return &converter{
		contractInfo: info,
	}
}

func (c *converter) Update(info *contract.Contract) {
	c.contractInfo = info
}

func (c *converter) CheckInfo() *contract.Contract {
	return c.contractInfo
}

// Convert the given watched event log into a types.Log for the given event
func (c *converter) Convert(watchedEvent core.WatchedEvent, event *types.Event) error {
	contract := bind.NewBoundContract(common.HexToAddress(c.contractInfo.Address), c.contractInfo.ParsedAbi, nil, nil, nil)
	values := make(map[string]interface{})

	for _, field := range event.Fields {
		var i interface{}
		values[field.Name] = i
	}

	log := helpers.ConvertToLog(watchedEvent)
	err := contract.UnpackLogIntoMap(values, event.Name, log)
	if err != nil {
		return err
	}

	strValues := make(map[string]string, len(values))

	for fieldName, input := range values {
		// Postgres cannot handle custom types, resolve to strings
		switch input.(type) {
		case *big.Int:
			var b *big.Int
			b = input.(*big.Int)
			strValues[fieldName] = b.String()
		case common.Address:
			var a common.Address
			a = input.(common.Address)
			strValues[fieldName] = a.String()
			c.contractInfo.AddTokenHolderAddress(a.String()) // cache address in a list of contract's token holder addresses
		case common.Hash:
			var h common.Hash
			h = input.(common.Hash)
			strValues[fieldName] = h.String()
		case string:
			strValues[fieldName] = input.(string)
		case bool:
			strValues[fieldName] = strconv.FormatBool(input.(bool))
		default:
			return errors.New("error: unhandled abi type")
		}
	}

	// Only hold onto logs that pass our address filter, if any
	// Persist log here and don't hold onto it
	if c.contractInfo.PassesEventFilter(strValues) {
		eventLog := types.Log{
			Id:     watchedEvent.LogID,
			Values: strValues,
			Block:  watchedEvent.BlockNumber,
			Tx:     watchedEvent.TxHash,
		}

		event.Logs[watchedEvent.LogID] = eventLog
	}

	return nil
}
