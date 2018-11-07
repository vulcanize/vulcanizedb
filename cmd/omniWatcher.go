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
	"bufio"
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/core"
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

	if len(contractEvents) == 0 || len(contractMethods) == 0 {
		var str string
		for str != "y" {
			reader := bufio.NewReader(os.Stdin)
			if len(contractEvents) == 0 && len(contractMethods) == 0 {
				fmt.Print("Warning: no events or methods specified.\n Proceed to watch every event and poll no methods? (Y/n)\n> ")
			} else if len(contractEvents) == 0 {
				fmt.Print("Warning: no events specified.\n Proceed to watch every event? (Y/n)\n> ")
			} else {
				fmt.Print("Warning: no methods specified.\n Proceed to poll no methods? (Y/n)\n> ")
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

	blockChain, db := setupBCandDB()

	t := transformer.NewTransformer(network, blockChain, db)

	contractAddresses = append(contractAddresses, contractAddress)
	for _, addr := range contractAddresses {
		t.SetEvents(addr, contractEvents)
		t.SetMethods(addr, contractMethods)
		t.SetEventAddrs(addr, eventAddrs)
		t.SetMethodAddrs(addr, methodAddrs)
		t.SetRange(addr, [2]int64{startingBlockNumber, endingBlockNumber})
	}

	err := t.Init()
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
	omniWatcherCmd.Flags().StringArrayVarP(&contractAddresses, "contract-addresses", "l", []string{}, "list of addresses to use; warning: watcher targets the same events and methods for each address")
	omniWatcherCmd.Flags().StringArrayVarP(&contractEvents, "contract-events", "e", []string{}, "Subset of events to watch; by default all events are watched")
	omniWatcherCmd.Flags().StringArrayVarP(&contractEvents, "contract-methods", "m", nil, "Subset of methods to poll; by default no methods are polled")
	omniWatcherCmd.Flags().StringArrayVarP(&eventAddrs, "event-filter-addresses", "f", []string{}, "Account addresses to persist event data for; default is to persist for all found token holder addresses")
	omniWatcherCmd.Flags().StringArrayVarP(&methodAddrs, "method-filter-addresses", "g", []string{}, "Account addresses to poll methods with; default is to poll with all found token holder addresses")
	omniWatcherCmd.Flags().StringVarP(&network, "network", "n", "", `Network the contract is deployed on; options: "ropsten", "kovan", and "rinkeby"; default is mainnet"`)
	omniWatcherCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "Block to begin watching- default is first block the contract exists")
	omniWatcherCmd.Flags().Int64VarP(&startingBlockNumber, "ending-block-number", "d", -1, "Block to end watching- default is most recent block")
}

func setupBCandDB() (core.BlockChain, *postgres.DB) {
	rawRpcClient, err := rpc.Dial(ipc)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to initialize rpc client\r\nerr: %v\r\n", err))
	}

	rpcClient := client.NewRpcClient(rawRpcClient, ipc)
	ethClient := ethclient.NewClient(rawRpcClient)
	cli := client.NewEthClient(ethClient)
	n := node.MakeNode(rpcClient)
	transactionConverter := vRpc.NewRpcTransactionConverter(ethClient)
	blockChain := geth.NewBlockChain(cli, n, transactionConverter)
	db, err := postgres.NewDB(databaseConfig, blockChain.Node())
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to initialize database\r\nerr: %v\r\n", err))
	}

	return blockChain, db
}
