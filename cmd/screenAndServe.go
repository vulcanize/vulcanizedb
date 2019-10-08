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

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/super_node"
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
		screenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(screenAndServeCmd)
}

func screenAndServe() {
	superNode, newNodeErr := newSuperNodeWithoutPairedGethNode()
	if newNodeErr != nil {
		log.Fatal(newNodeErr)
	}
	wg := &syn.WaitGroup{}
	quitChan := make(chan bool, 1)
	emptyPayloadChan := make(chan ipfs.IPLDPayload)
	superNode.ScreenAndServe(wg, emptyPayloadChan, quitChan)

	serverErr := startServers(superNode)
	if serverErr != nil {
		log.Fatal(serverErr)
	}
	wg.Wait()
}

func startServers(superNode super_node.NodeInterface) error {
	var ipcPath string
	ipcPath = viper.GetString("server.ipcPath")
	if ipcPath == "" {
		home, homeDirErr := os.UserHomeDir()
		if homeDirErr != nil {
			return homeDirErr
		}
		ipcPath = filepath.Join(home, ".vulcanize/vulcanize.ipc")
	}
	_, _, ipcErr := rpc.StartIPCEndpoint(ipcPath, superNode.APIs())
	if ipcErr != nil {
		return ipcErr
	}

	var wsEndpoint string
	wsEndpoint = viper.GetString("server.wsEndpoint")
	if wsEndpoint == "" {
		wsEndpoint = "127.0.0.1:8080"
	}
	var exposeAll = true
	var wsOrigins []string = nil
	_, _, wsErr := rpc.StartWSEndpoint(wsEndpoint, superNode.APIs(), []string{"vdb"}, wsOrigins, exposeAll)
	if wsErr != nil {
		return wsErr
	}
	return nil
}

func newSuperNodeWithoutPairedGethNode() (super_node.NodeInterface, error) {
	ipfsPath = viper.GetString("client.ipfsPath")
	if ipfsPath == "" {
		home, homeDirErr := os.UserHomeDir()
		if homeDirErr != nil {
			return nil, homeDirErr
		}
		ipfsPath = filepath.Join(home, ".ipfs")
	}
	ipfsInitErr := ipfs.InitIPFSPlugins()
	if ipfsInitErr != nil {
		return nil, ipfsInitErr
	}
	ipldFetcher, newFetcherErr := ipfs.NewIPLDFetcher(ipfsPath)
	if newFetcherErr != nil {
		return nil, newFetcherErr
	}
	db := utils.LoadPostgres(databaseConfig, core.Node{})
	return &super_node.Service{
		IPLDFetcher:       ipldFetcher,
		Retriever:         super_node.NewCIDRetriever(&db),
		Resolver:          ipfs.NewIPLDResolver(),
		Subscriptions:     make(map[common.Hash]map[rpc.ID]super_node.Subscription),
		SubscriptionTypes: make(map[common.Hash]config.Subscription),
		GethNode:          core.Node{},
	}, nil
}
