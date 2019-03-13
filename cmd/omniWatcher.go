// VulcanizeDB
// Copyright © 2019 Vulcanize

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

	ft "github.com/vulcanize/vulcanizedb/pkg/omni/full/transformer"
	lt "github.com/vulcanize/vulcanizedb/pkg/omni/light/transformer"
	st "github.com/vulcanize/vulcanizedb/pkg/omni/shared/transformer"
	"github.com/vulcanize/vulcanizedb/utils"
)

// omniWatcherCmd represents the omniWatcher command
var omniWatcherCmd = &cobra.Command{
	Use:   "omniWatcher",
	Short: "Watches events at the provided contract address using fully synced vDB",
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

var (
	network           string
	contractAddress   string
	contractAddresses []string
	contractEvents    []string
	contractMethods   []string
	eventArgs         []string
	methodArgs        []string
	methodPiping      bool
	mode              string
)

func omniWatcher() {
	if contractAddress == "" && len(contractAddresses) == 0 {
		log.Fatal("Contract address required")
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())

	var t st.Transformer
	switch mode {
	case "light":
		t = lt.NewTransformer(network, blockChain, &db)
	case "full":
		t = ft.NewTransformer(network, blockChain, &db)
	default:
		log.Fatal("Invalid mode")
	}

	contractAddresses = append(contractAddresses, contractAddress)
	for _, addr := range contractAddresses {
		t.SetEvents(addr, contractEvents)
		t.SetMethods(addr, contractMethods)
		t.SetEventArgs(addr, eventArgs)
		t.SetMethodArgs(addr, methodArgs)
		t.SetPiping(addr, methodPiping)
		t.SetStartingBlock(addr, startingBlockNumber)
	}

	err := t.Init()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to initialized transformer\r\nerr: %v\r\n", err))
	}

	for range ticker.C {
		t.Execute()
	}
}

func init() {
	rootCmd.AddCommand(omniWatcherCmd)

	omniWatcherCmd.Flags().StringVarP(&mode, "mode", "o", "light", "'light' or 'full' mode to work with either light synced or fully synced vDB (default is light)")
	omniWatcherCmd.Flags().StringVarP(&contractAddress, "contract-address", "a", "", "Single address to generate watchers for")
	omniWatcherCmd.Flags().StringArrayVarP(&contractAddresses, "contract-addresses", "l", []string{}, "list of addresses to use; warning: watcher targets the same events and methods for each address")
	omniWatcherCmd.Flags().StringArrayVarP(&contractEvents, "events", "e", []string{}, "Subset of events to watch; by default all events are watched")
	omniWatcherCmd.Flags().StringArrayVarP(&contractMethods, "methods", "m", nil, "Subset of methods to poll; by default no methods are polled")
	omniWatcherCmd.Flags().StringArrayVarP(&eventArgs, "event-args", "f", []string{}, "Argument values to filter event logs for; will only persist event logs that emit at least one of the value specified")
	omniWatcherCmd.Flags().StringArrayVarP(&methodArgs, "method-args", "g", []string{}, "Argument values to limit methods to; will only call methods with emitted values that were specified here")
	omniWatcherCmd.Flags().StringVarP(&network, "network", "n", "", `Network the contract is deployed on; options: "ropsten", "kovan", and "rinkeby"; default is mainnet"`)
	omniWatcherCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "Block to begin watching- default is first block the contract exists")
	omniWatcherCmd.Flags().BoolVarP(&methodPiping, "piping", "p", false, "Turn on method output piping: methods listed first will be polled first and their output used as input to subsequent methods")
}
