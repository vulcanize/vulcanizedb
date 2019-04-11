package level

import (
	"bytes"
	"github.com/ethereum/go-ethereum/core/state"
	state_wrapper "github.com/vulcanize/eth-block-extractor/pkg/wrappers/core/state"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/rlp"
)

var EmptyStorageTrieRoot = []byte{86, 232, 31, 23, 27, 204, 85, 166, 255, 131, 69, 230, 146, 192, 248, 110, 91, 72, 224, 27, 153, 108, 173, 192, 1, 98, 47, 181, 227, 99, 180, 33}

type IStorageTrieReader interface {
	GetStorageTrie(stateTrieLeafNode []byte) (storageTrieResults [][]byte, err error)
}

type StorageTrieReader struct {
	db      state_wrapper.GethStateDatabase
	decoder rlp.Decoder
}

func NewStorageTrieReader(db state_wrapper.GethStateDatabase, decoder rlp.Decoder) *StorageTrieReader {
	return &StorageTrieReader{
		db:      db,
		decoder: decoder,
	}
}

func (stc *StorageTrieReader) GetStorageTrie(stateTrieLeafNode []byte) (storageTrieResults [][]byte, err error) {
	trieDb := stc.db.TrieDB()
	var account state.Account
	err = stc.decoder.Decode(stateTrieLeafNode, &account)
	if err != nil {
		return storageTrieResults, err
	}
	// if storage trie root corresponds to empty storage trie, continue to next iteration in state trie
	if bytes.Equal(EmptyStorageTrieRoot, account.Root.Bytes()) {
		return storageTrieResults, err
	}
	// if storage trie root not empty, fetch and append root node
	storageRootNode, err := trieDb.Node(account.Root)
	if err != nil {
		return storageTrieResults, err
	}
	storageTrieResults = append(storageTrieResults, storageRootNode)
	storageTrie, err := stc.db.OpenTrie(account.Root)
	if err != nil {
		return storageTrieResults, err
	}
	storageTrieIterator := storageTrie.NodeIterator(account.Root.Bytes())
	for storageTrieIterator.Next(true) {
		if storageTrieIterator.Leaf() {
			storageTrieResults = append(storageTrieResults, storageTrieIterator.LeafBlob())
		} else {
			nextStorageHash := storageTrieIterator.Hash()
			nextStorageNode, err := trieDb.Node(nextStorageHash)
			if err != nil {
				return storageTrieResults, err
			}
			storageTrieResults = append(storageTrieResults, nextStorageNode)
		}
	}
	return storageTrieResults, err
}
