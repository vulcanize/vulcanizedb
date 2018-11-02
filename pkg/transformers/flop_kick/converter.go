// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
