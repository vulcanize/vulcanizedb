// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tend_test

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

// These test are pending either being able to emit a Tend event on a Ganache test chain or until the contracts are deployed to Kovan.
var _ = XDescribe("Integration tests", func() {
	It("Fetches Tend event logs from a local test chain", func() {
		ipcPath := test_config.TestClient.IPCPath

		rawRpcClient, err := rpc.Dial(ipcPath)
		Expect(err).NotTo(HaveOccurred())

		rpcClient := client.NewRpcClient(rawRpcClient, ipcPath)
		ethClient := ethclient.NewClient(rawRpcClient)
		blockChainClient := client.NewEthClient(ethClient)
		realNode := node.MakeNode(rpcClient)
		transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
		realBlockChain := geth.NewBlockChain(blockChainClient, realNode, transactionConverter)
		realFetcher := shared.NewFetcher(realBlockChain)
		topic0 := common.HexToHash(shared.TendSignature)
		topics := [][]common.Hash{{topic0}}

		result, err := realFetcher.FetchLogs(test_data.FlipAddress, topics, test_data.FlipKickBlockNumber)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(result) > 0).To(BeTrue())
		Expect(result[0].Address).To(Equal(test_data.EthFlipKickLog.Address))
		Expect(result[0].TxHash).To(Equal(test_data.EthFlipKickLog.TxHash))
		Expect(result[0].BlockNumber).To(Equal(test_data.EthFlipKickLog.BlockNumber))
		Expect(result[0].Topics).To(Equal(test_data.EthFlipKickLog.Topics))
		Expect(result[0].Index).To(Equal(test_data.EthFlipKickLog.Index))
	})

	It("unpacks an event log", func() {
		address := common.HexToAddress(test_data.FlipAddress)
		abi, err := geth.ParseAbi(shared.FlipperABI)
		Expect(err).NotTo(HaveOccurred())

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		entity := tend.TendEntity{}

		var eventLog = test_data.TendLog

		err = contract.UnpackLog(&entity, "Tend", eventLog)
		Expect(err).NotTo(HaveOccurred())

		expectedEntity := test_data.TendEntity
		Expect(entity.Id).To(Equal(expectedEntity.Id))
		Expect(entity.Lot).To(Equal(expectedEntity.Lot))
		Expect(entity.Bid).To(Equal(expectedEntity.Bid))
		Expect(entity.Guy).To(Equal(expectedEntity.Guy))
		Expect(entity.Tic).To(Equal(expectedEntity.Tic))
		Expect(entity.Era).To(Equal(expectedEntity.Era))
		Expect(entity.Raw).To(Equal(expectedEntity.Raw))
	})
})
