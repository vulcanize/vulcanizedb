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

package frob

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

type Converter interface {
	ToEntities(contractAbi string, ethLogs []types.Log) ([]FrobEntity, error)
	ToModels(entities []FrobEntity) ([]FrobModel, error)
}

type FrobConverter struct{}

func (FrobConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]FrobEntity, error) {
	var entities []FrobEntity
	for _, ethLog := range ethLogs {
		entity := FrobEntity{}
		address := ethLog.Address
		abi, err := geth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}
		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		err = contract.UnpackLog(&entity, "Frob", ethLog)
		entity.TransactionIndex = ethLog.TxIndex
		entity.Raw = ethLog
		entities = append(entities, entity)
	}

	return entities, nil
}

func (FrobConverter) ToModels(entities []FrobEntity) ([]FrobModel, error) {
	var models []FrobModel
	for _, entity := range entities {
		rawLog, err := json.Marshal(entity.Raw)
		if err != nil {
			return nil, err
		}
		model := FrobModel{
			Ilk:              common.BytesToAddress(entity.Ilk[:common.AddressLength]).String(),
			Urn:              common.BytesToAddress(entity.Urn[:common.AddressLength]).String(),
			Ink:              entity.Ink.String(),
			Art:              entity.Art.String(),
			Dink:             entity.Dink.String(),
			Dart:             entity.Dart.String(),
			IArt:             entity.IArt.String(),
			TransactionIndex: entity.TransactionIndex,
			Raw:              rawLog,
		}
		models = append(models, model)
	}
	return models, nil
}
