// VulcanizeDB
// Copyright © 2019 Vulcanize

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
	createOrUpdateHeaderCallCount          int
	createOrUpdateHeaderErr                error
	createOrUpdateHeaderPassedBlockNumbers []int64
	createOrUpdateHeaderReturnID           int64
	CreateTransactionsCalled               bool
	CreateTransactionsError                error
	getHeaderError                         error
	getHeaderReturnBlockHash               string
	missingBlockNumbers                    []int64
	headerExists                           bool
	GetHeaderPassedBlockNumber             int64
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

func (repository *MockHeaderRepository) CreateTransactions(headerID int64, transactions []core.TransactionModel) error {
	repository.CreateTransactionsCalled = true
	return repository.CreateTransactionsError
}

func (repository *MockHeaderRepository) GetHeader(blockNumber int64) (core.Header, error) {
	repository.GetHeaderPassedBlockNumber = blockNumber
	return core.Header{BlockNumber: blockNumber, Hash: repository.getHeaderReturnBlockHash}, repository.getHeaderError
}

func (repository *MockHeaderRepository) MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64, nodeID string) ([]int64, error) {
	return repository.missingBlockNumbers, nil
}

func (repository *MockHeaderRepository) SetGetHeaderError(err error) {
	repository.getHeaderError = err
}

func (repository *MockHeaderRepository) SetGetHeaderReturnBlockHash(hash string) {
	repository.getHeaderReturnBlockHash = hash
}

func (repository *MockHeaderRepository) AssertCreateOrUpdateHeaderCallCountAndPassedBlockNumbers(times int, blockNumbers []int64) {
	Expect(repository.createOrUpdateHeaderCallCount).To(Equal(times))
	Expect(repository.createOrUpdateHeaderPassedBlockNumbers).To(Equal(blockNumbers))
}
