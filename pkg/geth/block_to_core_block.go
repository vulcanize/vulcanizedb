package geth

import (
	"strings"

	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type Client interface {
	TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

func ToCoreBlock(gethBlock *types.Block, client Client) (core.Block, error) {
	transactions, err := convertTransactionsToCore(gethBlock, client)
	if err != nil {
		return core.Block{}, err
	}
	coreBlock := core.Block{
		Difficulty:   gethBlock.Difficulty().Int64(),
		ExtraData:    hexutil.Encode(gethBlock.Extra()),
		GasLimit:     gethBlock.GasLimit(),
		GasUsed:      gethBlock.GasUsed(),
		Hash:         gethBlock.Hash().Hex(),
		Miner:        strings.ToLower(gethBlock.Coinbase().Hex()),
		Nonce:        hexutil.Encode(gethBlock.Header().Nonce[:]),
		Number:       gethBlock.Number().Int64(),
		ParentHash:   gethBlock.ParentHash().Hex(),
		Size:         gethBlock.Size().String(),
		Time:         gethBlock.Time().Int64(),
		Transactions: transactions,
		UncleHash:    gethBlock.UncleHash().Hex(),
	}
	coreBlock.Reward = CalcBlockReward(coreBlock, gethBlock.Uncles())
	coreBlock.UnclesReward = CalcUnclesReward(coreBlock, gethBlock.Uncles())
	return coreBlock, nil
}

func convertTransactionsToCore(gethBlock *types.Block, client Client) ([]core.Transaction, error) {
	var g errgroup.Group
	coreTransactions := make([]core.Transaction, len(gethBlock.Transactions()))

	for gethTransactionIndex, gethTransaction := range gethBlock.Transactions() {
		//https://golang.org/doc/faq#closures_and_goroutines
		transaction := gethTransaction
		transactionIndex := uint(gethTransactionIndex)
		g.Go(func() error {
			from, err := client.TransactionSender(context.Background(), transaction, gethBlock.Hash(), transactionIndex)
			if err != nil {
				log.Println("transaction sender: ", err)
				return err
			}
			coreTransaction := transToCoreTrans(transaction, &from)
			coreTransaction, err = appendReceiptToTransaction(client, coreTransaction)
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

func appendReceiptToTransaction(client Client, transaction core.Transaction) (core.Transaction, error) {
	gethReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(transaction.Hash))
	if err != nil {
		return transaction, err
	}
	receipt := ReceiptToCoreReceipt(gethReceipt)
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
