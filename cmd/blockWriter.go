// Copyright © 2020 Vulcanize, Inc
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
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/utils"

	"github.com/spf13/cobra"
)

// blockWriterCmd represents the blockWriter command
var blockWriterCmd = &cobra.Command{
	Use:   "blockWriter",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("blockWriter called")
		blockWriter()
	},
}

var (
	blockNumber int
)

func blockWriter() {
	dbConfig := new(config.Database)
	dbConfig.Init()
	db := utils.LoadPostgres(*dbConfig, core.Node{})
	b, err := eth.NewEthBackend(&db)
	if err != nil {
		logrus.Fatal(err)
	}
	blk, err := b.BlockByNumber(context.Background(), rpc.BlockNumber(blockNumber))
	if err != nil {
		logrus.Fatal(err)
	}
	blkBytes, err := rlp.EncodeToBytes(blk)
	if err != nil {
		logrus.Fatal(err)
	}
	fileName := fmt.Sprintf("./block%d_rlp", blockNumber)
	logrus.Infof("writing block to file %s", fileName)
	f, err := os.Create(fileName)
	if err != nil {
		logrus.Fatal(err)
	}
	_, err = f.Write(blkBytes)
	if err != nil {
		logrus.Fatal(err)
	}
	f.Close()
}

func init() {
	rootCmd.AddCommand(blockWriterCmd)
	blockWriterCmd.Flags().IntVarP(&blockNumber, "block-number", "b", 0, "Block number to write to disk")
}
