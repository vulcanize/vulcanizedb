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

package ilk

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type Converter interface {
	ToModel(contractAddress string, contractAbi string, ethLog types.Log) (PitFileIlkModel, error)
}

type PitFileIlkConverter struct{}

func (PitFileIlkConverter) ToModel(contractAddress string, contractAbi string, ethLog types.Log) (entity PitFileIlkModel, err error) {
	ilk := string(bytes.Trim(ethLog.Topics[2].Bytes(), "\x00"))
	what := string(bytes.Trim(ethLog.Topics[3].Bytes(), "\x00"))
	itemByteLength := 32
	riskBytes := ethLog.Data[len(ethLog.Data)-itemByteLength:]
	risk := big.NewInt(0).SetBytes(riskBytes).String()

	raw, err := json.Marshal(ethLog)
	return PitFileIlkModel{
		Ilk:              ilk,
		What:             what,
		Data:             risk,
		TransactionIndex: ethLog.TxIndex,
		Raw:              raw,
	}, err
}
