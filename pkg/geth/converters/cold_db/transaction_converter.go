package cold_db

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"golang.org/x/sync/errgroup"
	"strings"
)

type ColdDbTransactionConverter struct{}

func NewColdDbTransactionConverter() *ColdDbTransactionConverter {
	return &ColdDbTransactionConverter{}
}

func (cdtc *ColdDbTransactionConverter) ConvertTransactionsToCore(gethBlock *types.Block) ([]core.Transaction, error) {
	var g errgroup.Group
	coreTransactions := make([]core.Transaction, len(gethBlock.Transactions()))

	for gethTransactionIndex, gethTransaction := range gethBlock.Transactions() {
		transaction := gethTransaction
		transactionIndex := uint(gethTransactionIndex)
		g.Go(func() error {
			signer := getSigner(transaction)
			sender, err := signer.Sender(transaction)
			if err != nil {
				return err
			}
			coreTransaction := transToCoreTrans(transaction, &sender)
			coreTransactions[transactionIndex] = coreTransaction
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return coreTransactions, err
	}
	return coreTransactions, nil
}

func getSigner(tx *types.Transaction) types.Signer {
	v, _, _ := tx.RawSignatureValues()
	if v.Sign() != 0 && tx.Protected() {
		return types.NewEIP155Signer(tx.ChainId())
	}
	return types.HomesteadSigner{}
}

func transToCoreTrans(transaction *types.Transaction, from *common.Address) core.Transaction {
	data := hexutil.Encode(transaction.Data())
	return core.Transaction{
		Hash:     transaction.Hash().Hex(),
		Nonce:    transaction.Nonce(),
		To:       strings.ToLower(addressToHex(transaction.To())),
		From:     strings.ToLower(addressToHex(from)),
		GasLimit: transaction.Gas(),
		GasPrice: transaction.GasPrice().Int64(),
		Value:    transaction.Value().String(),
		Data:     data,
	}
}

func addressToHex(to *common.Address) string {
	if to == nil {
		return ""
	}
	return to.Hex()
}
