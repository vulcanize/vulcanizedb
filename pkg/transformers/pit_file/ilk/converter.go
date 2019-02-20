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

package ilk

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type PitFileIlkConverter struct{}

func (PitFileIlkConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var models []interface{}
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}
		ilk := shared.GetHexWithoutPrefix(ethLog.Topics[2].Bytes())
		what := string(bytes.Trim(ethLog.Topics[3].Bytes(), "\x00"))
		dataBytes := ethLog.Data[len(ethLog.Data)-constants.DataItemLength:]
		data, err := getData(dataBytes, what)
		if err != nil {
			return nil, err
		}

		raw, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}
		model := PitFileIlkModel{
			Ilk:              ilk,
			What:             what,
			Data:             data,
			LogIndex:         ethLog.Index,
			TransactionIndex: ethLog.TxIndex,
			Raw:              raw,
		}
		models = append(models, model)
	}
	return models, nil
}

func getData(dataBytes []byte, what string) (string, error) {
	n := big.NewInt(0).SetBytes(dataBytes).String()
	if what == "spot" {
		return shared.ConvertToRay(n), nil
	} else if what == "line" {
		return shared.ConvertToWad(n), nil
	} else {
		return "", errors.New("unexpected payload for 'what'")
	}
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 4 {
		return errors.New("log missing topics")
	}
	if len(log.Data) < constants.DataItemLength {
		return errors.New("log missing data")
	}
	return nil
}
