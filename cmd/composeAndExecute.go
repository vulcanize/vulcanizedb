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
	"os"
	"plugin"
	syn "sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
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
    name     = "vulcanize_public"
    hostname = "localhost"
    user     = "vulcanize"
    password = "vulcanize"
    port     = 5432

[client]
    ipcPath  = "http://kovan0.vulcanize.io:8545"

[exporter]
    home     = "github.com/vulcanize/vulcanizedb"
    clone    = false
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
    [exporter.transformer2]
        path = "path/to/transformer2"
        type = "eth_event"
        repository = "github.com/account/repo"
        migrations = "db/migrations"
    [exporter.transformer3]
        path = "path/to/transformer3"
        type = "eth_event"
        repository = "github.com/account/repo"
        migrations = "db/migrations"
    [exporter.transformer4]
        path = "path/to/transformer4"
        type = "eth_storage"
        repository = "github.com/account2/repo2"
        migrations = "to/db/migrations"


Note: If any of the imported transformer need additional
config variables do not forget to include those as well

This information is used to write and build a go plugin with a transformer 
set composed from the transformer imports specified in the config file
This plugin is loaded and the set of transformer initializers is exported
from it and loaded into and executed over by the appropriate watcher. 

The type of watcher that the transformer works with is specified using the 
type variable for each transformer in the config. Currently there are watchers 
of event data from an eth node (eth_event) and storage data from an eth node 
(eth_storage).

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
	log.Info("generating plugin")
	generator, err := p2.NewGenerator(genConfig, databaseConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = generator.GenerateExporterPlugin()
	if err != nil {
		log.Debug("generating plugin failed")
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
	log.Info("opening plugin")
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		log.Debug("opening pluggin failed")
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
			// TODO Handle watcher errors in composeAndExecute
		}
	}
}

func watchEthStorage(w *watcher.StorageWatcher, wg *syn.WaitGroup) {
	defer wg.Done()
	// Execute over the TransformerInitializer set using the watcher
	log.Info("executing storage transformers")
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
	log.Info("configuring plugin")
	names := viper.GetStringSlice("exporter.transformerNames")
	transformers := make(map[string]config.Transformer)
	for _, name := range names {
		transformer := viper.GetStringMapString("exporter." + name)
		p, ok := transformer["path"]
		if !ok || p == "" {
			log.Fatal(fmt.Sprintf("%s transformer config is missing `path` value", name))
		}
		r, ok := transformer["repository"]
		if !ok || r == "" {
			log.Fatal(fmt.Sprintf("%s transformer config is missing `repository` value", name))
		}
		m, ok := transformer["migrations"]
		if !ok || m == "" {
			log.Fatal(fmt.Sprintf("%s transformer config is missing `migrations` value", name))
		}
		t, ok := transformer["type"]
		if !ok {
			log.Fatal(fmt.Sprintf("%s transformer config is missing `type` value", name))
		}
		transformerType := config.GetTransformerType(t)
		if transformerType == config.UnknownTransformerType {
			log.Fatal(errors.New(`unknown transformer type in exporter config accepted types are "eth_event", "eth_storage"`))
		}

		transformers[name] = config.Transformer{
			Path:           p,
			Type:           transformerType,
			RepositoryPath: r,
			MigrationPath:  m,
		}
	}

	genConfig = config.Plugin{
		Transformers: transformers,
		FilePath:     "$GOPATH/src/github.com/vulcanize/vulcanizedb/plugins",
		FileName:     viper.GetString("exporter.name"),
		Save:         viper.GetBool("exporter.save"),
		Home:         viper.GetString("exporter.home"),
		Clone:        viper.GetBool("exporter.clone"),
	}
}
