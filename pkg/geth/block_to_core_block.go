package geth

import (
	"strings"

	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"golang.org/x/net/context"
)

type GethClient interface {
	TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

func ToCoreBlock(gethBlock *types.Block, client GethClient) core.Block {
	transactions := convertTransactionsToCore(gethBlock, client)
	coreBlock := core.Block{
		Difficulty:   gethBlock.Difficulty().Int64(),
		ExtraData:    hexutil.Encode(gethBlock.Extra()),
		GasLimit:     gethBlock.GasLimit().Int64(),
		GasUsed:      gethBlock.GasUsed().Int64(),
		Hash:         gethBlock.Hash().Hex(),
		Miner:        strings.ToLower(gethBlock.Coinbase().Hex()),
		Nonce:        hexutil.Encode(gethBlock.Header().Nonce[:]),
		Number:       gethBlock.Number().Int64(),
		ParentHash:   gethBlock.ParentHash().Hex(),
		Size:         gethBlock.Size().Int64(),
		Time:         gethBlock.Time().Int64(),
		Transactions: transactions,
		UncleHash:    gethBlock.UncleHash().Hex(),
	}
	coreBlock.Reward = CalcBlockReward(coreBlock, gethBlock.Uncles())
	coreBlock.UnclesReward = CalcUnclesReward(coreBlock, gethBlock.Uncles())
	return coreBlock
}

func convertTransactionsToCore(gethBlock *types.Block, client GethClient) []core.Transaction {
	transactions := make([]core.Transaction, 0)
	for i, gethTransaction := range gethBlock.Transactions() {
		from, err := client.TransactionSender(context.Background(), gethTransaction, gethBlock.Hash(), uint(i))
		if err != nil {
			log.Println(err)
		}
		transaction := transToCoreTrans(gethTransaction, &from)
		transaction, err = appendReceiptToTransaction(client, transaction)
		if err != nil {
			log.Println(err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}

func appendReceiptToTransaction(client GethClient, transaction core.Transaction) (core.Transaction, error) {
	gethReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(transaction.Hash))
	if err != nil {
		log.Println(err)
		return transaction, err
	}
	receipt := ReceiptToCoreReceipt(gethReceipt)
	transaction.Receipt = receipt
	return transaction, err
}

func transToCoreTrans(transaction *types.Transaction, from *common.Address) core.Transaction {
	data := hexutil.Encode(transaction.Data())
	return core.Transaction{
		Hash:     transaction.Hash().Hex(),
		Nonce:    transaction.Nonce(),
		To:       strings.ToLower(addressToHex(transaction.To())),
		From:     strings.ToLower(addressToHex(from)),
		GasLimit: transaction.Gas().Int64(),
		GasPrice: transaction.GasPrice().Int64(),
		Value:    transaction.Value().String(),
		Data:     data,
	}
}

func addressToHex(to *common.Address) string {
	if to == nil {
		return ""
	} else {
		return to.Hex()
	}
}
