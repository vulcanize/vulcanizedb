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

package tusd

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/examples/generic/event_triggered"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers"
)

// Converter converts a raw event log into its corresponding entity
// and can subsequently convert the entity into a model

type GenericConverterInterface interface {
	ToBurnEntity(watchedEvent core.WatchedEvent) (*BurnEntity, error)
	ToBurnModel(entity *BurnEntity) *event_triggered.BurnModel
	ToMintEntity(watchedEvent core.WatchedEvent) (*MintEntity, error)
	ToMintModel(entity *MintEntity) *event_triggered.MintModel
}

type GenericConverter struct {
	config shared.ContractConfig
}

func NewGenericConverter(config shared.ContractConfig) (*GenericConverter, error) {
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

func (c GenericConverter) ToBurnModel(entity *BurnEntity) *event_triggered.BurnModel {
	burner := entity.Burner.String()
	tokens := entity.Value.String()

	return &event_triggered.BurnModel{
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

func (c GenericConverter) ToMintModel(entity *MintEntity) *event_triggered.MintModel {
	mintee := entity.To.String()
	minter := c.config.Owner
	tokens := entity.Amount.String()

	return &event_triggered.MintModel{
		TokenName:    c.config.Name,
		TokenAddress: c.config.Address,
		Mintee:       mintee,
		Minter:       minter,
		Tokens:       tokens,
		Block:        entity.Block,
		TxHash:       entity.TxHash,
	}
}
