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
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"github.com/vulcanize/vulcanizedb/utils"
)

// fullSyncCmd represents the fullSync command
var fullSyncCmd = &cobra.Command{
	Use:   "fullSync",
	Short: "Syncs VulcanizeDB with local ethereum node",
	Long: `Syncs VulcanizeDB with local ethereum node. Populates
Postgres with blocks, transactions, receipts, and logs.

./vulcanizedb fullSync --starting-block-number 0 --config public.toml

Expects ethereum node to be running and requires a .toml config:

  [database]
  name = "vulcanize_public"
  hostname = "localhost"
  port = 5432

  [client]
  ipcPath = "/Users/user/Library/Ethereum/geth.ipc"
`,
	Run: func(cmd *cobra.Command, args []string) {
		fullSync()
	},
}

func init() {
	rootCmd.AddCommand(fullSyncCmd)

	fullSyncCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "Block number to start syncing from")
}

func backFillAllBlocks(blockchain core.BlockChain, blockRepository datastore.BlockRepository, missingBlocksPopulated chan int, startingBlockNumber int64) {
	populated, err := history.PopulateMissingBlocks(blockchain, blockRepository, startingBlockNumber)
	if err != nil {
		log.Error("backfillAllBlocks: error in populateMissingBlocks: ", err)
	}
	missingBlocksPopulated <- populated
}

func fullSync() {
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	blockChain := getBlockChain()
	lastBlock, err := blockChain.LastBlock()
	if err != nil {
		log.Error("fullSync: Error getting last block: ", err)
	}
	if lastBlock.Int64() == 0 {
		log.Fatal("geth initial: state sync not finished")
	}
	if startingBlockNumber > lastBlock.Int64() {
		log.Fatal("fullSync: starting block number > current block number")
	}

	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	blockRepository := repositories.NewBlockRepository(&db)
	validator := history.NewBlockValidator(blockChain, blockRepository, validationWindow)
	missingBlocksPopulated := make(chan int)
	go backFillAllBlocks(blockChain, blockRepository, missingBlocksPopulated, startingBlockNumber)

	for {
		select {
		case <-ticker.C:
			window, err := validator.ValidateBlocks()
			if err != nil {
				log.Error("fullSync: error in validateBlocks: ", err)
			}
			log.Info(window.GetString())
		case <-missingBlocksPopulated:
			go backFillAllBlocks(blockChain, blockRepository, missingBlocksPopulated, startingBlockNumber)
		}
	}
}
