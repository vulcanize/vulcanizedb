package transformers_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/eth-block-extractor/pkg/transformers"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/db"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/ipfs"
	"io/ioutil"
	"log"
)

var _ = Describe("Eth block receipts transformer", func() {
	BeforeEach(func() {
		log.SetOutput(ioutil.Discard)
	})

	It("returns error if ending block number is less than starting block number", func() {
		transformer := transformers.NewEthBlockReceiptTransformer(db.NewMockDatabase(), ipfs.NewMockPublisher())

		err := transformer.Execute(1, 0)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(transformers.ErrInvalidRange))
	})

	It("fetches blocks' receipts from database", func() {
		mockDatabase := db.NewMockDatabase()
		transformer := transformers.NewEthBlockReceiptTransformer(mockDatabase, ipfs.NewMockPublisher())

		err := transformer.Execute(0, 1)

		Expect(err).NotTo(HaveOccurred())
		mockDatabase.AssertGetBlockReceiptsCalledWith([]int64{0, 1})
	})

	It("publishes block receipts", func() {
		mockDatabase := db.NewMockDatabase()
		fakeReceipts := types.Receipts{
			&types.Receipt{},
			&types.Receipt{},
		}
		mockDatabase.SetGetBlockReceiptsReturnReceipts(fakeReceipts)
		mockPublisher := ipfs.NewMockPublisher()
		transformer := transformers.NewEthBlockReceiptTransformer(mockDatabase, mockPublisher)

		err := transformer.Execute(0, 0)

		Expect(err).NotTo(HaveOccurred())
		mockPublisher.AssertWriteCalledWithInterfaces([]interface{}{fakeReceipts})
	})

	It("returns error if publishing block receipts fails", func() {
		mockDatabase := db.NewMockDatabase()
		fakeReceipts := types.Receipts{
			&types.Receipt{},
			&types.Receipt{},
		}
		mockDatabase.SetGetBlockReceiptsReturnReceipts(fakeReceipts)
		mockPublisher := ipfs.NewMockPublisher()
		mockPublisher.SetError(test_helpers.FakeError)
		transformer := transformers.NewEthBlockReceiptTransformer(mockDatabase, mockPublisher)

		err := transformer.Execute(0, 0)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})
})
