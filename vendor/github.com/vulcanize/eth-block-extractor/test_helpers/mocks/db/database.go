package db

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
)

type MockDatabase struct {
	computeBlockStateTrieErr                          error
	computeBlockStateTriePassedCurrentBlock           *types.Block
	computeBlockStateTriePassedParentBlock            *types.Block
	computeBlockStateTrieReturnHash                   common.Hash
	getBlockBodyByBlockNumberPassedBlockNumbers       []int64
	getBlockBodyByBlockNumberReturnBodies             []*types.Body
	getBlockByBlockNumberPassedNumbers                []int64
	getBlockByBlockNumberReturnBlock                  *types.Block
	getBlockHeaderByBlockNumberPassedBlockNumbers     []int64
	getBlockHeaderByBlockNumberReturnHeader           *types.Header
	getRawBlockHeaderByBlockNumberPassedBlockNumbers  []int64
	getRawBlockHeaderByBlockNumberReturnBytes         [][]byte
	getBlockReceiptsPassedBlockNumbers                []int64
	getBlockReceiptsReturnReceipts                    types.Receipts
	getStateAndStorageTrieNodesErr                    error
	getStateAndStorageTrieNodesPassedRoot             common.Hash
	getStateAndStorageTrieNodesReturnStateTrieBytes   [][]byte
	getStateAndStorageTrieNodesReturnStorageTrieBytes [][]byte
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		computeBlockStateTrieErr:                          nil,
		computeBlockStateTriePassedCurrentBlock:           nil,
		computeBlockStateTriePassedParentBlock:            nil,
		computeBlockStateTrieReturnHash:                   common.Hash{},
		getBlockBodyByBlockNumberPassedBlockNumbers:       nil,
		getBlockBodyByBlockNumberReturnBodies:             nil,
		getBlockByBlockNumberPassedNumbers:                nil,
		getBlockByBlockNumberReturnBlock:                  nil,
		getBlockHeaderByBlockNumberPassedBlockNumbers:     nil,
		getBlockHeaderByBlockNumberReturnHeader:           nil,
		getRawBlockHeaderByBlockNumberPassedBlockNumbers:  nil,
		getRawBlockHeaderByBlockNumberReturnBytes:         nil,
		getBlockReceiptsPassedBlockNumbers:                nil,
		getBlockReceiptsReturnReceipts:                    nil,
		getStateAndStorageTrieNodesErr:                    nil,
		getStateAndStorageTrieNodesPassedRoot:             common.Hash{},
		getStateAndStorageTrieNodesReturnStateTrieBytes:   nil,
		getStateAndStorageTrieNodesReturnStorageTrieBytes: nil,
	}
}

func (db *MockDatabase) SetComputeBlockStateTrieError(err error) {
	db.computeBlockStateTrieErr = err
}

func (db *MockDatabase) SetComputeBlockStateTrieReturnHash(hash common.Hash) {
	db.computeBlockStateTrieReturnHash = hash
}

func (db *MockDatabase) SetGetBlockBodyByBlockNumberReturnBody(bodies []*types.Body) {
	db.getBlockBodyByBlockNumberReturnBodies = bodies
}

func (db *MockDatabase) SetGetBlockByBlockNumberReturnBlock(returnBlock *types.Block) {
	db.getBlockByBlockNumberReturnBlock = returnBlock
}

func (db *MockDatabase) SetGetBlockHeaderByBlockNumberReturnHeader(header *types.Header) {
	db.getBlockHeaderByBlockNumberReturnHeader = header
}

func (db *MockDatabase) SetGetRawBlockHeaderByBlockNumberReturnBytes(returnBytes [][]byte) {
	db.getRawBlockHeaderByBlockNumberReturnBytes = returnBytes
}

func (db *MockDatabase) SetGetBlockReceiptsReturnReceipts(receipts types.Receipts) {
	db.getBlockReceiptsReturnReceipts = receipts
}

func (db *MockDatabase) SetGetStateAndStorageTrieNodesError(err error) {
	db.getStateAndStorageTrieNodesErr = err
}

func (db *MockDatabase) SetGetStateAndStorageTrieNodesReturnStateTrieBytes(returnBytes [][]byte) {
	db.getStateAndStorageTrieNodesReturnStateTrieBytes = returnBytes
}

func (db *MockDatabase) SetGetStateAndStorageTrieNodesReturnStorageTrieBytes(returnBytes [][]byte) {
	db.getStateAndStorageTrieNodesReturnStorageTrieBytes = returnBytes
}

