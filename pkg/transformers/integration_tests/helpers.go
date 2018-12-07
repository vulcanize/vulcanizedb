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
	_, err = headerRepository.CreateOrUpdateHeader(header)
	return header, err
}
