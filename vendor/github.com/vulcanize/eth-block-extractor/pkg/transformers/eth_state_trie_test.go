package transformers_test

import (
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/eth-block-extractor/pkg/transformers"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/db"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/ipfs"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/rlp"
)

var _ = Describe("Eth state trie transformer", func() {
	BeforeEach(func() {
		log.SetOutput(ioutil.Discard)
	})

	It("returns error if ending block number is less than starting block number", func() {
		transformer := transformers.NewEthStateTrieTransformer(db.NewMockDatabase(), ipfs.NewMockPublisher(), ipfs.NewMockPublisher())

		err := transformer.Execute(1, 0)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(transformers.ErrInvalidRange))
	})

	It("fetches block header for block", func() {
		mockDB := db.NewMockDatabase()
		mockDB.SetGetBlockHeaderByBlockNumberReturnHeader(&types.Header{})
		transformer := transformers.NewEthStateTrieTransformer(mockDB, ipfs.NewMockPublisher(), ipfs.NewMockPublisher())

		err := transformer.Execute(0, 0)

		Expect(err).NotTo(HaveOccurred())
		mockDB.AssertGetBlockHeaderByBlockNumberCalledWith([]int64{0})
	})

	It("fetches state and storage trie nodes with state root from decoded block header", func() {
		mockDB := db.NewMockDatabase()
		mockDB.SetGetBlockHeaderByBlockNumberReturnHeader(&types.Header{Root: test_helpers.FakeHash})
		transformer := transformers.NewEthStateTrieTransformer(mockDB, ipfs.NewMockPublisher(), ipfs.NewMockPublisher())

		err := transformer.Execute(0, 0)

		Expect(err).NotTo(HaveOccurred())
		mockDB.AssertGetStateTrieNodesCalledWith(test_helpers.FakeHash)
	})

	It("returns err if fetching state trie returns err", func() {
		mockDB := db.NewMockDatabase()
		mockDB.SetGetBlockHeaderByBlockNumberReturnHeader(&types.Header{})
		mockDB.SetGetStateAndStorageTrieNodesError(test_helpers.FakeError)
		transformer := transformers.NewEthStateTrieTransformer(mockDB, ipfs.NewMockPublisher(), ipfs.NewMockPublisher())

		err := transformer.Execute(0, 0)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(test_helpers.FakeError.Error()))
	})

	It("writes state trie nodes to ipfs", func() {
		mockDB := db.NewMockDatabase()
		mockDB.SetGetBlockHeaderByBlockNumberReturnHeader(&types.Header{})
		fakeStateTrieNodes := [][]byte{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 0}}
		mockDB.SetGetStateAndStorageTrieNodesReturnStateTrieBytes(fakeStateTrieNodes)
		mockDecoder := rlp.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{})
		mockStateTriePublisher := ipfs.NewMockPublisher()
		mockStateTriePublisher.SetReturnStrings([][]string{{"one"}, {"two"}})
		transformer := transformers.NewEthStateTrieTransformer(mockDB, mockStateTriePublisher, ipfs.NewMockPublisher())

		err := transformer.Execute(0, 0)

		Expect(err).NotTo(HaveOccurred())
		mockStateTriePublisher.AssertWriteCalledWithBytes(fakeStateTrieNodes)
	})

	It("writes storage trie nodes to ipfs", func() {
		mockDB := db.NewMockDatabase()
		mockDB.SetGetBlockHeaderByBlockNumberReturnHeader(&types.Header{})
		fakeStateTrieNodes := [][]byte{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 0}}
		mockDB.SetGetStateAndStorageTrieNodesReturnStorageTrieBytes(fakeStateTrieNodes)
		mockDecoder := rlp.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{})
		mockStorageTriePublisher := ipfs.NewMockPublisher()
		mockStorageTriePublisher.SetReturnStrings([][]string{{"one"}, {"two"}})
		transformer := transformers.NewEthStateTrieTransformer(mockDB, ipfs.NewMockPublisher(), mockStorageTriePublisher)

		err := transformer.Execute(0, 0)

		Expect(err).NotTo(HaveOccurred())
		mockStorageTriePublisher.AssertWriteCalledWithBytes(fakeStateTrieNodes)
	})
})
