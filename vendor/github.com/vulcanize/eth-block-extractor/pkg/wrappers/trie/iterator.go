package trie

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
)

type GethTrieNodeIterator interface {
	Hash() common.Hash
	Leaf() bool
	LeafBlob() []byte
	Next(bool) bool
}

type NodeIterator struct {
	iterator trie.NodeIterator
}

func NewNodeIterator(nodeIterator trie.NodeIterator) *NodeIterator {
	return &NodeIterator{iterator: nodeIterator}
}

func (ni *NodeIterator) Hash() common.Hash {
	return ni.iterator.Hash()
}

func (ni *NodeIterator) Leaf() bool {
	return ni.iterator.Leaf()
}

func (ni *NodeIterator) LeafBlob() []byte {
	return ni.iterator.LeafBlob()
}

func (ni *NodeIterator) Next(b bool) bool {
	return ni.iterator.Next(b)
}
