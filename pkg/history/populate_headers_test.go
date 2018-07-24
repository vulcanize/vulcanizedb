package history_test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
)

var _ = Describe("Populating headers", func() {

	var headerRepository *fakes.MockHeaderRepository

	BeforeEach(func() {
		headerRepository = fakes.NewMockHeaderRepository()
	})

	It("returns number of headers added", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(2))
		headerRepository.SetMissingBlockNumbers([]int64{2})

		headersAdded, err := history.PopulateMissingHeaders(blockChain, headerRepository, 1, []transformers.Transformer{})

		Expect(err).NotTo(HaveOccurred())
		Expect(headersAdded).To(Equal(1))
	})

	It("adds missing headers to the db", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(2))
		headerRepository.SetMissingBlockNumbers([]int64{2})

		_, err := history.PopulateMissingHeaders(blockChain, headerRepository, 1, []transformers.Transformer{})

		Expect(err).NotTo(HaveOccurred())
		headerRepository.AssertCreateOrUpdateHeaderCallCountAndPassedBlockNumbers(1, []int64{2})
	})

	It("executes passed transformers with created headers", func() {
		blockNumber := int64(54321)
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(blockNumber))
		headerRepository.SetMissingBlockNumbers([]int64{blockNumber})
		headerID := int64(12345)
		headerRepository.SetCreateOrUpdateHeaderReturnID(headerID)
		transformer := fakes.NewMockTransformer()

		_, err := history.PopulateMissingHeaders(blockChain, headerRepository, 1, []transformers.Transformer{transformer})

		Expect(err).NotTo(HaveOccurred())
		transformer.AssertExecuteCalledWith(core.Header{BlockNumber: blockNumber}, headerID)
	})

	It("returns error if executing transformer fails", func() {
		blockNumber := int64(54321)
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(blockNumber))
		headerRepository.SetMissingBlockNumbers([]int64{blockNumber})
		headerID := int64(12345)
		headerRepository.SetCreateOrUpdateHeaderReturnID(headerID)
		transformer := fakes.NewMockTransformer()
		transformer.SetExecuteErr(fakes.FakeError)

		_, err := history.PopulateMissingHeaders(blockChain, headerRepository, 1, []transformers.Transformer{transformer})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
