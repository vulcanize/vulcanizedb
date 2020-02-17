package backfill

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	storage2 "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/sirupsen/logrus"
)

func NewStorageValueLoader(bc core.BlockChain, db *postgres.DB, initializers []transformer.StorageTransformerInitializer, blockNumber int64) StorageValueLoader {
	return StorageValueLoader{
		bc:              bc,
		db:              db,
		HeaderRepo:      repositories.NewHeaderRepository(db),
		StorageDiffRepo: storage2.NewDiffRepository(db),
		initializers:    initializers,
		blockNumber:     blockNumber,
	}
}

type StorageValueLoader struct {
	bc              core.BlockChain
	db              *postgres.DB
	HeaderRepo      datastore.HeaderRepository
	StorageDiffRepo storage2.DiffRepository
	initializers    []transformer.StorageTransformerInitializer
	blockNumber     int64
}

func (r *StorageValueLoader) Run() error {
	addressToKeys, getKeysErr := r.getStorageKeys()
	if getKeysErr != nil {
		return getKeysErr
	}

	header, getHeaderErr := r.HeaderRepo.GetHeader(r.blockNumber)
	if getHeaderErr != nil {
		return getHeaderErr
	}

	for address, keys := range addressToKeys {
		persistStorageErr := r.getAndPersistStorageValues(address, keys, r.blockNumber, header.Hash)
		if persistStorageErr != nil {
			return persistStorageErr
		}
	}
	logrus.Infof("Persisted storage values for %v addresses.", len(addressToKeys))

	return nil
}

func (r *StorageValueLoader) getStorageKeys() (map[common.Address][]common.Hash, error) {
	addressToKeys := make(map[common.Address][]common.Hash)
	for _, i := range r.initializers {
		transformer := i(r.db)
		keysLookup, ok := transformer.GetStorageKeysLookup().(storage.KeysLookup)
		if !ok {
			return addressToKeys, fmt.Errorf("%v type incompatible. Should be a storage.KeysLookup", keysLookup)
		}
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

func (r *StorageValueLoader) getAndPersistStorageValues(address common.Address, keys []common.Hash, blockNumber int64, headerHash string) error {
	blockNumberBigInt := big.NewInt(blockNumber)
	keccakOfAddress := crypto.Keccak256Hash(address[:])
	logrus.Infof("Getting and persisting %v storage keys for address: %v, keccak hash of address: %v", len(keys), address.Hex(), keccakOfAddress.Hex())
	for _, key := range keys {
		value, getStorageErr := r.bc.GetStorageAt(address, key, blockNumberBigInt)
		if getStorageErr != nil {
			return getStorageErr
		}
		diff := types.RawDiff{
			HashedAddress: keccakOfAddress,
			BlockHash:     common.HexToHash(headerHash),
			BlockHeight:   int(blockNumber),
			StorageKey:    key,
			StorageValue:  common.BytesToHash(value),
		}

		diffId, createDiffErr := r.StorageDiffRepo.CreateStorageDiff(diff)
		if createDiffErr != nil {
			if createDiffErr == sql.ErrNoRows {
				return nil
			}
			return createDiffErr
		}

		markFromBackfillErr := r.StorageDiffRepo.MarkFromBackfill(diffId)
		if markFromBackfillErr != nil {
			return markFromBackfillErr
		}
	}
	return nil
}
