// Copyright © 2018 Rob Mulholand <rmulholand@8thlight.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vulcanize/vulcanizedb/pkg/config"
)

var (
	blockNumber         int64
	cfgFile             string
	computeState        bool
	databaseConfig      config.Database
	endingBlockNumber   int64
	ipc                 string
	ipfsPath            string
	levelDbPath         string
	startingBlockNumber int64
)

var rootCmd = &cobra.Command{
	Use:              "blockWatcher",
	PersistentPreRun: database,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func database(cmd *cobra.Command, args []string) {
	ipc = viper.GetString("client.ipcpath")
	levelDbPath = viper.GetString("client.leveldbpath")
	ipfsPath = viper.GetString("client.ipfspath")
	databaseConfig = config.Database{
		Name:     viper.GetString("database.name"),
		Hostname: viper.GetString("database.hostname"),
		Port:     viper.GetInt("database.port"),
	}
	viper.Set("database.config", databaseConfig)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "environments/public.toml", "config file location")
	rootCmd.PersistentFlags().String("database-name", "vulcanize_public", "database name")
	rootCmd.PersistentFlags().Int("database-port", 5432, "database port")
	rootCmd.PersistentFlags().String("database-hostname", "localhost", "database hostname")
	rootCmd.PersistentFlags().String("client-ipcPath", "", "location of geth.ipc file")
	rootCmd.PersistentFlags().String("client-ipfsPath", "", "location of ipfs directory")
	rootCmd.PersistentFlags().String("client-levelDbPath", "", "location of levelDb chaindata")

	viper.BindPFlag("database.name", rootCmd.PersistentFlags().Lookup("database-name"))
	viper.BindPFlag("database.port", rootCmd.PersistentFlags().Lookup("database-port"))
	viper.BindPFlag("database.hostname", rootCmd.PersistentFlags().Lookup("database-hostname"))
	viper.BindPFlag("client.ipcPath", rootCmd.PersistentFlags().Lookup("client-ipcPath"))
	viper.BindPFlag("client.ipfsPath", rootCmd.PersistentFlags().Lookup("client-ipfsPath"))
	viper.BindPFlag("client.levelDbPath", rootCmd.PersistentFlags().Lookup("client-levelDbPath"))

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".vulcanizedb")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s\n\n", viper.ConfigFileUsed())
	}
}
