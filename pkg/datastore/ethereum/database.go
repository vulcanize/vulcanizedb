package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/ethereum/level"
)

type Database interface {
	GetBlock(hash []byte, blockNumber int64) *types.Block
	GetBlockHash(blockNumber int64) []byte
	GetBlockReceipts(blockHash []byte, blockNumber int64) types.Receipts
	GetHeadBlockNumber() int64
}

func CreateDatabase(config DatabaseConfig) (Database, error) {
	switch config.Type {
	case Level:
		levelDBConnection, err := ethdb.NewLDBDatabase(config.Path, 128, 1024)
		if err != nil {
			return nil, err
		}
		levelDBReader := level.NewLevelDatabaseReader(levelDBConnection)
		levelDB := level.NewLevelDatabase(levelDBReader)
		return levelDB, nil
	default:
		return nil, fmt.Errorf("Unknown ethereum database: %s", config.Path)
	}
}
