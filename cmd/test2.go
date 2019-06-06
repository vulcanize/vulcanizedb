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
	"bytes"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
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
	payloadChan := make(chan ipfs.ResponsePayload, 8000)
	streamFilters := ipfs.StreamFilters{}
	streamFilters.HeaderFilter.FinalOnly = true
	streamFilters.ReceiptFilter.Topic0s = []string{
		"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		"0x930a61a57a70a73c2a503615b87e2e54fe5b9cdeacda518270b852296ab1a377",
	}
	streamFilters.StateFilter.Addresses = []string{"0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"}
	streamFilters.StorageFilter.Off = true
	//streamFilters.TrxFilter.Src = []string{"0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"}
	//streamFilters.TrxFilter.Dst = []string{"0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"}
	sub, err := str.Stream(payloadChan, streamFilters)
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
				var header types.Header
				err = rlp.Decode(bytes.NewBuffer(headerRlp), &header)
				if err != nil {
					println(err.Error())
					log.Error(err)
				}
				println("header")
				println(header.Hash().Hex())
				println(header.Number.Int64())
			}
			for _, trxRlp := range payload.TransactionsRlp {
				var trx types.Transaction
				buff := bytes.NewBuffer(trxRlp)
				stream := rlp.NewStream(buff, 0)
				err := trx.DecodeRLP(stream)
				if err != nil {
					println(err.Error())
					log.Error(err)
				}
				println("trx")
				println(trx.Hash().Hex())
			}
			for _, rctRlp := range payload.ReceiptsRlp {
				var rct types.Receipt
				buff := bytes.NewBuffer(rctRlp)
				stream := rlp.NewStream(buff, 0)
				err = rct.DecodeRLP(stream)
				if err != nil {
					println(err.Error())
					log.Error(err)
				}
				println("rct")
				for _, l := range rct.Logs {
					println("log")
					println(l.BlockHash.Hex())
					println(l.TxHash.Hex())
					println(l.Address.Hex())
				}
			}
			for _, stateRlp := range payload.StateNodesRlp {
				var acct state.Account
				err = rlp.Decode(bytes.NewBuffer(stateRlp), &acct)
				if err != nil {
					println(err.Error())
					log.Error(err)
				}
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
