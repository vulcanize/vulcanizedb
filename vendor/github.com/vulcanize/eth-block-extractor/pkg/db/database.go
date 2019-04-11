package db

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/vulcanize/eth-block-extractor/pkg/db/level"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/core"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/core/rawdb"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/core/state"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/rlp"
)

var ErrNoSuchDb = errors.New("no such database")

type ReadError struct {
	msg string
	err error
}

func (re ReadError) Error() string {
	return fmt.Sprintf("%s: %s", re.msg, re.err.Error())
}

type Database interface {
	ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) (common.Hash, error)
	GetBlockByBlockNumber(blockNumber int64) *types.Block
	GetBlockBodyByBlockNumber(blockNumber int64) *types.Body
	GetBlockHeaderByBlockNumber(blockNumber int64) *types.Header
	GetRawBlockHeaderByBlockNumber(blockNumber int64) []byte
	GetBlockReceipts(blockNumber int64) types.Receipts
	GetStateAndStorageTrieNodes(root common.Hash) (stateTrieNodes, storageTrieNodes [][]byte, err error)
}

func CreateDatabase(config DatabaseConfig) (Database, error) {
	switch config.Type {
	case Level:
		levelDBConnection, err := ethdb.NewLDBDatabase(config.Path, 128, 1024)
		if err != nil {
			return nil, ReadError{msg: "Failed to connect to LevelDB", err: err}
		}
		stateDatabase := state.NewDatabase(levelDBConnection)
		stateTrieReader := createStateTrieReader(stateDatabase)
		levelDBReader := rawdb.NewAccessorsChain(levelDBConnection)
		stateComputer, err := createStateComputer(levelDBConnection, stateDatabase)
		if err != nil {
			return nil, err
		}
		levelDB := level.NewLevelDatabase(levelDBReader, stateComputer, stateTrieReader)
		return levelDB, nil
	default:
		return nil, ReadError{msg: "Unknown database not implemented", err: ErrNoSuchDb}
	}
}

func createStateTrieReader(stateDatabase state.GethStateDatabase) level.IStateTrieReader {
	decoder := rlp.RlpDecoder{}
	storageTrieReader := level.NewStorageTrieReader(stateDatabase, decoder)
	return level.NewStateTrieReader(stateDatabase, storageTrieReader)
}

func createStateComputer(databaseConnection ethdb.Database, stateDatabase state.GethStateDatabase) (level.IStateComputer, error) {
	blockChain, err := core.NewBlockChain(databaseConnection)
	if err != nil {
		return nil, err
	}
	processor := core.NewStateProcessor(*blockChain)
	trieFactory := state.NewStateDBFactory()
	validator := core.NewBlockValidator(*blockChain)
	computer := level.NewStateComputer(blockChain, stateDatabase, processor, trieFactory, validator)
	return computer, nil
}
