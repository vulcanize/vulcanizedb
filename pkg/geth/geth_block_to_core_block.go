package geth

import (
	"strings"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/net/context"
)

type GethClient interface {
	TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

func GethBlockToCoreBlock(gethBlock *types.Block, client GethClient) core.Block {
	transactions := []core.Transaction{}
	for i, gethTransaction := range gethBlock.Transactions() {
		from, _ := client.TransactionSender(context.Background(), gethTransaction, gethBlock.Hash(), uint(i))
		transaction := gethTransToCoreTrans(gethTransaction, &from)
		transactions = append(transactions, transaction)
	}
	blockReward := CalcBlockReward(gethBlock, client)
	uncleReward := CalcUnclesReward(gethBlock)
	return core.Block{
		Difficulty:   gethBlock.Difficulty().Int64(),
		ExtraData:    hexutil.Encode(gethBlock.Extra()),
		GasLimit:     gethBlock.GasLimit().Int64(),
		GasUsed:      gethBlock.GasUsed().Int64(),
		Hash:         gethBlock.Hash().Hex(),
		Miner:        gethBlock.Coinbase().Hex(),
		Nonce:        hexutil.Encode(gethBlock.Header().Nonce[:]),
		Number:       gethBlock.Number().Int64(),
		ParentHash:   gethBlock.ParentHash().Hex(),
		Reward:       blockReward,
		Size:         gethBlock.Size().Int64(),
		Time:         gethBlock.Time().Int64(),
		Transactions: transactions,
		UncleHash:    gethBlock.UncleHash().Hex(),
		UnclesReward: uncleReward,
	}
}

func gethTransToCoreTrans(transaction *types.Transaction, from *common.Address) core.Transaction {
	data := hexutil.Encode(transaction.Data())
	return core.Transaction{
		Hash:     transaction.Hash().Hex(),
		Nonce:    transaction.Nonce(),
		To:       strings.ToLower(addressToHex(transaction.To())),
		From:     strings.ToLower(addressToHex(from)),
		GasLimit: transaction.Gas().Int64(),
		GasPrice: transaction.GasPrice().Int64(),
		Value:    transaction.Value().Int64(),
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
