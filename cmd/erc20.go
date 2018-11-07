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

package cmd

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered/dai"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/omni/constants"
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
	}

	// It is important that the event transformer is executed before the every_block transformer
	// because the events are used to generate the token holder address list that is used to
	// collect balances and allowances at every block
	transformers := append(dai.DaiEventTriggeredTransformerInitializer(), every_block.ERC20EveryBlockTransformerInitializer()...)

	err = watcher.AddTransformers(transformers, con)
	if err != nil {
		log.Fatal(err)
	}

	for range ticker.C {
		watcher.Execute()
	}
}

func init() {
	rootCmd.AddCommand(erc20Cmd)
}
