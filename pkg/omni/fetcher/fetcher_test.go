// Copyright 2018 Vulcanize
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

package fetcher_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/omni/fetcher"
)

var _ = Describe("Fetcher Test", func() {
	blockNumber := int64(6194634)
	var realFetcher fetcher.Fetcher

	BeforeEach(func() {
		infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
		rawRpcClient, err := rpc.Dial(infuraIPC)
		Expect(err).NotTo(HaveOccurred())
		rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
		ethClient := ethclient.NewClient(rawRpcClient)
		blockChainClient := client.NewEthClient(ethClient)
		node := node.MakeNode(rpcClient)
		transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
		blockChain := geth.NewBlockChain(blockChainClient, node, transactionConverter)
		realFetcher = fetcher.NewFetcher(blockChain)
	})

	Describe("Fetch big.Int test", func() {

		It("fetch totalSupply big.Int", func() {
			bigInt, err := realFetcher.FetchBigInt("totalSupply", constants.DaiAbiString, constants.DaiContractAddress, blockNumber, nil)
			Expect(err).NotTo(HaveOccurred())
			expectedBigInt := big.Int{}
			expectedBigInt.SetString("47327413946297204537985606", 10)
			Expect(bigInt.String()).To(Equal(expectedBigInt.String()))
		})

		It("returns an error if the call to the blockchain fails", func() {
			result, err := realFetcher.FetchBigInt("totalSupply", "", "", 0, nil)

			Expect(result).To(Equal(big.Int{}))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("totalSupply"))
		})
	})

	Describe("Fetch bool test", func() {

		It("fetch stopped boolean", func() {
			boo, err := realFetcher.FetchBool("stopped", constants.DaiAbiString, constants.DaiContractAddress, blockNumber, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(boo).To(Equal(false))
		})

		It("returns an error if the call to the blockchain fails", func() {
			boo, err := realFetcher.FetchBool("stopped", "", "", 0, nil)

			Expect(boo).To(Equal(false))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("stopped"))
		})
	})

	Describe("Fetch address test", func() {

		It("fetch owner address", func() {
			expectedAddr := common.HexToAddress("0x0000000000000000000000000000000000000000")
			addr, err := realFetcher.FetchAddress("owner", constants.DaiAbiString, constants.DaiContractAddress, blockNumber, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(addr).To(Equal(expectedAddr))
		})

		It("returns an error if the call to the blockchain fails", func() {
			addr, err := realFetcher.FetchAddress("owner", "", "", 0, nil)

			Expect(addr).To(Equal(common.Address{}))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("owner"))
		})
	})

	Describe("Fetch string test", func() {

		It("fetch name string", func() {
			expectedStr := "TrueUSD"
			str, err := realFetcher.FetchString("name", constants.TusdAbiString, constants.TusdContractAddress, blockNumber, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(str).To(Equal(expectedStr))
		})

		It("returns an error if the call to the blockchain fails", func() {
			str, err := realFetcher.FetchString("name", "", "", 0, nil)

			Expect(str).To(Equal(""))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name"))
		})
	})

	Describe("Fetch hash test", func() {

		It("fetch name hash", func() {
			expectedHash := common.HexToHash("0x44616920537461626c65636f696e2076312e3000000000000000000000000000")
			hash, err := realFetcher.FetchHash("name", constants.DaiAbiString, constants.DaiContractAddress, blockNumber, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(hash).To(Equal(expectedHash))
		})

		It("returns an error if the call to the blockchain fails", func() {
			hash, err := realFetcher.FetchHash("name", "", "", 0, nil)

			Expect(hash).To(Equal(common.Hash{}))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name"))
		})
	})

})
