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
	ToModel(ethLog types.Log) (DealModel, error)
}

type DealConverter struct{}

func NewDealConverter() DealConverter {
	return DealConverter{}
}

func (DealConverter) ToModel(ethLog types.Log) (DealModel, error) {
	err := validateLog(ethLog)
	if err != nil {
		return DealModel{}, err
	}

	bidId := ethLog.Topics[2].Big()
	raw, err := json.Marshal(ethLog)
	if err != nil {
		return DealModel{}, err
	}

	return DealModel{
		BidId:            bidId.String(),
		TransactionIndex: ethLog.TxIndex,
		Raw:              raw,
	}, nil
}

func validateLog(ethLog types.Log) error {
	if len(ethLog.Topics) < 3 {
		return errors.New("deal log does not contain expected topics")
	}
	return nil
}
