package transformers_test

import (
	"io/ioutil"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/eth-block-extractor/pkg/transformers"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/db"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/ipfs"
)

var _ = Describe("Eth block header transformer", func() {
	var mockDB *db.MockDatabase
	var mockPublisher *ipfs.MockPublisher

	Describe("Fetching one block header", func() {
		var fakeBytes []byte
		var blockNumber int64

		BeforeEach(func() {
			mockDB = db.NewMockDatabase()
			mockPublisher = ipfs.NewMockPublisher()
			blockNumber = 54321
			fakeBytes = []byte{6, 7, 8, 9, 0}
			mockDB.SetGetRawBlockHeaderByBlockNumberReturnBytes([][]byte{fakeBytes})
			mockPublisher.SetReturnStrings([][]string{{"cid_one", "cid_two"}})
			log.SetOutput(ioutil.Discard)
		})

		It("returns error if ending block number is less than starting block number", func() {
			transformer := transformers.NewEthBlockHeaderTransformer(mockDB, mockPublisher)

			err := transformer.Execute(1, 0)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(transformers.ErrInvalidRange))
		})

		It("Fetches RLP data from ethereum db", func() {
			transformer := transformers.NewEthBlockHeaderTransformer(mockDB, mockPublisher)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockDB.AssertGetRawBlockHeaderByBlockNumberCalledWith([]int64{blockNumber})
		})

		It("Persists block RLP data to IPFS", func() {
			transformer := transformers.NewEthBlockHeaderTransformer(mockDB, mockPublisher)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockPublisher.AssertWriteCalledWithBytes([][]byte{fakeBytes})
		})

		It("Returns err if persisting block RLP data to IPFS fails", func() {
			mockPublisher.SetError(test_helpers.FakeError)
			transformer := transformers.NewEthBlockHeaderTransformer(mockDB, mockPublisher)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(transformers.NewExecuteError(transformers.PutIpldErr, test_helpers.FakeError)))
			mockDB.AssertGetRawBlockHeaderByBlockNumberCalledWith([]int64{blockNumber})
			mockPublisher.AssertWriteCalledWithBytes([][]byte{fakeBytes})
		})
	})

	Describe("Fetching multiple blocks", func() {
		var fakeRlpBytes []byte
		var startingBlockNumber int64
		var endingBlockNumber int64

		BeforeEach(func() {
			mockDB = db.NewMockDatabase()
			mockPublisher = ipfs.NewMockPublisher()
			startingBlockNumber = 54321
			endingBlockNumber = 54322
			fakeRlpBytes = []byte{6, 7, 8, 9, 0}
			mockDB.SetGetRawBlockHeaderByBlockNumberReturnBytes([][]byte{fakeRlpBytes, fakeRlpBytes})
			fakeOutputString := []string{"cid_one", "cid_two"}
			mockPublisher.SetReturnStrings([][]string{fakeOutputString, fakeOutputString})
			log.SetOutput(ioutil.Discard)
		})

		It("Fetches block RLP data from ethereum db for every block in range", func() {
			transformer := transformers.NewEthBlockHeaderTransformer(mockDB, mockPublisher)

			err := transformer.Execute(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockDB.AssertGetRawBlockHeaderByBlockNumberCalledWith([]int64{startingBlockNumber, endingBlockNumber})
		})

		It("Persists block RLP data to IPFS for every block in range", func() {
			transformer := transformers.NewEthBlockHeaderTransformer(mockDB, mockPublisher)

			err := transformer.Execute(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockPublisher.AssertWriteCalledWithBytes([][]byte{fakeRlpBytes, fakeRlpBytes})
		})
	})

})
