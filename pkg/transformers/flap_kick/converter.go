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

package flap_kick

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"time"
)

type FlapKickConverter struct {
}

func (FlapKickConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]interface{}, error) {
	var entities []interface{}
	for _, ethLog := range ethLogs {
		entity := &FlapKickEntity{}
		address := ethLog.Address
		abi, err := geth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}
		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		err = contract.UnpackLog(entity, "Kick", ethLog)
		if err != nil {
			return nil, err
		}
		entity.Raw = ethLog
		entity.TransactionIndex = ethLog.TxIndex
		entity.LogIndex = ethLog.Index
		entities = append(entities, *entity)
	}
	return entities, nil
}

func (FlapKickConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var models []interface{}
	for _, entity := range entities {
		flapKickEntity, ok := entity.(FlapKickEntity)
		if !ok {
			return nil, fmt.Errorf("entity of type %T, not %T", entity, FlapKickEntity{})
		}

		if flapKickEntity.Id == nil {
			return nil, errors.New("FlapKick log ID cannot be nil.")
		}

		id := flapKickEntity.Id.String()
		lot := shared.BigIntToString(flapKickEntity.Lot)
		bid := shared.BigIntToString(flapKickEntity.Bid)
		gal := flapKickEntity.Gal.String()
		endValue := shared.BigIntToInt64(flapKickEntity.End)
		end := time.Unix(endValue, 0)
		rawLog, err := json.Marshal(flapKickEntity.Raw)
		if err != nil {
			return nil, err
		}

		model := FlapKickModel{
			BidId:            id,
			Lot:              lot,
			Bid:              bid,
			Gal:              gal,
			End:              end,
			TransactionIndex: flapKickEntity.TransactionIndex,
			LogIndex:         flapKickEntity.LogIndex,
			Raw:              rawLog,
		}
		models = append(models, model)
	}
	return models, nil
}
