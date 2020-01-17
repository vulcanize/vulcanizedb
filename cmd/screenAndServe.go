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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/super_node"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/config"
	"github.com/vulcanize/vulcanizedb/utils"
)

// screenAndServeCmd represents the screenAndServe command
var screenAndServeCmd = &cobra.Command{
	Use:   "screenAndServe",
	Short: "Serve super-node data requests to requesting clients",
	Long: ` It then opens up WS and IPC servers on top of the super-node ETH-IPLD index which 
relays relevant data to requesting clients. In this mode, the super-node can only relay data which it has
already indexed it does not stream out live data.`,
	Run: func(cmd *cobra.Command, args []string) {
		subCommand = cmd.CalledAs()
		logWithCommand = *log.WithField("SubCommand", subCommand)
		screenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(screenAndServeCmd)
}

func screenAndServe() {
	superNode, err := newSuperNodeWithoutPairedGethNode()
	if err != nil {
		logWithCommand.Fatal(err)
	}
	wg := &syn.WaitGroup{}
	quitChan := make(chan bool, 1)
	emptyPayloadChan := make(chan interface{})
	superNode.ScreenAndServe(wg, emptyPayloadChan, quitChan)

	if err := startServers(superNode); err != nil {
		logWithCommand.Fatal(err)
	}
	wg.Wait()
}

func startServers(superNode super_node.NodeInterface) error {
	var ipcPath string
	ipcPath = viper.GetString("server.ipcPath")
	if ipcPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		ipcPath = filepath.Join(home, ".vulcanize/vulcanize.ipc")
	}
	_, _, err := rpc.StartIPCEndpoint(ipcPath, superNode.APIs())
	if err != nil {
		return err
	}

	var wsEndpoint string
	wsEndpoint = viper.GetString("server.wsEndpoint")
	if wsEndpoint == "" {
		wsEndpoint = "127.0.0.1:8080"
	}
	var exposeAll = true
	var wsOrigins []string
	_, _, err = rpc.StartWSEndpoint(wsEndpoint, superNode.APIs(), []string{"vdb"}, wsOrigins, exposeAll)
	if err != nil {
		return err
	}
	return nil
}

func newSuperNodeWithoutPairedGethNode() (super_node.NodeInterface, error) {
	ipfsPath = viper.GetString("client.ipfsPath")
	if ipfsPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		ipfsPath = filepath.Join(home, ".ipfs")
	}
	if err := ipfs.InitIPFSPlugins(); err != nil {
		return nil, err
	}
	ipldFetcher, err := super_node.NewIPLDFetcher(config.Ethereum, ipfsPath)
	if err != nil {
		return nil, err
	}
	db := utils.LoadPostgres(databaseConfig, core.Node{})
	retriever, err := super_node.NewCIDRetriever(config.Ethereum, &db)
	if err != nil {
		return nil, err
	}
	resolver, err := super_node.NewIPLDResolver(config.Ethereum)
	if err != nil {
		return nil, err
	}
	return &super_node.Service{
		IPLDFetcher:       ipldFetcher,
		Retriever:         retriever,
		Resolver:          resolver,
		Subscriptions:     make(map[common.Hash]map[rpc.ID]super_node.Subscription),
		SubscriptionTypes: make(map[common.Hash]super_node.SubscriptionSettings),
		NodeInfo:          core.Node{},
	}, nil
}
