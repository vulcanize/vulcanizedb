package watched_contracts_test

import (
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
	"github.com/8thlight/vulcanizedb/pkg/watched_contracts"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The watched contract summary", func() {

	Context("when the given contract is not being watched", func() {
		It("returns an error", func() {
			repository := repositories.NewInMemory()

			contractSummary, err := watched_contracts.NewSummary(repository, "123")

			Expect(contractSummary).To(BeNil())
			Expect(err).NotTo(BeNil())
		})
	})

	Context("when the given contract is being watched", func() {
		It("returns the summary", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)

			contractSummary, err := watched_contracts.NewSummary(repository, "0x123")

			Expect(contractSummary).NotTo(BeNil())
			Expect(err).To(BeNil())
		})

		It("includes the contract hash in the summary", func() {
			repository := repositories.NewInMemory()
			watchedContract := core.WatchedContract{Hash: "0x123"}
			repository.CreateWatchedContract(watchedContract)

			contractSummary, _ := watched_contracts.NewSummary(repository, "0x123")

			Expect(contractSummary.ContractHash).To(Equal("0x123"))
		})
	})

})
