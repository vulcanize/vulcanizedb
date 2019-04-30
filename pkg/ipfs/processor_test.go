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
	"sync"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/test_helpers/mocks"
)

var _ = Describe("Processor", func() {

	Describe("Process", func() {
		It("Streams StatediffPayloads, converts them to IPLDPayloads, publishes IPLDPayloads, and indexes CIDPayloads", func() {
			wg := new(sync.WaitGroup)
			payloadChan := make(chan statediff.Payload, 800)
			processor := ipfs.Processor{
				Repository: &mocks.CIDRepository{
					ReturnErr: nil,
				},
				Publisher: &mocks.IPLDPublisher{
					ReturnCIDPayload: &test_helpers.MockCIDPayload,
					ReturnErr:        nil,
				},
				Streamer: &mocks.StateDiffStreamer{
					ReturnSub: &rpc.ClientSubscription{},
					StreamPayloads: []statediff.Payload{
						test_helpers.MockStatediffPayload,
					},
					ReturnErr: nil,
					WaitGroup: wg,
				},
				Converter: &mocks.PayloadConverter{
					ReturnIPLDPayload: &test_helpers.MockIPLDPayload,
					ReturnErr:         nil,
				},
				PayloadChan: payloadChan,
			}
			err := processor.Process(wg)
			Expect(err).ToNot(HaveOccurred())
			wg.Wait()
		})

		It("Fails if", func() {

		})
	})
})
