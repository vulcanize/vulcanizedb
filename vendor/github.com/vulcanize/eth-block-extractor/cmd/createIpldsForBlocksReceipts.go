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
	"github.com/spf13/cobra"
	"github.com/vulcanize/eth-block-extractor/pkg/db"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_receipts"
	"github.com/vulcanize/eth-block-extractor/pkg/transformers"
	"log"
)

// createIpldsForBlocksReceiptsCmd represents the createIpldsForBlocksReceipts command
var createIpldsForBlocksReceiptsCmd = &cobra.Command{
	Use:   "createIpldsForBlocksReceipts",
	Short: "Create IPLDs for receipts in a range of blocks",
	Long: `Fetch receipts in a range of blocks and persist them to IPFS as IPLDs. For example:

./eth-block-extractor createIpldsForBlocksReceipts --starting-block-number 5000000 --ending-block-number 5000123

The starting and ending block number flags specify the range of blocks for which to fetch and persist receipts.
Ending block number must be greater than or equal to starting block number.`,
	Run: func(cmd *cobra.Command, args []string) {
		createBlocksReceipts()
	},
}

func init() {
	rootCmd.AddCommand(createIpldsForBlocksReceiptsCmd)
	createIpldsForBlocksReceiptsCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "First block number to create IPLD for.")
	createIpldsForBlocksReceiptsCmd.Flags().Int64VarP(&endingBlockNumber, "ending-block-number", "e", 5900000, "Last block number to create IPLD for.")
}

func createBlocksReceipts() {
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
	dagPutter := eth_block_receipts.NewEthBlockReceiptDagPutter(ipfsNode)
	publisher := ipfs.NewIpfsPublisher(dagPutter)

	// execute transformer
	transformer := transformers.NewEthBlockReceiptTransformer(ethDB, publisher)
	err = transformer.Execute(startingBlockNumber, endingBlockNumber)
	if err != nil {
		log.Fatal("Error creating receipt IPLDs for block: ", err)
	}
}
