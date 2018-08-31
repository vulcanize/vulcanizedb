// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event_triggered

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/examples/generic/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

// Converter converts a raw event log into its corresponding entity
// and can subsequently convert the entity into a model

type GenericConverterInterface interface {
	ToBurnEntity(watchedEvent core.WatchedEvent) (*BurnEntity, error)
	ToBurnModel(entity *BurnEntity) *BurnModel
	ToMintEntity(watchedEvent core.WatchedEvent) (*MintEntity, error)
	ToMintModel(entity *MintEntity) *MintModel
}

type GenericConverter struct {
	config generic.ContractConfig
}

func NewGenericConverter(config generic.ContractConfig) (*GenericConverter, error) {
	var err error

	config.ParsedAbi, err = geth.ParseAbi(config.Abi)
	if err != nil {
		return nil, err
	}

	converter := &GenericConverter{
		config: config,
	}

	return converter, nil
}

func (c GenericConverter) ToBurnEntity(watchedEvent core.WatchedEvent) (*BurnEntity, error) {
	result := &BurnEntity{}
	contract := bind.NewBoundContract(common.HexToAddress(c.config.Address), c.config.ParsedAbi, nil, nil, nil)
	event := helpers.ConvertToLog(watchedEvent)
	err := contract.UnpackLog(result, constants.BurnEvent.String(), event)
	if err != nil {
		return result, err
	}
	result.TokenName = c.config.Name
	result.TokenAddress = common.HexToAddress(c.config.Address)
	result.Block = watchedEvent.BlockNumber
	result.TxHash = watchedEvent.TxHash

	return result, nil
}

func (c GenericConverter) ToBurnModel(entity *BurnEntity) *BurnModel {
	burner := entity.Burner.String()
	tokens := entity.Value.String()

	return &BurnModel{
		TokenName:    c.config.Name,
		TokenAddress: c.config.Address,
		Burner:       burner,
		Tokens:       tokens,
		Block:        entity.Block,
		TxHash:       entity.TxHash,
	}
}

func (c GenericConverter) ToMintEntity(watchedEvent core.WatchedEvent) (*MintEntity, error) {
	result := &MintEntity{}
	contract := bind.NewBoundContract(common.HexToAddress(c.config.Address), c.config.ParsedAbi, nil, nil, nil)
	event := helpers.ConvertToLog(watchedEvent)
	err := contract.UnpackLog(result, constants.MintEvent.String(), event)
	if err != nil {
		return result, err
	}
	result.TokenName = c.config.Name
	result.TokenAddress = common.HexToAddress(c.config.Address)
	result.Block = watchedEvent.BlockNumber
	result.TxHash = watchedEvent.TxHash

	return result, nil
}

func (c GenericConverter) ToMintModel(entity *MintEntity) *MintModel {
	mintee := entity.To.String()
	minter := c.config.Owner
	tokens := entity.Amount.String()

	return &MintModel{
		TokenName:    c.config.Name,
		TokenAddress: c.config.Address,
		Mintee:       mintee,
		Minter:       minter,
		Tokens:       tokens,
		Block:        entity.Block,
		TxHash:       entity.TxHash,
	}
}
