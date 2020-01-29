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
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
)

type MockReceiptRepository struct {
	createReceiptsAndLogsCalled         bool
	createReceiptsAndLogsPassedBlockID  int64
	createReceiptsAndLogsPassedReceipts []core.Receipt
	createReceiptsAndLogsReturnErr      error
}

func NewMockReceiptRepository() *MockReceiptRepository {
	return &MockReceiptRepository{
		createReceiptsAndLogsCalled:         false,
		createReceiptsAndLogsPassedBlockID:  0,
		createReceiptsAndLogsPassedReceipts: nil,
		createReceiptsAndLogsReturnErr:      nil,
	}
}

func (mrr *MockReceiptRepository) SetCreateReceiptsAndLogsReturnErr(err error) {
	mrr.createReceiptsAndLogsReturnErr = err
}

func (mrr *MockReceiptRepository) CreateReceiptsAndLogs(blockID int64, receipts []core.Receipt) error {
	mrr.createReceiptsAndLogsCalled = true
	mrr.createReceiptsAndLogsPassedBlockID = blockID
	mrr.createReceiptsAndLogsPassedReceipts = receipts
	return mrr.createReceiptsAndLogsReturnErr
}

func (mrr *MockReceiptRepository) CreateFullSyncReceiptInTx(blockID int64, receipt core.Receipt, tx *sqlx.Tx) (int64, error) {
	panic("implement me")
}

func (mrr *MockReceiptRepository) GetFullSyncReceipt(txHash string) (core.Receipt, error) {
	panic("implement me")
}

func (mrr *MockReceiptRepository) AssertCreateReceiptsAndLogsCalledWith(blockID int64, receipts []core.Receipt) {
	Expect(mrr.createReceiptsAndLogsCalled).To(BeTrue())
	Expect(mrr.createReceiptsAndLogsPassedBlockID).To(Equal(blockID))
	Expect(mrr.createReceiptsAndLogsPassedReceipts).To(Equal(receipts))
}

func (mrr *MockReceiptRepository) AssertCreateReceiptsAndLogsNotCalled() {
	Expect(mrr.createReceiptsAndLogsCalled).To(BeFalse())
}
