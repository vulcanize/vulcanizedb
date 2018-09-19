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

package cmd

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered/dai"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

// erc20Cmd represents the erc20 command
var erc20Cmd = &cobra.Command{
	Use:   "erc20",
	Short: "Fetches and persists token supply",
	Long: `Fetches transfer and approval events, totalSupply, allowances, and 
balances for the configured token from each block and persists it in Vulcanize DB.
vulcanizedb erc20 --config environments/public

Expects an ethereum node to be running and requires a .toml config file:

  [database]
  name = "vulcanize_public"
  hostname = "localhost"
  port = 5432

  [client]
  ipcPath = "/Users/user/Library/Ethereum/geth.ipc"
`,
	Run: func(cmd *cobra.Command, args []string) {
		watchERC20s()
	},
}

func watchERC20s() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	rawRpcClient, err := rpc.Dial(ipc)
	if err != nil {
		log.Fatal(err)
	}
	rpcClient := client.NewRpcClient(rawRpcClient, ipc)
	ethClient := ethclient.NewClient(rawRpcClient)
	client := client.NewEthClient(ethClient)
	node := node.MakeNode(rpcClient)
	transactionConverter := vRpc.NewRpcTransactionConverter(ethClient)
	blockChain := geth.NewBlockChain(client, node, transactionConverter)
	db, err := postgres.NewDB(databaseConfig, blockChain.Node())
	if err != nil {
		log.Fatal("Failed to initialize database.")
	}

	con := generic.DaiConfig
	con.Filters = constants.DaiERC20Filters
	watcher := shared.Watcher{
		DB:         *db,
		Blockchain: blockChain,
		Config:     con,
	}

	watcher.AddTransformers(dai.DaiEventTriggeredTransformerInitializer())
	watcher.AddTransformers(every_block.ERC20EveryBlockTransformerInitializers())
	for range ticker.C {
		watcher.Execute()
	}
}

func init() {
	rootCmd.AddCommand(erc20Cmd)
}
