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

package flip_kick

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type FlipKickConverter struct{}

func (FlipKickConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]interface{}, error) {
	var entities []interface{}
	for _, ethLog := range ethLogs {
		entity := &FlipKickEntity{}
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

func (FlipKickConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var models []interface{}
	for _, entity := range entities {
		flipKickEntity, ok := entity.(FlipKickEntity)
		if !ok {
			return nil, fmt.Errorf("entity of type %T, not %T", entity, FlipKickEntity{})
		}

		if flipKickEntity.Id == nil {
			return nil, errors.New("FlipKick log ID cannot be nil.")
		}

		id := flipKickEntity.Id.String()
		lot := shared.BigIntToString(flipKickEntity.Lot)
		bid := shared.BigIntToString(flipKickEntity.Bid)
		gal := flipKickEntity.Gal.String()
		endValue := shared.BigIntToInt64(flipKickEntity.End)
		end := time.Unix(endValue, 0)
		urn := common.BytesToAddress(flipKickEntity.Urn[:common.AddressLength]).String()
		tab := shared.BigIntToString(flipKickEntity.Tab)
		rawLog, err := json.Marshal(flipKickEntity.Raw)
		if err != nil {
			return nil, err
		}

		model := FlipKickModel{
			BidId:            id,
			Lot:              lot,
			Bid:              bid,
			Gal:              gal,
			End:              end,
			Urn:              urn,
			Tab:              tab,
			TransactionIndex: flipKickEntity.TransactionIndex,
			LogIndex:         flipKickEntity.LogIndex,
			Raw:              rawLog,
		}
		models = append(models, model)
	}
	return models, nil
}
