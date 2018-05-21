package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type LevelDatabase struct {
	reader Reader
}

func NewLevelDatabase(ldbReader Reader) *LevelDatabase {
	return &LevelDatabase{
		reader: ldbReader,
	}
}

func (l LevelDatabase) GetBlock(blockHash []byte, blockNumber int64) *types.Block {
	n := uint64(blockNumber)
	h := common.BytesToHash(blockHash)
	return l.reader.GetBlock(h, n)
}

func (l LevelDatabase) GetBlockHash(blockNumber int64) []byte {
	n := uint64(blockNumber)
	h := l.reader.GetCanonicalHash(n)
	return h.Bytes()
}

func (l LevelDatabase) GetBlockReceipts(blockHash []byte, blockNumber int64) types.Receipts {
	n := uint64(blockNumber)
	h := common.BytesToHash(blockHash)
	return l.reader.GetBlockReceipts(h, n)
}

func (l LevelDatabase) GetHeadBlockNumber() int64 {
	h := l.reader.GetHeadBlockHash()
	n := l.reader.GetBlockNumber(h)
	return int64(n)
}
