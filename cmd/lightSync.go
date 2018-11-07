// VulcanizeDB
// Copyright © 2018 Vulcanize

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
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
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
	missingBlocksPopulated <- history.PopulateMissingHeaders(blockchain, headerRepository, startingBlockNumber)
}

func lightSync() {
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	rawRpcClient, err := rpc.Dial(ipc)
	if err != nil {
		log.Fatal(err)
	}
	rpcClient := client.NewRpcClient(rawRpcClient, ipc)
	ethClient := ethclient.NewClient(rawRpcClient)
	client := client.NewEthClient(ethClient)
	node := node.MakeNode(rpcClient)
	transactionConverter := vRpc.NewRpcTransactionConverter(client)
	blockChain := geth.NewBlockChain(client, node, transactionConverter)

	lastBlock := blockChain.LastBlock().Int64()
	if lastBlock == 0 {
		log.Fatal("geth initial: state sync not finished")
	}
	if startingBlockNumber > lastBlock {
		log.Fatal("starting block number > current block number")
	}

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
