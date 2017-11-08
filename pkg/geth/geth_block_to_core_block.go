package geth

import (
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/net/context"
)

type GethClient interface {
	TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error)
}

func GethBlockToCoreBlock(gethBlock *types.Block, client GethClient) core.Block {
	transactions := []core.Transaction{}
	for i, gethTransaction := range gethBlock.Transactions() {
		from, _ := client.TransactionSender(context.Background(), gethTransaction, gethBlock.Hash(), uint(i))
		transaction := gethTransToCoreTrans(gethTransaction, &from)
		transactions = append(transactions, transaction)
	}
	return core.Block{
		Difficulty:   gethBlock.Difficulty().Int64(),
		GasLimit:     gethBlock.GasLimit().Int64(),
		GasUsed:      gethBlock.GasUsed().Int64(),
		Hash:         gethBlock.Hash().Hex(),
		Nonce:        hexutil.Encode(gethBlock.Header().Nonce[:]),
		Number:       gethBlock.Number().Int64(),
		ParentHash:   gethBlock.ParentHash().Hex(),
		Size:         gethBlock.Size().Int64(),
		Time:         gethBlock.Time().Int64(),
		Transactions: transactions,
		UncleHash:    gethBlock.UncleHash().Hex(),
	}
}

func gethTransToCoreTrans(transaction *types.Transaction, from *common.Address) core.Transaction {
	return core.Transaction{
		Hash:     transaction.Hash().Hex(),
		Data:     transaction.Data(),
		Nonce:    transaction.Nonce(),
		To:       addressToHex(transaction.To()),
		From:     addressToHex(from),
		GasLimit: transaction.Gas().Int64(),
		GasPrice: transaction.GasPrice().Int64(),
		Value:    transaction.Value().Int64(),
	}
}

func addressToHex(to *common.Address) string {
	if to == nil {
		return ""
	} else {
		return to.Hex()
	}
}
