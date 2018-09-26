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

package stability_fee

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Converter interface {
	ToModel(ethLog types.Log) (PitFileStabilityFeeModel, error)
}

type PitFileStabilityFeeConverter struct{}

func (PitFileStabilityFeeConverter) ToModel(ethLog types.Log) (PitFileStabilityFeeModel, error) {
	err := verifyLog(ethLog)
	if err != nil {
		return PitFileStabilityFeeModel{}, err
	}
	what := string(bytes.Trim(ethLog.Topics[2].Bytes(), "\x00"))
	data := common.HexToAddress(ethLog.Topics[3].String()).Hex()

	raw, err := json.Marshal(ethLog)
	return PitFileStabilityFeeModel{
		What:             what,
		Data:             data,
		TransactionIndex: ethLog.TxIndex,
		Raw:              raw,
	}, err
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 3 {
		return errors.New("log missing topics")
	}
	return nil
}
