package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/trie"
)

type GethStateDatabase interface {
	Database() state.Database
	OpenTrie(root common.Hash) (GethTrie, error)
	TrieDB() trie.GethTrieDatabase
}

type Database struct {
	db     state.Database
	trieDB trie.GethTrieDatabase
}

func NewDatabase(databaseConnection ethdb.Database) *Database {
	db := state.NewDatabase(databaseConnection)
	trieDB := trie.NewTrieDatabase(db.TrieDB())
	return &Database{db: db, trieDB: trieDB}
}

func (db Database) Database() state.Database {
	return db.db
}

func (db Database) OpenTrie(root common.Hash) (GethTrie, error) {
	stateTrie, err := db.db.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	return NewTrie(stateTrie), nil
}

func (db Database) TrieDB() trie.GethTrieDatabase {
	return db.trieDB
}
