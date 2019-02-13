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
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
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
            transformer1 = "path/to/transformer1"
            transformer2 = "path/to/transformer2"
            transformer3 = "path/to/transformer3"
            transformer4 = "path/to/transformer4"
	[exporter.types]
            transformer1 = "eth_event"
            transformer2 = "eth_event"
            transformer3 = "eth_event"
            transformer4 = "eth_storage"
    [exporter.repositories]
            transformers = "github.com/account/repo"
            transformer4 = "github.com/account2/repo2"
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
		ew := watcher.NewEventWatcher(&db, blockChain)
		ew.AddTransformers(ethEventInitializers)
		wg.Add(1)
		go watchEthEvents(&ew, &wg)
	}

	if len(ethStorageInitializers) > 0 {
		tailer := fs.FileTailer{Path: storageDiffsPath}
		sw := watcher.NewStorageWatcher(tailer, &db)
		sw.AddTransformers(ethStorageInitializers)
		wg.Add(1)
		go watchEthStorage(&sw, &wg)
	}
	wg.Wait()
}

type Exporter interface {
	Export() ([]transformer.TransformerInitializer, []transformer.StorageTransformerInitializer)
}

func init() {
	rootCmd.AddCommand(composeAndExecuteCmd)
	composeAndExecuteCmd.Flags().BoolVar(&recheckHeadersArg, "recheckHeaders", false, "checks headers that are already checked for each transformer.")
}

func watchEthEvents(w *watcher.EventWatcher, wg *syn.WaitGroup) {
	defer wg.Done()
	// Execute over the TransformerInitializer set using the watcher
	fmt.Println("executing event transformers")
	var recheck constants.TransformerExecution
	if recheckHeadersArg {
		recheck = constants.HeaderRecheck
	} else {
		recheck = constants.HeaderMissing
	}
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	for range ticker.C {
		err := w.Execute(recheck)
		if err != nil {
			// TODO Handle watcher errors in composeAndExecute
		}
	}
}

func watchEthStorage(w *watcher.StorageWatcher, wg *syn.WaitGroup) {
	defer wg.Done()
	// Execute over the TransformerInitializer set using the watcher
	fmt.Println("executing storage transformers")
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
	names := viper.GetStringSlice("exporter.transformerNames")
	transformers := make(map[string]config.Transformer)
	for _, name := range names {
		transformer := viper.GetStringMapString("exporter." + name)
		_, ok := transformer["path"]
		if !ok {
			log.Fatal(fmt.Sprintf("%s transformer config is missing `path` value", name))
		}
		_, ok = transformer["repository"]
		if !ok {
			log.Fatal(fmt.Sprintf("%s transformer config is missing `repository` value", name))
		}
		_, ok = transformer["migrations"]
		if !ok {
			log.Fatal(fmt.Sprintf("%s transformer config is missing `migrations` value", name))
		}
		ty, ok := transformer["type"]
		if !ok {
			log.Fatal(fmt.Sprintf("%s transformer config is missing `type` value", name))
		}

		transformerType := config.GetTransformerType(ty)
		if transformerType == config.UnknownTransformerType {
			log.Fatal(errors.New(`unknown transformer type in exporter config
accepted types are "eth_event", "eth_storage"`))
		}

		transformers[name] = config.Transformer{
			Path:           transformer["path"],
			Type:           transformerType,
			RepositoryPath: transformer["repository"],
			MigrationPath:  transformer["migrations"],
		}
	}

	genConfig = config.Plugin{
		Transformers: transformers,
		FilePath:     "$GOPATH/src/github.com/vulcanize/vulcanizedb/plugins",
		FileName:     viper.GetString("exporter.name"),
		Save:         viper.GetBool("exporter.save"),
	}
}
