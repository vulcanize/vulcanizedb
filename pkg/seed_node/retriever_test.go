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

package seed_node_test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/seed_node"
)

var (
	retriever  seed_node.CIDRetriever
	openFilter = config.Subscription{
		StartingBlock: big.NewInt(0),
		EndingBlock:   big.NewInt(1),
		HeaderFilter:  config.HeaderFilter{},
		TrxFilter:     config.TrxFilter{},
		ReceiptFilter: config.ReceiptFilter{},
		StateFilter:   config.StateFilter{},
		StorageFilter: config.StorageFilter{},
	}
	rctContractFilter = config.Subscription{
		StartingBlock: big.NewInt(0),
		EndingBlock:   big.NewInt(1),
		HeaderFilter: config.HeaderFilter{
			Off: true,
		},
		TrxFilter: config.TrxFilter{
			Off: true,
		},
		ReceiptFilter: config.ReceiptFilter{
			Contracts: []string{"0x0000000000000000000000000000000000000001"},
		},
		StateFilter: config.StateFilter{
			Off: true,
		},
		StorageFilter: config.StorageFilter{
			Off: true,
		},
	}
	rctTopicsFilter = config.Subscription{
		StartingBlock: big.NewInt(0),
		EndingBlock:   big.NewInt(1),
		HeaderFilter: config.HeaderFilter{
			Off: true,
		},
		TrxFilter: config.TrxFilter{
			Off: true,
		},
		ReceiptFilter: config.ReceiptFilter{
			Topic0s: []string{"0x0000000000000000000000000000000000000000000000000000000000000004"},
		},
		StateFilter: config.StateFilter{
			Off: true,
		},
		StorageFilter: config.StorageFilter{
			Off: true,
		},
	}
	rctTopicsAndContractFilter = config.Subscription{
		StartingBlock: big.NewInt(0),
		EndingBlock:   big.NewInt(1),
		HeaderFilter: config.HeaderFilter{
			Off: true,
		},
		TrxFilter: config.TrxFilter{
			Off: true,
		},
		ReceiptFilter: config.ReceiptFilter{
			Topic0s:   []string{"0x0000000000000000000000000000000000000000000000000000000000000004", "0x0000000000000000000000000000000000000000000000000000000000000005"},
			Contracts: []string{"0x0000000000000000000000000000000000000000"},
		},
		StateFilter: config.StateFilter{
			Off: true,
		},
		StorageFilter: config.StorageFilter{
			Off: true,
		},
	}
	rctContractsAndTopicFilter = config.Subscription{
		StartingBlock: big.NewInt(0),
		EndingBlock:   big.NewInt(1),
		HeaderFilter: config.HeaderFilter{
			Off: true,
		},
		TrxFilter: config.TrxFilter{
			Off: true,
		},
		ReceiptFilter: config.ReceiptFilter{
			Topic0s:   []string{"0x0000000000000000000000000000000000000000000000000000000000000005"},
			Contracts: []string{"0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000001"},
		},
		StateFilter: config.StateFilter{
			Off: true,
		},
		StorageFilter: config.StorageFilter{
			Off: true,
		},
	}
	rctsForAllCollectedTrxs = config.Subscription{
		StartingBlock: big.NewInt(0),
		EndingBlock:   big.NewInt(1),
		HeaderFilter: config.HeaderFilter{
			Off: true,
		},
		TrxFilter: config.TrxFilter{}, // Trx filter open so we will collect all trxs, therefore we will also collect all corresponding rcts despite rct filter
		ReceiptFilter: config.ReceiptFilter{
			Topic0s:   []string{"0x0000000000000000000000000000000000000000000000000000000000000006"}, // Topic isn't one of the topics we have
			Contracts: []string{"0x0000000000000000000000000000000000000002"},                         // Contract isn't one of the contracts we have
		},
		StateFilter: config.StateFilter{
			Off: true,
		},
		StorageFilter: config.StorageFilter{
			Off: true,
		},
	}
	rctsForSelectCollectedTrxs = config.Subscription{
		StartingBlock: big.NewInt(0),
		EndingBlock:   big.NewInt(1),
		HeaderFilter: config.HeaderFilter{
			Off: true,
		},
		TrxFilter: config.TrxFilter{
			Dst: []string{"0x0000000000000000000000000000000000000001"}, // We only filter for one of the trxs so we will only get the one corresponding receipt
		},
		ReceiptFilter: config.ReceiptFilter{
			Topic0s:   []string{"0x0000000000000000000000000000000000000000000000000000000000000006"}, // Topic isn't one of the topics we have
			Contracts: []string{"0x0000000000000000000000000000000000000002"},                         // Contract isn't one of the contracts we have
		},
		StateFilter: config.StateFilter{
			Off: true,
		},
		StorageFilter: config.StorageFilter{
			Off: true,
		},
	}
	stateFilter = config.Subscription{
		StartingBlock: big.NewInt(0),
		EndingBlock:   big.NewInt(1),
		HeaderFilter: config.HeaderFilter{
			Off: true,
		},
		TrxFilter: config.TrxFilter{
			Off: true,
		},
		ReceiptFilter: config.ReceiptFilter{
			Off: true,
		},
		StateFilter: config.StateFilter{
			Addresses: []string{mocks.Address.Hex()},
		},
		StorageFilter: config.StorageFilter{
			Off: true,
		},
	}
)

