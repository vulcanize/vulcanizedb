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

	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/utils"
)

// syncPublishScreenAndServeCmd represents the syncPublishScreenAndServe command
var syncPublishScreenAndServeCmd = &cobra.Command{
	Use:   "syncPublishScreenAndServe",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		syncPublishScreenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(syncPublishScreenAndServeCmd)
	syncPublishScreenAndServeCmd.Flags().StringVarP(&ipfsPath, "ipfs-path", "i", "~/.ipfs", "Path for configuring IPFS node")
	syncPublishScreenAndServeCmd.Flags().StringVarP(&vulcPath, "ipc-path", "p", "~/.vulcanize/vulcanize.ipc", "IPC path for the Vulcanize seed node server")
}

func syncPublishScreenAndServe() {
	blockChain, ethClient, rpcClient := getBlockChainAndClients()

	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	quitChan := make(chan bool)
	processor, err := ipfs.NewIPFSProcessor(ipfsPath, &db, ethClient, rpcClient, quitChan)
	if err != nil {
		log.Fatal(err)
	}

	wg := &syn.WaitGroup{}
	forwardPayloadChan := make(chan ipfs.IPLDPayload)
	forwardQuitChan := make(chan bool)
	err = processor.SyncAndPublish(wg, forwardPayloadChan, forwardQuitChan)
	if err != nil {
		log.Fatal(err)
	}
	processor.ScreenAndServe(wg, forwardPayloadChan, forwardQuitChan)
	_, _, err = rpc.StartIPCEndpoint(vulcPath, processor.APIs())
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
