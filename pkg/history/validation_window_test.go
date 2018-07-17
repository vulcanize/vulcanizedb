package history_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
)

var _ = Describe("", func() {
	It("creates a ValidationWindow equal to (HEAD-windowSize, HEAD)", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 1},
			{Number: 2},
			{Number: 3},
			{Number: 4},
			{Number: 5},
		})

		validationWindow := history.MakeValidationWindow(blockchain, 2)

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

	It("logs window message", func() {
		expectedMessage := &bytes.Buffer{}
		window := history.ValidationWindow{LowerBound: 5, UpperBound: 7}
		history.ParsedWindowTemplate.Execute(expectedMessage, window)
		actualMessage := &bytes.Buffer{}

		window.Log(actualMessage)

		Expect(actualMessage).To(Equal(expectedMessage))
	})
})
