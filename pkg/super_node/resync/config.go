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
	RESYNC_CHAIN            = "RESYNC_CHAIN"
	RESYNC_START            = "RESYNC_START"
	RESYNC_STOP             = "RESYNC_STOP"
	RESYNC_BATCH_SIZE       = "RESYNC_BATCH_SIZE"
	RESYNC_BATCH_NUMBER     = "RESYNC_BATCH_NUMBER"
	RESYNC_CLEAR_OLD_CACHE  = "RESYNC_CLEAR_OLD_CACHE"
	RESYNC_TYPE             = "RESYNC_TYPE"
	RESYNC_RESET_VALIDATION = "RESYNC_RESET_VALIDATION"
)

// Config holds the parameters needed to perform a resync
type Config struct {
	Chain           shared.ChainType // The type of resync to perform
	ResyncType      shared.DataType  // The type of data to resync
	ClearOldCache   bool             // Resync will first clear all the data within the range
	ResetValidation bool             // If true, resync will reset the validation level to 0 for the given range

	// DB info
	DB       *postgres.DB
	DBConfig config.Database
	IPFSPath string
	IPFSMode shared.IPFSMode

	HTTPClient  interface{}   // Note this client is expected to support the retrieval of the specified data type(s)
	NodeInfo    core.Node     // Info for the associated node
	Ranges      [][2]uint64   // The block height ranges to resync
	BatchSize   uint64        // BatchSize for the resync http calls (client has to support batch sizing)
	Timeout     time.Duration // HTTP connection timeout in seconds
	BatchNumber uint64

	Quit chan bool // Channel for shutting down
}

// NewReSyncConfig fills and returns a resync config from toml parameters
func NewReSyncConfig() (*Config, error) {
	c := new(Config)
	var err error

	viper.BindEnv("resync.start", RESYNC_START)
	viper.BindEnv("resync.stop", RESYNC_STOP)
	viper.BindEnv("resync.clearOldCache", RESYNC_CLEAR_OLD_CACHE)
	viper.BindEnv("resync.type", RESYNC_TYPE)
	viper.BindEnv("resync.chain", RESYNC_CHAIN)
	viper.BindEnv("ethereum.httpPath", shared.ETH_HTTP_PATH)
	viper.BindEnv("bitcoin.httpPath", shared.BTC_HTTP_PATH)
	viper.BindEnv("resync.batchSize", RESYNC_BATCH_SIZE)
	viper.BindEnv("resync.batchNumber", RESYNC_BATCH_NUMBER)
	viper.BindEnv("resync.resetValidation", RESYNC_RESET_VALIDATION)
	viper.BindEnv("resync.timeout", shared.HTTP_TIMEOUT)

	timeout := viper.GetInt("resync.timeout")
	if timeout < 15 {
		timeout = 15
	}
	c.Timeout = time.Second * time.Duration(timeout)

	start := uint64(viper.GetInt64("resync.start"))
	stop := uint64(viper.GetInt64("resync.stop"))
	c.Ranges = [][2]uint64{{start, stop}}
	c.ClearOldCache = viper.GetBool("resync.clearOldCache")
	c.ResetValidation = viper.GetBool("resync.resetValidation")

	c.IPFSMode, err = shared.GetIPFSMode()
	if err != nil {
		return nil, err
	}
	if c.IPFSMode == shared.LocalInterface || c.IPFSMode == shared.RemoteClient {
		c.IPFSPath, err = shared.GetIPFSPath()
		if err != nil {
			return nil, err
		}
	}
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
		ethHTTP := viper.GetString("ethereum.httpPath")
		c.NodeInfo, c.HTTPClient, err = shared.GetEthNodeAndClient(fmt.Sprintf("http://%s", ethHTTP))
		if err != nil {
			return nil, err
		}
	case shared.Bitcoin:
		btcHTTP := viper.GetString("bitcoin.httpPath")
		c.NodeInfo, c.HTTPClient = shared.GetBtcNodeAndClient(btcHTTP)
	}

	c.DBConfig.Init()
	db := utils.LoadPostgres(c.DBConfig, c.NodeInfo)
	c.DB = &db

	c.Quit = make(chan bool)
	c.BatchSize = uint64(viper.GetInt64("resync.batchSize"))
	c.BatchNumber = uint64(viper.GetInt64("resync.batchNumber"))
	return c, nil
}
