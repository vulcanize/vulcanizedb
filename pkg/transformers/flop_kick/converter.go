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

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type Converter interface {
	ToEntities(contractAbi string, ethLogs []types.Log) ([]Entity, error)
	ToModels(entities []Entity) ([]Model, error)
}

type FlopKickConverter struct{}

func (FlopKickConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]Entity, error) {
	var results []Entity
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

func (FlopKickConverter) ToModels(entities []Entity) ([]Model, error) {
	var results []Model
	for _, entity := range entities {
		endValue := shared.BigIntToInt64(entity.End)
		rawLogJson, err := json.Marshal(entity.Raw)
		if err != nil {
			return nil, err
		}

		model := Model{
			BidId:            shared.BigIntToString(entity.Id),
			Lot:              shared.BigIntToString(entity.Lot),
			Bid:              shared.BigIntToString(entity.Bid),
			Gal:              entity.Gal.String(),
			End:              time.Unix(endValue, 0),
			TransactionIndex: entity.TransactionIndex,
			LogIndex:         entity.LogIndex,
			Raw:              rawLogJson,
		}
		results = append(results, model)
	}

	return results, nil
}
