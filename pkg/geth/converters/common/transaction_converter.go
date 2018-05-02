package common

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type TransactionConverter interface {
	ConvertTransactionsToCore(gethBlock *types.Block) ([]core.Transaction, error)
}
