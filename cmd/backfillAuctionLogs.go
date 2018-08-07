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

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
)

// backfillAuctionLogsCmd represents the backfillAuctionLogs command
var backfillAuctionLogsCmd = &cobra.Command{
	Use:   "backfillAuctionLogs",
	Short: "Backfill auction event logs",
	Long: `Backfills auction event logs based on previously populated block Header records.
vulcanize backfillAuctionLogs --config environments/local.toml

This command expects a light sync to have been run, and the presence of header records in the Vulcanize database.`,
	Run: func(cmd *cobra.Command, args []string) {
		backfillAuctionLogs()
	},
}

func blockChain() *geth.BlockChain {
	rawRpcClient, err := rpc.Dial(ipc)

	if err != nil {
		log.Fatal(err)
	}
	rpcClient := client.NewRpcClient(rawRpcClient, ipc)
	ethClient := ethclient.NewClient(rawRpcClient)
	vdbEthClient := client.NewEthClient(ethClient)
	vdbNode := node.MakeNode(rpcClient)
	transactionConverter := vRpc.NewRpcTransactionConverter(ethClient)
	return geth.NewBlockChain(vdbEthClient, vdbNode, transactionConverter)
}

func backfillAuctionLogs() {
	blockChain := blockChain()
	db, err := postgres.NewDB(databaseConfig, blockChain.Node())
	if err != nil {
		log.Fatal("Failed to initialize database.")
	}

	watcher := shared.Watcher{
		DB:         *db,
		Blockchain: blockChain,
	}

	watcher.AddTransformers(transformers.TransformerInitializers())
	watcher.Execute()
}

func init() {
	rootCmd.AddCommand(backfillAuctionLogsCmd)
}
