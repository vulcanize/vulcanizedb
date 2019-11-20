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
	"github.com/makerdao/vulcanizedb/pkg/crypto"
	"github.com/makerdao/vulcanizedb/pkg/datastore/ethereum"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/eth/cold_import"
	"github.com/makerdao/vulcanizedb/pkg/eth/converters/cold_db"
	"github.com/makerdao/vulcanizedb/pkg/eth/converters/common"
	"github.com/makerdao/vulcanizedb/pkg/fs"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var coldImportCmd = &cobra.Command{
	Use:   "coldImport",
	Short: "Sync vulcanize from a cold instance of LevelDB.",
	Long: `Populate core vulcanize db data directly out of LevelDB, rather than over rpc calls. For example:

./vulcanizedb coldImport -s 0 -e 5000000

Geth must be synced over all of the desired blocks and must not be running in order to execute this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommand = cmd.CalledAs()
		LogWithCommand = *logrus.WithField("SubCommand", SubCommand)
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
		LogWithCommand.Fatal("Error connecting to ethereum db: ", err)
	}
	mostRecentBlockNumberInDb := ethDB.GetHeadBlockNumber()
	if syncAll {
		startingBlockNumber = 0
		endingBlockNumber = mostRecentBlockNumberInDb
	}
	if endingBlockNumber < startingBlockNumber {
		LogWithCommand.Fatal("Ending block number must be greater than starting block number for cold import.")
	}
	if endingBlockNumber > mostRecentBlockNumberInDb {
		LogWithCommand.Fatal("Ending block number is greater than most recent block in db: ", mostRecentBlockNumberInDb)
	}

	// init pg db
	genesisBlock := ethDB.GetBlockHash(0)
	reader := fs.FsReader{}
	parser := crypto.EthPublicKeyParser{}
	nodeBuilder := cold_import.NewColdImportNodeBuilder(reader, parser)
	coldNode, err := nodeBuilder.GetNode(genesisBlock, levelDbPath)
	if err != nil {
		LogWithCommand.Fatal("Error getting node: ", err)
	}
	pgDB := utils.LoadPostgres(databaseConfig, coldNode)

	// init cold importer deps
	blockRepository := repositories.NewBlockRepository(&pgDB)
	receiptRepository := repositories.FullSyncReceiptRepository{DB: &pgDB}
	transactionConverter := cold_db.NewColdDbTransactionConverter()
	blockConverter := common.NewBlockConverter(transactionConverter)

	// init and execute cold importer
	coldImporter := cold_import.NewColdImporter(ethDB, blockRepository, receiptRepository, blockConverter)
	err = coldImporter.Execute(startingBlockNumber, endingBlockNumber, coldNode.ID)
	if err != nil {
		LogWithCommand.Fatal("Error executing cold import: ", err)
	}
}
