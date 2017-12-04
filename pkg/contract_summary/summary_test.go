package contract_summary_test

import (
	"math/big"

	"github.com/8thlight/vulcanizedb/pkg/contract_summary"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/fakes"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func NewCurrentContractSummary(blockchain core.Blockchain, repository repositories.Repository, contractHash string) (contract_summary.ContractSummary, error) {
	return contract_summary.NewSummary(blockchain, repository, contractHash, nil)
}

var _ = Describe("The watched contract summary", func() {

	Context("when the given contract is not being watched", func() {
		It("returns an error", func() {
			repository := repositories.NewInMemory()
			blockchain := fakes.NewBlockchain()

			contractSummary, err := NewCurrentContractSummary(blockchain, repository, "0x123")

			Expect(contractSummary).To(Equal(contract_summary.ContractSummary{}))
			Expect(err).NotTo(BeNil())
		})
	})

	Context("when the given contract is being watched", func() {
		It("returns the summary", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)
			blockchain := fakes.NewBlockchain()

			contractSummary, err := NewCurrentContractSummary(blockchain, repository, "0x123")

			Expect(contractSummary).NotTo(Equal(contract_summary.ContractSummary{}))
			Expect(err).To(BeNil())
		})

		It("includes the contract hash in the summary", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)
			blockchain := fakes.NewBlockchain()

			contractSummary, _ := NewCurrentContractSummary(blockchain, repository, "0x123")

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

			contractSummary, _ := NewCurrentContractSummary(blockchain, repository, "0x123")

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

			contractSummary, _ := NewCurrentContractSummary(blockchain, repository, "0x123")

			Expect(contractSummary.LastTransaction.Hash).To(Equal("TRANSACTION2"))
		})

		It("gets contract state attribute for the contract from the blockchain", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)
			blockchain := fakes.NewBlockchain()
			blockchain.SetContractStateAttribute("0x123", nil, "foo", "bar")

			contractSummary, _ := NewCurrentContractSummary(blockchain, repository, "0x123")
			attribute := contractSummary.GetStateAttribute("foo")

			Expect(attribute).To(Equal("bar"))
		})

		It("gets contract state attribute for the contract from the blockchain at specific block height", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)
			blockchain := fakes.NewBlockchain()
			blockNumber := big.NewInt(1000)
			blockchain.SetContractStateAttribute("0x123", nil, "foo", "bar")
			blockchain.SetContractStateAttribute("0x123", blockNumber, "foo", "baz")

			contractSummary, _ := contract_summary.NewSummary(blockchain, repository, "0x123", blockNumber)
			attribute := contractSummary.GetStateAttribute("foo")

			Expect(attribute).To(Equal("baz"))
		})

		It("gets attributes for the contract from the blockchain", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)
			blockchain := fakes.NewBlockchain()
			blockchain.SetContractStateAttribute("0x123", nil, "foo", "bar")
			blockchain.SetContractStateAttribute("0x123", nil, "baz", "bar")

			contractSummary, _ := NewCurrentContractSummary(blockchain, repository, "0x123")

			Expect(contractSummary.Attributes).To(Equal(
				core.ContractAttributes{
					{Name: "baz", Type: "string"},
					{Name: "foo", Type: "string"},
				},
			))
		})
	})

})
