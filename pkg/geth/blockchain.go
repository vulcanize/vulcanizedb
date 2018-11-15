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

package geth

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/net/context"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	vulcCommon "github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
)

type BlockChain struct {
	client          core.EthClient
	blockConverter  vulcCommon.BlockConverter
	headerConverter vulcCommon.HeaderConverter
	node            core.Node
}

func NewBlockChain(client core.EthClient, node core.Node, converter vulcCommon.TransactionConverter) *BlockChain {
	return &BlockChain{
		client:          client,
		blockConverter:  vulcCommon.NewBlockConverter(converter),
		headerConverter: vulcCommon.HeaderConverter{},
		node:            node,
	}
}

func (blockChain *BlockChain) GetBlockByNumber(blockNumber int64) (block core.Block, err error) {
	gethBlock, err := blockChain.client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		return block, err
	}
	return blockChain.blockConverter.ToCoreBlock(gethBlock)
}

func (blockChain *BlockChain) GetHeaderByNumber(blockNumber int64) (header core.Header, err error) {
	gethHeader, err := blockChain.client.HeaderByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		return header, err
	}
	return blockChain.headerConverter.Convert(gethHeader)
}

func (blockChain *BlockChain) GetLogs(contract core.Contract, startingBlockNumber, endingBlockNumber *big.Int) ([]core.Log, error) {
	if endingBlockNumber == nil {
		endingBlockNumber = startingBlockNumber
	}
	contractAddress := common.HexToAddress(contract.Hash)
	fc := ethereum.FilterQuery{
		FromBlock: startingBlockNumber,
		ToBlock:   endingBlockNumber,
		Addresses: []common.Address{contractAddress},
	}
	gethLogs, err := blockChain.client.FilterLogs(context.Background(), fc)
	if err != nil {
		return []core.Log{}, err
	}
	logs := vulcCommon.ToCoreLogs(gethLogs)
	return logs, nil
}

func (blockChain *BlockChain) LastBlock() *big.Int {
	block, _ := blockChain.client.HeaderByNumber(context.Background(), nil)
	return block.Number
}

func (blockChain *BlockChain) Node() core.Node {
	return blockChain.node
}
