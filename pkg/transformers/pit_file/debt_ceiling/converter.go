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

package debt_ceiling

import (
	"encoding/json"
	"math/big"

	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type Converter interface {
	ToModel(ethLog types.Log) (PitFileDebtCeilingModel, error)
}

type PitFileDebtCeilingConverter struct{}

func (PitFileDebtCeilingConverter) ToModel(ethLog types.Log) (PitFileDebtCeilingModel, error) {
	err := verifyLog(ethLog)
	if err != nil {
		return PitFileDebtCeilingModel{}, err
	}
	what := common.HexToAddress(ethLog.Topics[1].String()).String()
	riskBytes := ethLog.Data[len(ethLog.Data)-shared.DataItemLength:]
	data := big.NewInt(0).SetBytes(riskBytes).String()

	raw, err := json.Marshal(ethLog)
	return PitFileDebtCeilingModel{
		What:             what,
		Data:             data,
		TransactionIndex: ethLog.TxIndex,
		Raw:              raw,
	}, err
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 2 {
		return errors.New("log missing topics")
	}
	if len(log.Data) < shared.DataItemLength {
		return errors.New("log missing data")
	}
	return nil
}
