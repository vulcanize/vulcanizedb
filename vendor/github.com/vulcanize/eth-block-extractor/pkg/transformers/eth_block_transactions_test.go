package transformers_test

import (
	"io/ioutil"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/eth-block-extractor/pkg/transformers"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/db"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/ipfs"
)

var _ = Describe("Eth block transactions transformer", func() {
	Describe("executing on a single block", func() {
		BeforeEach(func() {
			log.SetOutput(ioutil.Discard)
		})

		It("returns error if ending block number is less than starting block number", func() {
			transformer := transformers.NewEthBlockTransactionsTransformer(db.NewMockDatabase(), ipfs.NewMockPublisher())

			err := transformer.Execute(1, 0)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(transformers.ErrInvalidRange))
		})

		It("fetches rlp data for block body", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockBodyByBlockNumberReturnBody([]*types.Body{{}})
			mockPublisher := ipfs.NewMockPublisher()
			mockPublisher.SetReturnStrings([][]string{{"cid"}})
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDB, mockPublisher)
			blockNumber := int64(1234567)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockDB.AssertGetBlockBodyByBlockNumberCalledWith([]int64{blockNumber})
		})

		It("publishes block body data to IPFS", func() {
			mockDB := db.NewMockDatabase()
			fakeRawData := []*types.Body{{}}
			mockDB.SetGetBlockBodyByBlockNumberReturnBody(fakeRawData)
			mockPublisher := ipfs.NewMockPublisher()
			mockPublisher.SetReturnStrings([][]string{{"cid"}})
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDB, mockPublisher)
			blockNumber := int64(1234567)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockPublisher.AssertWriteCalledWithBodies(fakeRawData)
		})

		It("returns error if publishing data returns error", func() {
			mockDB := db.NewMockDatabase()
			fakeRawData := []*types.Body{{}}
			mockDB.SetGetBlockBodyByBlockNumberReturnBody(fakeRawData)
			mockPublisher := ipfs.NewMockPublisher()
			mockPublisher.SetError(test_helpers.FakeError)
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDB, mockPublisher)
			blockNumber := int64(1234567)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(transformers.NewExecuteError(transformers.PutIpldErr, test_helpers.FakeError)))
		})
	})

	Describe("executing on a range of blocks", func() {
		It("fetches rlp data for every block's body", func() {
			mockDatabase := db.NewMockDatabase()
			mockDatabase.SetGetBlockBodyByBlockNumberReturnBody([]*types.Body{{}, {}})
			mockPublisher := ipfs.NewMockPublisher()
			mockPublisher.SetReturnStrings([][]string{{"cid_one"}, {"cid_two"}})
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDatabase, mockPublisher)
			startingBlockNumber := int64(1234567)
			endingBlockNumber := int64(1234568)

			err := transformer.Execute(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockDatabase.AssertGetBlockBodyByBlockNumberCalledWith([]int64{startingBlockNumber, endingBlockNumber})
		})

		It("publishes every block body's data to IPFS", func() {
			mockDatabase := db.NewMockDatabase()
			fakeRawData := []*types.Body{{}, {}}
			mockDatabase.SetGetBlockBodyByBlockNumberReturnBody(fakeRawData)
			mockPublisher := ipfs.NewMockPublisher()
			mockPublisher.SetReturnStrings([][]string{{"cid_one"}, {"cid_two"}})
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDatabase, mockPublisher)
			startingBlockNumber := int64(1234567)
			endingBlockNumber := int64(1234568)

			err := transformer.Execute(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockPublisher.AssertWriteCalledWithBodies(fakeRawData)
		})
	})
})
