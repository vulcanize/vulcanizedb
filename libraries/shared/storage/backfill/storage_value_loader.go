package backfill

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	storage2 "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/sirupsen/logrus"
)

var (
	MaxRequestSize    = 400
	ErrNoTransformers = errors.New("storage value loader initialized without transformers")
	emptyStorageValue = common.BytesToHash([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
)

func NewStorageValueLoader(bc core.BlockChain, db *postgres.DB, initializers []storage.TransformerInitializer, startingBlock, endingBlock int64) StorageValueLoader {
	return StorageValueLoader{
		bc:               bc,
		db:               db,
		HeaderRepo:       repositories.NewHeaderRepository(db),
		StorageDiffRepo:  storage2.NewDiffRepository(db),
		storageByAddress: make(map[common.Address]chunksOfKeysToValues, len(initializers)),
		initializers:     initializers,
		startingBlock:    startingBlock,
		endingBlock:      endingBlock,
	}
}

type storageKey = common.Hash
type storageValue = common.Hash
type chunksOfKeysToValues = []map[storageKey]storageValue

type StorageValueLoader struct {
	bc               core.BlockChain
	db               *postgres.DB
	HeaderRepo       datastore.HeaderRepository
	StorageDiffRepo  storage2.DiffRepository
	storageByAddress map[common.Address]chunksOfKeysToValues
	initializers     []storage.TransformerInitializer
	startingBlock    int64
	endingBlock      int64
}

func (r *StorageValueLoader) Run() error {
	if r.storageByAddress == nil {
		return ErrNoTransformers
	}
	getKeysErr := r.addKeysToStorageByAddress()
	if getKeysErr != nil {
		return getKeysErr
	}
	headers, getHeadersErr := r.HeaderRepo.GetHeadersInRange(r.startingBlock, r.endingBlock)
	if getHeadersErr != nil {
		return getHeadersErr
	}

	for _, header := range headers {
		persistStorageErr := r.getAndPersistStorageValues(header.BlockNumber, header.Hash)
		if persistStorageErr != nil {
			return persistStorageErr
		}

	}
	logrus.Infof("Finished persisting storage values for %v addresses from block %v to %v.", len(r.storageByAddress), r.startingBlock, r.endingBlock)

	return nil
}

func (r *StorageValueLoader) addKeysToStorageByAddress() error {
	for _, i := range r.initializers {
		transformer := i(r.db)
		keysLookup := transformer.GetStorageKeysLookup()
		keys, getKeysErr := keysLookup.GetKeys()
		if getKeysErr != nil {
			return getKeysErr
		}
		address := transformer.GetContractAddress()
		chunkedKeys := chunkKeys(keys)
		for chunkIndex, chunk := range chunkedKeys {
			nextChunkOfKeysToValues := make(map[storageKey]storageValue, len(chunk))
			r.storageByAddress[address] = append(r.storageByAddress[address], nextChunkOfKeysToValues)
			for _, key := range chunk {
				// set default initial value to empty
				r.storageByAddress[address][chunkIndex][key] = emptyStorageValue
			}
		}
		logrus.Infof("Received %v storage keys for address:%v", len(keys), address.Hex())
	}

	return nil
}

func (r *StorageValueLoader) getAndPersistStorageValues(blockNumber int64, headerHashStr string) error {
	blockNumberBigInt := big.NewInt(blockNumber)
	blockHash := common.HexToHash(headerHashStr)

	for address, chunkedKeysToValues := range r.storageByAddress {
		for chunkIndex, currentKeysToValues := range chunkedKeysToValues {
			var keys []storageKey
			for key, _ := range currentKeysToValues {
				keys = append(keys, key)
			}
			logrus.WithFields(logrus.Fields{
				"Address":     address.Hex(),
				"BlockNumber": blockNumber,
			}).Infof("Getting and persisting %v storage values", len(keys))
			newKeysToValues, getStorageValuesErr := r.bc.BatchGetStorageAt(address, keys, blockNumberBigInt)
			if getStorageValuesErr != nil {
				return getStorageValuesErr
			}
			for key, newValue := range newKeysToValues {
				newValueHash := common.BytesToHash(newValue)
				// don't attempt insert if new value matches last known value
				if newValueHash != currentKeysToValues[key] {
					// update last known value to new value if changed
					diff := types.RawDiff{
						Address:      address,
						BlockHash:    blockHash,
						BlockHeight:  int(blockNumber),
						StorageKey:   key,
						StorageValue: newValueHash,
					}
					createDiffErr := r.StorageDiffRepo.CreateBackFilledStorageValue(diff)
					if createDiffErr != nil {
						return createDiffErr
					}
					r.storageByAddress[address][chunkIndex][key] = newValueHash
				}
			}
		}
	}

	return nil
}

func chunkKeys(keys []storageKey) [][]storageKey {
	result := make([][]storageKey, getNumberOfChunks(keys))
	for index, key := range keys {
		resultIndex := getChunkIndex(index)
		result[resultIndex] = append(result[resultIndex], key)
	}
	return result
}

func getNumberOfChunks(keys []storageKey) int {
	keysLength := len(keys)
	if keysLength%MaxRequestSize == 0 {
		return keysLength / MaxRequestSize
	}
	return keysLength/MaxRequestSize + 1
}

func getChunkIndex(index int) int {
	return index / MaxRequestSize
}
