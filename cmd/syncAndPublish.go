// Copyright © 2019 Vulcanize, Inc
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	syn "sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/utils"
)

// syncAndPublishCmd represents the syncAndPublish command
var syncAndPublishCmd = &cobra.Command{
	Use:   "syncAndPublish",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		syncAndPulish()
	},
}

var (
	ipfsPath string
)

func init() {
	rootCmd.AddCommand(syncAndPublishCmd)
	syncAndPublishCmd.Flags().StringVarP(&ipfsPath, "ipfs-path", "i", "", "Path for configuring IPFS node")
}

func syncAndPulish() {
	blockChain, ethClient, rpcClient := getBlockChainAndClients()

	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	quitChan := make(chan bool)
	indexer, err := ipfs.NewIPFSIndexer(ipfsPath, &db, ethClient, rpcClient, quitChan)
	if err != nil {
		log.Fatal(err)
	}

	wg := syn.WaitGroup{}
	indexer.Index(wg)
	wg.Wait()
}

func getBlockChainAndClients() (*geth.BlockChain, core.EthClient, core.RpcClient) {
	rawRpcClient, err := rpc.Dial(ipc)

	if err != nil {
		log.Fatal(err)
	}
	rpcClient := client.NewRpcClient(rawRpcClient, ipc)
	ethClient := ethclient.NewClient(rawRpcClient)
	vdbEthClient := client.NewEthClient(ethClient)
	vdbNode := node.MakeNode(rpcClient)
	transactionConverter := vRpc.NewRpcTransactionConverter(ethClient)
	blockChain := geth.NewBlockChain(vdbEthClient, rpcClient, vdbNode, transactionConverter)
	return blockChain, ethClient, rpcClient
}