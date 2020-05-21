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
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth/state"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
	"github.com/vulcanize/vulcanizedb/utils"
)

// stateSnapshotCmd represents the stateSnapshot command
var stateSnapshotCmd = &cobra.Command{
	Use:   "stateSnapshot",
	Short: "Used to create snapshot of eth state at provided height",
	Long: `Uses the statediff_getStateTrieAt endpoint
to retrieve all of the state and storage trie nodes at a provided block height.

It writes their IPLDs to public.blocks and their metadata to the eth.state_trie_cids and eth.storage_trie_cids tables
Everything is hash-linked up to the appropriate header in eth.header_cids`,
	Run: func(cmd *cobra.Command, args []string) {
		subCommand = cmd.CalledAs()
		logWithCommand = *log.WithField("SubCommand", subCommand)
		stateSnapshot()
	},
}

func stateSnapshot() {
	ethHTTP := viper.GetString("ethereum.httpPath")
	nodeInfo, httpClient, err := shared.GetEthNodeAndClient(fmt.Sprintf("http://%s", ethHTTP))
	if err != nil {
		logWithCommand.Fatal(err)
	}
	dbConfig := new(config.Database)
	dbConfig.Init()
	db := utils.LoadPostgres(*dbConfig, nodeInfo)
	snapshotBuilder := state.NewSnapsShotBuilder(&db, httpClient)
	height := viper.GetInt("trie.height")
	if err := snapshotBuilder.BuildSnapShotAt(uint64(height)); err != nil {
		logWithCommand.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(stateSnapshotCmd)

	// flags for all config variables
	stateSnapshotCmd.PersistentFlags().Int("trie-height", 0, "block height to fetch and write entire tries")

	stateSnapshotCmd.PersistentFlags().String("eth-http-path", "", "http url for ethereum node")
	stateSnapshotCmd.PersistentFlags().String("eth-node-id", "", "eth node id")
	stateSnapshotCmd.PersistentFlags().String("eth-client-name", "", "eth client name")
	stateSnapshotCmd.PersistentFlags().String("eth-genesis-block", "", "eth genesis block hash")
	stateSnapshotCmd.PersistentFlags().String("eth-network-id", "", "eth network id")

	// and their bindings
	viper.BindPFlag("trie.height", stateSnapshotCmd.PersistentFlags().Lookup("trie-height"))

	viper.BindPFlag("ethereum.httpPath", stateSnapshotCmd.PersistentFlags().Lookup("eth-http-path"))
	viper.BindPFlag("ethereum.nodeID", stateSnapshotCmd.PersistentFlags().Lookup("eth-node-id"))
	viper.BindPFlag("ethereum.clientName", stateSnapshotCmd.PersistentFlags().Lookup("eth-client-name"))
	viper.BindPFlag("ethereum.genesisBlock", stateSnapshotCmd.PersistentFlags().Lookup("eth-genesis-block"))
	viper.BindPFlag("ethereum.networkID", stateSnapshotCmd.PersistentFlags().Lookup("eth-network-id"))
}
