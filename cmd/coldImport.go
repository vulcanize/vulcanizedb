// Copyright © 2018 Rob Mulholand <rmulholand@8thlight.com>
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

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/ethereum"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/converters/cold_db"
	vulcCommon "github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
	"github.com/vulcanize/vulcanizedb/utils"
)

var coldImportCmd = &cobra.Command{
	Use:   "coldImport",
	Short: "Sync vulcanize from a cold instance of LevelDB.",
	Long: `Populate core vulcanize db data directly out of LevelDB, rather than over rpc calls. For example:

./vulcanizedb coldImport -s 0 -e 5000000

Geth must be synced over all of the desired blocks and must not be running in order to execute this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		coldImport()
	},
}

func init() {
	rootCmd.AddCommand(coldImportCmd)
	coldImportCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "Number for first block to cold import")
	coldImportCmd.Flags().Int64VarP(&endingBlockNumber, "ending-block-number", "e", 5500000, "Number for last block to cold import")
}

func coldImport() {
	if endingBlockNumber < startingBlockNumber {
		log.Fatal("Ending block number must be greater than starting block number for cold import.")
	}

	// init eth db
	ethDBConfig := ethereum.CreateDatabaseConfig(ethereum.Level, levelDbPath)
	ethDB, err := ethereum.CreateDatabase(ethDBConfig)
	if err != nil {
		log.Fatal("Error connecting to ethereum db: ", err)
	}

	// init pg db
	genesisBlockHash := common.BytesToHash(ethDB.GetBlockHash(0)).String()
	coldNode := core.Node{
		GenesisBlock: genesisBlockHash,
		NetworkID:    1,
		ID:           "LevelDbColdImport",
		ClientName:   "LevelDbColdImport",
	}
	pgDB := utils.LoadPostgres(databaseConfig, coldNode)

	// init cold importer deps
	blockRepository := repositories.BlockRepository{DB: &pgDB}
	receiptRepository := repositories.ReceiptRepository{DB: &pgDB}
	transactionconverter := cold_db.NewColdDbTransactionConverter()
	blockConverter := vulcCommon.NewBlockConverter(transactionconverter)

	// init and execute cold importer
	coldImporter := geth.NewColdImporter(ethDB, blockRepository, receiptRepository, blockConverter)
	err = coldImporter.Execute(startingBlockNumber, endingBlockNumber)
	if err != nil {
		log.Fatal("Error executing cold import: ", err)
	}
}
