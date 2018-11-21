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
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

type FrobConverter struct{}

func (FrobConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]interface{}, error) {
	var entities []interface{}
	log.Info("blah")
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
			Ilk:              string(bytes.Trim(frobEntity.Ilk[:], "\x00)")),
			Urn:              common.BytesToAddress(frobEntity.Urn[:]).String(),
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
