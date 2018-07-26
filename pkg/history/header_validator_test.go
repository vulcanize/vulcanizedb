package history_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
	"math/big"
)

var _ = Describe("Header validator", func() {
	It("attempts to create every header in the validation window", func() {
		headerRepository := fakes.NewMockHeaderRepository()
		headerRepository.SetMissingBlockNumbers([]int64{})
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(3))
		validator := history.NewHeaderValidator(blockChain, headerRepository, 2, []transformers.Transformer{})

		validator.ValidateHeaders()

		headerRepository.AssertCreateOrUpdateHeaderCallCountAndPassedBlockNumbers(3, []int64{1, 2, 3})
	})

	It("passes transformers for execution on new blocks", func() {
		headerRepository := fakes.NewMockHeaderRepository()
		headerRepository.SetMissingBlockNumbers([]int64{})
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(3))
		transformer := fakes.NewMockTransformer()
		validator := history.NewHeaderValidator(blockChain, headerRepository, 1, []transformers.Transformer{transformer})

		validator.ValidateHeaders()

		headerRepository.AssertCreateOrUpdateHeaderCallCountAndPassedBlockNumbers(2, []int64{2, 3})
		transformer.AssertExecuteCalledWith(core.Header{BlockNumber: 3}, 0)
	})
})
