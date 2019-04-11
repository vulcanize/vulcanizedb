package level_test

import (
	"github.com/ethereum/go-ethereum/core/state"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/eth-block-extractor/pkg/db/level"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	state_wrapper "github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/core/state"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/rlp"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/trie"
)

var _ = Describe("Storage trie reader", func() {
	It("decodes passed state trie leaf node into account", func() {
		db := state_wrapper.NewMockStateDatabase()
		trieDb := db.CreateFakeUnderlyingDatabase()
		db.SetReturnDatabase(trieDb)
		mockIteratror := trie.NewMockIterator(1)
		mockIteratror.SetIncludeLeaf()
		mockTrie := state_wrapper.NewMockTrie()
		mockTrie.SetReturnIterator(mockIteratror)
		db.SetReturnTrie(mockTrie)
		decoder := rlp.NewMockDecoder()
		decoder.SetReturnOut(&state.Account{})
		reader := level.NewStorageTrieReader(db, decoder)

		_, err := reader.GetStorageTrie(test_helpers.FakeTrieNode)

		Expect(err).NotTo(HaveOccurred())
		decoder.AssertDecodeCalledWith(test_helpers.FakeTrieNode, &state.Account{})
	})

	It("fetches node associated with storage root", func() {
		db := state_wrapper.NewMockStateDatabase()
		trieDb := db.CreateFakeUnderlyingDatabase()
		db.SetReturnDatabase(trieDb)
		mockIteratror := trie.NewMockIterator(0)
		mockTrie := state_wrapper.NewMockTrie()
		mockTrie.SetReturnIterator(mockIteratror)
		db.SetReturnTrie(mockTrie)
		decoder := rlp.NewMockDecoder()
		decoder.SetReturnOut(&state.Account{})
		reader := level.NewStorageTrieReader(db, decoder)

		storageTrieNodes, err := reader.GetStorageTrie(test_helpers.FakeTrieNode)

		Expect(err).NotTo(HaveOccurred())
		Expect(len(storageTrieNodes)).To(Equal(1))
	})

	It("returns nodes found traversing storage trie", func() {
		db := state_wrapper.NewMockStateDatabase()
		trieDb := db.CreateFakeUnderlyingDatabase()
		db.SetReturnDatabase(trieDb)
		mockIteratror := trie.NewMockIterator(1)
		mockIteratror.SetIncludeLeaf()
		mockTrie := state_wrapper.NewMockTrie()
		mockTrie.SetReturnIterator(mockIteratror)
		db.SetReturnTrie(mockTrie)
		decoder := rlp.NewMockDecoder()
		decoder.SetReturnOut(&state.Account{})
		reader := level.NewStorageTrieReader(db, decoder)

		storageTrieNodes, err := reader.GetStorageTrie(test_helpers.FakeTrieNode)

		Expect(err).NotTo(HaveOccurred())
		Expect(len(storageTrieNodes)).To(Equal(2))
	})
})
