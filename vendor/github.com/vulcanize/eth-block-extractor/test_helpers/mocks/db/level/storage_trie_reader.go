package level

import (
	. "github.com/onsi/gomega"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
)

type MockStorageTrieReader struct {
	getStorageTrieCalled bool
}

func NewMockStorageTrieReader() *MockStorageTrieReader {
	return &MockStorageTrieReader{
		getStorageTrieCalled: false,
	}
}

func (mstr *MockStorageTrieReader) GetStorageTrie(stateTrieLeafNode []byte) (storageTrieResults [][]byte, err error) {
	mstr.getStorageTrieCalled = true
	return test_helpers.FakeTrieNodes, nil
}

func (mstr *MockStorageTrieReader) AssertGetStorageTrieCalled() {
	Expect(mstr.getStorageTrieCalled).To(BeTrue())
}
