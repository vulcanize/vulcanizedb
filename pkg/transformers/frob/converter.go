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

package frob

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

type FrobConverter struct{}

func (FrobConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]interface{}, error) {
	var entities []interface{}
	for _, ethLog := range ethLogs {
		entity := FrobEntity{}
		address := ethLog.Address
		abi, err := geth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}
		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		err = contract.UnpackLog(&entity, "Frob", ethLog)
		if err != nil {
			return entities, err
		}
		entity.LogIndex = ethLog.Index
		entity.TransactionIndex = ethLog.TxIndex
		entity.Raw = ethLog
		entities = append(entities, entity)
	}

	return entities, nil
}

func (FrobConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var models []interface{}
	for _, entity := range entities {
		frobEntity, ok := entity.(FrobEntity)
		if !ok {
			return nil, fmt.Errorf("entity of type %T, not %T", entity, FrobEntity{})
		}

		rawLog, err := json.Marshal(frobEntity.Raw)
		if err != nil {
			return nil, err
		}
		model := FrobModel{
			Ilk:              shared.GetHexWithoutPrefix(frobEntity.Ilk[:]),
			Urn:              shared.GetHexWithoutPrefix(frobEntity.Urn[:]),
			Ink:              frobEntity.Ink.String(),
			Art:              frobEntity.Art.String(),
			Dink:             frobEntity.Dink.String(),
			Dart:             frobEntity.Dart.String(),
			IArt:             frobEntity.IArt.String(),
			LogIndex:         frobEntity.LogIndex,
			TransactionIndex: frobEntity.TransactionIndex,
			Raw:              rawLog,
		}
		models = append(models, model)
	}
	return models, nil
}
