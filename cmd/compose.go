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
	"errors"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	p2 "github.com/vulcanize/vulcanizedb/pkg/plugin"
)

// composeCmd represents the compose command
var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Composes transformer initializer plugin",
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
        rank = "0"
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
./vulcanizedb compose --config=./environments/config_name.toml`,
	Run: func(cmd *cobra.Command, args []string) {
		subCommand = cmd.CalledAs()
		logWithCommand = *log.WithField("SubCommand", subCommand)
		compose()
	},
}

func compose() {
	// Build plugin generator config
	prepConfig()

	// Generate code to build the plugin according to the config file
	logWithCommand.Info("generating plugin")
	generator, err := p2.NewGenerator(genConfig, databaseConfig)
	if err != nil {
		logWithCommand.Debug("initializing plugin generator failed")
		logWithCommand.Fatal(err)
	}
	err = generator.GenerateExporterPlugin()
	if err != nil {
		logWithCommand.Debug("generating plugin failed")
		logWithCommand.Fatal(err)
	}
	// TODO: Embed versioning info in the .so files so we know which version of vulcanizedb to run them with
	_, pluginPath, err := genConfig.GetPluginPaths()
	if err != nil {
		logWithCommand.Debug("getting plugin path failed")
		logWithCommand.Fatal(err)
	}
	fmt.Printf("Composed plugin %s", pluginPath)
	logWithCommand.Info("plugin .so file output to ", pluginPath)
}

func init() {
	rootCmd.AddCommand(composeCmd)
}

func prepConfig() {
	logWithCommand.Info("configuring plugin")
	names := viper.GetStringSlice("exporter.transformerNames")
	transformers := make(map[string]config.Transformer)
	for _, name := range names {
		transformer := viper.GetStringMapString("exporter." + name)
		p, pOK := transformer["path"]
		if !pOK || p == "" {
			logWithCommand.Fatal(name, " transformer config is missing `path` value")
		}
		r, rOK := transformer["repository"]
		if !rOK || r == "" {
			logWithCommand.Fatal(name, " transformer config is missing `repository` value")
		}
		m, mOK := transformer["migrations"]
		if !mOK || m == "" {
			logWithCommand.Fatal(name, " transformer config is missing `migrations` value")
		}
		mr, mrOK := transformer["rank"]
		if !mrOK || mr == "" {
			logWithCommand.Fatal(name, " transformer config is missing `rank` value")
		}
		rank, err := strconv.ParseUint(mr, 10, 64)
		if err != nil {
			logWithCommand.Fatal(name, " migration `rank` can't be converted to an unsigned integer")
		}
		t, tOK := transformer["type"]
		if !tOK {
			logWithCommand.Fatal(name, " transformer config is missing `type` value")
		}
		transformerType := config.GetTransformerType(t)
		if transformerType == config.UnknownTransformerType {
			logWithCommand.Fatal(errors.New(`unknown transformer type in exporter config accepted types are "eth_event", "eth_storage"`))
		}

		transformers[name] = config.Transformer{
			Path:           p,
			Type:           transformerType,
			RepositoryPath: r,
			MigrationPath:  m,
			MigrationRank:  rank,
		}
	}

	genConfig = config.Plugin{
		Transformers: transformers,
		FilePath:     "$GOPATH/src/github.com/vulcanize/vulcanizedb/plugins",
		FileName:     viper.GetString("exporter.name"),
		Save:         viper.GetBool("exporter.save"),
		Home:         viper.GetString("exporter.home"),
	}
}
