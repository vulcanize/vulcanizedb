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

package super_node_test

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/super_node"
	mocks2 "github.com/vulcanize/vulcanizedb/pkg/super_node/eth/mocks"
)

var _ = Describe("Service", func() {
	Describe("SyncAndPublish", func() {
		It("Streams statediff.Payloads, converts them to IPLDPayloads, publishes IPLDPayloads, and indexes CIDPayloads", func() {
			wg := new(sync.WaitGroup)
			payloadChan := make(chan interface{}, 1)
			quitChan := make(chan bool, 1)
			mockCidIndexer := &mocks2.CIDIndexer{
				ReturnErr: nil,
			}
			mockPublisher := &mocks2.IPLDPublisher{
				ReturnCIDPayload: mocks2.MockCIDPayload,
				ReturnErr:        nil,
			}
			mockStreamer := &mocks2.StateDiffStreamer{
				ReturnSub: &rpc.ClientSubscription{},
				StreamPayloads: []statediff.Payload{
					mocks2.MockStateDiffPayload,
				},
				ReturnErr: nil,
			}
			mockConverter := &mocks2.PayloadConverter{
				ReturnIPLDPayload: mocks2.MockIPLDPayload,
				ReturnErr:         nil,
			}
			processor := &super_node.Service{
				Indexer:        mockCidIndexer,
				Publisher:      mockPublisher,
				Streamer:       mockStreamer,
				Converter:      mockConverter,
				PayloadChan:    payloadChan,
				QuitChan:       quitChan,
				WorkerPoolSize: 1,
			}
			err := processor.SyncAndPublish(wg, nil, nil)
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(2 * time.Second)
			quitChan <- true
			wg.Wait()
			Expect(mockConverter.PassedStatediffPayload).To(Equal(mocks2.MockStateDiffPayload))
			Expect(len(mockCidIndexer.PassedCIDPayload)).To(Equal(1))
			Expect(mockCidIndexer.PassedCIDPayload[0]).To(Equal(mocks2.MockCIDPayload))
			Expect(mockPublisher.PassedIPLDPayload).To(Equal(mocks2.MockIPLDPayload))
			Expect(mockStreamer.PassedPayloadChan).To(Equal(payloadChan))
		})
	})
})
