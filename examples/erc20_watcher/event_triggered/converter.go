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
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/examples/generic/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

// Converter converts a raw event log into its corresponding entity
// and can subsequently convert the entity into a model

type ERC20ConverterInterface interface {
	ToTransferEntity(watchedEvent core.WatchedEvent) (*TransferEntity, error)
	ToTransferModel(entity TransferEntity) TransferModel
	ToApprovalEntity(watchedEvent core.WatchedEvent) (*ApprovalEntity, error)
	ToApprovalModel(entity ApprovalEntity) ApprovalModel
}

type ERC20Converter struct {
	config generic.ContractConfig
}

func NewERC20Converter(config generic.ContractConfig) ERC20Converter {
	return ERC20Converter{
		config: config,
	}
}

func (c ERC20Converter) ToTransferEntity(watchedEvent core.WatchedEvent) (*TransferEntity, error) {
	result := &TransferEntity{}
	contract := bind.NewBoundContract(common.HexToAddress(c.config.Address), c.config.ParsedAbi, nil, nil, nil)
	event := helpers.ConvertToLog(watchedEvent)
	err := contract.UnpackLog(result, "Transfer", event)
	if err != nil {
		return result, err
	}
	result.TokenName = c.config.Name
	result.TokenAddress = common.HexToAddress(c.config.Address)
	result.Block = watchedEvent.BlockNumber
	result.TxHash = watchedEvent.TxHash

	return result, nil
}

func (c ERC20Converter) ToTransferModel(transferEntity TransferEntity) TransferModel {
	to := transferEntity.Dst.String()
	from := transferEntity.Src.String()
	tokens := transferEntity.Wad.String()
	block := transferEntity.Block
	tx := transferEntity.TxHash
	return TransferModel{
		TokenName:    c.config.Name,
		TokenAddress: c.config.Address,
		To:           to,
		From:         from,
		Tokens:       tokens,
		Block:        block,
		TxHash:       tx,
	}
}

func (c ERC20Converter) ToApprovalEntity(watchedEvent core.WatchedEvent) (*ApprovalEntity, error) {
	result := &ApprovalEntity{}
	contract := bind.NewBoundContract(common.HexToAddress(c.config.Address), c.config.ParsedAbi, nil, nil, nil)
	event := helpers.ConvertToLog(watchedEvent)
	err := contract.UnpackLog(result, "Approval", event)
	if err != nil {
		return result, err
	}
	result.TokenName = c.config.Name
	result.TokenAddress = common.HexToAddress(c.config.Address)
	result.Block = watchedEvent.BlockNumber
	result.TxHash = watchedEvent.TxHash

	return result, nil
}

func (c ERC20Converter) ToApprovalModel(TransferEntity ApprovalEntity) ApprovalModel {
	tokenOwner := TransferEntity.Src.String()
	spender := TransferEntity.Guy.String()
	tokens := TransferEntity.Wad.String()
	block := TransferEntity.Block
	tx := TransferEntity.TxHash
	return ApprovalModel{
		TokenName:    c.config.Name,
		TokenAddress: c.config.Address,
		Owner:        tokenOwner,
		Spender:      spender,
		Tokens:       tokens,
		Block:        block,
		TxHash:       tx,
	}
}
