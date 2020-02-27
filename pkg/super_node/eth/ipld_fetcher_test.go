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

package eth_test

import (
	"bytes"
	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/go-block-format"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
)

var (
	mockHeaderData    = []byte{0, 1, 2, 3, 4}
	mockUncleData     = []byte{1, 2, 3, 4, 5}
	mockTrxData       = []byte{2, 3, 4, 5, 6}
	mockReceiptData   = []byte{3, 4, 5, 6, 7}
	mockStateData     = []byte{4, 5, 6, 7, 8}
	mockStorageData   = []byte{5, 6, 7, 8, 9}
	mockStorageData2  = []byte{6, 7, 8, 9, 1}
	mockHeaderBlock   = blocks.NewBlock(mockHeaderData)
	mockUncleBlock    = blocks.NewBlock(mockUncleData)
	mockTrxBlock      = blocks.NewBlock(mockTrxData)
	mockReceiptBlock  = blocks.NewBlock(mockReceiptData)
	mockStateBlock    = blocks.NewBlock(mockStateData)
	mockStorageBlock1 = blocks.NewBlock(mockStorageData)
	mockStorageBlock2 = blocks.NewBlock(mockStorageData2)
	mockBlocks        = []blocks.Block{mockHeaderBlock, mockUncleBlock, mockTrxBlock, mockReceiptBlock, mockStateBlock, mockStorageBlock1, mockStorageBlock2}
	mockBlockService  *mocks.MockIPFSBlockService
	mockCIDWrapper    = &eth.CIDWrapper{
		BlockNumber: big.NewInt(9000),
		Header: eth.HeaderModel{
			TotalDifficulty: "1337",
			CID:             mockHeaderBlock.Cid().String(),
		},
		Uncles: []eth.UncleModel{
			{
				CID: mockUncleBlock.Cid().String(),
			},
		},
		Transactions: []eth.TxModel{
			{
				CID: mockTrxBlock.Cid().String(),
			},
		},
		Receipts: []eth.ReceiptModel{
			{
				CID: mockReceiptBlock.Cid().String(),
			},
		},
		StateNodes: []eth.StateNodeModel{{
			CID:      mockStateBlock.Cid().String(),
			Leaf:     true,
			StateKey: "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
		}},
		StorageNodes: []eth.StorageNodeWithStateKeyModel{{
			CID:        mockStorageBlock1.Cid().String(),
			Leaf:       true,
			StateKey:   "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
			StorageKey: "0000000000000000000000000000000000000000000000000000000000000001",
		},
			{
				CID:        mockStorageBlock2.Cid().String(),
				Leaf:       true,
				StateKey:   "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
				StorageKey: "0000000000000000000000000000000000000000000000000000000000000002",
			}},
	}
)

var _ = Describe("Fetcher", func() {
	Describe("FetchCIDs", func() {
		BeforeEach(func() {
			mockBlockService = new(mocks.MockIPFSBlockService)
			err := mockBlockService.AddBlocks(mockBlocks)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(mockBlockService.Blocks)).To(Equal(7))
		})

		It("Fetches and returns IPLDs for the CIDs provided in the CIDWrapper", func() {
			fetcher := new(eth.IPLDFetcher)
			fetcher.BlockService = mockBlockService
			i, err := fetcher.Fetch(mockCIDWrapper)
			Expect(err).ToNot(HaveOccurred())
			iplds, ok := i.(eth.IPLDs)
			Expect(ok).To(BeTrue())
			Expect(iplds.TotalDifficulty).To(Equal(big.NewInt(1337)))
			Expect(iplds.BlockNumber).To(Equal(mockCIDWrapper.BlockNumber))
			Expect(iplds.Header).To(Equal(ipfs.BlockModel{
				Data: mockHeaderBlock.RawData(),
				CID:  mockHeaderBlock.Cid().String(),
			}))
			Expect(len(iplds.Uncles)).To(Equal(1))
			Expect(iplds.Uncles[0]).To(Equal(ipfs.BlockModel{
				Data: mockUncleBlock.RawData(),
				CID:  mockUncleBlock.Cid().String(),
			}))
			Expect(len(iplds.Transactions)).To(Equal(1))
			Expect(iplds.Transactions[0]).To(Equal(ipfs.BlockModel{
				Data: mockTrxBlock.RawData(),
				CID:  mockTrxBlock.Cid().String(),
			}))
			Expect(len(iplds.Receipts)).To(Equal(1))
			Expect(iplds.Receipts[0]).To(Equal(ipfs.BlockModel{
				Data: mockReceiptBlock.RawData(),
				CID:  mockReceiptBlock.Cid().String(),
			}))
			Expect(len(iplds.StateNodes)).To(Equal(1))
			Expect(iplds.StateNodes[0].StateTrieKey).To(Equal(common.HexToHash("0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")))
			Expect(iplds.StateNodes[0].Leaf).To(BeTrue())
			Expect(iplds.StateNodes[0].IPLD).To(Equal(ipfs.BlockModel{
				Data: mockStateBlock.RawData(),
				CID:  mockStateBlock.Cid().String(),
			}))
			Expect(len(iplds.StorageNodes)).To(Equal(2))
			for _, storage := range iplds.StorageNodes {
				Expect(storage.Leaf).To(BeTrue())
				Expect(storage.StateTrieKey).To(Equal(common.HexToHash("0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")))
				if bytes.Equal(storage.StorageTrieKey.Bytes(), common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001").Bytes()) {
					Expect(storage.IPLD).To(Equal(ipfs.BlockModel{
						Data: mockStorageBlock1.RawData(),
						CID:  mockStorageBlock1.Cid().String(),
					}))
				}
				if bytes.Equal(storage.StorageTrieKey.Bytes(), common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002").Bytes()) {
					Expect(storage.IPLD).To(Equal(ipfs.BlockModel{
						Data: mockStorageBlock2.RawData(),
						CID:  mockStorageBlock2.Cid().String(),
					}))
				}
			}
		})
	})
})
