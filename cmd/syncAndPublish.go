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
	"os"
	"path/filepath"
	syn "sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth"
	"github.com/vulcanize/vulcanizedb/pkg/eth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/eth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/eth/node"
	"github.com/vulcanize/vulcanizedb/pkg/super_node"
	"github.com/vulcanize/vulcanizedb/utils"
)

// syncAndPublishCmd represents the syncAndPublish command
var syncAndPublishCmd = &cobra.Command{
	Use:   "syncAndPublish",
	Short: "Syncs all Ethereum data into IPFS, indexing the CIDs",
	Long: `This command works alongside a modified geth node which streams
all block and state (diff) data over a websocket subscription. This process 
then converts the eth data to IPLD objects and publishes them to IPFS. Additionally,
it maintains a local index of the IPLD objects' CIDs in Postgres.`,
	Run: func(cmd *cobra.Command, args []string) {
		subCommand = cmd.CalledAs()
		logWithCommand = *log.WithField("SubCommand", subCommand)
		syncAndPublish()
	},
}

var ipfsPath string

func init() {
	rootCmd.AddCommand(syncAndPublishCmd)
}

func syncAndPublish() {
	superNode, newNodeErr := newSuperNode()
	if newNodeErr != nil {
		logWithCommand.Fatal(newNodeErr)
	}
	wg := &syn.WaitGroup{}
	syncAndPubErr := superNode.SyncAndPublish(wg, nil, nil)
	if syncAndPubErr != nil {
		logWithCommand.Fatal(syncAndPubErr)
	}
	if viper.GetBool("superNodeBackFill.on") && viper.GetString("superNodeBackFill.rpcPath") != "" {
		backfiller, newBackFillerErr := newBackFiller()
		if newBackFillerErr != nil {
			logWithCommand.Fatal(newBackFillerErr)
		}
		backfiller.FillGaps(wg, nil)
	}
	wg.Wait() // If an error was thrown, wg.Add was never called and this will fall through
}

func getBlockChainAndClient(path string) (*eth.BlockChain, core.RpcClient) {
	rawRPCClient, dialErr := rpc.Dial(path)
	if dialErr != nil {
		logWithCommand.Fatal(dialErr)
	}
	rpcClient := client.NewRpcClient(rawRPCClient, ipc)
	ethClient := ethclient.NewClient(rawRPCClient)
	vdbEthClient := client.NewEthClient(ethClient)
	vdbNode := node.MakeNode(rpcClient)
	transactionConverter := vRpc.NewRpcTransactionConverter(ethClient)
	blockChain := eth.NewBlockChain(vdbEthClient, rpcClient, vdbNode, transactionConverter)
	return blockChain, rpcClient
}

func newSuperNode() (super_node.NodeInterface, error) {
	blockChain, rpcClient := getBlockChainAndClient(ipc)
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	quitChan := make(chan bool)
	ipfsPath = viper.GetString("client.ipfsPath")
	if ipfsPath == "" {
		home, homeDirErr := os.UserHomeDir()
		if homeDirErr != nil {
			logWithCommand.Fatal(homeDirErr)
		}
		ipfsPath = filepath.Join(home, ".ipfs")
	}
	workers := viper.GetInt("client.workers")
	if workers < 1 {
		workers = 1
	}
	return super_node.NewSuperNode(ipfsPath, &db, rpcClient, quitChan, workers, blockChain.Node())
}

func newBackFiller() (super_node.BackFillInterface, error) {
	blockChain, archivalRPCClient := getBlockChainAndClient(viper.GetString("superNodeBackFill.rpcPath"))
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	freq := viper.GetInt("superNodeBackFill.frequency")
	var frequency time.Duration
	if freq <= 0 {
		frequency = time.Minute * 5
	} else {
		frequency = time.Duration(freq)
	}
	return super_node.NewBackFillService(ipfsPath, &db, archivalRPCClient, time.Minute*frequency, super_node.DefaultMaxBatchSize)
}
