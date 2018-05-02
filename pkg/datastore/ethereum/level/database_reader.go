package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
)

type Reader interface {
	GetBlock(hash common.Hash, number uint64) *types.Block
	GetBlockReceipts(hash common.Hash, number uint64) types.Receipts
	GetCanonicalHash(number uint64) common.Hash
}

type LevelDatabaseReader struct {
	core.DatabaseReader
}

func NewLevelDatabaseReader(reader core.DatabaseReader) *LevelDatabaseReader {
	return &LevelDatabaseReader{DatabaseReader: reader}
}

func (ldbr *LevelDatabaseReader) GetBlock(hash common.Hash, number uint64) *types.Block {
	return core.GetBlock(ldbr.DatabaseReader, hash, number)
}

func (ldbr *LevelDatabaseReader) GetBlockReceipts(hash common.Hash, number uint64) types.Receipts {
	return core.GetBlockReceipts(ldbr.DatabaseReader, hash, number)
}

func (ldbr *LevelDatabaseReader) GetCanonicalHash(number uint64) common.Hash {
	return core.GetCanonicalHash(ldbr.DatabaseReader, number)
}
