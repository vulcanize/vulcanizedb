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

package chop_lump

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"math/big"
)

type Converter interface {
	ToModels(ethLogs []types.Log) ([]CatFileChopLumpModel, error)
}

type CatFileChopLumpConverter struct{}

func (CatFileChopLumpConverter) ToModels(ethLogs []types.Log) ([]CatFileChopLumpModel, error) {
	var results []CatFileChopLumpModel
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}
		ilk := string(bytes.Trim(ethLog.Topics[2].Bytes(), "\x00"))
		what := string(bytes.Trim(ethLog.Topics[3].Bytes(), "\x00"))
		dataBytes := ethLog.Data[len(ethLog.Data)-shared.DataItemLength:]
		data := big.NewInt(0).SetBytes(dataBytes).String()

		raw, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}
		result := CatFileChopLumpModel{
			Ilk:              ilk,
			What:             what,
			Data:             data,
			TransactionIndex: ethLog.TxIndex,
			Raw:              raw,
		}
		results = append(results, result)
	}
	return results, nil
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 4 {
		return errors.New("log missing topics")
	}
	if len(log.Data) < shared.DataItemLength {
		return errors.New("log missing data")
	}
	return nil
}
