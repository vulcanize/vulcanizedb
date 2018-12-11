// Copyright © 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	shared2 "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

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

	repository := shared2.Repository{}
	fetcher := shared2.NewFetcher(blockChain)
	watcher := shared.NewWatcher(db, fetcher, repository)

	watcher.AddTransformers(transformers.TransformerInitializers())
	watcher.Execute()
}

func init() {
	rootCmd.AddCommand(backfillMakerLogsCmd)
}
