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
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/omni/light/transformer"
	"github.com/vulcanize/vulcanizedb/utils"
)

// omniWatcherCmd represents the omniWatcher command
var lightOmniWatcherCmd = &cobra.Command{
	Use:   "lightOmniWatcher",
	Short: "Watches events at the provided contract address using lightSynced vDB",
	Long: `Uses input contract address and event filters to watch events

Expects an ethereum node to be running
Expects lightSync to have been run and the presence of headers in the Vulcanize database
Requires a .toml config file:

  [database]
  name = "vulcanize_public"
  hostname = "localhost"
  port = 5432

  [client]
  ipcPath = "/Users/user/Library/Ethereum/geth.ipc"
`,
	Run: func(cmd *cobra.Command, args []string) {
		lightOmniWatcher()
	},
}

func lightOmniWatcher() {
	if contractAddress == "" && len(contractAddresses) == 0 {
		log.Fatal("Contract address required")
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	t := transformer.NewTransformer(network, blockChain, &db)

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
	rootCmd.AddCommand(lightOmniWatcherCmd)

	lightOmniWatcherCmd.Flags().StringVarP(&contractAddress, "contract-address", "a", "", "Single address to generate watchers for")
	lightOmniWatcherCmd.Flags().StringArrayVarP(&contractAddresses, "contract-addresses", "l", []string{}, "list of addresses to use; warning: watcher targets the same events and methods for each address")
	lightOmniWatcherCmd.Flags().StringArrayVarP(&contractEvents, "contract-events", "e", []string{}, "Subset of events to watch; by default all events are watched")
	lightOmniWatcherCmd.Flags().StringArrayVarP(&contractMethods, "contract-methods", "m", nil, "Subset of methods to poll; by default no methods are polled")
	lightOmniWatcherCmd.Flags().StringArrayVarP(&eventAddrs, "event-filter-addresses", "f", []string{}, "Account addresses to persist event data for; default is to persist for all found token holder addresses")
	lightOmniWatcherCmd.Flags().StringArrayVarP(&methodAddrs, "method-filter-addresses", "g", []string{}, "Account addresses to poll methods with; default is to poll with all found token holder addresses")
	lightOmniWatcherCmd.Flags().StringVarP(&network, "network", "n", "", `Network the contract is deployed on; options: "ropsten", "kovan", and "rinkeby"; default is mainnet"`)
	lightOmniWatcherCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "Block to begin watching- default is first block the contract exists")
	lightOmniWatcherCmd.Flags().Int64VarP(&endingBlockNumber, "ending-block-number", "d", -1, "Block to end watching- default is most recent block")
}
