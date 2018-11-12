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

package tend

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TendConverter struct{}

func (TendConverter) ToModels(ethLogs []types.Log) (results []interface{}, err error) {
	for _, ethLog := range ethLogs {
		err := validateLog(ethLog)
		if err != nil {
			return nil, err
		}

		bidId := ethLog.Topics[2].Big()
		guy := common.HexToAddress(ethLog.Topics[1].Hex()).String()
		lot := ethLog.Topics[3].Big().String()

		lastDataItemStartIndex := len(ethLog.Data) - 32
		lastItem := ethLog.Data[lastDataItemStartIndex:]
		last := big.NewInt(0).SetBytes(lastItem)
		bidValue := last.String()
		tic := "0"
		//TODO: it is likely that the tic value will need to be added to an emitted event,
		//so this will need to be updated at that point
		transactionIndex := ethLog.TxIndex
		logIndex := ethLog.Index

		rawLog, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}

		model := TendModel{
			BidId:            bidId.String(),
			Lot:              lot,
			Bid:              bidValue,
			Guy:              guy,
			Tic:              tic,
			LogIndex:         logIndex,
			TransactionIndex: transactionIndex,
			Raw:              rawLog,
		}
		results = append(results, model)
	}
	return results, err
}

func validateLog(ethLog types.Log) error {
	if len(ethLog.Data) <= 0 {
		return errors.New("tend log note data is empty")
	}

	if len(ethLog.Topics) < 4 {
		return errors.New("tend log does not contain expected topics")
	}

	return nil
}
