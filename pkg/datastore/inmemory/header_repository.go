package inmemory

import "github.com/vulcanize/vulcanizedb/pkg/core"

type HeaderRepository struct {
	memory *InMemory
}

func NewHeaderRepository(memory *InMemory) *HeaderRepository {
	return &HeaderRepository{memory: memory}
}

func (repository *HeaderRepository) CreateOrUpdateHeader(header core.Header) (int64, error) {
	repository.memory.headers[header.BlockNumber] = header
	return 0, nil
}

func (repository *HeaderRepository) GetHeader(blockNumber int64) (core.Header, error) {
	return repository.memory.headers[blockNumber], nil
}

func (repository *HeaderRepository) MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64, nodeID string) []int64 {
	missingNumbers := []int64{}
	for blockNumber := int64(startingBlockNumber); blockNumber <= endingBlockNumber; blockNumber++ {
		if _, ok := repository.memory.headers[blockNumber]; !ok {
			missingNumbers = append(missingNumbers, blockNumber)
		}
	}
	return missingNumbers
}
