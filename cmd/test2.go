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
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/i-norden/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/libraries/shared/streamer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// test2Cmd represents the test2 command
var test2Cmd = &cobra.Command{
	Use:   "test2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		test2()
	},
}

func test2() {
	rpcClient := getRpcClient()
	str := streamer.NewSeedStreamer(rpcClient)
	payloadChan := make(chan ipfs.ResponsePayload, 800)
	filter := ipfs.StreamFilters{}
	filter.HeaderFilter.FinalOnly = true
	filter.TrxFilter.Src = []string{"0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"}
	filter.TrxFilter.Dst = []string{"0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"}
	filter.ReceiptFilter.Topic0s = []string{}
	filter.StateFilter.Addresses = []string{"0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"}
	filter.StorageFilter.Off = true
	sub, err := str.Stream(payloadChan, filter)
	if err != nil {
		println(err.Error())
		log.Fatal(err)
	}
	for {
		select {
		case payload := <-payloadChan:
			if payload.Err != nil {
				log.Error(payload.Err)
			}
			for _, headerRlp := range payload.HeadersRlp {
				header := new(types.Header)
				err = rlp.DecodeBytes(headerRlp, header)
				println("header")
				println(header.TxHash.Hex())
				println(header.Number.Int64())
			}
			for _, trxRlp := range payload.TransactionsRlp {
				trx := new(types.Transaction)
				err = rlp.DecodeBytes(trxRlp, trx)
				println("trx")
				println(trx.Hash().Hex())
				println(trx.Value().Int64())
			}
			for _, rctRlp := range payload.ReceiptsRlp {
				rct := new(types.Receipt)
				err = rlp.DecodeBytes(rctRlp, rct)
				println("rct")
				println(rct.TxHash.Hex())
				println(rct.BlockNumber.Bytes())
			}
			for _, stateRlp := range payload.StateNodesRlp {
				acct := new(state.Account)
				err = rlp.DecodeBytes(stateRlp, acct)
				println("state")
				println(acct.Root.Hex())
				println(acct.Balance.Int64())
			}
		case err = <-sub.Err():
			println(err.Error())
			log.Fatal(err)
		}
	}
}

func init() {
	rootCmd.AddCommand(test2Cmd)
	test2Cmd.Flags().StringVarP(&vulcPath, "ipc-path", "p", "~/.vulcanize/vulcanize.ipc", "IPC path for the Vulcanize seed node server")
}

func getRpcClient() core.RpcClient {
	println(vulcPath)
	rawRpcClient, err := rpc.Dial(vulcPath)
	if err != nil {
		log.Fatal(err)
	}
	return client.NewRpcClient(rawRpcClient, vulcPath)
}
