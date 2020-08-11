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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/streamer"
	"github.com/makerdao/vulcanizedb/pkg/fs"
	"github.com/sirupsen/logrus"
)

type GethRpcStorageFetcher struct {
	statediffPayloadChan chan filters.Payload
	streamer             streamer.Streamer
	statusWriter         fs.StatusWriter
}

func NewGethRpcStorageFetcher(streamer streamer.Streamer, statediffPayloadChan chan filters.Payload, statusWriter fs.StatusWriter) GethRpcStorageFetcher {
	return GethRpcStorageFetcher{
		statediffPayloadChan: statediffPayloadChan,
		streamer:             streamer,
		statusWriter:         statusWriter,
	}
}

var (
	processingDiffsLogString = "processing %d storage diffs for account %s"
	addingDiffsLogString     = "adding storage diff to out channel. keccak of address: %v, block height: %v, storage key: %v, storage value: %v"
)

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
		case diffPayload := <-ethStatediffPayloadChan:
			logrus.Trace("received a statediff payload")
			fetcher.handleDiffPayload(diffPayload, out, errs)
		}
	}
}

func (fetcher GethRpcStorageFetcher) handleDiffPayload(payload filters.Payload, out chan<- types.RawDiff, errs chan<- error) {
	var stateDiff filters.StateDiff
	decodeErr := rlp.DecodeBytes(payload.StateDiffRlp, &stateDiff)
	if decodeErr != nil {
		errs <- fmt.Errorf("error decoding storage diff from geth payload: %w", decodeErr)
		return
	}

	for _, account := range stateDiff.UpdatedAccounts {
		logrus.Infof(processingDiffsLogString, len(account.Storage), common.Bytes2Hex(account.Key))
		for _, accountStorage := range account.Storage {
			rawDiff, formatErr := types.FromGethStateDiff(account, &stateDiff, accountStorage)
			if formatErr != nil {
				errs <- formatErr
				return
			}

			logrus.Tracef(addingDiffsLogString, rawDiff.HashedAddress.Hex(), rawDiff.BlockHeight, rawDiff.StorageKey.Hex(), rawDiff.StorageValue.Hex())
			out <- rawDiff
		}
	}
}
