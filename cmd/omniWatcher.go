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
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/omni/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

// omniWatcherCmd represents the omniWatcher command
var omniWatcherCmd = &cobra.Command{
	Use:   "omniWatcher",
	Short: "Watches events at the provided contract address",
	Long: `Uses input contract address and event filters to watch events

Expects an ethereum node to be running
Expects an archival node synced into vulcanizeDB
Requires a .toml config file:

  [database]
  name = "vulcanize_public"
  hostname = "localhost"
  port = 5432

  [client]
  ipcPath = "/Users/user/Library/Ethereum/geth.ipc"
`,
	Run: func(cmd *cobra.Command, args []string) {
		omniWatcher()
	},
}

func omniWatcher() {

	if contractAddress == "" {
		log.Fatal("Contract address required")
	}

	if contractEvents == nil {
		var str string
		for str != "y" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Warning: no events specified, proceeding to watch every event at address" + contractAddress + "? (Y/n)\n> ")
			resp, err := reader.ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}

			str = strings.ToLower(string(resp))
			if str == "n" {
				return
			}
		}
	}

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
		log.Fatal(fmt.Sprintf("Failed to initialize database\r\nerr: %v\r\n", err))
	}

	con := types.Config{
		DB:      db,
		BC:      blockChain,
		Network: network,
	}

	t := transformer.NewTransformer(&con)
	t.Set(contractAddress, contractEvents)

	err = t.Init()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to initialized generator\r\nerr: %v\r\n", err))
	}

	log.Fatal(t.Execute())
}

func init() {
	rootCmd.AddCommand(omniWatcherCmd)

	omniWatcherCmd.Flags().StringVarP(&contractAddress, "contract-address", "a", "", "Single address to generate watchers for")
	omniWatcherCmd.Flags().StringArrayVarP(&contractEvents, "contract-events", "e", []string{}, "Subset of events to watch- use only with single address")
	omniWatcherCmd.Flags().StringArrayVarP(&contractAddresses, "contract-addresses", "l", []string{}, "Addresses of the contracts to generate watchers for")
	omniWatcherCmd.Flags().StringVarP(&network, "network", "n", "", `Network the contract is deployed on; options: "ropsten", "kovan", and "rinkeby"; default is mainnet"`)
}
