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

package fetcher_test

import (
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/makerdao/vulcanizedb/libraries/shared/mocks"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Geth RPC Storage Fetcher", func() {
	var (
		streamer             *mocks.MockStoragediffStreamer
		statediffPayloadChan chan filters.Payload
		statediffFetcher     fetcher.GethRpcStorageFetcher
		storagediffChan      chan types.RawDiff
		subscription         *fakes.MockSubscription
		errorChan            chan error
		statusWriter         fakes.MockStatusWriter
		stateDiffPayloads    []filters.Payload
		badStateDiffPayloads = []filters.Payload{{}} //This empty payload is "bad" because it does not contain the required StateDiffRlp
	)

	Describe("StorageFetcher for the New Geth Patch", func() {
		// This tests fetching diff payloads from the updated simplified geth patch: https://github.com/makerdao/go-ethereum/tree/allow-state-diff-subscription
		//  - diffs are formatted with the FromGethStateDiff method
		BeforeEach(func() {
			subscription = &fakes.MockSubscription{Errs: make(chan error)}
			streamer = &mocks.MockStoragediffStreamer{ClientSubscription: subscription}
			statediffPayloadChan = make(chan filters.Payload, 1)
			statediffFetcher = fetcher.NewGethRpcStorageFetcher(streamer, statediffPayloadChan, &statusWriter)
			storagediffChan = make(chan types.RawDiff)
			errorChan = make(chan error)
			stateDiffPayloads = []filters.Payload{test_data.MockStatediffPayload}
		})

		It("adds errors to errors channel if the streamer fails to subscribe", func(done Done) {
			streamer.SetSubscribeError(fakes.FakeError)

			go func() {
				failedSub := func() {
					statediffFetcher.FetchStorageDiffs(storagediffChan, errorChan)
				}
				Expect(failedSub).To(Panic())
			}()

			Expect(<-errorChan).To(MatchError(fakes.FakeError))
			close(done)
		})

		It("streams StatediffPayloads from a Geth RPC subscription", func(done Done) {
			streamer.SetPayloads(stateDiffPayloads)

			go statediffFetcher.FetchStorageDiffs(storagediffChan, errorChan)

			streamedPayload := <-statediffPayloadChan
			Expect(streamedPayload).To(Equal(test_data.MockStatediffPayload))
			Expect(streamer.PassedPayloadChan).To(Equal(statediffPayloadChan))
			close(done)
		})

		Describe("when subscription established", func() {
			It("creates file for health check when connection established", func(done Done) {
				go statediffFetcher.FetchStorageDiffs(storagediffChan, errorChan)

				Eventually(func() bool {
					return statusWriter.WriteCalled
				}).Should(BeTrue())
				close(done)
			})

			It("adds error to errors channel if the subscription fails", func(done Done) {
				go statediffFetcher.FetchStorageDiffs(storagediffChan, errorChan)

				subscription.Errs <- fakes.FakeError

				Expect(<-errorChan).To(MatchError(fakes.FakeError))
				close(done)
			})

			It("adds errors to error channel if decoding the state diff RLP fails", func(done Done) {
				streamer.SetPayloads(badStateDiffPayloads)

				go statediffFetcher.FetchStorageDiffs(storagediffChan, errorChan)

				expectedErr := fmt.Errorf("error decoding storage diff from geth payload: %w", io.EOF)
				Expect(<-errorChan).To(MatchError(expectedErr))
				close(done)
			})

			It("adds parsed statediff payloads to the out channel", func(done Done) {
				streamer.SetPayloads(stateDiffPayloads)

				go statediffFetcher.FetchStorageDiffs(storagediffChan, errorChan)

				height := test_data.BlockNumber
				intHeight := int(height.Int64())

				expectedDiff1 := types.RawDiff{
					Address:      common.BytesToAddress(test_data.ContractLeafKey[:]),
					BlockHash:    common.HexToHash("0xfa40fbe2d98d98b3363a778d52f2bcd29d6790b9b3f3cab2b167fd12d3550f73"),
					BlockHeight:  intHeight,
					StorageKey:   common.BytesToHash(test_data.StorageKey),
					StorageValue: common.BytesToHash(test_data.SmallStorageValue),
				}
				expectedDiff2 := types.RawDiff{
					Address:      common.BytesToAddress(test_data.AnotherContractLeafKey[:]),
					BlockHash:    common.HexToHash("0xfa40fbe2d98d98b3363a778d52f2bcd29d6790b9b3f3cab2b167fd12d3550f73"),
					BlockHeight:  intHeight,
					StorageKey:   common.BytesToHash(test_data.StorageKey),
					StorageValue: common.BytesToHash(test_data.LargeStorageValue),
				}
				expectedDiff3 := types.RawDiff{
					Address:      common.BytesToAddress(test_data.AnotherContractLeafKey[:]),
					BlockHash:    common.HexToHash("0xfa40fbe2d98d98b3363a778d52f2bcd29d6790b9b3f3cab2b167fd12d3550f73"),
					BlockHeight:  intHeight,
					StorageKey:   common.BytesToHash(test_data.StorageKey),
					StorageValue: common.BytesToHash(test_data.SmallStorageValue),
				}

				diff1 := <-storagediffChan
				diff2 := <-storagediffChan
				diff3 := <-storagediffChan

				Expect(diff1).To(Equal(expectedDiff1))
				Expect(diff2).To(Equal(expectedDiff2))
				Expect(diff3).To(Equal(expectedDiff3))

				close(done)
			})

			It("adds errors to error channel if formatting the diff as a StateDiff object fails", func(done Done) {
				stateDiff := test_data.StateDiffWithBadStorageValue
				stateDiffRlp, err := rlp.EncodeToBytes(stateDiff)
				Expect(err).NotTo(HaveOccurred())
				payloadToReturn := filters.Payload{StateDiffRlp: stateDiffRlp}

				streamer.SetPayloads([]filters.Payload{payloadToReturn})

				go statediffFetcher.FetchStorageDiffs(storagediffChan, errorChan)

				Expect(<-errorChan).To(MatchError(rlp.ErrMoreThanOneValue))

				close(done)
			})
		})
	})
})
