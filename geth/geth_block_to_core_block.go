package geth

import (
	"github.com/8thlight/vulcanizedb/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func gethTransToCoreTrans(transaction *types.Transaction) core.Transaction {
	to := transaction.To()
	toHex := convertTo(to)
	return core.Transaction{
		Hash:     transaction.Hash().Hex(),
		Data:     transaction.Data(),
		Nonce:    transaction.Nonce(),
		To:       toHex,
		GasLimit: transaction.Gas().Int64(),
		GasPrice: transaction.GasPrice().Int64(),
		Value:    transaction.Value().Int64(),
	}
}

func GethBlockToCoreBlock(gethBlock *types.Block) core.Block {
	transactions := []core.Transaction{}
	for _, gethTransaction := range gethBlock.Transactions() {
		transactions = append(transactions, gethTransToCoreTrans(gethTransaction))
	}
	return core.Block{
		Number:       gethBlock.Number().Int64(),
		GasLimit:     gethBlock.GasLimit().Int64(),
		GasUsed:      gethBlock.GasUsed().Int64(),
		Time:         gethBlock.Time().Int64(),
		Transactions: transactions,
	}
}

func convertTo(to *common.Address) string {
	if to == nil {
		return ""
	} else {
		return to.Hex()
	}
}
