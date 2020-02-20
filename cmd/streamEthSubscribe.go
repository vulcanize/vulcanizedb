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
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/libraries/shared/streamer"
	"github.com/vulcanize/vulcanizedb/pkg/eth/client"
	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/super_node"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
)

// streamEthSubscriptionCmd represents the streamEthSubscription command
var streamEthSubscriptionCmd = &cobra.Command{
	Use:   "streamEthSubscription",
	Short: "This command is used to subscribe to the super node eth stream with the provided filters",
	Long: `This command is for demo and testing purposes and is used to subscribe to the super node with the provided subscription configuration parameters.
It does not do anything with the data streamed from the super node other than unpack it and print it out for demonstration purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		subCommand = cmd.CalledAs()
		logWithCommand = *log.WithField("SubCommand", subCommand)
		streamEthSubscription()
	},
}

func init() {
	rootCmd.AddCommand(streamEthSubscriptionCmd)
}

func streamEthSubscription() {
	// Prep the subscription config/filters to be sent to the server
	ethSubConfig, err := eth.NewEthSubscriptionConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new rpc client and a subscription streamer with that client
	rpcClient := getRPCClient()
	str := streamer.NewSuperNodeStreamer(rpcClient)

	// Buffered channel for reading subscription payloads
	payloadChan := make(chan super_node.SubscriptionPayload, 20000)

	// Subscribe to the super node service with the given config/filter parameters
	sub, err := str.Stream(payloadChan, ethSubConfig)
	if err != nil {
		logWithCommand.Fatal(err)
	}
	logWithCommand.Info("awaiting payloads")
	// Receive response payloads and print out the results
	for {
		select {
		case payload := <-payloadChan:
			if payload.Err != "" {
				logWithCommand.Error(payload.Err)
				continue
			}
			var ethData eth.IPLDs
			if err := rlp.DecodeBytes(payload.Data, &ethData); err != nil {
				logWithCommand.Error(err)
			}
			for _, headerRlp := range ethData.Headers {
				var header types.Header
				err = rlp.Decode(bytes.NewBuffer(headerRlp.Data), &header)
				if err != nil {
					logWithCommand.Error(err)
					continue
				}
				fmt.Printf("Header number %d, hash %s\n", header.Number.Int64(), header.Hash().Hex())
				fmt.Printf("header: %v\n", header)
			}
			for _, trxRlp := range ethData.Transactions {
				var trx types.Transaction
				buff := bytes.NewBuffer(trxRlp.Data)
				stream := rlp.NewStream(buff, 0)
				err := trx.DecodeRLP(stream)
				if err != nil {
					logWithCommand.Error(err)
					continue
				}
				fmt.Printf("Transaction with hash %s\n", trx.Hash().Hex())
				fmt.Printf("trx: %v\n", trx)
			}
			for _, rctRlp := range ethData.Receipts {
				var rct types.ReceiptForStorage
				buff := bytes.NewBuffer(rctRlp.Data)
				stream := rlp.NewStream(buff, 0)
				err = rct.DecodeRLP(stream)
				if err != nil {
					logWithCommand.Error(err)
					continue
				}
				fmt.Printf("Receipt with block hash %s, trx hash %s\n", rct.BlockHash.Hex(), rct.TxHash.Hex())
				fmt.Printf("rct: %v\n", rct)
				for _, l := range rct.Logs {
					if len(l.Topics) < 1 {
						logWithCommand.Error(fmt.Sprintf("log only has %d topics", len(l.Topics)))
						continue
					}
					fmt.Printf("Log for block hash %s, trx hash %s, address %s, and with topic0 %s\n",
						l.BlockHash.Hex(), l.TxHash.Hex(), l.Address.Hex(), l.Topics[0].Hex())
					fmt.Printf("log: %v\n", l)
				}
			}
			// This assumes leafs only
			for _, stateNode := range ethData.StateNodes {
				var acct state.Account
				err = rlp.Decode(bytes.NewBuffer(stateNode.IPLD.Data), &acct)
				if err != nil {
					logWithCommand.Error(err)
					continue
				}
				fmt.Printf("Account for key %s, and root %s, with balance %d\n",
					stateNode.StateTrieKey.Hex(), acct.Root.Hex(), acct.Balance.Int64())
				fmt.Printf("state account: %v\n", acct)
			}
			for _, storageNode := range ethData.StorageNodes {
				fmt.Printf("Storage for state key %s ", storageNode.StateTrieKey.Hex())
				fmt.Printf("with storage key %s\n", storageNode.StorageTrieKey.Hex())
				var i []interface{}
				err := rlp.DecodeBytes(storageNode.IPLD.Data, &i)
				if err != nil {
					logWithCommand.Error(err)
					continue
				}
				// if a value node
				if len(i) == 1 {
					valueBytes, ok := i[0].([]byte)
					if !ok {
						continue
					}
					fmt.Printf("Storage leaf key: %s, and value hash: %s\n",
						storageNode.StorageTrieKey.Hex(), common.BytesToHash(valueBytes).Hex())
				}
			}
		case err = <-sub.Err():
			logWithCommand.Fatal(err)
		}
	}
}

func getRPCClient() core.RPCClient {
	vulcPath := viper.GetString("superNode.ethSubscription.path")
	if vulcPath == "" {
		vulcPath = "ws://127.0.0.1:8080" // default to and try the default ws url if no path is provided
	}
	rawRPCClient, err := rpc.Dial(vulcPath)
	if err != nil {
		logWithCommand.Fatal(err)
	}
	return client.NewRPCClient(rawRPCClient, vulcPath)
}
