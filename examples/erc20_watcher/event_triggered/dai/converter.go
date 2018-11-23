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

package dai

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers"
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
