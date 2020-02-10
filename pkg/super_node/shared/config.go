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

package shared

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/eth"
	"github.com/vulcanize/vulcanizedb/pkg/eth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/eth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth/node"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/utils"
)

// SuperNodeConfig struct
type SuperNodeConfig struct {
	// Ubiquitous fields
	Chain    ChainType
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

// NewSuperNodeConfig is used to initialize a SuperNode config from a config .toml file
func NewSuperNodeConfig() (*SuperNodeConfig, error) {
	sn := new(SuperNodeConfig)
	sn.DBConfig = config.Database{
		Name:     viper.GetString("superNode.database.name"),
		Hostname: viper.GetString("superNode.database.hostname"),
		Port:     viper.GetInt("superNode.database.port"),
		User:     viper.GetString("superNode.database.user"),
		Password: viper.GetString("superNode.database.password"),
	}
	var err error
	sn.Chain, err = NewChainType(viper.GetString("superNode.chain"))
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
	sn.IPFSPath = ipfsPath
	sn.Serve = viper.GetBool("superNode.server.on")
	sn.Sync = viper.GetBool("superNode.sync.on")
	if sn.Sync {
		workers := viper.GetInt("superNode.sync.workers")
		if workers < 1 {
			workers = 1
		}
		sn.Workers = workers
		if sn.Chain == Ethereum {
			sn.NodeInfo, sn.WSClient, err = getEthNodeAndClient(sn.Chain, viper.GetString("superNode.sync.wsPath"))
		}
		if sn.Chain == Bitcoin {
			sn.NodeInfo = core.Node{
				ID:           "temporaryID",
				ClientName:   "omnicored",
				GenesisBlock: "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
				NetworkID:    "0xD9B4BEF9",
			}
			sn.WSClient = &rpcclient.ConnConfig{
				Host:     viper.GetString("superNode.sync.wsPath"),
				Endpoint: "ws",
			}
		}
	}
	if sn.Serve {
		wsPath := viper.GetString("superNode.server.wsPath")
		if wsPath == "" {
			wsPath = "ws://127.0.0.1:8546"
		}
		sn.WSEndpoint = wsPath
		ipcPath := viper.GetString("superNode.server.ipcPath")
		if ipcPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			ipcPath = filepath.Join(home, ".vulcanize/vulcanize.ipc")
		}
		sn.IPCEndpoint = ipcPath
		httpPath := viper.GetString("superNode.server.httpPath")
		if httpPath == "" {
			httpPath = "http://127.0.0.1:8547"
		}
		sn.HTTPEndpoint = httpPath
	}
	db := utils.LoadPostgres(sn.DBConfig, sn.NodeInfo)
	sn.DB = &db
	sn.Quit = make(chan bool)
	if viper.GetBool("superNode.backFill.on") {
		if err := sn.BackFillFields(); err != nil {
			return nil, err
		}
	}
	return sn, err
}

// BackFillFields is used to fill in the BackFill fields of the config
func (sn *SuperNodeConfig) BackFillFields() error {
	sn.BackFill = true
	var httpClient interface{}
	var err error
	if sn.Chain == Ethereum {
		_, httpClient, err = getEthNodeAndClient(sn.Chain, viper.GetString("superNode.backFill.httpPath"))
		if err != nil {
			return err
		}
	}
	if sn.Chain == Bitcoin {
		httpClient = &rpcclient.ConnConfig{
			Host: viper.GetString("superNode.backFill.httpPath"),
		}
	}
	sn.HTTPClient = httpClient
	freq := viper.GetInt("superNode.backFill.frequency")
	var frequency time.Duration
	if freq <= 0 {
		frequency = time.Minute * 5
	} else {
		frequency = time.Duration(freq)
	}
	sn.Frequency = frequency
	sn.BatchSize = uint64(viper.GetInt64("superNode.backFill.batchSize"))
	return nil
}

func getEthNodeAndClient(chain ChainType, path string) (core.Node, interface{}, error) {
	rawRPCClient, err := rpc.Dial(path)
	if err != nil {
		return core.Node{}, nil, err
	}
	rpcClient := client.NewRPCClient(rawRPCClient, path)
	ethClient := ethclient.NewClient(rawRPCClient)
	vdbEthClient := client.NewEthClient(ethClient)
	vdbNode := node.MakeNode(rpcClient)
	transactionConverter := vRpc.NewRPCTransactionConverter(ethClient)
	blockChain := eth.NewBlockChain(vdbEthClient, rpcClient, vdbNode, transactionConverter)
	return blockChain.Node(), rpcClient, nil
}
