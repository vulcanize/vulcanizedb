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

package price_feeds

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type PriceFeedConverter struct{}

func (converter PriceFeedConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var results []interface{}
	for _, log := range ethLogs {
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		model := PriceFeedModel{
			BlockNumber:       log.BlockNumber,
			MedianizerAddress: log.Address.String(),
			UsdValue:          shared.ConvertToWad(hexutil.Encode(log.Data)),
			LogIndex:          log.Index,
			TransactionIndex:  log.TxIndex,
			Raw:               raw,
		}
		results = append(results, model)
	}
	return results, nil
}
