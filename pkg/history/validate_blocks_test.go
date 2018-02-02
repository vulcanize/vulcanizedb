package history_test

import (
	"bytes"

	"io/ioutil"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"github.com/vulcanize/vulcanizedb/pkg/repositories/inmemory"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

var _ = Describe("Blocks validator", func() {

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
		window := history.ValidationWindow{1, 3}
		Expect(window.Size()).To(Equal(2))
	})

	It("calls create or update for all blocks within the window", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 4},
			{Number: 5},
			{Number: 6},
			{Number: 7},
		})
		repository := inmemory.NewInMemory()

		validator := history.NewBlockValidator(blockchain, repository, 2)
		window := validator.ValidateBlocks()
		Expect(window).To(Equal(history.ValidationWindow{5, 7}))
		Expect(repository.BlockCount()).To(Equal(2))
		Expect(repository.CreateOrUpdateBlockCallCount).To(Equal(2))
	})

	It("logs window message", func() {
		expectedMessage := &bytes.Buffer{}
		window := history.ValidationWindow{5, 7}
		history.ParsedWindowTemplate.Execute(expectedMessage, history.ValidationWindow{5, 7})

		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{})
		repository := inmemory.NewInMemory()
		validator := history.NewBlockValidator(blockchain, repository, 2)
		actualMessage := &bytes.Buffer{}
		validator.Log(actualMessage, window)
		Expect(actualMessage).To(Equal(expectedMessage))
	})

	It("generates a range of int64s", func() {
		numberOfBlocksCreated := history.MakeRange(0, 5)
		expected := []int64{0, 1, 2, 3, 4}

		Expect(numberOfBlocksCreated).To(Equal(expected))
	})

	It("returns the number of largest block", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 1},
			{Number: 2},
			{Number: 3},
		})
		maxBlockNumber := blockchain.LastBlock()

		Expect(maxBlockNumber.Int64()).To(Equal(int64(3)))
	})

})
