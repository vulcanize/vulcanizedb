package history_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"math/big"
)

var _ = Describe("Header validator", func() {
	var (
		headerRepository *fakes.MockHeaderRepository
		blockChain       *fakes.MockBlockChain
	)

	BeforeEach(func() {
		headerRepository = fakes.NewMockHeaderRepository()
		blockChain = fakes.NewMockBlockChain()
	})

	It("attempts to create every header in the validation window", func() {
		headerRepository.SetMissingBlockNumbers([]int64{})
		blockChain.SetLastBlock(big.NewInt(3))
		validator := history.NewHeaderValidator(blockChain, headerRepository, 2)

		_, err := validator.ValidateHeaders()
		Expect(err).NotTo(HaveOccurred())

		headerRepository.AssertCreateOrUpdateHeaderCallCountAndPassedBlockNumbers(3, []int64{1, 2, 3})
	})

	It("propagates header repository errors", func() {
		blockChain.SetLastBlock(big.NewInt(3))
		headerRepositoryError := errors.New("CreateOrUpdate")
		headerRepository.SetCreateOrUpdateHeaderReturnErr(headerRepositoryError)
		validator := history.NewHeaderValidator(blockChain, headerRepository, 2)

		_, err := validator.ValidateHeaders()
		Expect(err).To(MatchError(headerRepositoryError))
	})
})
