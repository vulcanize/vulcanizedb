package fakes

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockBlockRepository struct {
	createOrUpdateBlockCalled      bool
	createOrUpdateBlockPassedBlock core.Block
	createOrUpdateBlockReturnInt   int64
	createOrUpdateBlockReturnErr   error
}

func NewMockBlockRepository() *MockBlockRepository {
	return &MockBlockRepository{
		createOrUpdateBlockCalled:      false,
		createOrUpdateBlockPassedBlock: core.Block{},
		createOrUpdateBlockReturnInt:   0,
		createOrUpdateBlockReturnErr:   nil,
	}
}

func (mbr *MockBlockRepository) SetCreateOrUpdateBlockReturnVals(i int64, err error) {
	mbr.createOrUpdateBlockReturnInt = i
	mbr.createOrUpdateBlockReturnErr = err
}

func (mbr *MockBlockRepository) CreateOrUpdateBlock(block core.Block) (int64, error) {
	mbr.createOrUpdateBlockCalled = true
	mbr.createOrUpdateBlockPassedBlock = block
	return mbr.createOrUpdateBlockReturnInt, mbr.createOrUpdateBlockReturnErr
}

func (mbr *MockBlockRepository) GetBlock(blockNumber int64) (core.Block, error) {
	panic("implement me")
}

func (mbr *MockBlockRepository) MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64) []int64 {
	panic("implement me")
}

func (mbr *MockBlockRepository) SetBlocksStatus(chainHead int64) {
	panic("implement me")
}

func (mbr *MockBlockRepository) AssertCreateOrUpdateBlockCalledWith(block core.Block) {
	Expect(mbr.createOrUpdateBlockCalled).To(BeTrue())
	Expect(mbr.createOrUpdateBlockPassedBlock).To(Equal(block))
}
