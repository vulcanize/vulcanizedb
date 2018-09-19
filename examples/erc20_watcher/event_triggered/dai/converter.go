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

package dai

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	"github.com/vulcanize/vulcanizedb/examples/generic/helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

// Converter converts a raw event log into its corresponding entity
// and can subsequently convert the entity into a model

type ERC20ConverterInterface interface {
	ToTransferEntity(watchedEvent core.WatchedEvent) (*TransferEntity, error)
	ToTransferModel(entity *TransferEntity) *event_triggered.TransferModel
	ToApprovalEntity(watchedEvent core.WatchedEvent) (*ApprovalEntity, error)
	ToApprovalModel(entity *ApprovalEntity) *event_triggered.ApprovalModel
}

type ERC20Converter struct {
	config shared.ContractConfig
}

func NewERC20Converter(config shared.ContractConfig) (*ERC20Converter, error) {
	var err error

	config.ParsedAbi, err = geth.ParseAbi(config.Abi)
	if err != nil {
		return nil, err
	}

	converter := &ERC20Converter{
		config: config,
	}

	return converter, nil
}

func (c ERC20Converter) ToTransferEntity(watchedEvent core.WatchedEvent) (*TransferEntity, error) {
	result := &TransferEntity{}
	contract := bind.NewBoundContract(common.HexToAddress(c.config.Address), c.config.ParsedAbi, nil, nil, nil)
	event := helpers.ConvertToLog(watchedEvent)
	err := contract.UnpackLog(result, constants.TransferEvent.String(), event)
	if err != nil {
		return result, err
	}
	result.TokenName = c.config.Name
	result.TokenAddress = common.HexToAddress(c.config.Address)
	result.Block = watchedEvent.BlockNumber
	result.TxHash = watchedEvent.TxHash

	return result, nil
}

func (c ERC20Converter) ToTransferModel(entity *TransferEntity) *event_triggered.TransferModel {
	to := entity.Dst.String()
	from := entity.Src.String()
	tokens := entity.Wad.String()

	return &event_triggered.TransferModel{
		TokenName:    c.config.Name,
		TokenAddress: c.config.Address,
		To:           to,
		From:         from,
		Tokens:       tokens,
		Block:        entity.Block,
		TxHash:       entity.TxHash,
	}
}

func (c ERC20Converter) ToApprovalEntity(watchedEvent core.WatchedEvent) (*ApprovalEntity, error) {
	result := &ApprovalEntity{}
	contract := bind.NewBoundContract(common.HexToAddress(c.config.Address), c.config.ParsedAbi, nil, nil, nil)
	event := helpers.ConvertToLog(watchedEvent)
	err := contract.UnpackLog(result, constants.ApprovalEvent.String(), event)
	if err != nil {
		return result, err
	}
	result.TokenName = c.config.Name
	result.TokenAddress = common.HexToAddress(c.config.Address)
	result.Block = watchedEvent.BlockNumber
	result.TxHash = watchedEvent.TxHash

	return result, nil
}

func (c ERC20Converter) ToApprovalModel(entity *ApprovalEntity) *event_triggered.ApprovalModel {
	tokenOwner := entity.Src.String()
	spender := entity.Guy.String()
	tokens := entity.Wad.String()

	return &event_triggered.ApprovalModel{
		TokenName:    c.config.Name,
		TokenAddress: c.config.Address,
		Owner:        tokenOwner,
		Spender:      spender,
		Tokens:       tokens,
		Block:        entity.Block,
		TxHash:       entity.TxHash,
	}
}
