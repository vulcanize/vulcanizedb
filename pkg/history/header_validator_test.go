package history_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"math/big"
)

var _ = Describe("Header validator", func() {
	It("attempts to create every header in the validation window", func() {
		headerRepository := fakes.NewMockHeaderRepository()
		headerRepository.SetMissingBlockNumbers([]int64{})
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(3))
		validator := history.NewHeaderValidator(blockChain, headerRepository, 2)

		validator.ValidateHeaders()

		headerRepository.AssertCreateOrUpdateHeaderCallCountAndPassedBlockNumbers(3, []int64{1, 2, 3})
	})
})
