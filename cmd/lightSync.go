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
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"github.com/vulcanize/vulcanizedb/utils"
)

// lightSyncCmd represents the lightSync command
var lightSyncCmd = &cobra.Command{
	Use:   "lightSync",
	Short: "Syncs VulcanizeDB with local ethereum node's block headers",
	Long: `Syncs VulcanizeDB with local ethereum node. Populates
Postgres with block headers.

./vulcanizedb lightSync --starting-block-number 0 --config public.toml

Expects ethereum node to be running and requires a .toml config:

  [database]
  name = "vulcanize_public"
  hostname = "localhost"
  port = 5432

  [client]
  ipcPath = "/Users/user/Library/Ethereum/geth.ipc"
`,
	Run: func(cmd *cobra.Command, args []string) {
		lightSync()
	},
}

func init() {
	rootCmd.AddCommand(lightSyncCmd)
	lightSyncCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "Block number to start syncing from")
}

func backFillAllHeaders(blockchain core.BlockChain, headerRepository datastore.HeaderRepository, missingBlocksPopulated chan int, startingBlockNumber int64) {
	populated, err := history.PopulateMissingHeaders(blockchain, headerRepository, startingBlockNumber)
	if err != nil {
		log.Fatal("Error populating headers: ", err)
	}
	missingBlocksPopulated <- populated
}

func lightSync() {
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	blockChain := getBlockChain()
	validateArgs(blockChain)
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())

	headerRepository := repositories.NewHeaderRepository(&db)
	validator := history.NewHeaderValidator(blockChain, headerRepository, validationWindow)
	missingBlocksPopulated := make(chan int)
	go backFillAllHeaders(blockChain, headerRepository, missingBlocksPopulated, startingBlockNumber)

	for {
		select {
		case <-ticker.C:
			window := validator.ValidateHeaders()
			window.Log(os.Stdout)
		case <-missingBlocksPopulated:
			go backFillAllHeaders(blockChain, headerRepository, missingBlocksPopulated, startingBlockNumber)
		}
	}
}

func validateArgs(blockChain *geth.BlockChain) {
	lastBlock := blockChain.LastBlock().Int64()
	if lastBlock == 0 {
		log.Fatal("geth initial: state sync not finished")
	}
	if startingBlockNumber > lastBlock {
		log.Fatal("starting block number > current block number")
	}
}
