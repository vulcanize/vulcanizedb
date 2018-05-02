package fakes

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockTransactionConverter struct {
	convertTransactionsToCoreCalled             bool
	convertTransactionsToCorePassedBlock        *types.Block
	convertTransactionsToCoreReturnTransactions []core.Transaction
	convertTransactionsToCoreReturnError        error
}

func NewMockTransactionConverter() *MockTransactionConverter {
	return &MockTransactionConverter{
		convertTransactionsToCoreCalled:             false,
		convertTransactionsToCorePassedBlock:        nil,
		convertTransactionsToCoreReturnTransactions: nil,
		convertTransactionsToCoreReturnError:        nil,
	}
}

func (mtc *MockTransactionConverter) SetConvertTransactionsToCoreReturnVals(transactions []core.Transaction, err error) {
	mtc.convertTransactionsToCoreReturnTransactions = transactions
	mtc.convertTransactionsToCoreReturnError = err
}

func (mtc *MockTransactionConverter) ConvertTransactionsToCore(gethBlock *types.Block) ([]core.Transaction, error) {
	mtc.convertTransactionsToCoreCalled = true
	mtc.convertTransactionsToCorePassedBlock = gethBlock
	return mtc.convertTransactionsToCoreReturnTransactions, mtc.convertTransactionsToCoreReturnError
}

func (mtc *MockTransactionConverter) AssertConvertTransactionsToCoreCalledWith(gethBlock *types.Block) {
	Expect(mtc.convertTransactionsToCoreCalled).To(BeTrue())
	Expect(mtc.convertTransactionsToCorePassedBlock).To(Equal(gethBlock))
}
