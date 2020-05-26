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
	"github.com/jmoiron/sqlx"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/gomega"
)

type MockHeaderRepository struct {
	createOrUpdateHeaderCallCount          int
	createOrUpdateHeaderErr                error
	createOrUpdateHeaderPassedBlockNumbers []int64
	createOrUpdateHeaderReturnID           int64
	AllHeaders                             []core.Header
	CreateTransactionsCalled               bool
	CreateTransactionsError                error
	GetHeaderByBlockNumberError            error
	GetHeaderByBlockNumberReturnHash       string
	GetHeaderByBlockNumberReturnID         int64
	GetHeaderByIDError                     error
	GetHeaderByIDHeaderToReturn            core.Header
	missingBlockNumbers                    []int64
	headerExists                           bool
	GetHeaderPassedBlockNumber             int64
	GetHeadersInRangeStartingBlock         int64
	GetHeadersInRangeEndingBlock           int64
	MostRecentHeaderBlockNumber            int64
	MostRecentHeaderBlockNumberErr         error
}

func NewMockHeaderRepository() *MockHeaderRepository {
	return &MockHeaderRepository{}
}

func (mock *MockHeaderRepository) SetCreateOrUpdateHeaderReturnID(id int64) {
	mock.createOrUpdateHeaderReturnID = id
}

func (mock *MockHeaderRepository) SetCreateOrUpdateHeaderReturnErr(err error) {
	mock.createOrUpdateHeaderErr = err
}

func (mock *MockHeaderRepository) SetMissingBlockNumbers(blockNumbers []int64) {
	mock.missingBlockNumbers = blockNumbers
}

func (mock *MockHeaderRepository) CreateOrUpdateHeader(header core.Header) (int64, error) {
	mock.createOrUpdateHeaderCallCount++
	mock.createOrUpdateHeaderPassedBlockNumbers = append(mock.createOrUpdateHeaderPassedBlockNumbers, header.BlockNumber)
	return mock.createOrUpdateHeaderReturnID, mock.createOrUpdateHeaderErr
}

func (mock *MockHeaderRepository) CreateTransactions(headerID int64, transactions []core.TransactionModel) error {
	mock.CreateTransactionsCalled = true
	return mock.CreateTransactionsError
}

func (mock *MockHeaderRepository) CreateTransactionInTx(tx *sqlx.Tx, headerID int64, transaction core.TransactionModel) (int64, error) {
	panic("implement me")
}

func (mock *MockHeaderRepository) GetHeaderByBlockNumber(blockNumber int64) (core.Header, error) {
	mock.GetHeaderPassedBlockNumber = blockNumber
	return core.Header{
		Id:          mock.GetHeaderByBlockNumberReturnID,
		BlockNumber: blockNumber,
		Hash:        mock.GetHeaderByBlockNumberReturnHash,
	}, mock.GetHeaderByBlockNumberError
}

func (mock *MockHeaderRepository) GetHeaderByID(id int64) (core.Header, error) {
	return mock.GetHeaderByIDHeaderToReturn, mock.GetHeaderByIDError
}

func (mock *MockHeaderRepository) GetHeadersInRange(startingBlock, endingBlock int64) ([]core.Header, error) {
	mock.GetHeadersInRangeStartingBlock = startingBlock
	mock.GetHeadersInRangeEndingBlock = endingBlock
	return mock.AllHeaders, mock.GetHeaderByBlockNumberError
}

func (mock *MockHeaderRepository) MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64) ([]int64, error) {
	return mock.missingBlockNumbers, nil
}

func (mock *MockHeaderRepository) GetMostRecentHeaderBlockNumber() (int64, error) {
	return mock.MostRecentHeaderBlockNumber, mock.MostRecentHeaderBlockNumberErr
}

func (mock *MockHeaderRepository) AssertCreateOrUpdateHeaderCallCountAndPassedBlockNumbers(times int, blockNumbers []int64) {
	Expect(mock.createOrUpdateHeaderCallCount).To(Equal(times))
	Expect(mock.createOrUpdateHeaderPassedBlockNumbers).To(Equal(blockNumbers))
}
