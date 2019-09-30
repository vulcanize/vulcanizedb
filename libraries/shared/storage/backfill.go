// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package storage

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
)

// IBackFiller is the backfilling interface
type IBackFiller interface {
	BackFill(bfa BackFillerArgs) (map[common.Hash][]utils.StorageDiff, error)
}

// BatchClient is an interface to a batch-fetching geth rpc client; created to allow mock insertion
type BatchClient interface {
	BatchCall(batch []client.BatchElem) error
}

// BackFiller is the backfilling struct
type BackFiller struct {
	client BatchClient
}

// BackFillerArgs are used to pass configuration params to the backfiller
type BackFillerArgs struct {
	// mapping of hashed addresses to a list of the storage key hashes we want to collect at that address
	WantedStorage map[common.Hash][]common.Hash
	StartingBlock uint64
	EndingBlock   uint64
}

const method = "statediff_stateDiffAt"

// NewStorageBackFiller returns a IBackFiller
func NewStorageBackFiller(bc BatchClient) IBackFiller {
	return &BackFiller{
		client: bc,
	}
}

// BackFill uses the provided config to fetch and return the state diff at the specified blocknumber
// StateDiffAt(ctx context.Context, blockNumber uint64) (*Payload, error)
func (bf *BackFiller) BackFill(bfa BackFillerArgs) (map[common.Hash][]utils.StorageDiff, error) {
	results := make(map[common.Hash][]utils.StorageDiff, len(bfa.WantedStorage))
	if bfa.EndingBlock < bfa.StartingBlock {
		return nil, errors.New("backfill: ending block number needs to be greater than starting block number")
	}
	batch := make([]client.BatchElem, 0)
	for i := bfa.StartingBlock; i <= bfa.EndingBlock; i++ {
		batch = append(batch, client.BatchElem{
			Method: method,
			Args:   []interface{}{i},
			Result: new(statediff.Payload),
		})
	}
	batchErr := bf.client.BatchCall(batch)
	if batchErr != nil {
		return nil, batchErr
	}
	for _, batchElem := range batch {
		payload := batchElem.Result.(*statediff.Payload)
		if batchElem.Error != nil {
			return nil, batchElem.Error
		}
		block := new(types.Block)
		blockDecodeErr := rlp.DecodeBytes(payload.BlockRlp, block)
		if blockDecodeErr != nil {
			return nil, blockDecodeErr
		}
		stateDiff := new(statediff.StateDiff)
		stateDiffDecodeErr := rlp.DecodeBytes(payload.StateDiffRlp, stateDiff)
		if stateDiffDecodeErr != nil {
			return nil, stateDiffDecodeErr
		}
		accounts := utils.GetAccountsFromDiff(*stateDiff)
		for _, account := range accounts {
			if wantedHashedAddress(bfa.WantedStorage, common.BytesToHash(account.Key)) {
				logrus.Trace(fmt.Sprintf("iterating through %d Storage values on account", len(account.Storage)))
				for _, storage := range account.Storage {
					if wantedHashedStorageKey(bfa.WantedStorage[common.BytesToHash(account.Key)], storage.Key) {
						diff, formatErr := utils.FromGethStateDiff(account, stateDiff, storage)
						logrus.Trace("adding storage diff to out channel",
							"keccak of address: ", diff.HashedAddress.Hex(),
							"block height: ", diff.BlockHeight,
							"storage key: ", diff.StorageKey.Hex(),
							"storage value: ", diff.StorageValue.Hex())
						if formatErr != nil {
							return nil, formatErr
						}
						results[diff.HashedAddress] = append(results[diff.HashedAddress], diff)
					}
				}
			}
		}
	}
	return results, nil
}

func wantedHashedAddress(wantedStorage map[common.Hash][]common.Hash, hashedKey common.Hash) bool {
	for addrHash := range wantedStorage {
		if bytes.Equal(addrHash.Bytes(), hashedKey.Bytes()) {
			return true
		}
	}
	return false
}

func wantedHashedStorageKey(wantedKeys []common.Hash, keyBytes []byte) bool {
	for _, key := range wantedKeys {
		if bytes.Equal(key.Bytes(), keyBytes) {
			return true
		}
	}
	return false
}
