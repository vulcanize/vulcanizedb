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

package client

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type EthClient struct {
	client *ethclient.Client
}

func NewEthClient(client *ethclient.Client) EthClient {
	return EthClient{client: client}
}

func (client EthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return client.client.BlockByNumber(ctx, number)
}

func (client EthClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return client.client.CallContract(ctx, msg, blockNumber)
}

func (client EthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return client.client.FilterLogs(ctx, q)
}

func (client EthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return client.client.HeaderByNumber(ctx, number)
}

func (client EthClient) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	return client.client.TransactionSender(ctx, tx, block, index)
}

func (client EthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return client.client.TransactionReceipt(ctx, txHash)
}
