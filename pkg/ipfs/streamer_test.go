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

package ipfs_test

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs/mocks"
)

var _ = Describe("Streamer", func() {
	Describe("Stream", func() {
		It("Streams StatediffPayloads from a Geth RPC subscription", func() {
			mockStreamer := mocks.StateDiffStreamer{}
			mockStreamer.ReturnSub = &rpc.ClientSubscription{}
			mockStreamer.StreamPayloads = []statediff.Payload{
				mocks.MockStatediffPayload,
			}
			payloadChan := make(chan statediff.Payload, 1)
			sub, err := mockStreamer.Stream(payloadChan)
			Expect(err).ToNot(HaveOccurred())
			Expect(sub).To(Equal(&rpc.ClientSubscription{}))
			Expect(mockStreamer.PassedPayloadChan).To(Equal(payloadChan))
			streamedPayload := <-payloadChan
			Expect(streamedPayload).To(Equal(mocks.MockStatediffPayload))
		})
	})
})
