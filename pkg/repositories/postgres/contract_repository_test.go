package postgres

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

var _ = Describe("Creating contracts", func() {
	var repository repositories.ContractRepository
	var node core.Node
	BeforeEach(func() {
		node = core.Node{
			GenesisBlock: "GENESIS",
			NetworkId:    1,
			Id:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
			ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
		}
		repository = BuildRepository(node)
	})

	It("returns the contract when it exists", func() {
		repository.CreateContract(core.Contract{Hash: "x123"})

		contract, err := repository.FindContract("x123")
		Expect(err).NotTo(HaveOccurred())
		Expect(contract.Hash).To(Equal("x123"))

		Expect(repository.ContractExists("x123")).To(BeTrue())
		Expect(repository.ContractExists("x456")).To(BeFalse())
	})

	It("returns err if contract does not exist", func() {
		_, err := repository.FindContract("x123")
		Expect(err).To(HaveOccurred())
	})

	It("returns empty array when no transactions 'To' a contract", func() {
		repository.CreateContract(core.Contract{Hash: "x123"})
		contract, err := repository.FindContract("x123")
		Expect(err).ToNot(HaveOccurred())
		Expect(contract.Transactions).To(BeEmpty())
	})

	It("returns transactions 'To' a contract", func() {
		var blockRepository repositories.BlockRepository
		blockRepository = BuildRepository(node)
		block := core.Block{
			Number: 123,
			Transactions: []core.Transaction{
				{Hash: "TRANSACTION1", To: "x123", Value: "0"},
				{Hash: "TRANSACTION2", To: "x345", Value: "0"},
				{Hash: "TRANSACTION3", To: "x123", Value: "0"},
			},
		}
		blockRepository.CreateOrUpdateBlock(block)

		repository.CreateContract(core.Contract{Hash: "x123"})
		contract, err := repository.FindContract("x123")
		Expect(err).ToNot(HaveOccurred())
		sort.Slice(contract.Transactions, func(i, j int) bool {
			return contract.Transactions[i].Hash < contract.Transactions[j].Hash
		})
		Expect(contract.Transactions).To(
			Equal([]core.Transaction{
				{Hash: "TRANSACTION1", To: "x123", Value: "0"},
				{Hash: "TRANSACTION3", To: "x123", Value: "0"},
			}))
	})

	It("stores the ABI of the contract", func() {
		repository.CreateContract(core.Contract{
			Abi:  "{\"some\": \"json\"}",
			Hash: "x123",
		})
		contract, err := repository.FindContract("x123")
		Expect(err).ToNot(HaveOccurred())
		Expect(contract.Abi).To(Equal("{\"some\": \"json\"}"))
	})

	It("updates the ABI of the contract if hash already present", func() {
		repository.CreateContract(core.Contract{
			Abi:  "{\"some\": \"json\"}",
			Hash: "x123",
		})
		repository.CreateContract(core.Contract{
			Abi:  "{\"some\": \"different json\"}",
			Hash: "x123",
		})
		contract, err := repository.FindContract("x123")
		Expect(err).ToNot(HaveOccurred())
		Expect(contract.Abi).To(Equal("{\"some\": \"different json\"}"))
	})
})
