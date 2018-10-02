/*
 *  Copyright 2018 Vulcanize
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package bite

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type Converter interface {
	ToEntities(contractAbi string, ethLogs []types.Log) ([]BiteEntity, error)
	ToModels(biteEntities []BiteEntity) ([]BiteModel, error)
}

type BiteConverter struct{}

func (BiteConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]BiteEntity, error) {
	var entities []BiteEntity
	for _, ethLog := range ethLogs {
		entity := BiteEntity{}
		address := ethLog.Address
		abi, err := geth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)

		err = contract.UnpackLog(&entity, "Bite", ethLog)
		if err != nil {
			return nil, err
		}

		entity.Raw = ethLog
		entity.TransactionIndex = ethLog.TxIndex

		entities = append(entities, entity)
	}

	return entities, nil
}

func (converter BiteConverter) ToModels(entities []BiteEntity) ([]BiteModel, error) {
	var models []BiteModel
	for _, entity := range entities {
		id := entity.Id
		ilk := entity.Ilk[:]
		urn := entity.Urn[:]
		ink := entity.Ink
		art := entity.Art
		iArt := entity.IArt
		tab := entity.Tab
		flip := entity.Flip
		txIdx := entity.TransactionIndex
		rawLogJson, err := json.Marshal(entity.Raw)
		rawLogString := string(rawLogJson)
		if err != nil {
			return nil, err
		}

		model := BiteModel{
			Id:               shared.BigIntToString(id),
			Ilk:              ilk,
			Urn:              urn,
			Ink:              shared.BigIntToString(ink),
			Art:              shared.BigIntToString(art),
			IArt:             shared.BigIntToString(iArt),
			Tab:              shared.BigIntToString(tab),
			Flip:             shared.BigIntToString(flip),
			TransactionIndex: txIdx,
			Raw:              rawLogString,
		}
		models = append(models, model)
	}
	return models, nil
}
