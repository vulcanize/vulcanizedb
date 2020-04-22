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

package watcher

import (
	"context"
	"errors"
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/wasm"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/eth/client"
	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
	shared2 "github.com/vulcanize/vulcanizedb/pkg/watcher/shared"
	"github.com/vulcanize/vulcanizedb/utils"
)

// Config holds all of the parameters necessary for defining and running an instance of a watcher
type Config struct {
	// Subscription settings
	SubscriptionConfig shared.SubscriptionSettings
	// Database settings
	DBConfig config.Database
	// DB itself
	DB *postgres.DB
	// Subscription client
	Client interface{}
	// WASM instantiation paths and namespaces
	WASMFunctions []wasm.WasmFunction
	// File paths for trigger functions (sql files) that (can) use the instantiated wasm namespaces
	TriggerFunctions []string
	// Chain type used to specify what type of raw data we will be processing
	Chain shared.ChainType
	// Source type used to specify which streamer to use based on what API we will be interfacing with
	Source shared2.SourceType
	// Info for the node
	NodeInfo core.Node
}

func NewWatcherConfig() (*Config, error) {
	c := new(Config)
	var err error
	chain := viper.GetString("watcher.chain")
	c.Chain, err = shared.NewChainType(chain)
	if err != nil {
		return nil, err
	}
	switch c.Chain {
	case shared.Ethereum:
		c.SubscriptionConfig, err = eth.NewEthSubscriptionConfig()
		if err != nil {
			return nil, err
		}
	case shared.Bitcoin:
		c.SubscriptionConfig, err = btc.NewBtcSubscriptionConfig()
		if err != nil {
			return nil, err
		}
	case shared.Omni:
		return nil, errors.New("omni chain type currently not supported")
	default:
		return nil, fmt.Errorf("unexpected chain type %s", c.Chain.String())
	}
	sourcePath := viper.GetString("watcher.dataSource")
	if sourcePath == "" {
		sourcePath = "ws://127.0.0.1:8080" // default to and try the default ws url if no path is provided
	}
	sourceType := viper.GetString("watcher.dataPath")
	c.Source, err = shared2.NewSourceType(sourceType)
	if err != nil {
		return nil, err
	}
	switch c.Source {
	case shared2.Ethereum:
		return nil, errors.New("ethereum data source currently not supported")
	case shared2.Bitcoin:
		return nil, errors.New("bitcoin data source currently not supported")
	case shared2.VulcanizeDB:
		rawRPCClient, err := rpc.Dial(sourcePath)
		if err != nil {
			return nil, err
		}
		cli := client.NewRPCClient(rawRPCClient, sourcePath)
		var nodeInfo core.Node
		if err := cli.CallContext(context.Background(), &nodeInfo, "vdb_node"); err != nil {
			return nil, err
		}
		c.NodeInfo = nodeInfo
		c.Client = cli
	default:
		return nil, fmt.Errorf("unexpected data source type %s", c.Source.String())
	}
	wasmBinaries := viper.GetStringSlice("watcher.wasmBinaries")
	wasmNamespaces := viper.GetStringSlice("watcher.wasmNamespaces")
	if len(wasmBinaries) != len(wasmNamespaces) {
		return nil, fmt.Errorf("watcher config needs a namespace for every wasm binary\r\nhave %d binaries and %d namespaces", len(wasmBinaries), len(wasmNamespaces))
	}
	c.WASMFunctions = make([]wasm.WasmFunction, len(wasmBinaries))
	for i, bin := range wasmBinaries {
		c.WASMFunctions[i] = wasm.WasmFunction{
			BinaryPath: bin,
			Namespace:  wasmNamespaces[i],
		}
	}
	c.TriggerFunctions = viper.GetStringSlice("watcher.triggerFunctions")
	c.DBConfig = config.Database{
		Name:     viper.GetString("watcher.database.name"),
		Hostname: viper.GetString("watcher.database.hostname"),
		Port:     viper.GetInt("watcher.database.port"),
		User:     viper.GetString("watcher.database.user"),
		Password: viper.GetString("watcher.database.password"),
	}
	db := utils.LoadPostgres(c.DBConfig, c.NodeInfo)
	c.DB = &db
	return c, nil
}
