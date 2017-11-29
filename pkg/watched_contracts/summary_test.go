package watched_contracts_test

import (
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/fakes"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
	"github.com/8thlight/vulcanizedb/pkg/watched_contracts"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ bool = Describe("The watched contract summary", func() {

	Context("when the given contract is not being watched", func() {
		It("returns an error", func() {
			repository := repositories.NewInMemory()
			blockchain := fakes.NewBlockchain()

			contractSummary, err := watched_contracts.NewSummary(blockchain, repository, "123")

			Expect(contractSummary).To(BeNil())
			Expect(err).NotTo(BeNil())
		})
	})

	Context("when the given contract is being watched", func() {
		It("returns the summary", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)
			blockchain := fakes.NewBlockchain()

			contractSummary, err := watched_contracts.NewSummary(blockchain, repository, "0x123")

			Expect(contractSummary).NotTo(BeNil())
			Expect(err).To(BeNil())
		})

		It("includes the contract hash in the summary", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)
			blockchain := fakes.NewBlockchain()

			contractSummary, _ := watched_contracts.NewSummary(blockchain, repository, "0x123")

			Expect(contractSummary.ContractHash).To(Equal("0x123"))
		})

		It("sets the number of transactions", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)
			block := core.Block{
				Transactions: []core.Transaction{
					{To: "0x123"},
					{To: "0x123"},
				},
			}
			repository.CreateBlock(block)
			blockchain := fakes.NewBlockchain()

			contractSummary, _ := watched_contracts.NewSummary(blockchain, repository, "0x123")

			Expect(contractSummary.NumberOfTransactions).To(Equal(2))
		})

		It("sets the last transaction", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)
			block := core.Block{
				Transactions: []core.Transaction{
					{Hash: "TRANSACTION2", To: "0x123"},
					{Hash: "TRANSACTION1", To: "0x123"},
				},
			}
			repository.CreateBlock(block)
			blockchain := fakes.NewBlockchain()

			contractSummary, _ := watched_contracts.NewSummary(blockchain, repository, "0x123")

			Expect(contractSummary.LastTransaction.Hash).To(Equal("TRANSACTION2"))
		})

		It("gets contract state attribute for the contract from the blockchain", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)
			blockchain := fakes.NewBlockchain()
			blockchain.SetContractStateAttribute("0x123", "foo", "bar")

			contractSummary, _ := watched_contracts.NewSummary(blockchain, repository, "0x123")
			attribute := contractSummary.GetStateAttribute("foo")

			Expect(attribute).To(Equal("bar"))
		})

		It("gets attributes for the contract from the blockchain", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)
			blockchain := fakes.NewBlockchain()
			blockchain.SetContractStateAttribute("0x123", "foo", "bar")
			blockchain.SetContractStateAttribute("0x123", "baz", "bar")

			contractSummary, _ := watched_contracts.NewSummary(blockchain, repository, "0x123")

			Expect(contractSummary.Attributes).To(Equal(
				core.ContractAttributes{
					{Name: "baz", Type: "string"},
					{Name: "foo", Type: "string"},
				},
			))
		})
	})

})
