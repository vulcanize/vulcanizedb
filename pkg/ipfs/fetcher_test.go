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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/go-block-format"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/mocks"
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
	mockCIDWrapper    = ipfs.CIDWrapper{
		BlockNumber:  big.NewInt(9000),
		Headers:      []string{mockHeaderBlock.Cid().String()},
		Uncles:       []string{mockUncleBlock.Cid().String()},
		Transactions: []string{mockTrxBlock.Cid().String()},
		Receipts:     []string{mockReceiptBlock.Cid().String()},
		StateNodes: []ipfs.StateNodeCID{{
			CID:  mockStateBlock.Cid().String(),
			Leaf: true,
			Key:  "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
		}},
		StorageNodes: []ipfs.StorageNodeCID{{
			CID:      mockStorageBlock1.Cid().String(),
			Leaf:     true,
			StateKey: "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
			Key:      "0000000000000000000000000000000000000000000000000000000000000001",
		},
			{
				CID:      mockStorageBlock2.Cid().String(),
				Leaf:     true,
				StateKey: "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
				Key:      "0000000000000000000000000000000000000000000000000000000000000002",
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
			fetcher := new(ipfs.EthIPLDFetcher)
			fetcher.BlockService = mockBlockService
			iplds, err := fetcher.FetchCIDs(mockCIDWrapper)
			Expect(err).ToNot(HaveOccurred())
			Expect(iplds.BlockNumber).To(Equal(mockCIDWrapper.BlockNumber))
			Expect(len(iplds.Headers)).To(Equal(1))
			Expect(iplds.Headers[0]).To(Equal(mockHeaderBlock))
			Expect(len(iplds.Uncles)).To(Equal(1))
			Expect(iplds.Uncles[0]).To(Equal(mockUncleBlock))
			Expect(len(iplds.Transactions)).To(Equal(1))
			Expect(iplds.Transactions[0]).To(Equal(mockTrxBlock))
			Expect(len(iplds.Receipts)).To(Equal(1))
			Expect(iplds.Receipts[0]).To(Equal(mockReceiptBlock))
			Expect(len(iplds.StateNodes)).To(Equal(1))
			stateNode, ok := iplds.StateNodes[common.HexToHash("0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")]
			Expect(ok).To(BeTrue())
			Expect(stateNode).To(Equal(mockStateBlock))
			Expect(len(iplds.StorageNodes)).To(Equal(1))
			storageNodes := iplds.StorageNodes[common.HexToHash("0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")]
			Expect(len(storageNodes)).To(Equal(2))
			storageNode1, ok := storageNodes[common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")]
			Expect(ok).To(BeTrue())
			Expect(storageNode1).To(Equal(mockStorageBlock1))
			storageNode2, ok := storageNodes[common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")]
			Expect(storageNode2).To(Equal(mockStorageBlock2))
			Expect(ok).To(BeTrue())
		})
	})
})
