// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package dent

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type DentConverter struct{}

func NewDentConverter() DentConverter {
	return DentConverter{}
}

func (c DentConverter) ToModels(ethLogs []types.Log) (result []interface{}, err error) {
	for _, log := range ethLogs {
		err := validateLog(log)
		if err != nil {
			return nil, err
		}

		bidId := log.Topics[2].Big()
		lot := log.Topics[3].Big().String()
		bidValue := getBidValue(log)
		guy := common.HexToAddress(log.Topics[1].Hex()).String()

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
