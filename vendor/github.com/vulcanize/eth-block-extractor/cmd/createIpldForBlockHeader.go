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
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_header"
	"github.com/vulcanize/eth-block-extractor/pkg/transformers"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/rlp"
)

// createIpldForBlockHeaderCmd represents the createIpldForBlockHeader command
var createIpldForBlockHeaderCmd = &cobra.Command{
	Use:   "createIpldForBlockHeader",
	Short: "Create an IPLD object for a block.",
	Long: `Create an IPLD object for a block.

e.g. ./eth-block-extractor createIpldForBlockHeader -b 1234567

Under the hood, the command fetches the block header RLP data from LevelDB and
puts it in IPFS, converting the data as an 'eth-block'`,
	Run: func(cmd *cobra.Command, args []string) {
		createIpldForBlockHeader()
	},
}

func init() {
	rootCmd.AddCommand(createIpldForBlockHeaderCmd)
	createIpldForBlockHeaderCmd.Flags().Int64VarP(&blockNumber, "block-number", "b", 0, "Create IPLD for this block header.")
}

func createIpldForBlockHeader() {
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
	decoder := rlp.RlpDecoder{}
	dagPutter := eth_block_header.NewBlockHeaderDagPutter(*ipfsNode, decoder)
	publisher := ipfs.NewIpfsPublisher(dagPutter)

	// execute transformer
	transformer := transformers.NewEthBlockHeaderTransformer(ethDB, publisher)
	err = transformer.Execute(blockNumber, blockNumber)
	if err != nil {
		log.Fatal("Error executing transformer: ", err.Error())
	}
}
