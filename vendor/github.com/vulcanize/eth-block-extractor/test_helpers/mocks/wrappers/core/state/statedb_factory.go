package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	. "github.com/onsi/gomega"
	state_wrapper "github.com/vulcanize/eth-block-extractor/pkg/wrappers/core/state"
)

type MockStateDBFactory struct {
	passedDatabase state.Database
	passedRoot     common.Hash
	returnErr      error
	returnStateDB  state_wrapper.GethStateDB
}

func NewMockStateDBFactory() *MockStateDBFactory {
	return &MockStateDBFactory{
		passedDatabase: nil,
		passedRoot:     common.Hash{},
		returnErr:      nil,
		returnStateDB:  nil,
	}
}

func (mstf *MockStateDBFactory) SetStateDB(stateTrie state_wrapper.GethStateDB) {
	mstf.returnStateDB = stateTrie
}

func (mstf *MockStateDBFactory) SetReturnErr(err error) {
	mstf.returnErr = err
}

func (mstf *MockStateDBFactory) NewStateDB(root common.Hash, db state.Database) (state_wrapper.GethStateDB, error) {
	mstf.passedRoot = root
	mstf.passedDatabase = db
	return mstf.returnStateDB, mstf.returnErr
}

func (mstf *MockStateDBFactory) AssertNewStateTrieCalledWith(root common.Hash, db state.Database) {
	Expect(mstf.passedRoot).To(Equal(root))
	Expect(mstf.passedDatabase).To(Equal(db))
}
