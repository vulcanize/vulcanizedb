// VulcanizeDB
// Copyright © 2018 Vulcanize

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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
)

// backfillMakerLogsCmd represents the backfillMakerLogs command
var backfillMakerLogsCmd = &cobra.Command{
	Use:   "backfillMakerLogs",
	Short: "Backfill Maker event logs",
	Long: `Backfills Maker event logs based on previously populated block Header records.
This currently includes logs related to Multi-collateral Dai (frob), Auctions (flip-kick),
and Price Feeds (ETH/USD, MKR/USD, and REP/USD - LogValue).

vulcanizedb backfillMakerLogs --config environments/local.toml

This command expects a light sync to have been run, and the presence of header records in the Vulcanize database.`,
	Run: func(cmd *cobra.Command, args []string) {
		backfillMakerLogs()
	},
}

func backfillMakerLogs() {
	blockChain := getBlockChain()
	db, err := postgres.NewDB(databaseConfig, blockChain.Node())
	if err != nil {
		log.Fatal("Failed to initialize database.")
	}

	watcher := shared.NewWatcher(db, blockChain)

	watcher.AddTransformers(transformers.TransformerInitializers())
	err = watcher.Execute()
	if err != nil {
		// TODO Handle watcher error in backfillMakerLogs
	}
}

func init() {
	rootCmd.AddCommand(backfillMakerLogsCmd)
}
