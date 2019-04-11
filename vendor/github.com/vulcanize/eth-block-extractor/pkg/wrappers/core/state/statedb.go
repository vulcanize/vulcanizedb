package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type GethStateDB interface {
	Commit(deleteEmptyObjects bool) (common.Hash, error)
	StateDB() *state.StateDB
}

type StateDB struct {
	db *state.StateDB
}

func (st *StateDB) StateDB() *state.StateDB {
	return st.db
}

func (st *StateDB) Commit(deleteEmptyObjects bool) (common.Hash, error) {
	return st.db.Commit(deleteEmptyObjects)
}
