package fakes

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockBlockRepository struct {
	createOrUpdateBlockCallCount                 int
	createOrUpdateBlockCalled                    bool
	createOrUpdateBlockPassedBlock               core.Block
	createOrUpdateBlockPassedBlockNumbers        []int64
	createOrUpdateBlockReturnErr                 error
	createOrUpdateBlockReturnInt                 int64
	missingBlockNumbersCalled                    bool
	missingBlockNumbersPassedEndingBlockNumber   int64
	missingBlockNumbersPassedNodeId              string
	missingBlockNumbersPassedStartingBlockNumber int64
	missingBlockNumbersReturnArray               []int64
	setBlockStatusCalled                         bool
	setBlockStatusPassedChainHead                int64
}

func NewMockBlockRepository() *MockBlockRepository {
	return &MockBlockRepository{
		createOrUpdateBlockCallCount:                 0,
		createOrUpdateBlockCalled:                    false,
		createOrUpdateBlockPassedBlock:               core.Block{},
		createOrUpdateBlockPassedBlockNumbers:        nil,
		createOrUpdateBlockReturnErr:                 nil,
		createOrUpdateBlockReturnInt:                 0,
		missingBlockNumbersCalled:                    false,
		missingBlockNumbersPassedEndingBlockNumber:   0,
		missingBlockNumbersPassedNodeId:              "",
		missingBlockNumbersPassedStartingBlockNumber: 0,
		missingBlockNumbersReturnArray:               nil,
		setBlockStatusCalled:                         false,
		setBlockStatusPassedChainHead:                0,
	}
}

func (repository *MockBlockRepository) SetCreateOrUpdateBlockReturnVals(i int64, err error) {
	repository.createOrUpdateBlockReturnInt = i
	repository.createOrUpdateBlockReturnErr = err
}

func (repository *MockBlockRepository) SetMissingBlockNumbersReturnArray(returnArray []int64) {
	repository.missingBlockNumbersReturnArray = returnArray
}

func (repository *MockBlockRepository) CreateOrUpdateBlock(block core.Block) (int64, error) {
	repository.createOrUpdateBlockCallCount++
	repository.createOrUpdateBlockCalled = true
	repository.createOrUpdateBlockPassedBlock = block
	repository.createOrUpdateBlockPassedBlockNumbers = append(repository.createOrUpdateBlockPassedBlockNumbers, block.Number)
	return repository.createOrUpdateBlockReturnInt, repository.createOrUpdateBlockReturnErr
}

func (repository *MockBlockRepository) GetBlock(blockNumber int64) (core.Block, error) {
	return core.Block{Number: blockNumber}, nil
}

func (repository *MockBlockRepository) MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64, nodeId string) []int64 {
	repository.missingBlockNumbersCalled = true
	repository.missingBlockNumbersPassedStartingBlockNumber = startingBlockNumber
	repository.missingBlockNumbersPassedEndingBlockNumber = endingBlockNumber
	repository.missingBlockNumbersPassedNodeId = nodeId
	return repository.missingBlockNumbersReturnArray
}

func (repository *MockBlockRepository) SetBlocksStatus(chainHead int64) {
	repository.setBlockStatusCalled = true
	repository.setBlockStatusPassedChainHead = chainHead
}

func (repository *MockBlockRepository) AssertCreateOrUpdateBlockCallCountEquals(times int) {
	Expect(repository.createOrUpdateBlockCallCount).To(Equal(times))
}

func (repository *MockBlockRepository) AssertCreateOrUpdateBlocksCallCountAndBlockNumbersEquals(times int, blockNumbers []int64) {
	Expect(repository.createOrUpdateBlockCallCount).To(Equal(times))
	Expect(repository.createOrUpdateBlockPassedBlockNumbers).To(Equal(blockNumbers))
}

func (repository *MockBlockRepository) AssertCreateOrUpdateBlockCalledWith(block core.Block) {
	Expect(repository.createOrUpdateBlockCalled).To(BeTrue())
	Expect(repository.createOrUpdateBlockPassedBlock).To(Equal(block))
}

func (repository *MockBlockRepository) AssertMissingBlockNumbersCalledWith(startingBlockNumber int64, endingBlockNumber int64, nodeId string) {
	Expect(repository.missingBlockNumbersCalled).To(BeTrue())
	Expect(repository.missingBlockNumbersPassedStartingBlockNumber).To(Equal(startingBlockNumber))
	Expect(repository.missingBlockNumbersPassedEndingBlockNumber).To(Equal(endingBlockNumber))
	Expect(repository.missingBlockNumbersPassedNodeId).To(Equal(nodeId))
}

func (repository *MockBlockRepository) AssertSetBlockStatusCalledWith(chainHead int64) {
	Expect(repository.setBlockStatusCalled).To(BeTrue())
	Expect(repository.setBlockStatusPassedChainHead).To(Equal(chainHead))
}
