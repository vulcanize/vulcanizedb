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

package vat_heal

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"strconv"
)

type Converter interface {
	ToModels(ethLogs []types.Log) ([]VatHealModel, error)
}

type VatHealConverter struct{}

func (VatHealConverter) ToModels(ethLogs []types.Log) ([]VatHealModel, error) {
	var models []VatHealModel
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}

		urn := common.BytesToAddress(ethLog.Topics[1].Bytes()[:common.AddressLength])
		v := common.BytesToAddress(ethLog.Topics[2].Bytes()[:common.AddressLength])
		radInt, err := strconv.ParseInt(ethLog.Topics[3].Hex(), 0, 64)
		if err != nil {
			return nil, err
		}

		rawLogJson, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}

		model := VatHealModel{
			Urn:              urn.String(),
			V:                v.String(),
			Rad:              int(radInt),
			TransactionIndex: ethLog.TxIndex,
			Raw:              rawLogJson,
		}

		models = append(models, model)
	}

	return models, nil
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 4 {
		return errors.New("log missing topics")
	}
	return nil
}
