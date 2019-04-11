package state

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/trie"
)

type GethTrie interface {
	NodeIterator(startKey []byte) trie.GethTrieNodeIterator
}

type Trie struct {
	trie state.Trie
}

func NewTrie(trie state.Trie) GethTrie {
	return &Trie{trie: trie}
}

func (t *Trie) NodeIterator(startKey []byte) trie.GethTrieNodeIterator {
	iterator := t.trie.NodeIterator(startKey)
	return trie.NewNodeIterator(iterator)
}
