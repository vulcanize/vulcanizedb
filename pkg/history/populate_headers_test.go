package history_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/inmemory"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
)

var _ = Describe("Populating headers", func() {

	var inMemory *inmemory.InMemory
	var headerRepository *inmemory.HeaderRepository

	BeforeEach(func() {
		inMemory = inmemory.NewInMemory()
		headerRepository = inmemory.NewHeaderRepository(inMemory)
	})

	Describe("When 1 missing header", func() {

		It("returns number of headers added", func() {
			headers := []core.Header{
				{BlockNumber: 1},
				{BlockNumber: 2},
			}
			blockChain := fakes.NewMockBlockChainWithHeaders(headers)
			headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: 2})

			headersAdded := history.PopulateMissingHeaders(blockChain, headerRepository, 1)

			Expect(headersAdded).To(Equal(1))
		})
	})

	It("adds missing headers to the db", func() {
		headers := []core.Header{
			{BlockNumber: 1},
			{BlockNumber: 2},
		}
		blockChain := fakes.NewMockBlockChainWithHeaders(headers)
		dbHeader, _ := headerRepository.GetHeader(1)
		Expect(dbHeader.BlockNumber).To(BeZero())

		history.PopulateMissingHeaders(blockChain, headerRepository, 1)

		dbHeader, _ = headerRepository.GetHeader(1)
		Expect(dbHeader.BlockNumber).To(Equal(int64(1)))
	})
})
