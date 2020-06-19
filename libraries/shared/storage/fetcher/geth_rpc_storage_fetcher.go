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
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/streamer"
	"github.com/makerdao/vulcanizedb/pkg/fs"
	"github.com/sirupsen/logrus"
)

type GethPatchVersion int

const (
	OldGethPatch GethPatchVersion = iota
	NewGethPatch
)

type GethRpcStorageFetcher struct {
	statediffPayloadChan chan statediff.Payload
	streamer             streamer.Streamer
	gethVersion          GethPatchVersion
	statusWriter         fs.StatusWriter
}

func NewGethRpcStorageFetcher(streamer streamer.Streamer, statediffPayloadChan chan statediff.Payload, gethVersion GethPatchVersion, statusWriter fs.StatusWriter) GethRpcStorageFetcher {
	return GethRpcStorageFetcher{
		statediffPayloadChan: statediffPayloadChan,
		streamer:             streamer,
		gethVersion:          gethVersion,
		statusWriter:         statusWriter,
	}
}

func (fetcher GethRpcStorageFetcher) FetchStorageDiffs(out chan<- types.RawDiff, errs chan<- error) {
	ethStatediffPayloadChan := fetcher.statediffPayloadChan
	clientSubscription, clientSubErr := fetcher.streamer.Stream(ethStatediffPayloadChan)
	if clientSubErr != nil {
		errs <- clientSubErr
		panic(fmt.Sprintf("Error creating a geth client subscription: %v", clientSubErr))
	}
	logrus.Info("Successfully created a geth client subscription: ", clientSubscription)

	writeErr := fetcher.statusWriter.Write()
	if writeErr != nil {
		errs <- writeErr
	}

	for {
		select {
		case err := <-clientSubscription.Err():
			logrus.Errorf("error with client subscription: %s", err.Error())
			errs <- err
		case diff := <-ethStatediffPayloadChan:
			logrus.Trace("received a statediff")
			stateDiff := new(statediff.StateDiff)
			decodeErr := rlp.DecodeBytes(diff.StateDiffRlp, stateDiff)
			logrus.Tracef("received a statediff from block: %v", stateDiff.BlockNumber)
			if decodeErr != nil {
				logrus.Warn("Error decoding state diff into RLP: ", decodeErr)
				errs <- decodeErr
			}

			accounts := getAccountsFromDiff(*stateDiff)
			logrus.Trace(fmt.Sprintf("iterating through %d accounts on stateDiff for block %d", len(accounts), stateDiff.BlockNumber))
			for _, account := range accounts {
				logrus.Trace(fmt.Sprintf("iterating through %d Storage values on account", len(account.Storage)))
				for _, accountStorage := range account.Storage {
					diff, formatErr := fetcher.formatDiff(account, stateDiff, accountStorage)
					logrus.Tracef("adding storage diff to out channel. keccak of address: %v, block height: %v, storage key: %v, storage value: %v",
						diff.HashedAddress.Hex(), diff.BlockHeight, diff.StorageKey.Hex(), diff.StorageValue.Hex())
					if formatErr != nil {
						errs <- formatErr
					}

					out <- diff
				}
			}
		}

	}
}

func (fetcher GethRpcStorageFetcher) formatDiff(account statediff.AccountDiff, stateDiff *statediff.StateDiff, storage statediff.StorageDiff) (types.RawDiff, error) {
	if fetcher.gethVersion == OldGethPatch {
		return types.FromOldGethStateDiff(account, stateDiff, storage)
	} else {
		return types.FromNewGethStateDiff(account, stateDiff, storage)
	}
}

func getAccountsFromDiff(stateDiff statediff.StateDiff) []statediff.AccountDiff {
	accounts := append(stateDiff.CreatedAccounts, stateDiff.UpdatedAccounts...)
	return append(accounts, stateDiff.DeletedAccounts...)
}
