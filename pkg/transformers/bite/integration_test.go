/*
 *  Copyright 2018 Vulcanize
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package bite_test

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/bite"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Integration tests", func() {
	XIt("Fetches bite event logs from a local test chain", func() {
		ipcPath := test_config.TestClient.IPCPath

		rawRpcClient, err := rpc.Dial(ipcPath)
		Expect(err).NotTo(HaveOccurred())

		rpcClient := client.NewRpcClient(rawRpcClient, ipcPath)
		ethClient := ethclient.NewClient(rawRpcClient)
		blockChainClient := client.NewEthClient(ethClient)
		realNode := node.MakeNode(rpcClient)
		transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
		realBlockChain := geth.NewBlockChain(blockChainClient, rpcClient, realNode, transactionConverter)
		realFetcher := shared.NewFetcher(realBlockChain)
		topic0 := common.HexToHash(shared.BiteSignature)
		topics := [][]common.Hash{{topic0}}

		result, err := realFetcher.FetchLogs(bite.BiteConfig.ContractAddresses, topics, int64(26))
		Expect(err).NotTo(HaveOccurred())

		Expect(len(result) > 0).To(BeTrue())
		Expect(result[0].Address).To(Equal(common.HexToAddress(shared.CatContractAddress)))
		Expect(result[0].TxHash).To(Equal(test_data.EthBiteLog.TxHash))
		Expect(result[0].BlockNumber).To(Equal(test_data.EthBiteLog.BlockNumber))
		Expect(result[0].Topics).To(Equal(test_data.EthBiteLog.Topics))
		Expect(result[0].Index).To(Equal(test_data.EthBiteLog.Index))
	})

	It("unpacks an event log", func() {
		address := common.HexToAddress(shared.CatContractAddress)
		abi, err := geth.ParseAbi(shared.CatABI)
		Expect(err).NotTo(HaveOccurred())

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		entity := &bite.BiteEntity{}

		var eventLog = test_data.EthBiteLog

		err = contract.UnpackLog(entity, "Bite", eventLog)
		Expect(err).NotTo(HaveOccurred())

		expectedEntity := test_data.BiteEntity
		Expect(entity.Art).To(Equal(expectedEntity.Art))
		Expect(entity.Ilk).To(Equal(expectedEntity.Ilk))
		Expect(entity.Ink).To(Equal(expectedEntity.Ink))
	})
})
