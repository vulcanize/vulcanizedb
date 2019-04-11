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
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_state_trie"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_storage_trie"
	"github.com/vulcanize/eth-block-extractor/pkg/transformers"
)

// createIpldsForStateTrieCmd represents the createIpldsForStateTrie command
var createIpldsForStateTrieCmd = &cobra.Command{
	Use:   "createIpldsForStateTrie",
	Short: "Create iplds for every ethereum state trie node",
	Long: `Create iplds for every ethereum state trie node. For example:

./eth-block-extractor createIpldsForStateTrie

Note that this operation is very expensive in terms of both cpu and disk,
as it is reconstructing the entire ethereum state trie in the same fashion
as an archive node.`,
	Run: func(cmd *cobra.Command, args []string) {
		createIpldsForStateTrie()
	},
}

func init() {
	rootCmd.AddCommand(createIpldsForStateTrieCmd)
	createIpldsForStateTrieCmd.Flags().BoolVarP(&computeState, "compute-state", "c", false, "Flag indicating state must be computed (non-archive node).")
	createIpldsForStateTrieCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "First block number to create IPLD for.")
	createIpldsForStateTrieCmd.Flags().Int64VarP(&endingBlockNumber, "ending-block-number", "e", 5900000, "Last block number to create IPLD for.")
}

func createIpldsForStateTrie() {
	if computeState && startingBlockNumber != 0 {
		log.Println("Computing state trie must begin at genesis block. Ignoring passed starting block number.")
	}

	// init eth db
	databaseConfig := db.CreateDatabaseConfig(db.Level, levelDbPath)
	database, err := db.CreateDatabase(databaseConfig)
	if err != nil {
		log.Fatal("Error connecting to the ethereum db: ", err)
	}

	// init ipfs publishers
	adder, err := ipfs.InitIPFSNode(ipfsPath)
	if err != nil {
		log.Fatal("Error connecting to ipfs: ", err)
	}
	stateTrieDagPutter := eth_state_trie.NewStateTrieDagPutter(adder)
	stateTriePublisher := ipfs.NewIpfsPublisher(stateTrieDagPutter)
	storageTrieDagPutter := eth_storage_trie.NewStorageTrieDagPutter(adder)
	storageTriePublisher := ipfs.NewIpfsPublisher(storageTrieDagPutter)

	// init and execute transformer
	if computeState {
		transformer := transformers.NewComputeEthStateTrieTransformer(database, stateTriePublisher, storageTriePublisher)
		err = transformer.Execute(endingBlockNumber)
	} else {
		transformer := transformers.NewEthStateTrieTransformer(database, stateTriePublisher, storageTriePublisher)
		err = transformer.Execute(startingBlockNumber, endingBlockNumber)
	}
	if err != nil {
		log.Fatal("Error executing transformer: ", err)
	}
}
