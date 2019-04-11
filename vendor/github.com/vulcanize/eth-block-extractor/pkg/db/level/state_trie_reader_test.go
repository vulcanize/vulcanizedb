package level_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/eth-block-extractor/pkg/db/level"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	level_wrapper "github.com/vulcanize/eth-block-extractor/test_helpers/mocks/db/level"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/core/state"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/trie"
)

var _ = Describe("State trie reader", func() {
	It("fetches node associated with state root", func() {
		db := state.NewMockStateDatabase()
		trieDB := db.CreateFakeUnderlyingDatabase()
		db.SetReturnDatabase(trieDB)
		mockIteratror := trie.NewMockIterator(0)
		mockTrie := state.NewMockTrie()
		mockTrie.SetReturnIterator(mockIteratror)
		db.SetReturnTrie(mockTrie)
		mockStorageTrieReader := level_wrapper.NewMockStorageTrieReader()
		reader := level.NewStateTrieReader(db, mockStorageTrieReader)

		stateTrieNodes, _, err := reader.GetStateAndStorageTrieNodes(test_helpers.FakeHash)

		Expect(err).NotTo(HaveOccurred())
		Expect(stateTrieNodes).To(ContainElement(test_helpers.FakeTrieNode))
	})

	It("returns nodes found traversing state trie", func() {
		db := state.NewMockStateDatabase()
		trieDB := db.CreateFakeUnderlyingDatabase()
		db.SetReturnDatabase(trieDB)
		mockIteratror := trie.NewMockIterator(2)
		mockTrie := state.NewMockTrie()
		mockTrie.SetReturnIterator(mockIteratror)
		db.SetReturnTrie(mockTrie)
		mockStorageTrieReader := level_wrapper.NewMockStorageTrieReader()
		reader := level.NewStateTrieReader(db, mockStorageTrieReader)

		stateTrieNodes, storageTrieNodes, err := reader.GetStateAndStorageTrieNodes(test_helpers.FakeHash)

		Expect(err).NotTo(HaveOccurred())
		Expect(len(stateTrieNodes)).To(Equal(3))
		Expect(len(storageTrieNodes)).To(BeZero())
	})

	It("invokes storage trie reader for state trie leaf nodes", func() {
		db := state.NewMockStateDatabase()
		trieDB := db.CreateFakeUnderlyingDatabase()
		db.SetReturnDatabase(trieDB)
		mockIteratror := trie.NewMockIterator(1)
		mockIteratror.SetIncludeLeaf()
		mockTrie := state.NewMockTrie()
		mockTrie.SetReturnIterator(mockIteratror)
		db.SetReturnTrie(mockTrie)
		mockStorageTrieReader := level_wrapper.NewMockStorageTrieReader()
		reader := level.NewStateTrieReader(db, mockStorageTrieReader)

		_, _, err := reader.GetStateAndStorageTrieNodes(test_helpers.FakeHash)

		Expect(err).NotTo(HaveOccurred())
		mockStorageTrieReader.AssertGetStorageTrieCalled()
	})
})
