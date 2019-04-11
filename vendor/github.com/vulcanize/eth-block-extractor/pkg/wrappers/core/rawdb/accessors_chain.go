package rawdb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
)

type IAccessorsChain interface {
	GetBlock(hash common.Hash, number uint64) *types.Block
	GetBlockReceipts(hash common.Hash, number uint64) types.Receipts
	GetBody(hash common.Hash, number uint64) *types.Body
	GetCanonicalHash(number uint64) common.Hash
	GetHeader(hash common.Hash, number uint64) *types.Header
	GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue
}

type AccessorsChain struct {
	ethDbConnection ethdb.Database
}

func NewAccessorsChain(databaseConnection ethdb.Database) *AccessorsChain {
	return &AccessorsChain{ethDbConnection: databaseConnection}
}

func (accessor *AccessorsChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return rawdb.ReadBlock(accessor.ethDbConnection, hash, number)
}

func (accessor *AccessorsChain) GetBlockReceipts(hash common.Hash, number uint64) types.Receipts {
	return rawdb.ReadReceipts(accessor.ethDbConnection, hash, number)
}

func (accessor *AccessorsChain) GetBody(hash common.Hash, number uint64) *types.Body {
	return rawdb.ReadBody(accessor.ethDbConnection, hash, number)
}

func (accessor *AccessorsChain) GetCanonicalHash(number uint64) common.Hash {
	return rawdb.ReadCanonicalHash(accessor.ethDbConnection, number)
}

func (accessor *AccessorsChain) GetHeader(hash common.Hash, number uint64) *types.Header {
	return rawdb.ReadHeader(accessor.ethDbConnection, hash, number)
}

func (accessor *AccessorsChain) GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue {
	return rawdb.ReadHeaderRLP(accessor.ethDbConnection, hash, number)
}
