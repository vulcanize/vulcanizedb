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

type MockReceiptRepository struct {
	createReceiptsAndLogsCalled         bool
	createReceiptsAndLogsPassedBlockId  int64
	createReceiptsAndLogsPassedReceipts []core.Receipt
	createReceiptsAndLogsReturnErr      error
}

func NewMockReceiptRepository() *MockReceiptRepository {
	return &MockReceiptRepository{
		createReceiptsAndLogsCalled:         false,
		createReceiptsAndLogsPassedBlockId:  0,
		createReceiptsAndLogsPassedReceipts: nil,
		createReceiptsAndLogsReturnErr:      nil,
	}
}

func (mrr *MockReceiptRepository) SetCreateReceiptsAndLogsReturnErr(err error) {
	mrr.createReceiptsAndLogsReturnErr = err
}

func (mrr *MockReceiptRepository) CreateReceiptsAndLogs(blockId int64, receipts []core.Receipt) error {
	mrr.createReceiptsAndLogsCalled = true
	mrr.createReceiptsAndLogsPassedBlockId = blockId
	mrr.createReceiptsAndLogsPassedReceipts = receipts
	return mrr.createReceiptsAndLogsReturnErr
}

func (mrr *MockReceiptRepository) CreateReceipt(blockId int64, receipt core.Receipt) (int64, error) {
	panic("implement me")
}

func (mrr *MockReceiptRepository) GetReceipt(txHash string) (core.Receipt, error) {
	panic("implement me")
}

func (mrr *MockReceiptRepository) AssertCreateReceiptsAndLogsCalledWith(blockId int64, receipts []core.Receipt) {
	Expect(mrr.createReceiptsAndLogsCalled).To(BeTrue())
	Expect(mrr.createReceiptsAndLogsPassedBlockId).To(Equal(blockId))
	Expect(mrr.createReceiptsAndLogsPassedReceipts).To(Equal(receipts))
}

func (mrr *MockReceiptRepository) AssertCreateReceiptsAndLogsNotCalled() {
	Expect(mrr.createReceiptsAndLogsCalled).To(BeFalse())
}
