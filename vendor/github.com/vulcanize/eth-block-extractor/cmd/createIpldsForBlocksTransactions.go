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

	"github.com/spf13/cobra"

	"github.com/vulcanize/eth-block-extractor/pkg/db"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_transactions"
	"github.com/vulcanize/eth-block-extractor/pkg/transformers"
)

// createIpldsForBlocksTransactionsCmd represents the createIpldsForBlocksTransactions command
var createIpldsForBlocksTransactionsCmd = &cobra.Command{
	Use:   "createIpldsForBlocksTransactions",
	Short: "Create IPLDs for every transaction in a range of blocks",
	Long: `Create an IPLD for every transaction in a range of blocks. For example:

./eth-block-extractor createIpldsForBlocksTransactions --config environments/public.toml --starting-block-number 5000000 --ending-block-number 5000100

The starting and ending block numbers specify the range of blocks for which to create transaction IPLDs, and are required.`,
	Run: func(cmd *cobra.Command, args []string) {
		createIpldsForBlocksTransactions()
	},
}

func init() {
	rootCmd.AddCommand(createIpldsForBlocksTransactionsCmd)
	createIpldsForBlocksTransactionsCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "First block number to create IPLD for.")
	createIpldsForBlocksTransactionsCmd.Flags().Int64VarP(&endingBlockNumber, "ending-block-number", "e", 5900000, "Last block number to create IPLD for.")
}

func createIpldsForBlocksTransactions() {
	// init eth db
	ethDBConfig := db.CreateDatabaseConfig(db.Level, levelDbPath)
	ethDB, err := db.CreateDatabase(ethDBConfig)
	if err != nil {
		log.Fatal("Error connecting to ethereum db: ", err)
	}

	// init ipfs publisher
	ipfsNode, err := ipfs.InitIPFSNode(ipfsPath)
	if err != nil {
		log.Fatal("Error connecting to IPFS: ", err)
	}
	dagPutter := eth_block_transactions.NewBlockTransactionsDagPutter(*ipfsNode)
	publisher := ipfs.NewIpfsPublisher(dagPutter)

	// execute transformer
	transformer := transformers.NewEthBlockTransactionsTransformer(ethDB, publisher)
	err = transformer.Execute(startingBlockNumber, endingBlockNumber)
	if err != nil {
		log.Fatal("Error executing transformer: ", err.Error())
	}
}
