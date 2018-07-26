// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds/pep"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds/pip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds/rep"
	"github.com/vulcanize/vulcanizedb/utils"
)

// syncPriceFeedsCmd represents the syncPriceFeeds command
var syncPriceFeedsCmd = &cobra.Command{
	Use:   "syncPriceFeeds",
	Short: "Sync block headers with price feed data",
	Long: `Sync Ethereum block headers and price feed data. For example:

./vulcanizedb syncPriceFeeds --config <config.toml> --starting-block-number <block-number>

Price feed data will be updated when price feed contracts log value events.`,
	Run: func(cmd *cobra.Command, args []string) {
		syncPriceFeeds()
	},
}

func init() {
	rootCmd.AddCommand(syncPriceFeedsCmd)
	syncPriceFeedsCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "block number at which to start tracking price feeds")
}

func backFillPriceFeeds(blockchain core.BlockChain, headerRepository datastore.HeaderRepository, missingBlocksPopulated chan int, startingBlockNumber int64, transformers []transformers.Transformer) {
	populated, err := history.PopulateMissingHeaders(blockchain, headerRepository, startingBlockNumber, transformers)
	if err != nil {
		log.Fatal("Error populating headers: ", err)
	}
	missingBlocksPopulated <- populated
}

func syncPriceFeeds() {
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
	missingBlocksPopulated := make(chan int)
	transformers := []transformers.Transformer{
		pep.NewPepTransformer(blockChain, &db),
		pip.NewPipTransformer(blockChain, &db),
		rep.NewRepTransformer(blockChain, &db),
	}
	validator := history.NewHeaderValidator(blockChain, headerRepository, validationWindow, transformers)
	go backFillPriceFeeds(blockChain, headerRepository, missingBlocksPopulated, startingBlockNumber, transformers)

	for {
		select {
		case <-ticker.C:
			window := validator.ValidateHeaders()
			window.Log(os.Stdout)
		case <-missingBlocksPopulated:
			go backFillPriceFeeds(blockChain, headerRepository, missingBlocksPopulated, startingBlockNumber, transformers)
		}
	}
}
