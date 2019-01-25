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
	"log"
	"os"
	"plugin"
	"time"

	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/autogen"
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
    filePath = "~/go/src/github.com/vulcanize/vulcanizedb/plugins"
	fileName = "exporter"
    [exporter.transformers]
		transformerImport1 = "github.com/path_to/transformerInitializer1"
		transformerImport2 = "github.com/path_to/transformerInitializer2"

Note: If any of the imported transformer need additional
config variables do not forget to include those as well

This information is used to write and build a .so with an arbitrary transformer 
set composed from the transformer imports specified in the config file
This .so is loaded as a plugin and the set of transformer initializers is 
loaded into and executed over by a generic watcher`,
	Run: func(cmd *cobra.Command, args []string) {
		composeAndExecute()
	},
}

func composeAndExecute() {
	generator := autogen.NewGenerator(autogenConfig)
	err := generator.GenerateTransformerPlugin()
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())

	_, pluginPath, err := autogen.GetPaths(autogenConfig)
	if err != nil {
		log.Fatal(err)
	}

	plug, err := plugin.Open(pluginPath)
	if err != nil {
		log.Fatal(err)
	}

	symExporter, err := plug.Lookup("Exporter")
	if err != nil {
		log.Fatal(err)
	}

	exporter, ok := symExporter.(Exporter)
	if !ok {
		fmt.Println("plugged-in symbol not of type Exporter")
		os.Exit(1)
	}

	initializers := exporter.Export()
	w := watcher.NewWatcher(&db, blockChain)
	w.AddTransformers(initializers)

	for range ticker.C {
		err := w.Execute()
		if err != nil {
			// TODO Handle watcher errors in composeAndExecute
		}
	}
}

type Exporter interface {
	Export() []transformer.TransformerInitializer
}

func init() {
	rootCmd.AddCommand(composeAndExecuteCmd)
	composeAndExecuteCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "Block number to start transformer execution from")
}
