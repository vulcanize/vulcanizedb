// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package fakes

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockHeaderRepository struct {
	createOrUpdateBlockNumbersCallCount          int
	createOrUpdateBlockNumbersPassedBlockNumbers []int64
	missingBlockNumbers                          []int64
}

func NewMockHeaderRepository() *MockHeaderRepository {
	return &MockHeaderRepository{}
}

func (repository *MockHeaderRepository) SetMissingBlockNumbers(blockNumbers []int64) {
	repository.missingBlockNumbers = blockNumbers
}

func (repository *MockHeaderRepository) CreateOrUpdateHeader(header core.Header) (int64, error) {
	repository.createOrUpdateBlockNumbersCallCount++
	repository.createOrUpdateBlockNumbersPassedBlockNumbers = append(repository.createOrUpdateBlockNumbersPassedBlockNumbers, header.BlockNumber)
	return 0, nil
}

func (*MockHeaderRepository) GetHeader(blockNumber int64) (core.Header, error) {
	return core.Header{BlockNumber: blockNumber}, nil
}

func (repository *MockHeaderRepository) MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64, nodeID string) []int64 {
	return repository.missingBlockNumbers
}

func (repository *MockHeaderRepository) AssertCreateOrUpdateHeaderCallCountAndPassedBlockNumbers(times int, blockNumbers []int64) {
	Expect(repository.createOrUpdateBlockNumbersCallCount).To(Equal(times))
	Expect(repository.createOrUpdateBlockNumbersPassedBlockNumbers).To(Equal(blockNumbers))
}
