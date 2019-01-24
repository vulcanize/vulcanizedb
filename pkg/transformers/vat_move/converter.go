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

package vat_move

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type VatMoveConverter struct{}

func (VatMoveConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var models []interface{}
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return []interface{}{}, err
		}

		src := common.BytesToAddress(ethLog.Topics[1].Bytes())
		dst := common.BytesToAddress(ethLog.Topics[2].Bytes())
		rad := ethLog.Topics[3].Big()
		raw, err := json.Marshal(ethLog)
		if err != nil {
			return []interface{}{}, err
		}

		models = append(models, VatMoveModel{
			Src:              src.String(),
			Dst:              dst.String(),
			Rad:              rad.String(),
			LogIndex:         ethLog.Index,
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
