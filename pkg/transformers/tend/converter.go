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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Converter interface {
	Convert(contractAddress string, contractAbi string, ethLog types.Log) (TendModel, error)
}

type TendConverter struct{}

func (c TendConverter) Convert(contractAddress string, contractAbi string, ethLog types.Log) (TendModel, error) {
	entity := TendModel{}
	entity.Guy = common.HexToAddress(ethLog.Topics[1].Hex()).String()
	entity.BidId = ethLog.Topics[2].Big().String()
	entity.Lot = ethLog.Topics[3].Big().String()

	itemByteLength := 32
	lastDataItemStartIndex := len(ethLog.Data) - itemByteLength
	lastItem := ethLog.Data[lastDataItemStartIndex:]
	last := big.NewInt(0).SetBytes(lastItem)
	entity.Bid = last.String()

	entity.Tic = "0" //TODO: how do we get the bid tic?
	entity.TransactionIndex = ethLog.TxIndex
	rawJson, err := json.Marshal(ethLog)
	if err != nil {
		return TendModel{}, err
	}

	entity.Raw = string(rawJson)
	return entity, err
}
