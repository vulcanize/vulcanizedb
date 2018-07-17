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

	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/pkg/crypto"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/ethereum"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
	"github.com/vulcanize/vulcanizedb/pkg/geth/cold_import"
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
	coldImportCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "BlockNumber for first block to cold import.")
	coldImportCmd.Flags().Int64VarP(&endingBlockNumber, "ending-block-number", "e", 5500000, "BlockNumber for last block to cold import.")
	coldImportCmd.Flags().BoolVarP(&syncAll, "all", "a", false, "Option to sync all missing blocks.")
}

func coldImport() {
	// init eth db
	ethDBConfig := ethereum.CreateDatabaseConfig(ethereum.Level, levelDbPath)
	ethDB, err := ethereum.CreateDatabase(ethDBConfig)
	if err != nil {
		log.Fatal("Error connecting to ethereum db: ", err)
	}
	mostRecentBlockNumberInDb := ethDB.GetHeadBlockNumber()
	if syncAll {
		startingBlockNumber = 0
		endingBlockNumber = mostRecentBlockNumberInDb
	}
	if endingBlockNumber < startingBlockNumber {
		log.Fatal("Ending block number must be greater than starting block number for cold import.")
	}
	if endingBlockNumber > mostRecentBlockNumberInDb {
		log.Fatal("Ending block number is greater than most recent block in db: ", mostRecentBlockNumberInDb)
	}

	// init pg db
	genesisBlock := ethDB.GetBlockHash(0)
	reader := fs.FsReader{}
	parser := crypto.EthPublicKeyParser{}
	nodeBuilder := cold_import.NewColdImportNodeBuilder(reader, parser)
	coldNode, err := nodeBuilder.GetNode(genesisBlock, levelDbPath)
	if err != nil {
		log.Fatal("Error getting node: ", err)
	}
	pgDB := utils.LoadPostgres(databaseConfig, coldNode)

	// init cold importer deps
	blockRepository := repositories.NewBlockRepository(&pgDB)
	receiptRepository := repositories.ReceiptRepository{DB: &pgDB}
	transactionconverter := cold_db.NewColdDbTransactionConverter()
	blockConverter := vulcCommon.NewBlockConverter(transactionconverter)

	// init and execute cold importer
	coldImporter := cold_import.NewColdImporter(ethDB, blockRepository, receiptRepository, blockConverter)
	err = coldImporter.Execute(startingBlockNumber, endingBlockNumber, coldNode.ID)
	if err != nil {
		log.Fatal("Error executing cold import: ", err)
	}
}
