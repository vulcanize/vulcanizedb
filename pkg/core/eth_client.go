package core

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EthClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}
