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
	"github.com/sirupsen/logrus"

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

func (fetcher GethRpcStorageFetcher) FetchStorageDiffs(out chan<- utils.StorageDiff, errs chan<- error) {
	ethStatediffPayloadChan := fetcher.statediffPayloadChan
	clientSubscription, clientSubErr := fetcher.streamer.Stream(ethStatediffPayloadChan)
	if clientSubErr != nil {
		errs <- clientSubErr
		panic(fmt.Sprintf("Error creating a geth client subscription: %v", clientSubErr))
	}
	logrus.Info("Successfully created a geth client subscription: ", clientSubscription)

	for {
		diff := <-ethStatediffPayloadChan
		logrus.Trace("received a statediff")
		stateDiff := new(statediff.StateDiff)
		decodeErr := rlp.DecodeBytes(diff.StateDiffRlp, stateDiff)
		if decodeErr != nil {
			logrus.Warn("Error decoding state diff into RLP: ", decodeErr)
			errs <- decodeErr
		}

		accounts := utils.GetAccountsFromDiff(*stateDiff)
		logrus.Trace(fmt.Sprintf("iterating through %d accounts on stateDiff for block %d", len(accounts), stateDiff.BlockNumber))
		for _, account := range accounts {
			logrus.Trace(fmt.Sprintf("iterating through %d Storage values on account", len(account.Storage)))
			for _, storage := range account.Storage {
				diff, formatErr := utils.FromGethStateDiff(account, stateDiff, storage)
				logrus.Trace("adding storage diff to out channel",
					"keccak of address: ", diff.HashedAddress.Hex(),
					"block height: ", diff.BlockHeight,
					"storage key: ", diff.StorageKey.Hex(),
					"storage value: ", diff.StorageValue.Hex())
				if formatErr != nil {
					errs <- formatErr
				}

				out <- diff
			}
		}
	}
}
