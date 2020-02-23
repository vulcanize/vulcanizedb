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
	"github.com/vulcanize/vulcanizedb/pkg/eth/fakes"
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
	var statediffFetcher fetcher.GethRPCStorageFetcher
	var storagediffChan chan utils.StorageDiffInput
	var errorChan chan error

	BeforeEach(func() {
		streamer = MockStoragediffStreamer{}
		statediffFetcher = fetcher.NewGethRPCStorageFetcher(&streamer)
		storagediffChan = make(chan utils.StorageDiffInput)
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

		streamedPayload := <-statediffFetcher.StatediffPayloadChan
		Expect(streamedPayload).To(Equal(test_data.MockStatediffPayload))
		Expect(streamer.PassedPayloadChan).To(Equal(statediffFetcher.StatediffPayloadChan))
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
		createdExpectedStorageDiff := utils.StorageDiffInput{
			HashedAddress: common.BytesToHash(test_data.ContractLeafKey[:]),
			BlockHash:     common.HexToHash("0xfa40fbe2d98d98b3363a778d52f2bcd29d6790b9b3f3cab2b167fd12d3550f73"),
			BlockHeight:   intHeight,
			StorageKey:    common.BytesToHash(test_data.StorageKey),
			StorageValue:  common.BytesToHash(test_data.SmallStorageValue),
		}
		updatedExpectedStorageDiff := utils.StorageDiffInput{
			HashedAddress: common.BytesToHash(test_data.AnotherContractLeafKey[:]),
			BlockHash:     common.HexToHash("0xfa40fbe2d98d98b3363a778d52f2bcd29d6790b9b3f3cab2b167fd12d3550f73"),
			BlockHeight:   intHeight,
			StorageKey:    common.BytesToHash(test_data.StorageKey),
			StorageValue:  common.BytesToHash(test_data.LargeStorageValue),
		}
		deletedExpectedStorageDiff := utils.StorageDiffInput{
			HashedAddress: common.BytesToHash(test_data.AnotherContractLeafKey[:]),
			BlockHash:     common.HexToHash("0xfa40fbe2d98d98b3363a778d52f2bcd29d6790b9b3f3cab2b167fd12d3550f73"),
			BlockHeight:   intHeight,
			StorageKey:    common.BytesToHash(test_data.StorageKey),
			StorageValue:  common.BytesToHash(test_data.SmallStorageValue),
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
		accountDiffs := test_data.CreatedAccountDiffs
		accountDiffs[0].Storage = []statediff.StorageDiff{test_data.StorageWithBadValue}

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