func (db *MockDatabase) ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) (common.Hash, error) {
	db.computeBlockStateTriePassedCurrentBlock = currentBlock
	db.computeBlockStateTriePassedParentBlock = parentBlock
	return db.computeBlockStateTrieReturnHash, db.computeBlockStateTrieErr
}

func (db *MockDatabase) GetBlockBodyByBlockNumber(blockNumber int64) *types.Body {
	db.getBlockBodyByBlockNumberPassedBlockNumbers = append(db.getBlockBodyByBlockNumberPassedBlockNumbers, blockNumber)
	returnBytes := db.getBlockBodyByBlockNumberReturnBodies[0]
	db.getBlockBodyByBlockNumberReturnBodies = db.getBlockBodyByBlockNumberReturnBodies[1:]
	return returnBytes
}

func (db *MockDatabase) GetBlockByBlockNumber(blockNumber int64) *types.Block {
	db.getBlockByBlockNumberPassedNumbers = append(db.getBlockByBlockNumberPassedNumbers, blockNumber)
	return db.getBlockByBlockNumberReturnBlock
}

func (db *MockDatabase) GetBlockHeaderByBlockNumber(blockNumber int64) *types.Header {
	db.getBlockHeaderByBlockNumberPassedBlockNumbers = append(db.getBlockHeaderByBlockNumberPassedBlockNumbers, blockNumber)
	return db.getBlockHeaderByBlockNumberReturnHeader
}

func (db *MockDatabase) GetRawBlockHeaderByBlockNumber(blockNumber int64) []byte {
	db.getRawBlockHeaderByBlockNumberPassedBlockNumbers = append(db.getRawBlockHeaderByBlockNumberPassedBlockNumbers, blockNumber)
	returnBytes := db.getRawBlockHeaderByBlockNumberReturnBytes[0]
	db.getRawBlockHeaderByBlockNumberReturnBytes = db.getRawBlockHeaderByBlockNumberReturnBytes[1:]
	return returnBytes
}

func (db *MockDatabase) GetBlockReceipts(blockNumber int64) types.Receipts {
	db.getBlockReceiptsPassedBlockNumbers = append(db.getBlockReceiptsPassedBlockNumbers, blockNumber)
	return db.getBlockReceiptsReturnReceipts
}

func (db *MockDatabase) GetStateAndStorageTrieNodes(root common.Hash) ([][]byte, [][]byte, error) {
	db.getStateAndStorageTrieNodesPassedRoot = root
	return db.getStateAndStorageTrieNodesReturnStateTrieBytes, db.getStateAndStorageTrieNodesReturnStorageTrieBytes, db.getStateAndStorageTrieNodesErr
}

func (db *MockDatabase) AssertComputeBlockStateTrieCalledWith(currentBlock *types.Block, parentBlock *types.Block) {
	Expect(db.computeBlockStateTriePassedCurrentBlock).To(Equal(currentBlock))
	Expect(db.computeBlockStateTriePassedParentBlock).To(Equal(parentBlock))
}

func (db *MockDatabase) AssertGetBlockBodyByBlockNumberCalledWith(blockNumbers []int64) {
	Expect(db.getBlockBodyByBlockNumberPassedBlockNumbers).To(Equal(blockNumbers))
}

func (db *MockDatabase) AssertGetBlockByBlockNumberCalledwith(blockNumbers []int64) {
	for i := 0; i < len(blockNumbers); i++ {
		Expect(db.getBlockByBlockNumberPassedNumbers).To(ContainElement(blockNumbers[i]))
	}
	for i := 0; i < len(db.getBlockByBlockNumberPassedNumbers); i++ {
		Expect(blockNumbers).To(ContainElement(db.getBlockByBlockNumberPassedNumbers[i]))
	}
}

func (db *MockDatabase) AssertGetBlockHeaderByBlockNumberCalledWith(blockNumbers []int64) {
	Expect(db.getBlockHeaderByBlockNumberPassedBlockNumbers).To(Equal(blockNumbers))
}

func (db *MockDatabase) AssertGetRawBlockHeaderByBlockNumberCalledWith(blockNumbers []int64) {
	Expect(db.getRawBlockHeaderByBlockNumberPassedBlockNumbers).To(Equal(blockNumbers))
}

func (db *MockDatabase) AssertGetBlockReceiptsCalledWith(blockNumbers []int64) {
	Expect(db.getBlockReceiptsPassedBlockNumbers).To(Equal(blockNumbers))
}

func (db *MockDatabase) AssertGetStateTrieNodesCalledWith(root common.Hash) {
	Expect(db.getStateAndStorageTrieNodesPassedRoot).To(Equal(root))
}
