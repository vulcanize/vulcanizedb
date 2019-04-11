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

// createIpldsForBlockReceiptsCmd represents the createIpldsForBlockReceipts command
var createIpldsForBlockReceiptsCmd = &cobra.Command{
	Use:   "createIpldsForBlockReceipts",
	Short: "Create IPLDs for a blocks receipts",
	Long: `Fetch a block's receipts and persist them to IPFS as IPLDs'. For example:

./eth-block-extractor createIpldsForBlockReceipts --block-number 5000000

The block number flag specifies the block for which to fetch and persist receipts.`,
	Run: func(cmd *cobra.Command, args []string) {
		createBlockReceipts()
	},
}

func init() {
	rootCmd.AddCommand(createIpldsForBlockReceiptsCmd)
	createIpldsForBlockReceiptsCmd.Flags().Int64VarP(&blockNumber, "block-number", "b", 0, "block for which to create receipt IPLDs")
}

func createBlockReceipts() {
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
	err = transformer.Execute(blockNumber, blockNumber)
	if err != nil {
		log.Fatal("Error creating receipt IPLDs for block: ", err)
	}
}
