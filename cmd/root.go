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

package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth"
	"github.com/vulcanize/vulcanizedb/pkg/eth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/eth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/eth/node"
)

var (
	cfgFile              string
	databaseConfig       config.Database
	genConfig            config.Plugin
	subscriptionConfig   config.Subscription
	ipc                  string
	levelDbPath          string
	queueRecheckInterval time.Duration
	startingBlockNumber  int64
	storageDiffsPath     string
	syncAll              bool
	endingBlockNumber    int64
	recheckHeadersArg    bool
	subCommand           string
	logWithCommand       log.Entry
	storageDiffsSource   string
)

const (
	pollingInterval  = 7 * time.Second
	validationWindow = 15
)

var rootCmd = &cobra.Command{
	Use:              "vulcanizedb",
	PersistentPreRun: initFuncs,
}

func Execute() {
	log.Info("----- Starting vDB -----")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initFuncs(cmd *cobra.Command, args []string) {
	setViperConfigs()
	logLvlErr := logLevel()
	if logLvlErr != nil {
		log.Fatal("Could not set log level: ", logLvlErr)
	}

}

func setViperConfigs() {
	ipc = viper.GetString("client.ipcpath")
	levelDbPath = viper.GetString("client.leveldbpath")
	storageDiffsPath = viper.GetString("filesystem.storageDiffsPath")
	storageDiffsSource = viper.GetString("storageDiffs.source")
	databaseConfig = config.Database{
		Name:     viper.GetString("database.name"),
		Hostname: viper.GetString("database.hostname"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
	}
	viper.Set("database.config", databaseConfig)
}

func logLevel() error {
	lvl, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		return err
	}
	log.SetLevel(lvl)
	if lvl > log.InfoLevel {
		log.SetReportCaller(true)
	}
	log.Info("Log level set to ", lvl.String())
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)
	// When searching for env variables, replace dots in config keys with underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file location")
	rootCmd.PersistentFlags().String("database-name", "vulcanize_public", "database name")
	rootCmd.PersistentFlags().Int("database-port", 5432, "database port")
	rootCmd.PersistentFlags().String("database-hostname", "localhost", "database hostname")
	rootCmd.PersistentFlags().String("database-user", "", "database user")
	rootCmd.PersistentFlags().String("database-password", "", "database password")
	rootCmd.PersistentFlags().String("client-ipcPath", "", "location of geth.ipc file")
	rootCmd.PersistentFlags().String("client-levelDbPath", "", "location of levelDb chaindata")
	rootCmd.PersistentFlags().String("filesystem-storageDiffsPath", "", "location of storage diffs csv file")
	rootCmd.PersistentFlags().String("storageDiffs-source", "csv", "where to get the state diffs: csv or geth")
	rootCmd.PersistentFlags().String("exporter-name", "exporter", "name of exporter plugin")
	rootCmd.PersistentFlags().String("log-level", log.InfoLevel.String(), "Log level (trace, debug, info, warn, error, fatal, panic")

	viper.BindPFlag("database.name", rootCmd.PersistentFlags().Lookup("database-name"))
	viper.BindPFlag("database.port", rootCmd.PersistentFlags().Lookup("database-port"))
	viper.BindPFlag("database.hostname", rootCmd.PersistentFlags().Lookup("database-hostname"))
	viper.BindPFlag("database.user", rootCmd.PersistentFlags().Lookup("database-user"))
	viper.BindPFlag("database.password", rootCmd.PersistentFlags().Lookup("database-password"))
	viper.BindPFlag("client.ipcPath", rootCmd.PersistentFlags().Lookup("client-ipcPath"))
	viper.BindPFlag("client.levelDbPath", rootCmd.PersistentFlags().Lookup("client-levelDbPath"))
	viper.BindPFlag("filesystem.storageDiffsPath", rootCmd.PersistentFlags().Lookup("filesystem-storageDiffsPath"))
	viper.BindPFlag("storageDiffs.source", rootCmd.PersistentFlags().Lookup("storageDiffs-source"))
	viper.BindPFlag("exporter.fileName", rootCmd.PersistentFlags().Lookup("exporter-name"))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		noConfigError := "No config file passed with --config flag"
		fmt.Println("Error: ", noConfigError)
		log.Fatal(noConfigError)
	}

	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s\n\n", viper.ConfigFileUsed())
	} else {
		invalidConfigError := "Couldn't read config file"
		formattedError := fmt.Sprintf("%s: %s", invalidConfigError, err.Error())
		log.Fatal(formattedError)
	}
}

func getBlockChain() *eth.BlockChain {
	rpcClient, ethClient := getClients()
	vdbEthClient := client.NewEthClient(ethClient)
	vdbNode := node.MakeNode(rpcClient)
	transactionConverter := vRpc.NewRPCTransactionConverter(ethClient)
	return eth.NewBlockChain(vdbEthClient, rpcClient, vdbNode, transactionConverter)
}

func getClients() (client.RPCClient, *ethclient.Client) {
	rawRPCClient, err := rpc.Dial(ipc)

	if err != nil {
		logWithCommand.Fatal(err)
	}
	rpcClient := client.NewRPCClient(rawRPCClient, ipc)
	ethClient := ethclient.NewClient(rawRPCClient)

	return rpcClient, ethClient
}

func getWSClient() core.RPCClient {
	wsRPCpath := viper.GetString("client.wsPath")
	if wsRPCpath == "" {
		logWithCommand.Fatal(errors.New("getWSClient() was called but no ws rpc path is provided"))
	}
	wsRPCClient, dialErr := rpc.Dial(wsRPCpath)
	if dialErr != nil {
		logWithCommand.Fatal(dialErr)
	}
	return client.NewRPCClient(wsRPCClient, wsRPCpath)
}
