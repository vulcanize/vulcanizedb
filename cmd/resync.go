// Copyright © 2020 Vulcanize, Inc
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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/resync"
)

// resyncCmd represents the resync command
var resyncCmd = &cobra.Command{
	Use:   "resync",
	Short: "Resync historical data",
	Long:  `Use this command to fill in sections of missing data in the super node`,
	Run: func(cmd *cobra.Command, args []string) {
		subCommand = cmd.CalledAs()
		logWithCommand = *log.WithField("SubCommand", subCommand)
		rsyncCmdCommand()
	},
}

func rsyncCmdCommand() {
	rConfig, err := resync.NewReSyncConfig()
	if err != nil {
		logWithCommand.Fatal(err)
	}
	if err := ipfs.InitIPFSPlugins(); err != nil {
		logWithCommand.Fatal(err)
	}
	rService, err := resync.NewResyncService(rConfig)
	if err != nil {
		logWithCommand.Fatal(err)
	}
	if err := rService.Resync(); err != nil {
		logWithCommand.Fatal(err)
	}
	logWithCommand.Infof("%s %s resync finished", rConfig.Chain.String(), rConfig.ResyncType.String())
}

func init() {
	rootCmd.AddCommand(resyncCmd)

	// flags
	resyncCmd.PersistentFlags().String("ipfs-path", "", "ipfs repository path")

	resyncCmd.PersistentFlags().String("resync-chain", "", "which chain to support, options are currently Ethereum or Bitcoin.")
	resyncCmd.PersistentFlags().String("resync-type", "", "which type of data to resync")
	resyncCmd.PersistentFlags().Int("resync-start", 0, "block height to start resync")
	resyncCmd.PersistentFlags().Int("resync-stop", 0, "block height to stop resync")
	resyncCmd.PersistentFlags().Int("resync-batch-size", 0, "data fetching batch size")
	resyncCmd.PersistentFlags().Int("resync-batch-number", 0, "how many goroutines to fetch data concurrently")
	resyncCmd.PersistentFlags().Bool("resync-clear-old", false, "if true, clear out old data of the provided type within the resync range before resyncing")

	resyncCmd.PersistentFlags().String("btc-http-path", "", "http url for bitcoin node")
	resyncCmd.PersistentFlags().String("btc-password", "", "password for btc node")
	resyncCmd.PersistentFlags().String("btc-username", "", "username for btc node")
	resyncCmd.PersistentFlags().String("btc-node-id", "", "btc node id")
	resyncCmd.PersistentFlags().String("btc-client-name", "", "btc client name")
	resyncCmd.PersistentFlags().String("btc-genesis-block", "", "btc genesis block hash")
	resyncCmd.PersistentFlags().String("btc-network-id", "", "btc network id")

	resyncCmd.PersistentFlags().String("eth-http-path", "", "http url for ethereum node")

	// and their bindings
	viper.BindPFlag("ipfs.path", resyncCmd.PersistentFlags().Lookup("ipfs-path"))

	viper.BindPFlag("resync.chain", resyncCmd.PersistentFlags().Lookup("resync-chain"))
	viper.BindPFlag("resync.type", resyncCmd.PersistentFlags().Lookup("resync-type"))
	viper.BindPFlag("resync.start", resyncCmd.PersistentFlags().Lookup("resync-start"))
	viper.BindPFlag("resync.stop", resyncCmd.PersistentFlags().Lookup("resync-stop"))
	viper.BindPFlag("resync.batchSize", resyncCmd.PersistentFlags().Lookup("resync-batch-size"))
	viper.BindPFlag("resync.batchNumber", resyncCmd.PersistentFlags().Lookup("resync-batch-number"))
	viper.BindPFlag("resync.clearOldCache", resyncCmd.PersistentFlags().Lookup("resync-clear-old"))

	viper.BindPFlag("bitcoin.httpPath", resyncCmd.PersistentFlags().Lookup("btc-http-path"))
	viper.BindPFlag("bitcoin.pass", resyncCmd.PersistentFlags().Lookup("btc-password"))
	viper.BindPFlag("bitcoin.user", resyncCmd.PersistentFlags().Lookup("btc-username"))
	viper.BindPFlag("bitcoin.nodeID", resyncCmd.PersistentFlags().Lookup("btc-node-id"))
	viper.BindPFlag("bitcoin.clientName", resyncCmd.PersistentFlags().Lookup("btc-client-name"))
	viper.BindPFlag("bitcoin.genesisBlock", resyncCmd.PersistentFlags().Lookup("btc-genesis-block"))
	viper.BindPFlag("bitcoin.networkID", resyncCmd.PersistentFlags().Lookup("btc-network-id"))

	viper.BindPFlag("ethereum.httpPath", resyncCmd.PersistentFlags().Lookup("eth-http-path"))
}
