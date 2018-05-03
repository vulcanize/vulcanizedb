package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
)

type Reader interface {
	GetBlock(hash common.Hash, number uint64) *types.Block
	GetBlockNumber(hash common.Hash) uint64
	GetBlockReceipts(hash common.Hash, number uint64) types.Receipts
	GetCanonicalHash(number uint64) common.Hash
	GetHeadBlockHash() common.Hash
}

type LevelDatabaseReader struct {
	reader core.DatabaseReader
}

func NewLevelDatabaseReader(reader core.DatabaseReader) *LevelDatabaseReader {
	return &LevelDatabaseReader{reader: reader}
}

func (ldbr *LevelDatabaseReader) GetBlock(hash common.Hash, number uint64) *types.Block {
	return core.GetBlock(ldbr.reader, hash, number)
}

func (ldbr *LevelDatabaseReader) GetBlockNumber(hash common.Hash) uint64 {
	return core.GetBlockNumber(ldbr.reader, hash)
}

func (ldbr *LevelDatabaseReader) GetBlockReceipts(hash common.Hash, number uint64) types.Receipts {
	return core.GetBlockReceipts(ldbr.reader, hash, number)
}

func (ldbr *LevelDatabaseReader) GetCanonicalHash(number uint64) common.Hash {
	return core.GetCanonicalHash(ldbr.reader, number)
}

func (ldbr *LevelDatabaseReader) GetHeadBlockHash() common.Hash {
	return core.GetHeadBlockHash(ldbr.reader)
}
