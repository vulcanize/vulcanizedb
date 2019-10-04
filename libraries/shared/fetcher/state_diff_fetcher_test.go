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
	"bytes"
	"encoding/json"
	"errors"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
)

type mockClient struct {
	MappedStateDiffAt map[uint64][]byte
}

// SetReturnDiffAt method to set what statediffs the mock client returns
func (mc *mockClient) SetReturnDiffAt(height uint64, diffPayload statediff.Payload) error {
	if mc.MappedStateDiffAt == nil {
		mc.MappedStateDiffAt = make(map[uint64][]byte)
	}
	by, err := json.Marshal(diffPayload)
	if err != nil {
		return err
	}
	mc.MappedStateDiffAt[height] = by
	return nil
}

// BatchCall mockClient method to simulate batch call to geth
func (mc *mockClient) BatchCall(batch []client.BatchElem) error {
	if mc.MappedStateDiffAt == nil {
		return errors.New("mockclient needs to be initialized with statediff payloads and errors")
	}
	for _, batchElem := range batch {
		if len(batchElem.Args) != 1 {
			return errors.New("expected batch elem to contain single argument")
		}
		blockHeight, ok := batchElem.Args[0].(uint64)
		if !ok {
			return errors.New("expected batch elem argument to be a uint64")
		}
		err := json.Unmarshal(mc.MappedStateDiffAt[blockHeight], batchElem.Result)
		if err != nil {
			return err
		}
	}
	return nil
}

var _ = Describe("StateDiffFetcher", func() {
	Describe("FetchStateDiffsAt", func() {
		var (
			mc               *mockClient
			stateDiffFetcher fetcher.IStateDiffFetcher
		)
		BeforeEach(func() {
			mc = new(mockClient)
			setDiffAtErr1 := mc.SetReturnDiffAt(test_data.BlockNumber.Uint64(), test_data.MockStatediffPayload)
			Expect(setDiffAtErr1).ToNot(HaveOccurred())
			setDiffAtErr2 := mc.SetReturnDiffAt(test_data.BlockNumber2.Uint64(), test_data.MockStatediffPayload2)
			Expect(setDiffAtErr2).ToNot(HaveOccurred())
			stateDiffFetcher = fetcher.NewStateDiffFetcher(mc)
		})
		It("Batch calls statediff_stateDiffAt", func() {
			blockHeights := []uint64{
				test_data.BlockNumber.Uint64(),
				test_data.BlockNumber2.Uint64(),
			}
			stateDiffPayloads, fetchErr := stateDiffFetcher.FetchStateDiffsAt(blockHeights)
			Expect(fetchErr).ToNot(HaveOccurred())
			Expect(len(stateDiffPayloads)).To(Equal(2))
			// Can only rlp encode the slice of diffs as part of a struct
			// Rlp encoding allows us to compare content of the slices when the order in the slice may vary
			expectedPayloadsStruct := struct {
				payloads []*statediff.Payload
			}{
				[]*statediff.Payload{
					&test_data.MockStatediffPayload,
					&test_data.MockStatediffPayload2,
				},
			}
			expectedPayloadsBytes, rlpErr1 := rlp.EncodeToBytes(expectedPayloadsStruct)
			Expect(rlpErr1).ToNot(HaveOccurred())
			receivedPayloadsStruct := struct {
				payloads []*statediff.Payload
			}{
				stateDiffPayloads,
			}
			receivedPayloadsBytes, rlpErr2 := rlp.EncodeToBytes(receivedPayloadsStruct)
			Expect(rlpErr2).ToNot(HaveOccurred())
			Expect(bytes.Equal(expectedPayloadsBytes, receivedPayloadsBytes)).To(BeTrue())
		})
	})
})
