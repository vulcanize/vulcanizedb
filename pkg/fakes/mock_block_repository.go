package fakes

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockBlockRepository struct {
	createOrUpdateBlockCalled                    bool
	createOrUpdateBlockPassedBlock               core.Block
	createOrUpdateBlockReturnInt                 int64
	createOrUpdateBlockReturnErr                 error
	missingBlockNumbersCalled                    bool
	missingBlockNumbersPassedStartingBlockNumber int64
	missingBlockNumbersPassedEndingBlockNumber   int64
	missingBlockNumbersPassedNodeId              string
	missingBlockNumbersReturnArray               []int64
	setBlockStatusCalled                         bool
	setBlockStatusPassedChainHead                int64
}

func NewMockBlockRepository() *MockBlockRepository {
	return &MockBlockRepository{
		createOrUpdateBlockCalled:                    false,
		createOrUpdateBlockPassedBlock:               core.Block{},
		createOrUpdateBlockReturnInt:                 0,
		createOrUpdateBlockReturnErr:                 nil,
		missingBlockNumbersCalled:                    false,
		missingBlockNumbersPassedStartingBlockNumber: 0,
		missingBlockNumbersPassedEndingBlockNumber:   0,
		missingBlockNumbersPassedNodeId:              "",
		missingBlockNumbersReturnArray:               nil,
		setBlockStatusCalled:                         false,
		setBlockStatusPassedChainHead:                0,
	}
}

func (mbr *MockBlockRepository) SetCreateOrUpdateBlockReturnVals(i int64, err error) {
	mbr.createOrUpdateBlockReturnInt = i
	mbr.createOrUpdateBlockReturnErr = err
}

func (mbr *MockBlockRepository) SetMissingBlockNumbersReturnArray(returnArray []int64) {
	mbr.missingBlockNumbersReturnArray = returnArray
}

func (mbr *MockBlockRepository) CreateOrUpdateBlock(block core.Block) (int64, error) {
	mbr.createOrUpdateBlockCalled = true
	mbr.createOrUpdateBlockPassedBlock = block
	return mbr.createOrUpdateBlockReturnInt, mbr.createOrUpdateBlockReturnErr
}

func (mbr *MockBlockRepository) GetBlock(blockNumber int64) (core.Block, error) {
	panic("implement me")
}

func (mbr *MockBlockRepository) MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64, nodeId string) []int64 {
	mbr.missingBlockNumbersCalled = true
	mbr.missingBlockNumbersPassedStartingBlockNumber = startingBlockNumber
	mbr.missingBlockNumbersPassedEndingBlockNumber = endingBlockNumber
	mbr.missingBlockNumbersPassedNodeId = nodeId
	return mbr.missingBlockNumbersReturnArray
}

func (mbr *MockBlockRepository) SetBlocksStatus(chainHead int64) {
	mbr.setBlockStatusCalled = true
	mbr.setBlockStatusPassedChainHead = chainHead
}

func (mbr *MockBlockRepository) AssertCreateOrUpdateBlockCalledWith(block core.Block) {
	Expect(mbr.createOrUpdateBlockCalled).To(BeTrue())
	Expect(mbr.createOrUpdateBlockPassedBlock).To(Equal(block))
}

func (mbr *MockBlockRepository) AssertMissingBlockNumbersCalledWith(startingBlockNumber int64, endingBlockNumber int64, nodeId string) {
	Expect(mbr.missingBlockNumbersCalled).To(BeTrue())
	Expect(mbr.missingBlockNumbersPassedStartingBlockNumber).To(Equal(startingBlockNumber))
	Expect(mbr.missingBlockNumbersPassedEndingBlockNumber).To(Equal(endingBlockNumber))
	Expect(mbr.missingBlockNumbersPassedNodeId).To(Equal(nodeId))
}

func (mbr *MockBlockRepository) AssertSetBlockStatusCalledWith(chainHead int64) {
	Expect(mbr.setBlockStatusCalled).To(BeTrue())
	Expect(mbr.setBlockStatusPassedChainHead).To(Equal(chainHead))
}
