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

package flop_kick

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type FlopKickConverter struct{}

func (FlopKickConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]interface{}, error) {
	var results []interface{}
	for _, ethLog := range ethLogs {
		entity := Entity{}
		address := ethLog.Address
		abi, err := geth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)

		err = contract.UnpackLog(&entity, "Kick", ethLog)
		if err != nil {
			return nil, err
		}
		entity.Raw = ethLog
		entity.TransactionIndex = ethLog.TxIndex
		entity.LogIndex = ethLog.Index
		results = append(results, entity)
	}
	return results, nil
}

func (FlopKickConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var results []interface{}
	for _, entity := range entities {
		flopKickEntity, ok := entity.(Entity)
		if !ok {
			return nil, fmt.Errorf("entity of type %T, not %T", entity, Entity{})
		}

		endValue := shared.BigIntToInt64(flopKickEntity.End)
		rawLogJson, err := json.Marshal(flopKickEntity.Raw)
		if err != nil {
			return nil, err
		}

		model := Model{
			BidId:            shared.BigIntToString(flopKickEntity.Id),
			Lot:              shared.BigIntToString(flopKickEntity.Lot),
			Bid:              shared.BigIntToString(flopKickEntity.Bid),
			Gal:              flopKickEntity.Gal.String(),
			End:              time.Unix(endValue, 0),
			TransactionIndex: flopKickEntity.TransactionIndex,
			LogIndex:         flopKickEntity.LogIndex,
			Raw:              rawLogJson,
		}
		results = append(results, model)
	}

	return results, nil
}
