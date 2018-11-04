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
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/libraries/shared"
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
	if contractAddress == "" && len(contractAddresses) == 0 {
		log.Fatal("Contract address required")
	}

	if !methodsOn && !eventsOn {
		log.Fatal("Method polling and event watching turned off- nothing to do!")
	}

	if len(contractEvents) == 0 || len(contractMethods) == 0 {
		var str string
		for str != "y" {
			reader := bufio.NewReader(os.Stdin)
			if len(contractEvents) == 0 && len(contractMethods) == 0 {
				fmt.Print("Warning: no events or methods specified, proceed to watch every event and method? (Y/n)\n> ")
			} else if len(contractEvents) == 0 {
				fmt.Print("Warning: no events specified, proceed to watch every event? (Y/n)\n> ")
			} else {
				fmt.Print("Warning: no methods specified, proceed to watch every method? (Y/n)\n> ")
			}
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

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	rawRpcClient, err := rpc.Dial(ipc)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to initialize rpc client\r\nerr: %v\r\n", err))
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

	contractAddresses = append(contractAddresses, contractAddress)
	for _, addr := range contractAddresses {
		t.SetEvents(addr, contractEvents)
		t.SetMethods(addr, contractMethods)
	}

	err = t.Init()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to initialized transformer\r\nerr: %v\r\n", err))
	}

	w := shared.Watcher{}
	w.AddTransformer(t)

	for range ticker.C {
		w.Execute()
	}
}

func init() {
	rootCmd.AddCommand(omniWatcherCmd)

	omniWatcherCmd.Flags().StringVarP(&contractAddress, "contract-address", "a", "", "Single address to generate watchers for")
	omniWatcherCmd.Flags().StringArrayVarP(&contractAddresses, "contract-addresses", "l", []string{}, "List of addresses to generate watchers for")
	omniWatcherCmd.Flags().BoolVarP(&eventsOn, "events-on", "o", true, "Set to false to turn off watching of any event")
	omniWatcherCmd.Flags().BoolVarP(&methodsOn, "methods-on", "p", true, "Set to false to turn off polling of any method")
	omniWatcherCmd.Flags().StringVarP(&contractAddress, "methods-off", "a", "", "Single address to generate watchers for")
	omniWatcherCmd.Flags().StringArrayVarP(&contractEvents, "contract-events", "e", []string{}, "Subset of events to watch; by default all events are watched")
	omniWatcherCmd.Flags().StringArrayVarP(&contractEvents, "contract-methods", "m", []string{}, "Subset of methods to watch; by default all methods are watched")
	omniWatcherCmd.Flags().StringVarP(&network, "network", "n", "", `Network the contract is deployed on; options: "ropsten", "kovan", and "rinkeby"; default is mainnet"`)
}
