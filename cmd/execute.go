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
	"plugin"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/statediff"
	"github.com/makerdao/vulcanizedb/libraries/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/streamer"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/libraries/shared/watcher"
	"github.com/makerdao/vulcanizedb/pkg/fs"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "executes a precomposed transformer initializer plugin",
	Long: `This command needs a config .toml file of form:

[database]
    name     = "vulcanize_public"
    hostname = "localhost"
    user     = "vulcanize"
    password = "vulcanize"
    port     = 5432

[client]
    ipcPath  = "/Users/user/Library/Ethereum/geth.ipc"

[exporter]
    name     = "exampleTransformerExporter"

Note: If any of the plugin transformer need additional
configuration variables include them in the .toml file as well

The exporter.name is the name (without extension) of the plugin to be loaded.
The plugin file needs to be located in the /plugins directory and this command assumes 
the db migrations remain from when the plugin was composed. Additionally, the plugin 
must have been composed by the same version of vulcanizedb or else it will not be compatible.

Specify config location when executing the command:
./vulcanizedb execute --config=./environments/config_name.toml`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommand = cmd.CalledAs()
		LogWithCommand = *logrus.WithField("SubCommand", SubCommand)
		execute()
	},
}

func execute() {
	// Build plugin generator config
	configErr := prepConfig()
	if configErr != nil {
		LogWithCommand.Fatalf("failed to prepare config: %s", configErr.Error())
	}
	executeTransformers()
}

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.Flags().BoolVarP(&recheckHeadersArg, "recheck-headers", "r", false, "whether to re-check headers for watched events")
	executeCmd.Flags().DurationVarP(&queueRecheckInterval, "queue-recheck-interval", "q", 5*time.Minute, "interval duration for rechecking queued storage diffs (ex: 5m30s)")
	executeCmd.Flags().DurationVarP(&retryInterval, "retry-interval", "i", 7*time.Second, "interval duration between retries on execution error")
	executeCmd.Flags().IntVarP(&maxUnexpectedErrors, "max-unexpected-errs", "m", 5, "maximum number of unexpected errors to allow (with retries) before exiting")
}

func executeTransformers() {
	// Get the plugin path and load the plugin
	_, pluginPath, pathErr := genConfig.GetPluginPaths()
	if pathErr != nil {
		LogWithCommand.Fatalf("failed to get plugin paths: %s", pathErr.Error())
	}

	LogWithCommand.Info("linking plugin ", pluginPath)
	plug, openErr := plugin.Open(pluginPath)
	if openErr != nil {
		LogWithCommand.Fatalf("linking plugin failed: %s", openErr.Error())
	}

	// Load the `Exporter` symbol from the plugin
	LogWithCommand.Info("loading transformers from plugin")
	symExporter, lookupErr := plug.Lookup("Exporter")
	if lookupErr != nil {
		LogWithCommand.Fatalf("loading Exporter symbol failed: %s", lookupErr.Error())
	}

	// Assert that the symbol is of type Exporter
	exporter, ok := symExporter.(Exporter)
	if !ok {
		LogWithCommand.Fatal("plugged-in symbol not of type Exporter")
	}

	// Use the Exporters export method to load the EventTransformerInitializer, StorageTransformerInitializer, and ContractTransformerInitializer sets
	ethEventInitializers, ethStorageInitializers, ethContractInitializers := exporter.Export()

	// Setup bc and db objects
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())

	// Execute over transformer sets returned by the exporter
	// Use WaitGroup to wait on both goroutines
	var wg sync.WaitGroup
	if len(ethEventInitializers) > 0 {
		ew := watcher.NewEventWatcher(&db, blockChain, maxUnexpectedErrors, retryInterval)
		addErr := ew.AddTransformers(ethEventInitializers)
		if addErr != nil {
			LogWithCommand.Fatalf("failed to add event transformer initializers to watcher: %s", addErr.Error())
		}
		wg.Add(1)
		go watchEthEvents(&ew, &wg)
	}

	if len(ethStorageInitializers) > 0 {
		switch storageDiffsSource {
		case "geth":
			logrus.Debug("fetching storage diffs from geth pub sub")
			rpcClient, _ := getClients()
			stateDiffStreamer := streamer.NewStateDiffStreamer(rpcClient)
			payloadChan := make(chan statediff.Payload)
			storageFetcher := fetcher.NewGethRpcStorageFetcher(&stateDiffStreamer, payloadChan)
			sw := watcher.NewStorageWatcher(storageFetcher, &db)
			sw.AddTransformers(ethStorageInitializers)
			wg.Add(1)
			go watchEthStorage(&sw, &wg)
		default:
			logrus.Debug("fetching storage diffs from csv")
			tailer := fs.FileTailer{Path: storageDiffsPath}
			storageFetcher := fetcher.NewCsvTailStorageFetcher(tailer)
			sw := watcher.NewStorageWatcher(storageFetcher, &db)
			sw.AddTransformers(ethStorageInitializers)
			wg.Add(1)
			go watchEthStorage(&sw, &wg)
		}
	}

	if len(ethContractInitializers) > 0 {
		gw := watcher.NewContractWatcher(&db, blockChain)
		gw.AddTransformers(ethContractInitializers)
		wg.Add(1)
		go watchEthContract(&gw, &wg)
	}
	wg.Wait()
}

type Exporter interface {
	Export() ([]transformer.EventTransformerInitializer, []transformer.StorageTransformerInitializer, []transformer.ContractTransformerInitializer)
}

func watchEthEvents(w *watcher.EventWatcher, wg *sync.WaitGroup) {
	defer wg.Done()
	// Execute over the EventTransformerInitializer set using the watcher
	LogWithCommand.Info("executing event transformers")
	var recheck constants.TransformerExecution
	if recheckHeadersArg {
		recheck = constants.HeaderRecheck
	} else {
		recheck = constants.HeaderUnchecked
	}
	err := w.Execute(recheck)
	if err != nil {
		LogWithCommand.Fatalf("error executing event watcher: %s", err.Error())
	}
}

func watchEthStorage(w watcher.IStorageWatcher, wg *sync.WaitGroup) {
	defer wg.Done()
	// Execute over the StorageTransformerInitializer set using the storage watcher
	LogWithCommand.Info("executing storage transformers")
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	err := w.Execute(queueRecheckInterval)
	if err != nil {
		LogWithCommand.Fatalf("error executing storage watcher: %s", err.Error())
	}
}

func watchEthContract(w *watcher.ContractWatcher, wg *sync.WaitGroup) {
	defer wg.Done()
	// Execute over the ContractTransformerInitializer set using the contract watcher
	LogWithCommand.Info("executing contract_watcher transformers")
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	for range ticker.C {
		w.Execute()
	}
}
