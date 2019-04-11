package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type GethStateDBFactory interface {
	NewStateDB(root common.Hash, db state.Database) (GethStateDB, error)
}

type StateDBFactory struct {
}

func NewStateDBFactory() *StateDBFactory {
	return &StateDBFactory{}
}

func (stf *StateDBFactory) NewStateDB(root common.Hash, db state.Database) (GethStateDB, error) {
	stateDb, err := state.New(root, db)
	if err != nil {
		return nil, err
	}
	return &StateDB{db: stateDb}, nil
}
