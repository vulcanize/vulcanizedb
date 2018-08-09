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

package frob

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

type Converter interface {
	ToEntity(contractAddress string, contractAbi string, ethLog types.Log) (FrobEntity, error)
	ToModel(flipKick FrobEntity) FrobModel
}

type FrobConverter struct {
}

func (FrobConverter) ToEntity(contractAddress string, contractAbi string, ethLog types.Log) (FrobEntity, error) {
	entity := FrobEntity{}
	address := common.HexToAddress(contractAddress)
	abi, err := geth.ParseAbi(contractAbi)
	if err != nil {
		return entity, err
	}
	contract := bind.NewBoundContract(address, abi, nil, nil, nil)
	err = contract.UnpackLog(&entity, "Frob", ethLog)
	return entity, err
}

func (FrobConverter) ToModel(frob FrobEntity) FrobModel {
	return FrobModel{
		Ilk: frob.Ilk[:],
		Lad: frob.Lad[:],
		Gem: frob.Gem.String(),
		Ink: frob.Ink.String(),
		Art: frob.Art.String(),
		Era: frob.Era.String(),
	}
}
