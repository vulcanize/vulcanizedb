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
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds/pep"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds/pip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds/rep"
	"github.com/vulcanize/vulcanizedb/utils"
)

// syncPriceFeedsCmd represents the syncPriceFeeds command
var syncPriceFeedsCmd = &cobra.Command{
	Use:   "syncPriceFeeds",
	Short: "Sync block headers with price feed data",
	Long: `Sync Ethereum block headers and price feed data. For example:

./vulcanizedb syncPriceFeeds --config <config.toml> --starting-block-number <block-number>

Price feed data will be updated when price feed contracts log value events.`,
	Run: func(cmd *cobra.Command, args []string) {
		syncPriceFeeds()
	},
}

func init() {
	rootCmd.AddCommand(syncPriceFeedsCmd)
	syncPriceFeedsCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "block number at which to start tracking price feeds")
}

func backFillPriceFeeds(blockchain core.BlockChain, headerRepository datastore.HeaderRepository, missingBlocksPopulated chan int, startingBlockNumber int64, transformers []transformers.Transformer) {
	populated, err := history.PopulateMissingHeaders(blockchain, headerRepository, startingBlockNumber, transformers)
	if err != nil {
		log.Fatal("Error populating headers: ", err)
	}
	missingBlocksPopulated <- populated
}

func syncPriceFeeds() {
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	blockChain := getBlockChain()
	validateArgs(blockChain)
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())

	transformers := []transformers.Transformer{
		pep.NewPepTransformer(blockChain, &db, pepContractAddress),
		pip.NewPipTransformer(blockChain, &db, pipContractAddress),
		rep.NewRepTransformer(blockChain, &db, repContractAddress),
	}
	headerRepository := repositories.NewHeaderRepository(&db)
	missingBlocksPopulated := make(chan int)
	validator := history.NewHeaderValidator(blockChain, headerRepository, validationWindow, transformers)
	go backFillPriceFeeds(blockChain, headerRepository, missingBlocksPopulated, startingBlockNumber, transformers)

	for {
		select {
		case <-ticker.C:
			window := validator.ValidateHeaders()
			window.Log(os.Stdout)
		case <-missingBlocksPopulated:
			go backFillPriceFeeds(blockChain, headerRepository, missingBlocksPopulated, startingBlockNumber, transformers)
		}
	}
}
