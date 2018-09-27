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

package deal

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
)

type Converter interface {
	ToModels(ethLog []types.Log) ([]DealModel, error)
}

type DealConverter struct{}

func NewDealConverter() DealConverter {
	return DealConverter{}
}

func (DealConverter) ToModels(ethLogs []types.Log) (result []DealModel, err error) {
	for _, log := range ethLogs {
		err := validateLog(log)
		if err != nil {
			return nil, err
		}

		bidId := log.Topics[2].Big()
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}

		model := DealModel{
			BidId:            bidId.String(),
			ContractAddress:  log.Address.Hex(),
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		}
		result = append(result, model)
	}

	return result, nil
}

func validateLog(ethLog types.Log) error {
	if len(ethLog.Topics) < 3 {
		return errors.New("deal log does not contain expected topics")
	}
	return nil
}
