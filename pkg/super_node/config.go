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

package super_node

import (
	"os"
	"path/filepath"
	"time"

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

// Config struct
type Config struct {
	// Ubiquitous fields
	Chain    shared.ChainType
	IPFSPath string
	DB       *postgres.DB
	DBConfig config.Database
	Quit     chan bool
	// Server fields
	Serve        bool
	WSEndpoint   string
	HTTPEndpoint string
	IPCEndpoint  string
	// Sync params
	Sync     bool
	Workers  int
	WSClient interface{}
	NodeInfo core.Node
	// Backfiller params
	BackFill   bool
	HTTPClient interface{}
	Frequency  time.Duration
	BatchSize  uint64
}

// NewSuperNodeConfig is used to initialize a SuperNode config from a .toml file
// Separate chain supernode instances need to be ran with separate ipfs path in order to avoid lock contention on the ipfs repository lockfile
func NewSuperNodeConfig() (*Config, error) {
	c := new(Config)
	var err error

	chain := viper.GetString("superNode.chain")
	c.Chain, err = shared.NewChainType(chain)
	if err != nil {
		return nil, err
	}

	ipfsPath := viper.GetString("superNode.ipfsPath")
	if ipfsPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		ipfsPath = filepath.Join(home, ".ipfs")
	}
	c.IPFSPath = ipfsPath

	c.Sync = viper.GetBool("superNode.sync")
	if c.Sync {
		workers := viper.GetInt("superNode.workers")
		if workers < 1 {
			workers = 1
		}
		c.Workers = workers
		switch c.Chain {
		case shared.Ethereum:
			c.NodeInfo, c.WSClient, err = getEthNodeAndClient(viper.GetString("ethereum.wsPath"))
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
			c.WSClient = &rpcclient.ConnConfig{
				Host:         viper.GetString("bitcoin.wsPath"),
				HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
				DisableTLS:   true, // Bitcoin core does not provide TLS by default
				Pass:         viper.GetString("bitcoin.pass"),
				User:         viper.GetString("bitcoin.user"),
			}
		}
	}

	c.Serve = viper.GetBool("superNode.server")
	if c.Serve {
		wsPath := viper.GetString("superNode.wsPath")
		if wsPath == "" {
			wsPath = "ws://127.0.0.1:8546"
		}
		c.WSEndpoint = wsPath
		ipcPath := viper.GetString("superNode.ipcPath")
		if ipcPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			ipcPath = filepath.Join(home, ".vulcanize/vulcanize.ipc")
		}
		c.IPCEndpoint = ipcPath
		httpPath := viper.GetString("superNode.httpPath")
		if httpPath == "" {
			httpPath = "http://127.0.0.1:8545"
		}
		c.HTTPEndpoint = httpPath
	}

	c.DBConfig = config.Database{
		Name:     viper.GetString("database.name"),
		Hostname: viper.GetString("database.hostname"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
	}

	db := utils.LoadPostgres(c.DBConfig, c.NodeInfo)
	c.DB = &db
	c.Quit = make(chan bool)
	if viper.GetBool("superNode.backFill") {
		if err := c.BackFillFields(chain); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// BackFillFields is used to fill in the BackFill fields of the config
func (sn *Config) BackFillFields(chain string) error {
	sn.BackFill = true
	var httpClient interface{}
	var err error
	switch sn.Chain {
	case shared.Ethereum:
		_, httpClient, err = getEthNodeAndClient(viper.GetString("ethereum.httpPath"))
		if err != nil {
			return err
		}
	case shared.Bitcoin:
		httpClient = &rpcclient.ConnConfig{
			Host:         viper.GetString("bitcoin.httpPath"),
			HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
			DisableTLS:   true, // Bitcoin core does not provide TLS by default
			Pass:         viper.GetString("bitcoin.pass"),
			User:         viper.GetString("bitcoin.user"),
		}
	}
	sn.HTTPClient = httpClient
	freq := viper.GetInt("superNode.frequency")
	var frequency time.Duration
	if freq <= 0 {
		frequency = time.Second * 30
	} else {
		frequency = time.Second * time.Duration(freq)
	}
	sn.Frequency = frequency
	sn.BatchSize = uint64(viper.GetInt64("superNode.batchSize"))
	return nil
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
