package common

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"strings"
)

type BlockConverter struct {
	transactionConverter TransactionConverter
}

func NewBlockConverter(transactionConverter TransactionConverter) BlockConverter {
	return BlockConverter{transactionConverter: transactionConverter}
}

func (bc BlockConverter) ToCoreBlock(gethBlock *types.Block) (core.Block, error) {
	transactions, err := bc.transactionConverter.ConvertTransactionsToCore(gethBlock)
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
