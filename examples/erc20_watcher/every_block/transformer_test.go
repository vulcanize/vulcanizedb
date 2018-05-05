package every_block_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/examples/mocks"
	"math/big"
)

//allow for setting configuration OR using a default config?
//handle errors

var config = erc20_watcher.DaiConfig

var _ = Describe("Everyblock transformer", func() {
	var fetcher mocks.Fetcher
	var repository mocks.TotalSupplyRepository
	var transformer every_block.Transformer
	var initialSupply = "27647235749155415536952630"
	var initialSupplyPlusOne = "27647235749155415536952631"
	var initialSupplyPlusTwo = "27647235749155415536952632"
	var initialSupplyPlusThree = "27647235749155415536952633"

	var testContractConfig = erc20_watcher.ContractConfig{
		Address:    "testAddress",
		Abi:        "testAbi",
		FirstBlock: 111,
		LastBlock:  112,
		Name:       "A test contract",
	}

	BeforeEach(func() {
		fetcher = mocks.Fetcher{}
		fetcher.SetSupply(initialSupply)
		repository = mocks.TotalSupplyRepository{}
		repository.SetMissingBlocks([]int64{config.FirstBlock})
		//setting the mock repository to return the first block as the missing blocks

		transformer = every_block.Transformer{
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		transformer.SetConfiguration(erc20_watcher.DaiConfig)
	})

	It("fetches and persists the total supply of a token for a single block", func() {
		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		Expect(len(fetcher.FetchedBlocks)).To(Equal(1))
		Expect(fetcher.FetchedBlocks).To(ConsistOf(config.FirstBlock))
		Expect(fetcher.Abi).To(Equal(config.Abi))
		Expect(fetcher.ContractAddress).To(Equal(config.Address))

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
		repository.SetMissingBlocks(missingBlocks)
		transformer.Execute()

		Expect(len(fetcher.FetchedBlocks)).To(Equal(3))
		Expect(fetcher.FetchedBlocks).To(ConsistOf(config.FirstBlock, config.FirstBlock+1, config.FirstBlock+2))
		Expect(fetcher.Abi).To(Equal(config.Abi))
		Expect(fetcher.ContractAddress).To(Equal(config.Address))

		Expect(len(repository.TotalSuppliesCreated)).To(Equal(3))
		Expect(repository.TotalSuppliesCreated[0].Value).To(Equal(initialSupplyPlusOne))
		Expect(repository.TotalSuppliesCreated[1].Value).To(Equal(initialSupplyPlusTwo))
		Expect(repository.TotalSuppliesCreated[2].Value).To(Equal(initialSupplyPlusThree))
	})

	It("uses the set contract configuration", func() {
		repository.SetMissingBlocks([]int64{testContractConfig.FirstBlock})
		transformer.SetConfiguration(testContractConfig)
		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		Expect(fetcher.FetchedBlocks).To(ConsistOf(testContractConfig.FirstBlock))
		Expect(fetcher.Abi).To(Equal(testContractConfig.Abi))
		Expect(fetcher.ContractAddress).To(Equal(testContractConfig.Address))

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

	It("returns an error if the call to get missing blocks fails", func() {
		failureRepository := mocks.FailureRepository{}
		failureRepository.SetMissingBlocksFail(true)
		transformer = every_block.Transformer{
			Fetcher:    &fetcher,
			Repository: &failureRepository,
		}
		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("TestError"))
		Expect(err.Error()).To(ContainSubstring("fetching missing blocks"))
	})

	It("returns an error if the call to the blockchain fails", func() {
		fetcher := every_block.NewFetcher(mocks.FailureBlockchain{})
		transformer = every_block.Transformer{
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("TestError"))
		Expect(err.Error()).To(ContainSubstring("supply"))
	})

	It("returns an error if the call to save the token_supply fails", func() {
		failureRepository := mocks.FailureRepository{}
		failureRepository.SetMissingBlocks([]int64{config.FirstBlock})
		failureRepository.SetCreateFail(true)

		transformer = every_block.Transformer{
			Fetcher:    &fetcher,
			Repository: &failureRepository,
		}
		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("TestError"))
		Expect(err.Error()).To(ContainSubstring("supply"))
	})
})
