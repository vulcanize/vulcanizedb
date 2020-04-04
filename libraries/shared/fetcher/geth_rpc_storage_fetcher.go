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
package fetcher

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/streamer"
)

const (
	PayloadChanBufferSize = 20000 // the max eth sub buffer size
)

type GethRPCStorageFetcher struct {
	StatediffPayloadChan chan statediff.Payload
	streamer             streamer.Streamer
}

func NewGethRPCStorageFetcher(streamer streamer.Streamer) GethRPCStorageFetcher {
	return GethRPCStorageFetcher{
		StatediffPayloadChan: make(chan statediff.Payload, PayloadChanBufferSize),
		streamer:             streamer,
	}
}

func (fetcher GethRPCStorageFetcher) FetchStorageDiffs(out chan<- utils.StorageDiffInput, errs chan<- error) {
	ethStatediffPayloadChan := fetcher.StatediffPayloadChan
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
			logrus.Trace(fmt.Sprintf("iterating through %d Storage values on account with key %s", len(account.Storage), common.BytesToHash(account.LeafKey).Hex()))
			for _, storage := range account.Storage {
				diff, formatErr := utils.FromGethStateDiff(account, stateDiff, storage)
				if formatErr != nil {
					logrus.Error("failed to format utils.StorageDiff from storage with key: ", common.BytesToHash(storage.LeafKey), "from account with key: ", common.BytesToHash(account.LeafKey))
					errs <- formatErr
					continue
				}
				logrus.Trace("adding storage diff to out channel",
					"keccak of address: ", diff.HashedAddress.Hex(),
					"block height: ", diff.BlockHeight,
					"storage key: ", diff.StorageKey.Hex(),
					"storage value: ", diff.StorageValue.Hex())

				out <- diff
			}
		}
	}
}
