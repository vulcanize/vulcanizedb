package history_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"math/big"
)

var _ = Describe("Validation window", func() {
	It("creates a ValidationWindow equal to (HEAD-windowSize, HEAD)", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(5))

		validationWindow := history.MakeValidationWindow(blockChain, 2)

		Expect(validationWindow.LowerBound).To(Equal(int64(3)))
		Expect(validationWindow.UpperBound).To(Equal(int64(5)))
	})

	It("returns the window size", func() {
		window := history.ValidationWindow{LowerBound: 1, UpperBound: 3}

		Expect(window.Size()).To(Equal(2))
	})

	It("generates a range of int64s", func() {
		numberOfBlocksCreated := history.MakeRange(0, 5)
		expected := []int64{0, 1, 2, 3, 4, 5}

		Expect(numberOfBlocksCreated).To(Equal(expected))
	})
})
