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

package vat_move

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Converter interface {
	ToModels(ethLog []types.Log) ([]VatMoveModel, error)
}

type VatMoveConverter struct{}

func (VatMoveConverter) ToModels(ethLogs []types.Log) ([]VatMoveModel, error) {
	var models []VatMoveModel
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return []VatMoveModel{}, err
		}

		src := common.BytesToAddress(ethLog.Topics[1].Bytes())
		dst := common.BytesToAddress(ethLog.Topics[2].Bytes())
		rad := ethLog.Topics[3].Big()
		raw, err := json.Marshal(ethLog)
		if err != nil {
			return []VatMoveModel{}, err
		}

		models = append(models, VatMoveModel{
			Src:              src.String(),
			Dst:              dst.String(),
			Rad:              rad.String(),
			TransactionIndex: ethLog.TxIndex,
			Raw:              raw,
		})
	}

	return models, nil
}

func verifyLog(ethLog types.Log) error {
	if len(ethLog.Data) <= 0 {
		return errors.New("log data is empty")
	}
	if len(ethLog.Topics) < 4 {
		return errors.New("log missing topics")
	}
	return nil
}
