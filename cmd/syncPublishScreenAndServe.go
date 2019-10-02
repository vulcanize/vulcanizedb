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

	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// syncPublishScreenAndServeCmd represents the syncPublishScreenAndServe command
var syncPublishScreenAndServeCmd = &cobra.Command{
	Use:   "syncPublishScreenAndServe",
	Short: "Syncs all Ethereum data into IPFS, indexing the CIDs, and uses this to serve data requests to requesting clients",
	Long: `This command works alongside a modified geth node which streams
all block and state (diff) data over a websocket subscription. This process 
then converts the eth data to IPLD objects and publishes them to IPFS. Additionally,
it maintains a local index of the IPLD objects' CIDs in Postgres. It then opens up a server which 
relays relevant data to requesting clients.`,
	Run: func(cmd *cobra.Command, args []string) {
		syncPublishScreenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(syncPublishScreenAndServeCmd)
}

func syncPublishScreenAndServe() {
	superNode, err := newSuperNode()
	if err != nil {
		log.Fatal(err)
	}

	wg := &syn.WaitGroup{}
	forwardPayloadChan := make(chan ipfs.IPLDPayload, 20000)
	forwardQuitChan := make(chan bool, 1)
	err = superNode.SyncAndPublish(wg, forwardPayloadChan, forwardQuitChan)
	if err != nil {
		log.Fatal(err)
	}
	superNode.ScreenAndServe(forwardPayloadChan, forwardQuitChan)
	if viper.GetBool("backfill.on") && viper.GetString("backfill.ipcPath") != "" {
		backfiller := newBackFiller(superNode.GetPublisher())
		if err != nil {
			log.Fatal(err)
		}
		backfiller.FillGaps(wg, nil)
	}

	var ipcPath string
	ipcPath = viper.GetString("server.ipcPath")
	if ipcPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		ipcPath = filepath.Join(home, ".vulcanize/vulcanize.ipc")
	}
	_, _, err = rpc.StartIPCEndpoint(ipcPath, superNode.APIs())
	if err != nil {
		log.Fatal(err)
	}

	var wsEndpoint string
	wsEndpoint = viper.GetString("server.wsEndpoint")
	if wsEndpoint == "" {
		wsEndpoint = "127.0.0.1:8080"
	}
	var exposeAll = true
	var wsOrigins []string = nil
	_, _, err = rpc.StartWSEndpoint(wsEndpoint, superNode.APIs(), []string{"vulcanizedb"}, wsOrigins, exposeAll)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
