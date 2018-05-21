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
