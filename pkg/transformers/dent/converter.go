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

package dent

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Converter interface {
	ToModels(ethLogs []types.Log) ([]DentModel, error)
}

type DentConverter struct{}

func NewDentConverter() DentConverter {
	return DentConverter{}
}

func (c DentConverter) ToModels(ethLogs []types.Log) (result []DentModel, err error) {
	for _, log := range ethLogs {
		err := validateLog(log)
		if err != nil {
			return nil, err
		}

		bidId := log.Topics[2].Big()
		lot := log.Topics[3].Big().String()
		bidValue := getBidValue(log)
		guy := common.HexToAddress(log.Topics[1].Hex()).String()
		tic := "0"
		//TODO: it is likely that the tic value will need to be added to an emitted event,
		//so this will need to be updated at that point

		logIndex := log.Index
		transactionIndex := log.TxIndex

		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}

		model := DentModel{
			BidId:            bidId.String(),
			Lot:              lot,
			Bid:              bidValue,
			Guy:              guy,
			Tic:              tic,
			LogIndex:         logIndex,
			TransactionIndex: transactionIndex,
			Raw:              raw,
		}
		result = append(result, model)
	}
	return result, err
}

func validateLog(ethLog types.Log) error {
	if len(ethLog.Data) <= 0 {
		return errors.New("dent log data is empty")
	}

	if len(ethLog.Topics) < 4 {
		return errors.New("dent log does not contain expected topics")
	}

	return nil
}

func getBidValue(ethLog types.Log) string {
	itemByteLength := 32
	lastDataItemStartIndex := len(ethLog.Data) - itemByteLength
	lastItem := ethLog.Data[lastDataItemStartIndex:]
	lastValue := big.NewInt(0).SetBytes(lastItem)

	return lastValue.String()
}