var _ = Describe("Retriever", func() {
	BeforeEach(func() {
		db, err = seed_node.SetupDB()
		Expect(err).ToNot(HaveOccurred())
		repo = seed_node.NewCIDRepository(db)
		err = repo.Index(mocks.MockCIDPayload)
		Expect(err).ToNot(HaveOccurred())
		retriever = seed_node.NewCIDRetriever(db)
	})
	AfterEach(func() {
		seed_node.TearDownDB(db)
	})

	Describe("RetrieveCIDs", func() {
		It("Retrieves all CIDs for the given blocknumber when provided an open filter", func() {
			cidWrapper, err := retriever.RetrieveCIDs(openFilter, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(cidWrapper.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper.Headers)).To(Equal(1))
			Expect(cidWrapper.Headers).To(Equal(mocks.MockCIDWrapper.Headers))
			Expect(len(cidWrapper.Transactions)).To(Equal(2))
			Expect(seed_node.ListContainsString(cidWrapper.Transactions, mocks.MockCIDWrapper.Transactions[0])).To(BeTrue())
			Expect(seed_node.ListContainsString(cidWrapper.Transactions, mocks.MockCIDWrapper.Transactions[1])).To(BeTrue())
			Expect(len(cidWrapper.Receipts)).To(Equal(2))
			Expect(seed_node.ListContainsString(cidWrapper.Receipts, mocks.MockCIDWrapper.Receipts[0])).To(BeTrue())
			Expect(seed_node.ListContainsString(cidWrapper.Receipts, mocks.MockCIDWrapper.Receipts[1])).To(BeTrue())
			Expect(len(cidWrapper.StateNodes)).To(Equal(2))
			for _, stateNode := range cidWrapper.StateNodes {
				if stateNode.CID == "mockStateCID1" {
					Expect(stateNode.Key).To(Equal(mocks.ContractLeafKey.Hex()))
					Expect(stateNode.Leaf).To(Equal(true))
				}
				if stateNode.CID == "mockStateCID2" {
					Expect(stateNode.Key).To(Equal(mocks.AnotherContractLeafKey.Hex()))
					Expect(stateNode.Leaf).To(Equal(true))
				}
			}
			Expect(len(cidWrapper.StorageNodes)).To(Equal(1))
			Expect(cidWrapper.StorageNodes).To(Equal(mocks.MockCIDWrapper.StorageNodes))
		})
	})

	Describe("RetrieveCIDs", func() {
		It("Applies filters from the provided config.Subscription", func() {
			cidWrapper1, err := retriever.RetrieveCIDs(rctContractFilter, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(cidWrapper1.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper1.Headers)).To(Equal(0))
			Expect(len(cidWrapper1.Transactions)).To(Equal(0))
			Expect(len(cidWrapper1.StateNodes)).To(Equal(0))
			Expect(len(cidWrapper1.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper1.Receipts)).To(Equal(1))
			Expect(cidWrapper1.Receipts[0]).To(Equal("mockRctCID2"))

			cidWrapper2, err := retriever.RetrieveCIDs(rctTopicsFilter, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(cidWrapper2.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper2.Headers)).To(Equal(0))
			Expect(len(cidWrapper2.Transactions)).To(Equal(0))
			Expect(len(cidWrapper2.StateNodes)).To(Equal(0))
			Expect(len(cidWrapper2.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper2.Receipts)).To(Equal(1))
			Expect(cidWrapper2.Receipts[0]).To(Equal("mockRctCID1"))

			cidWrapper3, err := retriever.RetrieveCIDs(rctTopicsAndContractFilter, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(cidWrapper3.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper3.Headers)).To(Equal(0))
			Expect(len(cidWrapper3.Transactions)).To(Equal(0))
			Expect(len(cidWrapper3.StateNodes)).To(Equal(0))
			Expect(len(cidWrapper3.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper3.Receipts)).To(Equal(1))
			Expect(cidWrapper3.Receipts[0]).To(Equal("mockRctCID1"))

			cidWrapper4, err := retriever.RetrieveCIDs(rctContractsAndTopicFilter, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(cidWrapper4.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper4.Headers)).To(Equal(0))
			Expect(len(cidWrapper4.Transactions)).To(Equal(0))
			Expect(len(cidWrapper4.StateNodes)).To(Equal(0))
			Expect(len(cidWrapper4.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper4.Receipts)).To(Equal(1))
			Expect(cidWrapper4.Receipts[0]).To(Equal("mockRctCID2"))

			cidWrapper5, err := retriever.RetrieveCIDs(rctsForAllCollectedTrxs, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(cidWrapper5.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper5.Headers)).To(Equal(0))
			Expect(len(cidWrapper5.Transactions)).To(Equal(2))
			Expect(seed_node.ListContainsString(cidWrapper5.Transactions, "mockTrxCID1")).To(BeTrue())
			Expect(seed_node.ListContainsString(cidWrapper5.Transactions, "mockTrxCID2")).To(BeTrue())
			Expect(len(cidWrapper5.StateNodes)).To(Equal(0))
			Expect(len(cidWrapper5.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper5.Receipts)).To(Equal(2))
			Expect(seed_node.ListContainsString(cidWrapper5.Receipts, "mockRctCID1")).To(BeTrue())
			Expect(seed_node.ListContainsString(cidWrapper5.Receipts, "mockRctCID2")).To(BeTrue())

			cidWrapper6, err := retriever.RetrieveCIDs(rctsForSelectCollectedTrxs, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(cidWrapper6.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper6.Headers)).To(Equal(0))
			Expect(len(cidWrapper6.Transactions)).To(Equal(1))
			Expect(cidWrapper6.Transactions[0]).To(Equal("mockTrxCID2"))
			Expect(len(cidWrapper6.StateNodes)).To(Equal(0))
			Expect(len(cidWrapper6.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper6.Receipts)).To(Equal(1))
			Expect(cidWrapper6.Receipts[0]).To(Equal("mockRctCID2"))

			cidWrapper7, err := retriever.RetrieveCIDs(stateFilter, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(cidWrapper7.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper7.Headers)).To(Equal(0))
			Expect(len(cidWrapper7.Transactions)).To(Equal(0))
			Expect(len(cidWrapper7.Receipts)).To(Equal(0))
			Expect(len(cidWrapper7.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper7.StateNodes)).To(Equal(1))
			Expect(cidWrapper7.StateNodes[0]).To(Equal(ipfs.StateNodeCID{
				Leaf: true,
				Key:  mocks.ContractLeafKey.Hex(),
				CID:  "mockStateCID1",
			}))
		})
	})

	Describe("RetrieveFirstBlockNumber", func() {
		It("Gets the number of the first block that has data in the database", func() {
			num, err := retriever.RetrieveFirstBlockNumber()
			Expect(err).ToNot(HaveOccurred())
			Expect(num).To(Equal(int64(1)))
		})
	})

	Describe("RetrieveLastBlockNumber", func() {
		It("Gets the number of the latest block that has data in the database", func() {
			num, err := retriever.RetrieveLastBlockNumber()
			Expect(err).ToNot(HaveOccurred())
			Expect(num).To(Equal(int64(1)))
		})
	})
})
