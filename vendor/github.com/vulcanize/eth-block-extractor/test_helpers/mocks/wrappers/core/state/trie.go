package state

import (
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/trie"
)

type MockTrie struct {
	iterator trie.GethTrieNodeIterator
}

func NewMockTrie() *MockTrie {
	return &MockTrie{}
}

func (mt *MockTrie) SetReturnIterator(iterator trie.GethTrieNodeIterator) {
	mt.iterator = iterator
}

func (mt *MockTrie) NodeIterator(startKey []byte) trie.GethTrieNodeIterator {
	return mt.iterator
}
