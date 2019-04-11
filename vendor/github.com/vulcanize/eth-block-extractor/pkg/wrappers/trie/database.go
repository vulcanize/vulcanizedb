package trie

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
)

type GethTrieDatabase interface {
	Node(hash common.Hash) ([]byte, error)
}

type Database struct {
	db *trie.Database
}

func NewTrieDatabase(db *trie.Database) *Database {
	return &Database{db: db}
}

func (td *Database) Node(hash common.Hash) ([]byte, error) {
	return td.db.Node(hash)
}
