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

package converter

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

// ConverterInterface is used to convert watched event logs to
// custom logs containing event input name => value maps
type ConverterInterface interface {
	Convert(watchedEvent core.WatchedEvent, event types.Event) (*types.Log, error)
	Update(info *contract.Contract)
}

// Converter is the underlying struct for the ConverterInterface
type Converter struct {
	ContractInfo *contract.Contract
}

// Update configures the converter for a specific contract
func (c *Converter) Update(info *contract.Contract) {
	c.ContractInfo = info
}

// Convert converts the given watched event log into a types.Log for the given event
func (c *Converter) Convert(watchedEvent core.WatchedEvent, event types.Event) (*types.Log, error) {
	boundContract := bind.NewBoundContract(common.HexToAddress(c.ContractInfo.Address), c.ContractInfo.ParsedAbi, nil, nil, nil)
	values := make(map[string]interface{})
	log := helpers.ConvertToLog(watchedEvent)
	err := boundContract.UnpackLogIntoMap(values, event.Name, log)
	if err != nil {
		return nil, err
	}

	strValues := make(map[string]string, len(values))
	seenAddrs := make([]interface{}, 0, len(values))
	seenHashes := make([]interface{}, 0, len(values))
	for fieldName, input := range values {
		// Postgres cannot handle custom types, resolve to strings
		switch input.(type) {
		case *big.Int:
			b := input.(*big.Int)
			strValues[fieldName] = b.String()
		case common.Address:
			a := input.(common.Address)
			strValues[fieldName] = a.String()
			seenAddrs = append(seenAddrs, a)
		case common.Hash:
			h := input.(common.Hash)
			strValues[fieldName] = h.String()
			seenHashes = append(seenHashes, h)
		case string:
			strValues[fieldName] = input.(string)
		case bool:
			strValues[fieldName] = strconv.FormatBool(input.(bool))
		case []byte:
			b := input.([]byte)
			strValues[fieldName] = hexutil.Encode(b)
			if len(b) == 32 { // collect byte arrays of size 32 as hashes
				seenHashes = append(seenHashes, common.HexToHash(strValues[fieldName]))
			}
		case byte:
			b := input.(byte)
			strValues[fieldName] = string(b)
		default:
			return nil, fmt.Errorf("error: unhandled abi type %T", input)
		}
	}

	// Only hold onto logs that pass our address filter, if any
	if c.ContractInfo.PassesEventFilter(strValues) {
		eventLog := &types.Log{
			ID:     watchedEvent.LogID,
			Values: strValues,
			Block:  watchedEvent.BlockNumber,
			Tx:     watchedEvent.TxHash,
		}

		// Cache emitted values if their caching is turned on
		if c.ContractInfo.EmittedAddrs != nil {
			c.ContractInfo.AddEmittedAddr(seenAddrs...)
		}
		if c.ContractInfo.EmittedHashes != nil {
			c.ContractInfo.AddEmittedHash(seenHashes...)
		}

		return eventLog, nil
	}

	return nil, nil
}
