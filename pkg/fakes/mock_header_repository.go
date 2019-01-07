package fakes

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockHeaderRepository struct {
	createOrUpdateHeaderCallCount          int
	createOrUpdateHeaderErr                error
	createOrUpdateHeaderPassedBlockNumbers []int64
	createOrUpdateHeaderReturnID           int64
	missingBlockNumbers                    []int64
	headerExists                           bool
}

func NewMockHeaderRepository() *MockHeaderRepository {
	return &MockHeaderRepository{}
}

func (repository *MockHeaderRepository) SetCreateOrUpdateHeaderReturnID(id int64) {
	repository.createOrUpdateHeaderReturnID = id
}

func (repository *MockHeaderRepository) SetCreateOrUpdateHeaderReturnErr(err error) {
	repository.createOrUpdateHeaderErr = err
}

func (repository *MockHeaderRepository) SetMissingBlockNumbers(blockNumbers []int64) {
	repository.missingBlockNumbers = blockNumbers
}

func (repository *MockHeaderRepository) CreateOrUpdateHeader(header core.Header) (int64, error) {
	repository.createOrUpdateHeaderCallCount++
	repository.createOrUpdateHeaderPassedBlockNumbers = append(repository.createOrUpdateHeaderPassedBlockNumbers, header.BlockNumber)
	return repository.createOrUpdateHeaderReturnID, repository.createOrUpdateHeaderErr
}

func (*MockHeaderRepository) GetHeader(blockNumber int64) (core.Header, error) {
	return core.Header{BlockNumber: blockNumber}, nil
}

func (repository *MockHeaderRepository) MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64, nodeID string) []int64 {
	return repository.missingBlockNumbers
}

func (repository *MockHeaderRepository) HeaderExists(blockNumber int64) (bool, error) {
	return repository.headerExists, nil
}

func (repository *MockHeaderRepository) SetHeaderExists(headerExists bool) {
	repository.headerExists = headerExists
}

func (repository *MockHeaderRepository) AssertCreateOrUpdateHeaderCallCountAndPassedBlockNumbers(times int, blockNumbers []int64) {
	Expect(repository.createOrUpdateHeaderCallCount).To(Equal(times))
	Expect(repository.createOrUpdateHeaderPassedBlockNumbers).To(Equal(blockNumbers))
}
