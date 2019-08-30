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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

type MockStoragediffStreamer struct {
	subscribeError    error
	PassedPayloadChan chan statediff.Payload
	streamPayloads    []statediff.Payload
}

func (streamer *MockStoragediffStreamer) Stream(statediffPayloadChan chan statediff.Payload) (*rpc.ClientSubscription, error) {
	clientSubscription := rpc.ClientSubscription{}
	streamer.PassedPayloadChan = statediffPayloadChan

	go func() {
		for _, payload := range streamer.streamPayloads {
			streamer.PassedPayloadChan <- payload
		}
	}()

	return &clientSubscription, streamer.subscribeError
}

func (streamer *MockStoragediffStreamer) SetSubscribeError(err error) {
	streamer.subscribeError = err
}

func (streamer *MockStoragediffStreamer) SetPayloads(payloads []statediff.Payload) {
	streamer.streamPayloads = payloads
}

var _ = Describe("Geth RPC Storage Fetcher", func() {
	var streamer MockStoragediffStreamer
	var statediffPayloadChan chan statediff.Payload
	var statediffFetcher fetcher.GethRpcStorageFetcher
	var storagediffChan chan utils.StorageDiff
	var errorChan chan error

	BeforeEach(func() {
		streamer = MockStoragediffStreamer{}
		statediffPayloadChan = make(chan statediff.Payload, 1)
		statediffFetcher = fetcher.NewGethRpcStorageFetcher(&streamer, statediffPayloadChan)
		storagediffChan = make(chan utils.StorageDiff)
		errorChan = make(chan error)
	})

	It("adds errors to error channel if the RPC subscription fails and panics", func(done Done) {
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
		streamer.SetPayloads([]statediff.Payload{test_data.MockStatediffPayload})

		go statediffFetcher.FetchStorageDiffs(storagediffChan, errorChan)

		streamedPayload := <-statediffPayloadChan
		Expect(streamedPayload).To(Equal(test_data.MockStatediffPayload))
		Expect(streamer.PassedPayloadChan).To(Equal(statediffPayloadChan))
		close(done)
	})

	It("adds errors to error channel if decoding the state diff RLP fails", func(done Done) {
		badStatediffPayload := statediff.Payload{}
		streamer.SetPayloads([]statediff.Payload{badStatediffPayload})

		go statediffFetcher.FetchStorageDiffs(storagediffChan, errorChan)

		Expect(<-errorChan).To(MatchError("EOF"))

		close(done)
	})

	It("adds parsed statediff payloads to the rows channel", func(done Done) {
		streamer.SetPayloads([]statediff.Payload{test_data.MockStatediffPayload})

		go statediffFetcher.FetchStorageDiffs(storagediffChan, errorChan)

		height := test_data.BlockNumber
		intHeight := int(height.Int64())
		createdExpectedStorageDiff := utils.StorageDiff{
			KeccakOfContractAddress: common.BytesToHash(test_data.ContractLeafKey[:]),
			BlockHash:               common.HexToHash("0xfa40fbe2d98d98b3363a778d52f2bcd29d6790b9b3f3cab2b167fd12d3550f73"),
			BlockHeight:             intHeight,
			StorageKey:              common.BytesToHash(test_data.StorageKey),
			StorageValue:            common.BytesToHash(test_data.SmallStorageValue),
		}
		updatedExpectedStorageDiff := utils.StorageDiff{
			KeccakOfContractAddress: common.BytesToHash(test_data.AnotherContractLeafKey[:]),
			BlockHash:               common.HexToHash("0xfa40fbe2d98d98b3363a778d52f2bcd29d6790b9b3f3cab2b167fd12d3550f73"),
			BlockHeight:             intHeight,
			StorageKey:              common.BytesToHash(test_data.StorageKey),
			StorageValue:            common.BytesToHash(test_data.LargeStorageValue),
		}
		deletedExpectedStorageDiff := utils.StorageDiff{
			KeccakOfContractAddress: common.BytesToHash(test_data.AnotherContractLeafKey[:]),
			BlockHash:               common.HexToHash("0xfa40fbe2d98d98b3363a778d52f2bcd29d6790b9b3f3cab2b167fd12d3550f73"),
			BlockHeight:             intHeight,
			StorageKey:              common.BytesToHash(test_data.StorageKey),
			StorageValue:            common.BytesToHash(test_data.SmallStorageValue),
		}

		createdStateDiff := <-storagediffChan
		updatedStateDiff := <-storagediffChan
		deletedStateDiff := <-storagediffChan

		Expect(createdStateDiff).To(Equal(createdExpectedStorageDiff))
		Expect(updatedStateDiff).To(Equal(updatedExpectedStorageDiff))
		Expect(deletedStateDiff).To(Equal(deletedExpectedStorageDiff))

		close(done)
	})

	It("adds errors to error channel if formatting the diff as a StateDiff object fails", func(done Done) {
		badStorageDiffs := []statediff.StorageDiff{{
			Key:   test_data.StorageKey,
			Value: []byte{1, 2, 3},
			// this storage value will fail to be decoded as an RLP with the following error message:
			// "input contains more than one value"
			Path:  test_data.StoragePath,
			Proof: [][]byte{},
		}}

		accountDiffs := test_data.CreatedAccountDiffs
		accountDiffs[0].Storage = badStorageDiffs

		stateDiff := statediff.StateDiff{
			BlockNumber:     test_data.BlockNumber,
			BlockHash:       common.HexToHash(test_data.BlockHash),
			CreatedAccounts: accountDiffs,
		}

		stateDiffRlp, err := rlp.EncodeToBytes(stateDiff)
		Expect(err).NotTo(HaveOccurred())

		badStatediffPayload := statediff.Payload{
			StateDiffRlp: stateDiffRlp,
		}
		streamer.SetPayloads([]statediff.Payload{badStatediffPayload})

		go statediffFetcher.FetchStorageDiffs(storagediffChan, errorChan)

		Expect(<-errorChan).To(MatchError("rlp: input contains more than one value"))

		close(done)
	})
})
