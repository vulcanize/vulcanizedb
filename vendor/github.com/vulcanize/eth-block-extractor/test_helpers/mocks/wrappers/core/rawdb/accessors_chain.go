package rawdb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	. "github.com/onsi/gomega"
)

type MockAccessorsChain struct {
	getBlockPassedHash                                common.Hash
	getBlockPassedNumber                              uint64
	getBlockReturnBlock                               *types.Block
	getBlockReceiptsPassedHash                        common.Hash
	getBlockReceiptsPassedNumber                      uint64
	getBodyRLPPassedHash                              common.Hash
	getBodyRLPPassedNumber                            uint64
	getCanonicalHashPassedNumber                      uint64
	getCanonicalHashReturnHash                        common.Hash
	getHeaderPassedHash                               common.Hash
	getHeaderPassedNumber                             uint64
	getHeaderRLPPassedHash                            common.Hash
	getHeaderRLPPassedNumber                          uint64
	getStateAndStorageTrieNodesPassedRoot             common.Hash
	getStateAndStorageTrieNodesReturnErr              error
	getStateAndStorageTrieNodesReturnStateTrieBytes   [][]byte
	getStateAndStorageTrieNodesReturnStorageTrieBytes [][]byte
}

func NewMockAccessorsChain() *MockAccessorsChain {
	return &MockAccessorsChain{
		getBlockPassedHash:                                common.Hash{},
		getBlockPassedNumber:                              0,
		getBlockReturnBlock:                               nil,
		getBodyRLPPassedHash:                              common.Hash{},
		getBodyRLPPassedNumber:                            0,
		getCanonicalHashPassedNumber:                      0,
		getCanonicalHashReturnHash:                        common.Hash{},
		getHeaderRLPPassedHash:                            common.Hash{},
		getHeaderRLPPassedNumber:                          0,
		getStateAndStorageTrieNodesPassedRoot:             common.Hash{},
		getStateAndStorageTrieNodesReturnErr:              nil,
		getStateAndStorageTrieNodesReturnStateTrieBytes:   nil,
		getStateAndStorageTrieNodesReturnStorageTrieBytes: nil,
	}
}

func (accessor *MockAccessorsChain) SetGetBlockReturnBlock(returnBlock *types.Block) {
	accessor.getBlockReturnBlock = returnBlock
}

func (accessor *MockAccessorsChain) SetGetCanonicalHashReturnHash(hash common.Hash) {
	accessor.getCanonicalHashReturnHash = hash
}

func (accessor *MockAccessorsChain) SetGetStateTrieNodesReturnStateTrieBytes(returnBytes [][]byte) {
	accessor.getStateAndStorageTrieNodesReturnStateTrieBytes = returnBytes
}

func (accessor *MockAccessorsChain) SetGetStateTrieNodesReturnStorageTrieBytes(returnBytes [][]byte) {
	accessor.getStateAndStorageTrieNodesReturnStorageTrieBytes = returnBytes
}

func (accessor *MockAccessorsChain) SetGetStateTrieNodesReturnErr(err error) {
	accessor.getStateAndStorageTrieNodesReturnErr = err
}

func (accessor *MockAccessorsChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	accessor.getBlockPassedHash = hash
	accessor.getBlockPassedNumber = number
	return accessor.getBlockReturnBlock
}

func (accessor *MockAccessorsChain) GetBlockReceipts(hash common.Hash, number uint64) types.Receipts {
	accessor.getBlockReceiptsPassedHash = hash
	accessor.getBlockReceiptsPassedNumber = number
	return nil
}

func (accessor *MockAccessorsChain) GetBody(hash common.Hash, number uint64) *types.Body {
	accessor.getBodyRLPPassedHash = hash
	accessor.getBodyRLPPassedNumber = number
	return nil
}

func (accessor *MockAccessorsChain) GetCanonicalHash(number uint64) common.Hash {
	accessor.getCanonicalHashPassedNumber = number
	return accessor.getCanonicalHashReturnHash
}

func (accessor *MockAccessorsChain) GetHeader(hash common.Hash, number uint64) *types.Header {
	accessor.getHeaderPassedHash = hash
	accessor.getHeaderPassedNumber = number
	return nil
}

func (accessor *MockAccessorsChain) GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue {
	accessor.getHeaderRLPPassedHash = hash
	accessor.getHeaderRLPPassedNumber = number
	return nil
}

func (accessor *MockAccessorsChain) GetStateAndStorageTrieNodes(root common.Hash) ([][]byte, [][]byte, error) {
	accessor.getStateAndStorageTrieNodesPassedRoot = root
	return accessor.getStateAndStorageTrieNodesReturnStateTrieBytes, accessor.getStateAndStorageTrieNodesReturnStorageTrieBytes, accessor.getStateAndStorageTrieNodesReturnErr
}

func (accessor *MockAccessorsChain) AssertGetBlockCalledWith(hash common.Hash, number uint64) {
	Expect(accessor.getBlockPassedHash).To(Equal(hash))
	Expect(accessor.getBlockPassedNumber).To(Equal(number))
}

func (accessor *MockAccessorsChain) AssertGetBlockReceiptsCalledWith(hash common.Hash, number uint64) {
	Expect(accessor.getBlockReceiptsPassedHash).To(Equal(hash))
	Expect(accessor.getBlockReceiptsPassedNumber).To(Equal(number))
}

func (accessor *MockAccessorsChain) AssertGetBodyRLPCalledWith(hash common.Hash, number uint64) {
	Expect(accessor.getBodyRLPPassedHash).To(Equal(hash))
	Expect(accessor.getBodyRLPPassedNumber).To(Equal(number))
}

func (accessor *MockAccessorsChain) AssertGetCanonicalHashCalledWith(number uint64) {
	Expect(accessor.getCanonicalHashPassedNumber).To(Equal(number))
}

func (accessor *MockAccessorsChain) AssertGetHeaderCalledWith(hash common.Hash, number uint64) {
	Expect(accessor.getHeaderPassedHash).To(Equal(hash))
	Expect(accessor.getHeaderPassedNumber).To(Equal(number))
}

func (accessor *MockAccessorsChain) AssertGetHeaderRLPCalledWith(hash common.Hash, number uint64) {
	Expect(accessor.getHeaderRLPPassedHash).To(Equal(hash))
	Expect(accessor.getHeaderRLPPassedNumber).To(Equal(number))
}

func (accessor *MockAccessorsChain) AssertGetStateTrieNodesCalledWith(root common.Hash) {
	Expect(accessor.getStateAndStorageTrieNodesPassedRoot).To(Equal(root))
}
