package eth_state_trie_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_state_trie"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/ipfs"
)

var _ = Describe("Ethereum state trie node dag putter", func() {
	It("adds passed state trie node to ipfs", func() {
		mockAdder := ipfs.NewMockAdder()
		dagPutter := eth_state_trie.NewStateTrieDagPutter(mockAdder)

		_, err := dagPutter.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).NotTo(HaveOccurred())
		mockAdder.AssertAddCalled(1, &eth_state_trie.EthStateTrieNode{})
	})

	It("returns error if adding to ipfs fails", func() {
		mockAdder := ipfs.NewMockAdder()
		mockAdder.SetError(test_helpers.FakeError)
		dagPutter := eth_state_trie.NewStateTrieDagPutter(mockAdder)

		_, err := dagPutter.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})
})
