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

package integration_tests

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

var ipc string

func getClients(ipc string) (client.RpcClient, *ethclient.Client, error) {
	raw, err := rpc.Dial(ipc)
	if err != nil {
		return client.RpcClient{}, &ethclient.Client{}, err
	}
	return client.NewRpcClient(raw, ipc), ethclient.NewClient(raw), nil
}

func getBlockChain(rpcClient client.RpcClient, ethClient *ethclient.Client) (core.BlockChain, error) {
	client := client.NewEthClient(ethClient)
	node := node.MakeNode(rpcClient)
	transactionConverter := rpc2.NewRpcTransactionConverter(client)
	blockChain := geth.NewBlockChain(client, rpcClient, node, transactionConverter)
	return blockChain, nil
}

// Persist the header for a given block to postgres. Returns the header if successful.
func persistHeader(db *postgres.DB, blockNumber int64, blockChain core.BlockChain) (core.Header, error) {
	header, err := blockChain.GetHeaderByNumber(blockNumber)
	if err != nil {
		return core.Header{}, err
	}
	headerRepository := repositories.NewHeaderRepository(db)
	id, err := headerRepository.CreateOrUpdateHeader(header)
	header.Id = id
	return header, err
}
