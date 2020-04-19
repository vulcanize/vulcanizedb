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
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
	"github.com/vulcanize/vulcanizedb/utils"
)

// Env variables
const (
	SUPERNODE_CHAIN            = "SUPERNODE_CHAIN"
	SUPERNODE_SYNC             = "SUPERNODE_SYNC"
	SUPERNODE_WORKERS          = "SUPERNODE_WORKERS"
	SUPERNODE_SERVER           = "SUPERNODE_SERVER"
	SUPERNODE_WS_PATH          = "SUPERNODE_WS_PATH"
	SUPERNODE_IPC_PATH         = "SUPERNODE_IPC_PATH"
	SUPERNODE_HTTP_PATH        = "SUPERNODE_HTTP_PATH"
	SUPERNODE_BACKFILL         = "SUPERNODE_BACKFILL"
	SUPERNODE_FREQUENCY        = "SUPERNODE_FREQUENCY"
	SUPERNODE_BATCH_SIZE       = "SUPERNODE_BATCH_SIZE"
	SUPERNODE_BATCH_NUMBER     = "SUPERNODE_BATCH_NUMBER"
	SUPERNODE_VALIDATION_LEVEL = "SUPERNODE_VALIDATION_LEVEL"
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
	BackFill        bool
	HTTPClient      interface{}
	Frequency       time.Duration
	BatchSize       uint64
	BatchNumber     uint64
	ValidationLevel int
	Timeout         time.Duration // HTTP connection timeout in seconds
}

// NewSuperNodeConfig is used to initialize a SuperNode config from a .toml file
// Separate chain supernode instances need to be ran with separate ipfs path in order to avoid lock contention on the ipfs repository lockfile
func NewSuperNodeConfig() (*Config, error) {
	c := new(Config)
	var err error

	viper.BindEnv("superNode.chain", SUPERNODE_CHAIN)
	viper.BindEnv("superNode.sync", SUPERNODE_SYNC)
	viper.BindEnv("superNode.workers", SUPERNODE_WORKERS)
	viper.BindEnv("ethereum.wsPath", shared.ETH_WS_PATH)
	viper.BindEnv("bitcoin.wsPath", shared.BTC_WS_PATH)
	viper.BindEnv("superNode.server", SUPERNODE_SERVER)
	viper.BindEnv("superNode.wsPath", SUPERNODE_WS_PATH)
	viper.BindEnv("superNode.ipcPath", SUPERNODE_IPC_PATH)
	viper.BindEnv("superNode.httpPath", SUPERNODE_HTTP_PATH)
	viper.BindEnv("superNode.backFill", SUPERNODE_BACKFILL)
	viper.BindEnv("superNode.timeout", shared.HTTP_TIMEOUT)

	timeout := viper.GetInt("superNode.timeout")
	if timeout < 15 {
		timeout = 15
	}
	c.Timeout = time.Second * time.Duration(timeout)

	chain := viper.GetString("superNode.chain")
	c.Chain, err = shared.NewChainType(chain)
	if err != nil {
		return nil, err
	}

	c.IPFSPath, err = shared.GetIPFSPath()
	if err != nil {
		return nil, err
	}

	c.Sync = viper.GetBool("superNode.sync")
	if c.Sync {
		workers := viper.GetInt("superNode.workers")
		if workers < 1 {
			workers = 1
		}
		c.Workers = workers
		switch c.Chain {
		case shared.Ethereum:
			ethWS := viper.GetString("ethereum.wsPath")
			c.NodeInfo, c.WSClient, err = shared.GetEthNodeAndClient(fmt.Sprintf("ws://%s", ethWS))
			if err != nil {
				return nil, err
			}
		case shared.Bitcoin:
			btcWS := viper.GetString("bitcoin.wsPath")
			c.NodeInfo, c.WSClient = shared.GetBtcNodeAndClient(btcWS)
		}
	}

	c.Serve = viper.GetBool("superNode.server")
	if c.Serve {
		wsPath := viper.GetString("superNode.wsPath")
		if wsPath == "" {
			wsPath = "127.0.0.1:8080"
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
			httpPath = "127.0.0.1:8081"
		}
		c.HTTPEndpoint = httpPath
	}

	c.BackFill = viper.GetBool("superNode.backFill")
	if c.BackFill {
		if err := c.BackFillFields(); err != nil {
			return nil, err
		}
	}

	c.DBConfig.Init()
	db := utils.LoadPostgres(c.DBConfig, c.NodeInfo)
	c.DB = &db
	c.Quit = make(chan bool)

	return c, nil
}

// BackFillFields is used to fill in the BackFill fields of the config
func (c *Config) BackFillFields() error {
	var err error

	viper.BindEnv("ethereum.httpPath", shared.ETH_HTTP_PATH)
	viper.BindEnv("bitcoin.httpPath", shared.BTC_HTTP_PATH)
	viper.BindEnv("superNode.frequency", SUPERNODE_FREQUENCY)
	viper.BindEnv("superNode.batchSize", SUPERNODE_BATCH_SIZE)
	viper.BindEnv("superNode.batchNumber", SUPERNODE_BATCH_NUMBER)
	viper.BindEnv("superNode.validationLevel", SUPERNODE_VALIDATION_LEVEL)

	switch c.Chain {
	case shared.Ethereum:
		ethHTTP := viper.GetString("ethereum.httpPath")
		c.NodeInfo, c.HTTPClient, err = shared.GetEthNodeAndClient(fmt.Sprintf("http://%s", ethHTTP))
		if err != nil {
			return err
		}
	case shared.Bitcoin:
		btcHTTP := viper.GetString("bitcoin.httpPath")
		c.NodeInfo, c.HTTPClient = shared.GetBtcNodeAndClient(btcHTTP)
	}

	freq := viper.GetInt("superNode.frequency")
	var frequency time.Duration
	if freq <= 0 {
		frequency = time.Second * 30
	} else {
		frequency = time.Second * time.Duration(freq)
	}
	c.Frequency = frequency
	c.BatchSize = uint64(viper.GetInt64("superNode.batchSize"))
	c.BatchNumber = uint64(viper.GetInt64("superNode.batchNumber"))
	c.ValidationLevel = viper.GetInt("superNode.validationLevel")
	return nil
}
