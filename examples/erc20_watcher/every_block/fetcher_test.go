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

package every_block_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/examples/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"math/big"
)

var _ = Describe("ERC20 Fetcher", func() {
	blockNumber := int64(5502914)

	infuraIPC := "https://mainnet.infura.io/J5Vd2fRtGsw0zZ0Ov3BL"
	realBlockchain := geth.NewBlockchain(infuraIPC)
	realFetcher := every_block.NewFetcher(realBlockchain)

	fakeBlockchain := &mocks.Blockchain{}
	testFetcher := every_block.NewFetcher(fakeBlockchain)
	testAbi := "testAbi"
	testContractAddress := "testContractAddress"

	errorBlockchain := &mocks.FailureBlockchain{}
	errorFetcher := every_block.NewFetcher(errorBlockchain)

	Describe("FetchSupplyOf", func() {
		It("fetches data from the blockchain with the correct arguments", func() {
			_, err := testFetcher.FetchSupplyOf(testAbi, testContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeBlockchain.FetchedAbi).To(Equal(testAbi))
			Expect(fakeBlockchain.FetchedContractAddress).To(Equal(testContractAddress))
			Expect(fakeBlockchain.FetchedMethod).To(Equal("totalSupply"))
			Expect(fakeBlockchain.FetchedMethodArg).To(BeNil())
			expectedResult := big.Int{}
			expected := &expectedResult
			Expect(fakeBlockchain.FetchedResult).To(Equal(&expected))
			Expect(fakeBlockchain.FetchedBlockNumber).To(Equal(blockNumber))
		})

		It("fetches a token's total supply at the given block height", func() {
			result, err := realFetcher.FetchSupplyOf(constants.DaiAbiString, constants.DaiContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			expectedResult := big.Int{}
			expectedResult.SetString("27647235749155415536952630", 10)
			Expect(result).To(Equal(expectedResult))
		})

		It("returns an error if the call to the blockchain fails", func() {
			result, err := errorFetcher.FetchSupplyOf("", "", 0)

			Expect(result.String()).To(Equal("0"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("totalSupply"))
			Expect(err.Error()).To(ContainSubstring(mocks.TestError.Error()))
		})
	})
})
