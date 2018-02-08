package inmemory

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

func (repository *InMemory) GetReceipt(txHash string) (core.Receipt, error) {
	if receipt, ok := repository.receipts[txHash]; ok {
		return receipt, nil
	}
	return core.Receipt{}, repositories.ErrReceiptDoesNotExist(txHash)
}
