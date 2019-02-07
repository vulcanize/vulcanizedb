// Copyright © 2019 Vulcanize, Inc
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"plugin"
	syn "sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
	p2 "github.com/vulcanize/vulcanizedb/pkg/plugin"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"
	"github.com/vulcanize/vulcanizedb/utils"
)

// executePluginCmd represents the execute command
var composeAndExecuteCmd = &cobra.Command{
	Use:   "composeAndExecute",
	Short: "Composes, loads, and executes transformer initializer plugin",
	Long: `This command needs a config .toml file of form:

[database]
    name = "vulcanize_public"
    hostname = "localhost"
    user = "vulcanize"
    password = "vulcanize"
    port = 5432

[client]
    ipcPath = "http://kovan0.vulcanize.io:8545"

[exporter]
    name = "exporter"
    [exporter.transformers]
            transformer1 = "github.com/path/to/transformer1"
            transformer2 = "github.com/path/to/transformer2"
            transformer3 = "github.com/path/to/transformer3"
            transformer4 = "github.com/different/path/to/transformer1"
	[exporter.types]
            transformer1 = "eth_event"
            transformer2 = "eth_event"
            transformer3 = "eth_event"
            transformer4 = "eth_storage"
    [exporter.repositories]
            transformers = "github.com/path/to"
            transformer4 = "github.com/different/path"
    [exporter.migrations]
            transformers = "db/migrations"
            transformer4 = "to/db/migrations"

Note: If any of the imported transformer need additional
config variables do not forget to include those as well

This information is used to write and build a go plugin with a transformer 
set composed from the transformer imports specified in the config file
This plugin is loaded and the set of transformer initializers is exported
from it and loaded into and executed over by the appropriate watcher. 

The type of watcher that the transformer works with is specified using the 
exporter.types config variable as shown above. Currently there are watchers 
of event data from an eth node (eth_event) and storage data from an eth node 
(eth_storage). Soon there will be watchers for ipfs (ipfs_event and ipfs_storage).

Transformers of different types can be ran together in the same command using a 
single config file or in separate command instances using different config files

Specify config location when executing the command:
./vulcanizedb composeAndExecute --config=./environments/config_name.toml`,
	Run: func(cmd *cobra.Command, args []string) {
		composeAndExecute()
	},
}

func composeAndExecute() {
	// Build plugin generator config
	prepConfig()

	// Generate code to build the plugin according to the config file
	fmt.Println("generating plugin")
	generator, err := p2.NewGenerator(genConfig, databaseConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = generator.GenerateExporterPlugin()
	if err != nil {
		fmt.Fprint(os.Stderr, "generating plugin failed")
		log.Fatal(err)
	}

	// Get the plugin path and load the plugin
	_, pluginPath, err := genConfig.GetPluginPaths()
	if err != nil {
		log.Fatal(err)
	}
	if !genConfig.Save {
		defer helpers.ClearFiles(pluginPath)
	}
	fmt.Println("opening plugin")
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		fmt.Fprint(os.Stderr, "opening pluggin failed")
		log.Fatal(err)
	}

	// Load the `Exporter` symbol from the plugin
	fmt.Println("loading transformers from plugin")
	symExporter, err := plug.Lookup("Exporter")
	if err != nil {
		fmt.Fprint(os.Stderr, "loading Exporter symbol failed")
		log.Fatal(err)
	}

	// Assert that the symbol is of type Exporter
	exporter, ok := symExporter.(Exporter)
	if !ok {
		fmt.Fprint(os.Stderr, "plugged-in symbol not of type Exporter")
		os.Exit(1)
	}

	// Use the Exporters export method to load the TransformerInitializer and StorageTransformerInitializer sets
	ethEventInitializers, ethStorageInitializers := exporter.Export()

	// Setup bc and db objects
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())

	// Execute over transformer sets returned by the exporter
	// Use WaitGroup to wait on both goroutines
	var wg syn.WaitGroup
	if len(ethEventInitializers) > 0 {
		w := watcher.NewWatcher(&db, blockChain)
		w.AddTransformers(ethEventInitializers)
		wg.Add(1)
		go watchEthEvents(&w, &wg)
	}

	if len(ethStorageInitializers) > 0 {
		tailer := fs.FileTailer{Path: storageDiffsPath}
		w := watcher.NewStorageWatcher(tailer, &db)
		w.AddTransformers(ethStorageInitializers)
		wg.Add(1)
		go watchEthStorage(&w, &wg)
	}
	wg.Wait()
}

type Exporter interface {
	Export() ([]transformer.TransformerInitializer, []transformer.StorageTransformerInitializer)
}

func init() {
	rootCmd.AddCommand(composeAndExecuteCmd)
	composeAndExecuteCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "Block number to start transformer execution from")
}

func watchEthEvents(w *watcher.Watcher, wg *syn.WaitGroup) {
	defer wg.Done()
	// Execute over the TransformerInitializer set using the watcher
	fmt.Println("executing transformers")
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	for range ticker.C {
		err := w.Execute()
		if err != nil {
			// TODO Handle watcher errors in composeAndExecute
		}
	}
}

func watchEthStorage(w *watcher.StorageWatcher, wg *syn.WaitGroup) {
	defer wg.Done()
	// Execute over the TransformerInitializer set using the watcher
	fmt.Println("executing transformers")
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	for range ticker.C {
		err := w.Execute()
		if err != nil {
			// TODO Handle watcher errors in composeAndExecute
		}
	}
}

func prepConfig() {
	fmt.Println("configuring plugin")
	types := viper.GetStringMapString("exporter.types")
	genTypes := map[string]config.PluginType{}
	for transformerName, transformerType := range types {
		genType := config.GetPluginType(transformerType)
		if genType == config.UnknownTransformerType {
			log.Fatal(errors.New(`unknown transformer type in exporter config
accepted types are "eth_event", "eth_storage", "ipfs_event" and "ipfs_storage"`))
		}
		genTypes[transformerName] = genType
	}
	genConfig = config.Plugin{
		FilePath:     "$GOPATH/src/github.com/vulcanize/vulcanizedb/plugins",
		FileName:     viper.GetString("exporter.name"),
		Save:         viper.GetBool("exporter.save"),
		Initializers: viper.GetStringMapString("exporter.transformers"),
		Dependencies: viper.GetStringMapString("exporter.repositories"),
		Migrations:   viper.GetStringMapString("exporter.migrations"),
		Types:        genTypes,
	}
}
