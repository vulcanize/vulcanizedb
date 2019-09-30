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
	"os"
	"plugin"
	syn "sync"
	"time"

	"github.com/ethereum/go-ethereum/statediff"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/streamer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
	p2 "github.com/vulcanize/vulcanizedb/pkg/plugin"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"
	"github.com/vulcanize/vulcanizedb/utils"
)

// composeAndExecuteCmd represents the composeAndExecute command
var composeAndExecuteCmd = &cobra.Command{
	Use:   "composeAndExecute",
	Short: "Composes, loads, and executes transformer initializer plugin",
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
    home     = "github.com/vulcanize/vulcanizedb"
    name     = "exampleTransformerExporter"
    save     = false
    transformerNames = [
        "transformer1",
        "transformer2",
        "transformer3",
        "transformer4",
    ]
    [exporter.transformer1]
        path = "path/to/transformer1"
        type = "eth_event"
        repository = "github.com/account/repo"
        migrations = "db/migrations"
        rank = "0"
    [exporter.transformer2]
        path = "path/to/transformer2"
        type = "eth_contract"
        repository = "github.com/account/repo"
        migrations = "db/migrations"
        rank = "2"
    [exporter.transformer3]
        path = "path/to/transformer3"
        type = "eth_event"
        repository = "github.com/account/repo"
        migrations = "db/migrations"
        rank = "0"
    [exporter.transformer4]
        path = "path/to/transformer4"
        type = "eth_storage"
        repository = "github.com/account2/repo2"
        migrations = "to/db/migrations"
        rank = "1"


Note: If any of the plugin transformer need additional
configuration variables include them in the .toml file as well

This information is used to write and build a go plugin with a transformer 
set composed from the transformer imports specified in the config file
This plugin is loaded and the set of transformer initializers is exported
from it and loaded into and executed over by the appropriate watcher. 

The type of watcher that the transformer works with is specified using the 
type variable for each transformer in the config. Currently there are watchers 
of event data from an eth node (eth_event) and storage data from an eth node 
(eth_storage), and a more generic interface for accepting contract_watcher pkg
based transformers which can perform both event watching and public method 
polling (eth_contract).

Transformers of different types can be ran together in the same command using a 
single config file or in separate command instances using different config files

Specify config location when executing the command:
./vulcanizedb composeAndExecute --config=./environments/config_name.toml`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommand = cmd.CalledAs()
		LogWithCommand = *log.WithField("SubCommand", SubCommand)
		composeAndExecute()
	},
}

func composeAndExecute() {
	// Build plugin generator config
	prepConfig()

	// Generate code to build the plugin according to the config file
	LogWithCommand.Info("generating plugin")
	generator, err := p2.NewGenerator(genConfig, databaseConfig)
	if err != nil {
		LogWithCommand.Fatal(err)
	}
	err = generator.GenerateExporterPlugin()
	if err != nil {
		LogWithCommand.Debug("generating plugin failed")
		LogWithCommand.Fatal(err)
	}

	// Get the plugin path and load the plugin
	_, pluginPath, err := genConfig.GetPluginPaths()
	if err != nil {
		LogWithCommand.Fatal(err)
	}
	if !genConfig.Save {
		defer helpers.ClearFiles(pluginPath)
	}
	LogWithCommand.Info("linking plugin ", pluginPath)
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		LogWithCommand.Debug("linking plugin failed")
		LogWithCommand.Fatal(err)
	}

	// Load the `Exporter` symbol from the plugin
	LogWithCommand.Info("loading transformers from plugin")
	symExporter, err := plug.Lookup("Exporter")
	if err != nil {
		LogWithCommand.Debug("loading Exporter symbol failed")
		LogWithCommand.Fatal(err)
	}

	// Assert that the symbol is of type Exporter
	exporter, ok := symExporter.(Exporter)
	if !ok {
		LogWithCommand.Debug("plugged-in symbol not of type Exporter")
		os.Exit(1)
	}

	// Use the Exporters export method to load the EventTransformerInitializer, StorageTransformerInitializer, and ContractTransformerInitializer sets
	ethEventInitializers, ethStorageInitializers, ethContractInitializers := exporter.Export()

	// Setup bc and db objects
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())

	// Execute over transformer sets returned by the exporter
	// Use WaitGroup to wait on both goroutines
	var wg syn.WaitGroup
	if len(ethEventInitializers) > 0 {
		ew := watcher.NewEventWatcher(&db, blockChain)
		err := ew.AddTransformers(ethEventInitializers)
		if err != nil {
			LogWithCommand.Fatalf("failed to add event transformer initializers to watcher: %s", err.Error())
		}
		wg.Add(1)
		go watchEthEvents(&ew, &wg)
	}

	if len(ethStorageInitializers) > 0 {
		switch storageDiffsSource {
		case "geth":
			log.Debug("fetching storage diffs from geth pub sub")
			rpcClient, _ := getClients()
			stateDiffStreamer := streamer.NewStateDiffStreamer(rpcClient)
			payloadChan := make(chan statediff.Payload)
			storageFetcher := fetcher.NewGethRpcStorageFetcher(&stateDiffStreamer, payloadChan)
			sw := watcher.NewStorageWatcher(storageFetcher, &db)
			sw.AddTransformers(ethStorageInitializers)
			wg.Add(1)
			go watchEthStorage(&sw, &wg)
		default:
			log.Debug("fetching storage diffs from csv")
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

func init() {
	rootCmd.AddCommand(composeAndExecuteCmd)
	composeAndExecuteCmd.Flags().BoolVarP(&recheckHeadersArg, "recheck-headers", "r", false, "whether to re-check headers for watched events")
	composeAndExecuteCmd.Flags().DurationVarP(&queueRecheckInterval, "queue-recheck-interval", "q", 5*time.Minute, "interval duration for rechecking queued storage diffs (ex: 5m30s)")
}
