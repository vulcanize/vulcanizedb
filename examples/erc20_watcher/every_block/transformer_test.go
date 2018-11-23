// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package every_block_test

import (
	"math/big"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/examples/mocks"
	"github.com/vulcanize/vulcanizedb/examples/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
)

var testContractConfig = shared.ContractConfig{
	Address:    constants.DaiContractAddress,
	Abi:        constants.DaiAbiString,
	FirstBlock: int64(4752008),
	LastBlock:  int64(5750050),
	Name:       "Dai",
}

var config = testContractConfig

var _ = Describe("Everyblock transformer", func() {
	var getter mocks.Getter
	var repository mocks.ERC20TokenRepository
	var transformer every_block.ERC20Transformer
	var blockChain *fakes.MockBlockChain
	var initialSupply = "27647235749155415536952630"
	var initialSupplyPlusOne = "27647235749155415536952631"
	var initialSupplyPlusTwo = "27647235749155415536952632"
	var initialSupplyPlusThree = "27647235749155415536952633"
	var defaultLastBlock = big.Int{}

	BeforeEach(func() {
		blockChain = fakes.NewMockBlockChain()
		blockChain.SetLastBlock(&defaultLastBlock)
		getter = mocks.NewGetter(blockChain)
		getter.Fetcher.SetSupply(initialSupply)
		repository = mocks.ERC20TokenRepository{}
		repository.SetMissingSupplyBlocks([]int64{config.FirstBlock})
		db := test_helpers.CreateNewDatabase()
		rt := erc20_watcher.NewTokenHolderRetriever(db, config.Address)
		//setting the mock repository to return the first block as the missing blocks

		transformer = every_block.ERC20Transformer{
			Getter:     &getter,
			Repository: &repository,
			Retriever:  rt,
			Config:     config,
		}
		transformer.SetConfiguration(config)
	})

	It("fetches and persists the total supply of a token for a single block", func() {
		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		Expect(len(getter.Fetcher.FetchedBlocks)).To(Equal(1))
		Expect(getter.Fetcher.FetchedBlocks).To(ConsistOf(config.FirstBlock))
		Expect(getter.Fetcher.Abi).To(Equal(config.Abi))
		Expect(getter.Fetcher.ContractAddress).To(Equal(config.Address))

		Expect(repository.StartingBlock).To(Equal(config.FirstBlock))
		Expect(repository.EndingBlock).To(Equal(config.LastBlock))
		Expect(len(repository.TotalSuppliesCreated)).To(Equal(1))
		expectedSupply := big.Int{}
		expectedSupply.SetString(initialSupply, 10)
		expectedSupply.Add(&expectedSupply, big.NewInt(1))

		Expect(repository.TotalSuppliesCreated[0].Value).To(Equal(expectedSupply.String()))
	})

	It("retrieves the total supply for every missing block", func() {
		missingBlocks := []int64{
			config.FirstBlock,
			config.FirstBlock + 1,
			config.FirstBlock + 2,
		}
		repository.SetMissingSupplyBlocks(missingBlocks)
		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		Expect(len(getter.Fetcher.FetchedBlocks)).To(Equal(3))
		Expect(getter.Fetcher.FetchedBlocks).To(ConsistOf(config.FirstBlock, config.FirstBlock+1, config.FirstBlock+2))
		Expect(getter.Fetcher.Abi).To(Equal(config.Abi))
		Expect(getter.Fetcher.ContractAddress).To(Equal(config.Address))

		Expect(len(repository.TotalSuppliesCreated)).To(Equal(3))
		Expect(repository.TotalSuppliesCreated[0].Value).To(Equal(initialSupplyPlusOne))
		Expect(repository.TotalSuppliesCreated[1].Value).To(Equal(initialSupplyPlusTwo))
		Expect(repository.TotalSuppliesCreated[2].Value).To(Equal(initialSupplyPlusThree))
	})

	It("uses the set contract configuration", func() {
		repository.SetMissingSupplyBlocks([]int64{testContractConfig.FirstBlock})
		transformer.SetConfiguration(testContractConfig)
		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		Expect(getter.Fetcher.FetchedBlocks).To(ConsistOf(testContractConfig.FirstBlock))
		Expect(getter.Fetcher.Abi).To(Equal(testContractConfig.Abi))
		Expect(getter.Fetcher.ContractAddress).To(Equal(testContractConfig.Address))

		Expect(repository.StartingBlock).To(Equal(testContractConfig.FirstBlock))
		Expect(repository.EndingBlock).To(Equal(testContractConfig.LastBlock))
		Expect(len(repository.TotalSuppliesCreated)).To(Equal(1))
		expectedTokenSupply := every_block.TokenSupply{
			Value:        initialSupplyPlusOne,
			TokenAddress: testContractConfig.Address,
			BlockNumber:  testContractConfig.FirstBlock,
		}
		Expect(repository.TotalSuppliesCreated[0]).To(Equal(expectedTokenSupply))
	})

	It("uses the most recent block if the Config.LastBlock is -1", func() {
		testContractConfig.LastBlock = -1
		transformer.SetConfiguration(testContractConfig)

		randomBlockNumber := rand.Int63()
		numberToString := strconv.FormatInt(randomBlockNumber, 10)
		mostRecentBlock := big.Int{}
		mostRecentBlock.SetString(numberToString, 10)

		blockChain.SetLastBlock(&mostRecentBlock)

		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		Expect(repository.EndingBlock).To(Equal(randomBlockNumber))
	})

	It("returns an error if the call to get missing blocks fails", func() {
		failureRepository := mocks.FailureRepository{}
		failureRepository.SetMissingSupplyBlocksFail(true)
		transformer = every_block.ERC20Transformer{
			Getter:     &getter,
			Repository: &failureRepository,
		}
		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
		Expect(err.Error()).To(ContainSubstring("fetching missing blocks"))
	})

	It("returns an error if the call to the blockChain fails", func() {
		failureBlockchain := fakes.NewMockBlockChain()
		failureBlockchain.SetLastBlock(&defaultLastBlock)
		failureBlockchain.SetFetchContractDataErr(fakes.FakeError)
		getter := every_block.NewGetter(failureBlockchain)
		transformer = every_block.ERC20Transformer{
			Getter:     &getter,
			Repository: &repository,
		}
		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
		Expect(err.Error()).To(ContainSubstring("supply"))
	})

	It("returns an error if the call to save the token_supply fails", func() {
		failureRepository := mocks.FailureRepository{}
		failureRepository.SetMissingSupplyBlocks([]int64{config.FirstBlock})
		failureRepository.SetCreateSupplyFail(true)

		transformer = every_block.ERC20Transformer{
			Getter:     &getter,
			Repository: &failureRepository,
		}
		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
		Expect(err.Error()).To(ContainSubstring("supply"))
	})
})
