package rpc

import (
	"context"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/sync/errgroup"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	vulcCommon "github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
)

type RpcTransactionConverter struct {
	client core.EthClient
}

func NewRpcTransactionConverter(client core.EthClient) *RpcTransactionConverter {
	return &RpcTransactionConverter{client: client}
}

func (rtc *RpcTransactionConverter) ConvertTransactionsToCore(gethBlock *types.Block) ([]core.Transaction, error) {
	var g errgroup.Group
	coreTransactions := make([]core.Transaction, len(gethBlock.Transactions()))

	for gethTransactionIndex, gethTransaction := range gethBlock.Transactions() {
		//https://golang.org/doc/faq#closures_and_goroutines
		transaction := gethTransaction
		transactionIndex := uint(gethTransactionIndex)
		g.Go(func() error {
			from, err := rtc.client.TransactionSender(context.Background(), transaction, gethBlock.Hash(), transactionIndex)
			if err != nil {
				log.Println("transaction sender: ", err)
				return err
			}
			coreTransaction := transToCoreTrans(transaction, &from)
			coreTransaction, err = rtc.appendReceiptToTransaction(coreTransaction)
			if err != nil {
				log.Println("receipt: ", err)
				return err
			}
			coreTransactions[transactionIndex] = coreTransaction
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Println("transactions: ", err)
		return coreTransactions, err
	}
	return coreTransactions, nil
}

func (rtc *RpcTransactionConverter) appendReceiptToTransaction(transaction core.Transaction) (core.Transaction, error) {
	gethReceipt, err := rtc.client.TransactionReceipt(context.Background(), common.HexToHash(transaction.Hash))
	if err != nil {
		return transaction, err
	}
	receipt := vulcCommon.ToCoreReceipt(gethReceipt)
	transaction.Receipt = receipt
	return transaction, nil
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
