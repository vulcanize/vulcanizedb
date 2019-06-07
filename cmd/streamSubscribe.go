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
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/libraries/shared/streamer"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// streamSubscribeCmd represents the streamSubscribe command
var streamSubscribeCmd = &cobra.Command{
	Use:   "streamSubscribe",
	Short: "This command is used to subscribe to the seed node stream with the provided filters",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		streamSubscribe()
	},
}

func init() {
	rootCmd.AddCommand(streamSubscribeCmd)
}

func streamSubscribe() {
	// Prep the subscription config/filters to be sent to the server
	subscriptionConfig()

	// Create a new rpc client and a subscription streamer with that client
	rpcClient := getRpcClient()
	str := streamer.NewSeedStreamer(rpcClient)

	// Buffered channel for reading subscription payloads
	payloadChan := make(chan ipfs.ResponsePayload, 20000)

	// Subscribe to the seed node service with the given config/filter parameters
	sub, err := str.Stream(payloadChan, subConfig)
	if err != nil {
		println(err.Error())
		log.Fatal(err)
	}

	// Receive response payloads and print out the results
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
				println("Header")
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
				println("Trx")
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
				println("Rct")
				for _, l := range rct.Logs {
					println("log")
					println(l.BlockHash.Hex())
					println(l.TxHash.Hex())
					println(l.Address.Hex())
				}
			}
			for key, stateRlp := range payload.StateNodesRlp {
				var acct state.Account
				err = rlp.Decode(bytes.NewBuffer(stateRlp), &acct)
				if err != nil {
					println(err.Error())
					log.Error(err)
				}
				println("State")
				print("key: ")
				println(key.Hex())
				print("root: ")
				println(acct.Root.Hex())
				print("balance: ")
				println(acct.Balance.Int64())
			}
			for stateKey, mappedRlp := range payload.StorageNodesRlp {
				println("Storage")
				print("state key: ")
				println(stateKey.Hex())
				for storageKey, storageRlp := range mappedRlp {
					println("Storage")
					print("key: ")
					println(storageKey.Hex())
					var i []interface{}
					err := rlp.DecodeBytes(storageRlp, i)
					if err != nil {
						println(err.Error())
						log.Error(err)
					}
					print("bytes: ")
					println(storageRlp)
				}
			}
		case err = <-sub.Err():
			println(err.Error())
			log.Fatal(err)
		}
	}
}

func subscriptionConfig() {
	log.Info("loading subscription config")
	vulcPath = viper.GetString("subscription.path")
	subConfig = config.Subscription{
		// Below default to false, which means we do not backfill by default
		BackFill:     viper.GetBool("subscription.backfill"),
		BackFillOnly: viper.GetBool("subscription.backfillOnly"),

		// Below default to 0
		// 0 start means we start at the beginning and 0 end means we continue indefinitely
		StartingBlock: viper.GetInt64("subscription.startingBlock"),
		EndingBlock:   viper.GetInt64("subscription.endingBlock"),

		// Below default to false, which means we get all headers by default
		HeaderFilter: config.HeaderFilter{
			Off:       viper.GetBool("subscription.headerFilter.off"),
			FinalOnly: viper.GetBool("subscription.headerFilter.finalOnly"),
		},

		// Below defaults to false and two slices of length 0
		// Which means we get all transactions by default
		TrxFilter: config.TrxFilter{
			Off: viper.GetBool("subscription.trxFilter.off"),
			Src: viper.GetStringSlice("subscription.trxFilter.src"),
			Dst: viper.GetStringSlice("subscription.trxFilter.dst"),
		},

		// Below defaults to false and one slice of length 0
		// Which means we get all receipts by default
		ReceiptFilter: config.ReceiptFilter{
			Off:     viper.GetBool("subscription.receiptFilter.off"),
			Topic0s: viper.GetStringSlice("subscription.receiptFilter.topic0s"),
		},

		// Below defaults to two false, and a slice of length 0
		// Which means we get all state leafs by default, but no intermediate nodes
		StateFilter: config.StateFilter{
			Off:               viper.GetBool("subscription.stateFilter.off"),
			IntermediateNodes: viper.GetBool("subscription.stateFilter.intermediateNodes"),
			Addresses:         viper.GetStringSlice("subscription.stateFilter.addresses"),
		},

		// Below defaults to two false, and two slices of length 0
		// Which means we get all storage leafs by default, but no intermediate nodes
		StorageFilter: config.StorageFilter{
			Off:               viper.GetBool("subscription.storageFilter.off"),
			IntermediateNodes: viper.GetBool("subscription.storageFilter.intermediateNodes"),
			Addresses:         viper.GetStringSlice("subscription.storageFilter.addresses"),
			StorageKeys:       viper.GetStringSlice("subscription.storageFilter.storageKeys"),
		},
	}
}

func getRpcClient() core.RpcClient {
	rawRpcClient, err := rpc.Dial(vulcPath)
	if err != nil {
		log.Fatal(err)
	}
	return client.NewRpcClient(rawRpcClient, vulcPath)
}
