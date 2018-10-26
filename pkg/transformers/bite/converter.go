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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type BiteConverter struct{}

func (BiteConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]interface{}, error) {
	var entities []interface{}
	for _, ethLog := range ethLogs {
		entity := &BiteEntity{}
		address := ethLog.Address
		abi, err := geth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)

		err = contract.UnpackLog(entity, "Bite", ethLog)
		if err != nil {
			return nil, err
		}

		entity.Raw = ethLog
		entity.LogIndex = ethLog.Index
		entity.TransactionIndex = ethLog.TxIndex

		entities = append(entities, *entity)
	}

	return entities, nil
}

func (converter BiteConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var models []interface{}
	for _, entity := range entities {
		biteEntity, ok := entity.(BiteEntity)
		if !ok {
			return nil, fmt.Errorf("entity of type %T, not %T", entity, BiteEntity{})
		}

		ilk := string(bytes.Trim(biteEntity.Ilk[:], "\x00"))
		urn := common.BytesToAddress(biteEntity.Urn[:]).String()
		ink := biteEntity.Ink
		art := biteEntity.Art
		iArt := biteEntity.IArt
		tab := biteEntity.Tab
		flip := biteEntity.Flip
		logIdx := biteEntity.LogIndex
		txIdx := biteEntity.TransactionIndex
		rawLogJson, err := json.Marshal(biteEntity.Raw)
		rawLogString := string(rawLogJson)
		if err != nil {
			return nil, err
		}

		model := BiteModel{
			Ilk:              ilk,
			Urn:              urn,
			Ink:              shared.BigIntToString(ink),
			Art:              shared.BigIntToString(art),
			IArt:             shared.BigIntToString(iArt),
			Tab:              shared.BigIntToString(tab),
			NFlip:            shared.BigIntToString(flip),
			LogIndex:         logIdx,
			TransactionIndex: txIdx,
			Raw:              rawLogString,
		}
		models = append(models, model)
	}
	return models, nil
}
