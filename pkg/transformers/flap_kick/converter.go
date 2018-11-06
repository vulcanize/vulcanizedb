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

package flap_kick

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"time"
)

type FlapKickConverter struct {
}

func (FlapKickConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]interface{}, error) {
	var entities []interface{}
	for _, ethLog := range ethLogs {
		entity := &FlapKickEntity{}
		address := ethLog.Address
		abi, err := geth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}
		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		err = contract.UnpackLog(entity, "Kick", ethLog)
		if err != nil {
			return nil, err
		}
		entity.Raw = ethLog
		entity.TransactionIndex = ethLog.TxIndex
		entity.LogIndex = ethLog.Index
		entities = append(entities, *entity)
	}
	return entities, nil
}

func (FlapKickConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var models []interface{}
	for _, entity := range entities {
		flapKickEntity, ok := entity.(FlapKickEntity)
		if !ok {
			return nil, fmt.Errorf("entity of type %T, not %T", entity, FlapKickEntity{})
		}

		if flapKickEntity.Id == nil {
			return nil, errors.New("FlapKick log ID cannot be nil.")
		}

		id := flapKickEntity.Id.String()
		lot := shared.BigIntToString(flapKickEntity.Lot)
		bid := shared.BigIntToString(flapKickEntity.Bid)
		gal := flapKickEntity.Gal.String()
		endValue := shared.BigIntToInt64(flapKickEntity.End)
		end := time.Unix(endValue, 0)
		rawLog, err := json.Marshal(flapKickEntity.Raw)
		if err != nil {
			return nil, err
		}

		model := FlapKickModel{
			BidId:            id,
			Lot:              lot,
			Bid:              bid,
			Gal:              gal,
			End:              end,
			TransactionIndex: flapKickEntity.TransactionIndex,
			LogIndex:         flapKickEntity.LogIndex,
			Raw:              rawLog,
		}
		models = append(models, model)
	}
	return models, nil
}
