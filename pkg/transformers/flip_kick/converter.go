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

package flip_kick

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/utilities"
)

type Converter interface {
	ToEntity(contractAddress string, contractAbi string, ethLog types.Log) (*FlipKickEntity, error)
	ToModel(flipKick FlipKickEntity) (FlipKickModel, error)
}

type FlipKickConverter struct{}

func (FlipKickConverter) ToEntity(contractAddress string, contractAbi string, ethLog types.Log) (*FlipKickEntity, error) {
	entity := &FlipKickEntity{}
	address := common.HexToAddress(contractAddress)
	abi, err := geth.ParseAbi(contractAbi)
	if err != nil {
		return entity, err
	}

	contract := bind.NewBoundContract(address, abi, nil, nil, nil)

	err = contract.UnpackLog(entity, "FlipKick", ethLog)
	if err != nil {
		return entity, err
	}
	entity.Raw = ethLog
	return entity, nil
}

func (FlipKickConverter) ToModel(flipKick FlipKickEntity) (FlipKickModel, error) {
	//TODO: Confirm if the following values can be/ever will be nil

	if flipKick.Id == nil {
		return FlipKickModel{}, errors.New("FlipKick log ID cannot be nil.")
	}

	id := flipKick.Id.String()
	mom := strings.ToLower(flipKick.Mom.String())
	vat := strings.ToLower(flipKick.Vat.String())
	ilk := strings.ToLower(common.ToHex(flipKick.Ilk[:]))
	lot := utilities.ConvertNilToEmptyString(flipKick.Lot.String())
	bid := utilities.ConvertNilToEmptyString(flipKick.Bid.String())
	guy := strings.ToLower(flipKick.Guy.String())
	gal := strings.ToLower(flipKick.Gal.String())
	endValue := utilities.ConvertNilToZeroTimeValue(flipKick.End)
	end := time.Unix(endValue, 0)
	eraValue := utilities.ConvertNilToZeroTimeValue(flipKick.Era)
	era := time.Unix(eraValue, 0)
	lad := strings.ToLower(flipKick.Lad.String())
	tab := utilities.ConvertNilToEmptyString(flipKick.Tab.String())
	rawLogJson, err := json.Marshal(flipKick.Raw)
	if err != nil {
		return FlipKickModel{}, err
	}
	rawLogString := string(rawLogJson)

	return FlipKickModel{
		Id:  id,
		Mom: mom,
		Vat: vat,
		Ilk: ilk,
		Lot: lot,
		Bid: bid,
		Guy: guy,
		Gal: gal,
		End: end,
		Era: era,
		Lad: lad,
		Tab: tab,
		Raw: rawLogString,
	}, nil
}
