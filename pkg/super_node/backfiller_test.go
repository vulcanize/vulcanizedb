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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/super_node"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
	mocks2 "github.com/vulcanize/vulcanizedb/pkg/super_node/shared/mocks"
)

var _ = Describe("BackFiller", func() {
	Describe("FillGaps", func() {
		It("Periodically checks for and fills in gaps in the super node's data", func() {
			mockCidRepo := &mocks.CIDIndexer{
				ReturnErr: nil,
			}
			mockPublisher := &mocks.IterativeIPLDPublisher{
				ReturnCIDPayload: []*eth.CIDPayload{mocks.MockCIDPayload, mocks.MockCIDPayload},
				ReturnErr:        nil,
			}
			mockConverter := &mocks.IterativePayloadConverter{
				ReturnIPLDPayload: []eth.IPLDPayload{mocks.MockIPLDPayload, mocks.MockIPLDPayload},
				ReturnErr:         nil,
			}
			mockRetriever := &mocks2.MockCIDRetriever{
				FirstBlockNumberToReturn: 0,
				GapsToRetrieve: []shared.Gap{
					{
						Start: 100, Stop: 101,
					},
				},
			}
			mockFetcher := &mocks2.IPLDFetcher{
				PayloadsToReturn: map[uint64]shared.RawChainData{
					100: mocks.MockStateDiffPayload,
					101: mocks.MockStateDiffPayload,
				},
			}
			backfiller := &super_node.BackFillService{
				Indexer:           mockCidRepo,
				Publisher:         mockPublisher,
				Converter:         mockConverter,
				Fetcher:           mockFetcher,
				Retriever:         mockRetriever,
				GapCheckFrequency: time.Second * 2,
				BatchSize:         super_node.DefaultMaxBatchSize,
			}
			wg := &sync.WaitGroup{}
			quitChan := make(chan bool, 1)
			backfiller.FillGaps(wg, quitChan)
			time.Sleep(time.Second * 3)
			quitChan <- true
			Expect(len(mockCidRepo.PassedCIDPayload)).To(Equal(2))
			Expect(mockCidRepo.PassedCIDPayload[0]).To(Equal(mocks.MockCIDPayload))
			Expect(mockCidRepo.PassedCIDPayload[1]).To(Equal(mocks.MockCIDPayload))
			Expect(len(mockPublisher.PassedIPLDPayload)).To(Equal(2))
			Expect(mockPublisher.PassedIPLDPayload[0]).To(Equal(mocks.MockIPLDPayload))
			Expect(mockPublisher.PassedIPLDPayload[1]).To(Equal(mocks.MockIPLDPayload))
			Expect(len(mockConverter.PassedStatediffPayload)).To(Equal(2))
			Expect(mockConverter.PassedStatediffPayload[0]).To(Equal(mocks.MockStateDiffPayload))
			Expect(mockConverter.PassedStatediffPayload[1]).To(Equal(mocks.MockStateDiffPayload))
			Expect(mockRetriever.CalledTimes).To(Equal(1))
			Expect(len(mockFetcher.CalledAtBlockHeights)).To(Equal(1))
			Expect(mockFetcher.CalledAtBlockHeights[0]).To(Equal([]uint64{100, 101}))
		})

		It("Works for single block `ranges`", func() {
			mockCidRepo := &mocks.CIDIndexer{
				ReturnErr: nil,
			}
			mockPublisher := &mocks.IterativeIPLDPublisher{
				ReturnCIDPayload: []*eth.CIDPayload{mocks.MockCIDPayload},
				ReturnErr:        nil,
			}
			mockConverter := &mocks.IterativePayloadConverter{
				ReturnIPLDPayload: []eth.IPLDPayload{mocks.MockIPLDPayload},
				ReturnErr:         nil,
			}
			mockRetriever := &mocks2.MockCIDRetriever{
				FirstBlockNumberToReturn: 0,
				GapsToRetrieve: []shared.Gap{
					{
						Start: 100, Stop: 100,
					},
				},
			}
			mockFetcher := &mocks2.IPLDFetcher{
				PayloadsToReturn: map[uint64]shared.RawChainData{
					100: mocks.MockStateDiffPayload,
				},
			}
			backfiller := &super_node.BackFillService{
				Indexer:           mockCidRepo,
				Publisher:         mockPublisher,
				Converter:         mockConverter,
				Fetcher:           mockFetcher,
				Retriever:         mockRetriever,
				GapCheckFrequency: time.Second * 2,
				BatchSize:         super_node.DefaultMaxBatchSize,
			}
			wg := &sync.WaitGroup{}
			quitChan := make(chan bool, 1)
			backfiller.FillGaps(wg, quitChan)
			time.Sleep(time.Second * 3)
			quitChan <- true
			Expect(len(mockCidRepo.PassedCIDPayload)).To(Equal(1))
			Expect(mockCidRepo.PassedCIDPayload[0]).To(Equal(mocks.MockCIDPayload))
			Expect(len(mockPublisher.PassedIPLDPayload)).To(Equal(1))
			Expect(mockPublisher.PassedIPLDPayload[0]).To(Equal(mocks.MockIPLDPayload))
			Expect(len(mockConverter.PassedStatediffPayload)).To(Equal(1))
			Expect(mockConverter.PassedStatediffPayload[0]).To(Equal(mocks.MockStateDiffPayload))
			Expect(mockRetriever.CalledTimes).To(Equal(1))
			Expect(len(mockFetcher.CalledAtBlockHeights)).To(Equal(1))
			Expect(mockFetcher.CalledAtBlockHeights[0]).To(Equal([]uint64{100}))
		})

		It("Finds beginning gap", func() {
			mockCidRepo := &mocks.CIDIndexer{
				ReturnErr: nil,
			}
			mockPublisher := &mocks.IterativeIPLDPublisher{
				ReturnCIDPayload: []*eth.CIDPayload{mocks.MockCIDPayload, mocks.MockCIDPayload},
				ReturnErr:        nil,
			}
			mockConverter := &mocks.IterativePayloadConverter{
				ReturnIPLDPayload: []eth.IPLDPayload{mocks.MockIPLDPayload, mocks.MockIPLDPayload},
				ReturnErr:         nil,
			}
			mockRetriever := &mocks2.MockCIDRetriever{
				FirstBlockNumberToReturn: 3,
				GapsToRetrieve:           []shared.Gap{},
			}
			mockFetcher := &mocks2.IPLDFetcher{
				PayloadsToReturn: map[uint64]shared.RawChainData{
					1: mocks.MockStateDiffPayload,
					2: mocks.MockStateDiffPayload,
				},
			}
			backfiller := &super_node.BackFillService{
				Indexer:           mockCidRepo,
				Publisher:         mockPublisher,
				Converter:         mockConverter,
				Fetcher:           mockFetcher,
				Retriever:         mockRetriever,
				GapCheckFrequency: time.Second * 2,
				BatchSize:         super_node.DefaultMaxBatchSize,
			}
			wg := &sync.WaitGroup{}
			quitChan := make(chan bool, 1)
			backfiller.FillGaps(wg, quitChan)
			time.Sleep(time.Second * 3)
			quitChan <- true
			Expect(len(mockCidRepo.PassedCIDPayload)).To(Equal(2))
			Expect(mockCidRepo.PassedCIDPayload[0]).To(Equal(mocks.MockCIDPayload))
			Expect(mockCidRepo.PassedCIDPayload[1]).To(Equal(mocks.MockCIDPayload))
			Expect(len(mockPublisher.PassedIPLDPayload)).To(Equal(2))
			Expect(mockPublisher.PassedIPLDPayload[0]).To(Equal(mocks.MockIPLDPayload))
			Expect(mockPublisher.PassedIPLDPayload[1]).To(Equal(mocks.MockIPLDPayload))
			Expect(len(mockConverter.PassedStatediffPayload)).To(Equal(2))
			Expect(mockConverter.PassedStatediffPayload[0]).To(Equal(mocks.MockStateDiffPayload))
			Expect(mockConverter.PassedStatediffPayload[1]).To(Equal(mocks.MockStateDiffPayload))
			Expect(mockRetriever.CalledTimes).To(Equal(1))
			Expect(len(mockFetcher.CalledAtBlockHeights)).To(Equal(1))
			Expect(mockFetcher.CalledAtBlockHeights[0]).To(Equal([]uint64{0, 1, 2}))
		})
	})
})
