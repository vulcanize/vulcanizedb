package backfill

import (
	"database/sql"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	storage2 "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/sirupsen/logrus"
)

func NewStorageValueLoader(bc core.BlockChain, db *postgres.DB, initializers []storage.TransformerInitializer, startingBlock, endingBlock int64) StorageValueLoader {
	return StorageValueLoader{
		bc:              bc,
		db:              db,
		HeaderRepo:      repositories.NewHeaderRepository(db),
		StorageDiffRepo: storage2.NewDiffRepository(db),
		initializers:    initializers,
		startingBlock:   startingBlock,
		endingBlock:     endingBlock,
	}
}

type StorageValueLoader struct {
	bc              core.BlockChain
	db              *postgres.DB
	HeaderRepo      datastore.HeaderRepository
	StorageDiffRepo storage2.DiffRepository
	initializers    []storage.TransformerInitializer
	startingBlock   int64
	endingBlock     int64
}

func (r *StorageValueLoader) Run() error {
	addressToKeys, getKeysErr := r.getStorageKeys()
	if getKeysErr != nil {
		return getKeysErr
	}
	headers, getHeadersErr := r.HeaderRepo.GetHeadersInRange(r.startingBlock, r.endingBlock)
	if getHeadersErr != nil {
		return getHeadersErr
	}

	for _, header := range headers {
		for address, keys := range addressToKeys {
			persistStorageErr := r.getAndPersistStorageValues(address, keys, header.BlockNumber, header.Hash)
			if persistStorageErr != nil {
				return persistStorageErr
			}
		}
	}
	logrus.Infof("Finished persisting storage values for %v addresses from block %v to %v.", len(addressToKeys), r.startingBlock, r.endingBlock)

	return nil
}

func (r *StorageValueLoader) getStorageKeys() (map[common.Address][]common.Hash, error) {
	addressToKeys := make(map[common.Address][]common.Hash, len(r.initializers))
	for _, i := range r.initializers {
		transformer := i(r.db)
		keysLookup := transformer.GetStorageKeysLookup()
		keys, getKeysErr := keysLookup.GetKeys()
		if getKeysErr != nil {
			return addressToKeys, getKeysErr
		}
		address := transformer.GetContractAddress()
		addressToKeys[address] = keys
		logrus.Infof("Received %v storage keys for address:%v", len(keys), address.Hex())
	}

	return addressToKeys, nil
}

func (r *StorageValueLoader) getAndPersistStorageValues(address common.Address, keys []common.Hash, blockNumber int64, headerHashStr string) error {
	blockNumberBigInt := big.NewInt(blockNumber)
	blockHash := common.HexToHash(headerHashStr)
	keccakOfAddress := crypto.Keccak256Hash(address[:])
	logrus.WithFields(logrus.Fields{
		"Address":       address.Hex(),
		"HashedAddress": keccakOfAddress.Hex(),
		"BlockNumber":   blockNumber,
	}).Infof("Getting and persisting %v storage values", len(keys))
	storageValues, getStorageValuesErr := r.bc.BatchGetStorageAt(address, keys, blockNumberBigInt)
	if getStorageValuesErr != nil {
		return getStorageValuesErr
	}
	for storageKey, storageValue := range storageValues {
		diff := types.RawDiff{
			HashedAddress: keccakOfAddress,
			BlockHash:     blockHash,
			BlockHeight:   int(blockNumber),
			StorageKey:    storageKey,
			StorageValue:  common.BytesToHash(storageValue),
		}
		createDiffErr := r.StorageDiffRepo.CreateBackFilledStorageValue(diff)
		if createDiffErr != nil {
			if createDiffErr == sql.ErrNoRows {
				return nil
			}
			return createDiffErr
		}
	}
	return nil
}
