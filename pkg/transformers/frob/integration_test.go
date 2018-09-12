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

package frob_test

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Integration tests", func() {
	XIt("Fetches frob event logs from a local test chain", func() {
		ipcPath := test_config.TestClient.IPCPath

		rawRpcClient, err := rpc.Dial(ipcPath)
		Expect(err).NotTo(HaveOccurred())

		rpcClient := client.NewRpcClient(rawRpcClient, ipcPath)
		ethClient := ethclient.NewClient(rawRpcClient)
		blockChainClient := client.NewEthClient(ethClient)
		realNode := node.MakeNode(rpcClient)
		transactionConverter := vRpc.NewRpcTransactionConverter(ethClient)
		realBlockChain := geth.NewBlockChain(blockChainClient, realNode, transactionConverter)
		realFetcher := shared.NewFetcher(realBlockChain)
		topic0 := common.HexToHash(shared.FrobSignature)
		topics := [][]common.Hash{{topic0}}

		result, err := realFetcher.FetchLogs(shared.PitContractAddress, topics, int64(12))
		Expect(err).NotTo(HaveOccurred())

		Expect(len(result) > 0).To(BeTrue())
		Expect(result[0].Address).To(Equal(common.HexToAddress(shared.PitContractAddress)))
		Expect(result[0].TxHash).To(Equal(test_data.EthFrobLog.TxHash))
		Expect(result[0].BlockNumber).To(Equal(test_data.EthFrobLog.BlockNumber))
		Expect(result[0].Topics).To(Equal(test_data.EthFrobLog.Topics))
		Expect(result[0].Index).To(Equal(test_data.EthFrobLog.Index))
	})

	It("unpacks an event log", func() {
		address := common.HexToAddress(shared.PitContractAddress)
		abi, err := geth.ParseAbi(shared.PitABI)
		Expect(err).NotTo(HaveOccurred())

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		entity := &frob.FrobEntity{}

		var eventLog = test_data.EthFrobLog

		err = contract.UnpackLog(entity, "Frob", eventLog)
		Expect(err).NotTo(HaveOccurred())

		expectedEntity := test_data.FrobEntity
		Expect(entity.Art).To(Equal(expectedEntity.Art))
		Expect(entity.IArt).To(Equal(expectedEntity.IArt))
		Expect(entity.Ilk).To(Equal(expectedEntity.Ilk))
		Expect(entity.Ink).To(Equal(expectedEntity.Ink))
		Expect(entity.Urn).To(Equal(expectedEntity.Urn))
	})
})
