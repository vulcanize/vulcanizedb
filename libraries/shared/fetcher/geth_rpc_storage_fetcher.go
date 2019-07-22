// Copyright 2019 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fetcher

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/streamer"
)

type GethRpcStorageFetcher struct {
	statediffPayloadChan chan statediff.Payload
	streamer             streamer.Streamer
}

func NewGethRpcStorageFetcher(streamer streamer.Streamer, statediffPayloadChan chan statediff.Payload) GethRpcStorageFetcher {
	return GethRpcStorageFetcher{
		statediffPayloadChan: statediffPayloadChan,
		streamer:             streamer,
	}
}

func (fetcher *GethRpcStorageFetcher) FetchStorageDiffs(out chan<- utils.StorageDiffRow, errs chan<- error) {
	ethStatediffPayloadChan := fetcher.statediffPayloadChan
	_, err := fetcher.streamer.Stream(ethStatediffPayloadChan)
	if err != nil {
		errs <- err
	}

	for {
		diff := <-ethStatediffPayloadChan
		stateDiff := new(statediff.StateDiff)
		err = rlp.DecodeBytes(diff.StateDiffRlp, stateDiff)
		if err != nil {
			errs <- err
		}
		accounts := getAccountDiffs(*stateDiff)

		for _, account := range accounts {
			for _, storage := range account.Storage {
				out <- utils.StorageDiffRow{
					Contract:     common.BytesToAddress(account.Key),
					BlockHash:    stateDiff.BlockHash,
					BlockHeight:  int(stateDiff.BlockNumber.Int64()),
					StorageKey:   common.BytesToHash(storage.Key),
					StorageValue: common.BytesToHash(storage.Value),
				}
			}
		}
	}
}

func getAccountDiffs(stateDiff statediff.StateDiff) []statediff.AccountDiff {
	accounts := append(stateDiff.CreatedAccounts, stateDiff.UpdatedAccounts...)
	return append(accounts, stateDiff.DeletedAccounts...)
}
