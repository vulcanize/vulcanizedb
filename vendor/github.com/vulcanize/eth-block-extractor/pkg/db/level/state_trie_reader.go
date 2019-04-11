package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/core/state"
)

type IStateTrieReader interface {
	GetStateAndStorageTrieNodes(stateRoot common.Hash) (stateTrieNodes, storageTrieNodes [][]byte, err error)
}

type StateTrieReader struct {
	db                state.GethStateDatabase
	storageTrieReader IStorageTrieReader
}

func NewStateTrieReader(db state.GethStateDatabase, storageTrieReader IStorageTrieReader) *StateTrieReader {
	return &StateTrieReader{
		db:                db,
		storageTrieReader: storageTrieReader,
	}
}

func (str *StateTrieReader) GetStateAndStorageTrieNodes(stateRoot common.Hash) (stateTrieNodes, storageTrieNodes [][]byte, err error) {
	trieDb := str.db.TrieDB()
	// fetch and append state root node
	stateRootNode, err := trieDb.Node(stateRoot)
	if err != nil {
		return stateTrieNodes, storageTrieNodes, err
	}
	stateTrieNodes = append(stateTrieNodes, stateRootNode)

	// fetch and append remaining nodes in the state trie
	stateTrie, err := str.db.OpenTrie(stateRoot)
	if err != nil {
		return stateTrieNodes, storageTrieNodes, err
	}
	stateTrieIterator := stateTrie.NodeIterator(nil)
	for stateTrieIterator.Next(true) {
		if stateTrieIterator.Leaf() {
			node := stateTrieIterator.LeafBlob()
			stateTrieNodes = append(stateTrieNodes, node)
			// fetch and append storage trie nodes for state trie leaf (account snapshot)
			accountStorageTrieNodes, err := str.storageTrieReader.GetStorageTrie(node)
			if err != nil {
				return stateTrieNodes, storageTrieNodes, err
			}
			storageTrieNodes = append(storageTrieNodes, accountStorageTrieNodes...)
		} else {
			nodeKey := stateTrieIterator.Hash()
			node, err := str.db.TrieDB().Node(nodeKey)
			if err != nil {
				return stateTrieNodes, storageTrieNodes, err
			}
			stateTrieNodes = append(stateTrieNodes, node)
		}
	}
	return stateTrieNodes, storageTrieNodes, nil
}
