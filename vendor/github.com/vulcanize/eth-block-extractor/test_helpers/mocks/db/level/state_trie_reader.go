package level

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/gomega"
)

type MockStateTrieReader struct {
	passedRoot common.Hash
}

func NewMockStateTrieReader() *MockStateTrieReader {
	return &MockStateTrieReader{}
}

func (mstr *MockStateTrieReader) GetStateAndStorageTrieNodes(stateRoot common.Hash) (stateTrieNodes, storageTrieNodes [][]byte, err error) {
	mstr.passedRoot = stateRoot
	return nil, nil, nil
}

func (mstr *MockStateTrieReader) AssertGetStateAndStorageTrieNodesCalledWith(root common.Hash) {
	Expect(mstr.passedRoot).To(Equal(root))
}
