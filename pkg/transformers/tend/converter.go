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
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/utilities"
)

type Converter interface {
	ToEntity(contractAddress string, contractAbi string, ethLog types.Log) (TendEntity, error)
	ToModel(entity TendEntity) (TendModel, error)
}

type TendConverter struct{}

func (c TendConverter) ToEntity(contractAddress string, contractAbi string, ethLog types.Log) (TendEntity, error) {
	entity := TendEntity{}
	address := common.HexToAddress(contractAddress)
	abi, err := geth.ParseAbi(contractAbi)

	if err != nil {
		return entity, err
	}

	contract := bind.NewBoundContract(address, abi, nil, nil, nil)
	err = contract.UnpackLog(&entity, "Tend", ethLog)
	if err != nil {
		return entity, err
	}
	entity.TransactionIndex = ethLog.TxIndex
	entity.Raw = ethLog
	return entity, nil
}

func (c TendConverter) ToModel(entity TendEntity) (TendModel, error) {
	if entity.Id == nil {
		return TendModel{}, errors.New("Tend log ID cannot be nil.")
	}

	rawJson, err := json.Marshal(entity.Raw)
	if err != nil {
		return TendModel{}, err
	}
	era := utilities.ConvertNilToZeroTimeValue(entity.Era)
	return TendModel{
		Id:               utilities.ConvertNilToEmptyString(entity.Id.String()),
		Lot:              utilities.ConvertNilToEmptyString(entity.Lot.String()),
		Bid:              utilities.ConvertNilToEmptyString(entity.Bid.String()),
		Guy:              entity.Guy[:],
		Tic:              utilities.ConvertNilToEmptyString(entity.Tic.String()),
		Era:              time.Unix(era, 0),
		TransactionIndex: entity.TransactionIndex,
		Raw:              string(rawJson),
	}, nil
}
