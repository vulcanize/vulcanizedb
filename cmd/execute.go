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
	"fmt"
	"os"
	"plugin"
	syn "sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
	"github.com/vulcanize/vulcanizedb/utils"
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
    ipcPath  = "http://kovan0.vulcanize.io:8545"

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
		execute()
	},
}

func execute() {
	// Build plugin generator config
	prepConfig()

	// Get the plugin path and load the plugin
	_, pluginPath, err := genConfig.GetPluginPaths()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Executing plugin %s", pluginPath)
	log.Info("linking plugin", pluginPath)
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		log.Debug("linking plugin failed")
		log.Fatal(err)
	}

	// Load the `Exporter` symbol from the plugin
	log.Info("loading transformers from plugin")
	symExporter, err := plug.Lookup("Exporter")
	if err != nil {
		log.Debug("loading Exporter symbol failed")
		log.Fatal(err)
	}

	// Assert that the symbol is of type Exporter
	exporter, ok := symExporter.(Exporter)
	if !ok {
		log.Debug("plugged-in symbol not of type Exporter")
		os.Exit(1)
	}

	// Use the Exporters export method to load the EventTransformerInitializer and StorageTransformerInitializer sets
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

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.Flags().BoolVar(&recheckHeadersArg, "recheckHeaders", false, "checks headers that are already checked for each transformer.")
}

type Exporter interface {
	Export() ([]transformer.EventTransformerInitializer, []transformer.StorageTransformerInitializer)
}

func watchEthEvents(w *watcher.EventWatcher, wg *syn.WaitGroup) {
	defer wg.Done()
	// Execute over the EventTransformerInitializer set using the watcher
	log.Info("executing event transformers")
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
			// TODO Handle watcher errors in execute
		}
	}
}

func watchEthStorage(w *watcher.StorageWatcher, wg *syn.WaitGroup) {
	defer wg.Done()
	// Execute over the StorageTransformerInitializer set using the watcher
	log.Info("executing storage transformers")
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	for range ticker.C {
		err := w.Execute()
		if err != nil {
			// TODO Handle watcher errors in execute
		}
	}
}
