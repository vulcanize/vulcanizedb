// Copyright © 2019 Vulcanize, Inc
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
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		test()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func test() {
	_, _, rpcClient := getBlockChainAndClients()
	streamer := ipfs.NewStateDiffStreamer(rpcClient)
	payloadChan := make(chan statediff.Payload, 800)
	sub, err := streamer.Stream(payloadChan)
	if err != nil {
		println(err.Error())
		log.Fatal(err)
	}
	for {
		select {
		case payload := <-payloadChan:
			fmt.Printf("blockRlp: %v\r\nstateDiffRlp: %v\r\nerror: %v\r\n", payload.BlockRlp, payload.StateDiffRlp, payload.Err)
			var block types.Block
			err := rlp.DecodeBytes(payload.BlockRlp, &block)
			if err != nil {
				log.Fatal(err)
			}
			var stateDiff statediff.StateDiff
			err = rlp.DecodeBytes(payload.StateDiffRlp, &stateDiff)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("block: %v\r\nstateDiff: %v\r\n", block, stateDiff)
			fmt.Printf("block number: %d\r\n", block.Number())
		case err = <-sub.Err():
			println(err.Error())
			log.Fatal(err)
		}
	}
}
