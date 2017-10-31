package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Transaction struct {
	Hash     string
	Data     []byte
	Nonce    uint64
	To       string
	GasLimit int64
	GasPrice int64
	Value    int64
}

func gethTransToCoreTrans(transaction *types.Transaction) Transaction {
	to := transaction.To()
	toHex := convertTo(to)
	return Transaction{
		Hash:     transaction.Hash().Hex(),
		Data:     transaction.Data(),
		Nonce:    transaction.Nonce(),
		To:       toHex,
		GasLimit: transaction.Gas().Int64(),
		GasPrice: transaction.GasPrice().Int64(),
		Value:    transaction.Value().Int64(),
	}
}

func convertTo(to *common.Address) string {
	if to == nil {
		return ""
	} else {
		return to.Hex()
	}
}
