// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package resync

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/eth/client"
	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth/node"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
	"github.com/vulcanize/vulcanizedb/utils"
)

// Config holds the parameters needed to perform a resync
type Config struct {
	Chain         shared.ChainType // The type of resync to perform
	ResyncType    shared.DataType  // The type of data to resync
	ClearOldCache bool             // Resync will first clear all the data within the range

	// DB info
	DB       *postgres.DB
	DBConfig config.Database
	IPFSPath string

	HTTPClient  interface{} // Note this client is expected to support the retrieval of the specified data type(s)
	NodeInfo    core.Node   // Info for the associated node
	Ranges      [][2]uint64 // The block height ranges to resync
	BatchSize   uint64      // BatchSize for the resync http calls (client has to support batch sizing)
	BatchNumber uint64

	Quit chan bool // Channel for shutting down
}

// NewReSyncConfig fills and returns a resync config from toml parameters
func NewReSyncConfig() (*Config, error) {
	c := new(Config)
	var err error
	start := uint64(viper.GetInt64("resync.start"))
	stop := uint64(viper.GetInt64("resync.stop"))
	c.Ranges = [][2]uint64{{start, stop}}
	ipfsPath := viper.GetString("resync.ipfsPath")
	if ipfsPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		ipfsPath = filepath.Join(home, ".ipfs")
	}
	c.IPFSPath = ipfsPath
	c.DBConfig = config.Database{
		Name:     viper.GetString("database.name"),
		Hostname: viper.GetString("database.hostname"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
	}
	c.ClearOldCache = viper.GetBool("resync.clearOldCache")
	resyncType := viper.GetString("resync.type")
	c.ResyncType, err = shared.GenerateResyncTypeFromString(resyncType)
	if err != nil {
		return nil, err
	}
	chain := viper.GetString("resync.chain")
	c.Chain, err = shared.NewChainType(chain)
	if err != nil {
		return nil, err
	}
	if ok, err := shared.SupportedResyncType(c.ResyncType, c.Chain); !ok {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("chain type %s does not support data type %s", c.Chain.String(), c.ResyncType.String())
	}

	switch c.Chain {
	case shared.Ethereum:
		c.NodeInfo, c.HTTPClient, err = getEthNodeAndClient(fmt.Sprintf("http://%s", viper.GetString("ethereum.httpPath")))
		if err != nil {
			return nil, err
		}
	case shared.Bitcoin:
		c.NodeInfo = core.Node{
			ID:           viper.GetString("bitcoin.nodeID"),
			ClientName:   viper.GetString("bitcoin.clientName"),
			GenesisBlock: viper.GetString("bitcoin.genesisBlock"),
			NetworkID:    viper.GetString("bitcoin.networkID"),
		}
		// For bitcoin we load in node info from the config because there is no RPC endpoint to retrieve this from the node
		c.HTTPClient = &rpcclient.ConnConfig{
			Host:         viper.GetString("bitcoin.httpPath"),
			HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
			DisableTLS:   true, // Bitcoin core does not provide TLS by default
			Pass:         viper.GetString("bitcoin.pass"),
			User:         viper.GetString("bitcoin.user"),
		}
	}
	db := utils.LoadPostgres(c.DBConfig, c.NodeInfo)
	c.DB = &db
	c.Quit = make(chan bool)
	c.BatchSize = uint64(viper.GetInt64("resync.batchSize"))
	c.BatchNumber = uint64(viper.GetInt64("resync.batchNumber"))
	return c, nil
}

func getEthNodeAndClient(path string) (core.Node, interface{}, error) {
	rawRPCClient, err := rpc.Dial(path)
	if err != nil {
		return core.Node{}, nil, err
	}
	rpcClient := client.NewRPCClient(rawRPCClient, path)
	vdbNode := node.MakeNode(rpcClient)
	return vdbNode, rpcClient, nil
}
