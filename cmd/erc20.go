// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"log"
	"time"
)

// erc20Cmd represents the erc20 command
var erc20Cmd = &cobra.Command{
	Use:   "erc20",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		watchERC20s()
	},
}

func watchERC20s() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	blockchain := geth.NewBlockchain(ipc)
	db, err := postgres.NewDB(databaseConfig, blockchain.Node())
	if err != nil {
		log.Fatal("Failed to initialize database.")
	}
	watcher := shared.Watcher{
		DB:         *db,
		Blockchain: blockchain,
	}

	watcher.AddTransformers(every_block.TransformerInitializers())
	for range ticker.C {
		watcher.Execute()
	}
}

func init() {
	rootCmd.AddCommand(erc20Cmd)

	//var contractWatcherConfig string
	//contractViper := viper.New()
	//fmt.Println(contractViper)
	//
	//erc20Cmd.Flags().StringVar(&contractWatcherConfig, "contractWatcherConfig", "contract_watcher/config.toml", "config for desired ERC20 contracts to watch")
	//viper.BindPFlag("contract.config", erc20Cmd.Flags().Lookup("contractWatcherConfig"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// erc20Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// erc20Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
