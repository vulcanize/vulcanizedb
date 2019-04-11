package trie

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
)

type MockIterator struct {
	includeLeaf    bool
	returnHash     common.Hash
	timesToIterate int
}

func NewMockIterator(timesToIterate int) *MockIterator {
	return &MockIterator{
		includeLeaf:    false,
		returnHash:     common.Hash{},
		timesToIterate: timesToIterate,
	}
}

func (mi *MockIterator) SetReturnHash(hash common.Hash) {
	mi.returnHash = hash
}

func (mi *MockIterator) SetIncludeLeaf() {
	mi.includeLeaf = true
}

func (mi *MockIterator) Leaf() bool {
	if mi.includeLeaf {
		mi.includeLeaf = false
		return true
	}
	return false
}

func (mi *MockIterator) LeafBlob() []byte {
	return test_helpers.FakeTrieNode
}

func (mi *MockIterator) Next(bool) bool {
	if mi.timesToIterate > 0 {
		mi.timesToIterate--
		return true
	}
	return false
}
func (mi *MockIterator) Hash() common.Hash {
	return mi.returnHash
}
