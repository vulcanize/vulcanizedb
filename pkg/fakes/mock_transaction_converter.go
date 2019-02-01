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
