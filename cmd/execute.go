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
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
	"plugin"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	e "github.com/vulcanize/vulcanizedb/libraries/shared/executor"
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
		LogWithCommand = *log.WithField("SubCommand", SubCommand)
		execute()
	},
}

func execute() {
	// Build plugin generator config
	prepConfig()

	// Get the plugin path
	_, pluginPath, err := genConfig.GetPluginPaths()
	if err != nil {
		LogWithCommand.Fatal(err)
	}

	executePlugin(pluginPath)
}

func executePlugin(pluginPath string) {
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())

	LogWithCommand.Info("linking plugin", pluginPath)

	plug, linkErr := plugin.Open(pluginPath)
	if linkErr != nil {
		LogWithCommand.Info("linking plugin failed")
		LogWithCommand.Fatal(linkErr)
	}

	tailer := fs.FileTailer{Path: storageDiffsPath}
	storageFetcher := fetcher.NewCsvTailStorageFetcher(tailer)
	sw := watcher.NewStorageWatcher(storageFetcher, &db, pollingInterval, queueRecheckInterval)
	ew := watcher.NewEventWatcher(&db, blockChain, recheckHeadersArg, pollingInterval)
	cw := watcher.NewContractWatcher(&db, blockChain, pollingInterval)
	executor := e.NewExecutor(plug, &ew, sw, &cw)

	LogWithCommand.Info("loading transformers from plugin")
	loadErr := executor.LoadTransformerSets()
	if loadErr != nil {
		LogWithCommand.Fatal(loadErr)
	}

	LogWithCommand.Info("executing transformers")
	executor.ExecuteTransformerSets()
}

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.Flags().BoolVarP(&recheckHeadersArg, "recheck-headers", "r", false, "whether to re-check headers for watched events")
	executeCmd.Flags().DurationVarP(&queueRecheckInterval, "queue-recheck-interval", "q", 5*time.Minute, "interval duration for rechecking queued storage diffs (ex: 5m30s)")
}
