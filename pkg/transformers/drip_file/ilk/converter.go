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

	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type Converter interface {
	ToModel(ethLog types.Log) (DripFileIlkModel, error)
}

type DripFileIlkConverter struct{}

func (DripFileIlkConverter) ToModel(ethLog types.Log) (DripFileIlkModel, error) {
	err := verifyLog(ethLog)
	if err != nil {
		return DripFileIlkModel{}, err
	}
	ilk := string(bytes.Trim(ethLog.Topics[2].Bytes(), "\x00"))
	vow := string(bytes.Trim(ethLog.Topics[3].Bytes(), "\x00"))
	taxBytes := ethLog.Data[len(ethLog.Data)-shared.DataItemLength:]
	tax := big.NewInt(0).SetBytes(taxBytes).String()
	raw, err := json.Marshal(ethLog)
	return DripFileIlkModel{
		Ilk:              ilk,
		Vow:              vow,
		Tax:              tax,
		TransactionIndex: ethLog.TxIndex,
		Raw:              raw,
	}, err
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 4 {
		return errors.New("log missing topics")
	}
	if len(log.Data) < shared.DataItemLength {
		return errors.New("log missing data")
	}
	return nil
}
