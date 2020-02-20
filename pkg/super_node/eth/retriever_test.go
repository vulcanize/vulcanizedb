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
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	eth2 "github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

var (
	openFilter = &eth.SubscriptionSettings{
		Start:         big.NewInt(0),
		End:           big.NewInt(1),
		HeaderFilter:  eth.HeaderFilter{},
		TxFilter:      eth.TxFilter{},
		ReceiptFilter: eth.ReceiptFilter{},
		StateFilter:   eth.StateFilter{},
		StorageFilter: eth.StorageFilter{},
	}
	rctContractFilter = &eth.SubscriptionSettings{
		Start: big.NewInt(0),
		End:   big.NewInt(1),
		HeaderFilter: eth.HeaderFilter{
			Off: true,
		},
		TxFilter: eth.TxFilter{
			Off: true,
		},
		ReceiptFilter: eth.ReceiptFilter{
			Contracts: []string{mocks.AnotherAddress.String()},
		},
		StateFilter: eth.StateFilter{
			Off: true,
		},
		StorageFilter: eth.StorageFilter{
			Off: true,
		},
	}
	rctTopicsFilter = &eth.SubscriptionSettings{
		Start: big.NewInt(0),
		End:   big.NewInt(1),
		HeaderFilter: eth.HeaderFilter{
			Off: true,
		},
		TxFilter: eth.TxFilter{
			Off: true,
		},
		ReceiptFilter: eth.ReceiptFilter{
			Topics: [][]string{{"0x0000000000000000000000000000000000000000000000000000000000000004"}},
		},
		StateFilter: eth.StateFilter{
			Off: true,
		},
		StorageFilter: eth.StorageFilter{
			Off: true,
		},
	}
	rctTopicsAndContractFilter = &eth.SubscriptionSettings{
		Start: big.NewInt(0),
		End:   big.NewInt(1),
		HeaderFilter: eth.HeaderFilter{
			Off: true,
		},
		TxFilter: eth.TxFilter{
			Off: true,
		},
		ReceiptFilter: eth.ReceiptFilter{
			Topics: [][]string{
				{"0x0000000000000000000000000000000000000000000000000000000000000004"},
				{"0x0000000000000000000000000000000000000000000000000000000000000006"},
			},
			Contracts: []string{mocks.Address.String()},
		},
		StateFilter: eth.StateFilter{
			Off: true,
		},
		StorageFilter: eth.StorageFilter{
			Off: true,
		},
	}
	rctTopicsAndContractFilterFail = &eth.SubscriptionSettings{
		Start: big.NewInt(0),
		End:   big.NewInt(1),
		HeaderFilter: eth.HeaderFilter{
			Off: true,
		},
		TxFilter: eth.TxFilter{
			Off: true,
		},
		ReceiptFilter: eth.ReceiptFilter{
			Topics: [][]string{
				{"0x0000000000000000000000000000000000000000000000000000000000000004"},
				{"0x0000000000000000000000000000000000000000000000000000000000000007"}, // This topic won't match on the mocks.Address.String() contract receipt
			},
			Contracts: []string{mocks.Address.String()},
		},
		StateFilter: eth.StateFilter{
			Off: true,
		},
		StorageFilter: eth.StorageFilter{
			Off: true,
		},
	}
	rctContractsAndTopicFilter = &eth.SubscriptionSettings{
		Start: big.NewInt(0),
		End:   big.NewInt(1),
		HeaderFilter: eth.HeaderFilter{
			Off: true,
		},
		TxFilter: eth.TxFilter{
			Off: true,
		},
		ReceiptFilter: eth.ReceiptFilter{
			Topics:    [][]string{{"0x0000000000000000000000000000000000000000000000000000000000000005"}},
			Contracts: []string{mocks.Address.String(), mocks.AnotherAddress.String()},
		},
		StateFilter: eth.StateFilter{
			Off: true,
		},
		StorageFilter: eth.StorageFilter{
			Off: true,
		},
	}
	rctsForAllCollectedTrxs = &eth.SubscriptionSettings{
		Start: big.NewInt(0),
		End:   big.NewInt(1),
		HeaderFilter: eth.HeaderFilter{
			Off: true,
		},
		TxFilter: eth.TxFilter{}, // Trx filter open so we will collect all trxs, therefore we will also collect all corresponding rcts despite rct filter
		ReceiptFilter: eth.ReceiptFilter{
			MatchTxs:  true,
			Topics:    [][]string{{"0x0000000000000000000000000000000000000000000000000000000000000006"}}, // Topic0 isn't one of the topic0s we have
			Contracts: []string{"0x0000000000000000000000000000000000000002"},                             // Contract isn't one of the contracts we have
		},
		StateFilter: eth.StateFilter{
			Off: true,
		},
		StorageFilter: eth.StorageFilter{
			Off: true,
		},
	}
	rctsForSelectCollectedTrxs = &eth.SubscriptionSettings{
		Start: big.NewInt(0),
		End:   big.NewInt(1),
		HeaderFilter: eth.HeaderFilter{
			Off: true,
		},
		TxFilter: eth.TxFilter{
			Dst: []string{mocks.AnotherAddress.String()}, // We only filter for one of the trxs so we will only get the one corresponding receipt
		},
		ReceiptFilter: eth.ReceiptFilter{
			MatchTxs:  true,
			Topics:    [][]string{{"0x0000000000000000000000000000000000000000000000000000000000000006"}}, // Topic0 isn't one of the topic0s we have
			Contracts: []string{"0x0000000000000000000000000000000000000002"},                             // Contract isn't one of the contracts we have
		},
		StateFilter: eth.StateFilter{
			Off: true,
		},
		StorageFilter: eth.StorageFilter{
			Off: true,
		},
	}
	stateFilter = &eth.SubscriptionSettings{
		Start: big.NewInt(0),
		End:   big.NewInt(1),
		HeaderFilter: eth.HeaderFilter{
			Off: true,
		},
		TxFilter: eth.TxFilter{
			Off: true,
		},
		ReceiptFilter: eth.ReceiptFilter{
			Off: true,
		},
		StateFilter: eth.StateFilter{
			Addresses: []string{mocks.Address.Hex()},
		},
		StorageFilter: eth.StorageFilter{
			Off: true,
		},
	}
)

var _ = Describe("Retriever", func() {
	var (
		db        *postgres.DB
		repo      *eth2.CIDIndexer
		retriever *eth2.CIDRetriever
	)
	BeforeEach(func() {
		var err error
		db, err = shared.SetupDB()
		Expect(err).ToNot(HaveOccurred())
		repo = eth2.NewCIDIndexer(db)
		retriever = eth2.NewCIDRetriever(db)
	})
	AfterEach(func() {
		eth.TearDownDB(db)
	})

	Describe("Retrieve", func() {
		BeforeEach(func() {
			err := repo.Index(mocks.MockCIDPayload)
			Expect(err).ToNot(HaveOccurred())
		})
		It("Retrieves all CIDs for the given blocknumber when provided an open filter", func() {
			cids, empty, err := retriever.Retrieve(openFilter, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(empty).ToNot(BeTrue())
			cidWrapper, ok := cids.(*eth.CIDWrapper)
			Expect(ok).To(BeTrue())
			Expect(cidWrapper.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper.Headers)).To(Equal(1))
			expectedHeaderCIDs := mocks.MockCIDWrapper.Headers
			expectedHeaderCIDs[0].ID = cidWrapper.Headers[0].ID
			expectedHeaderCIDs[0].NodeID = cidWrapper.Headers[0].NodeID
			Expect(cidWrapper.Headers).To(Equal(expectedHeaderCIDs))
			Expect(len(cidWrapper.Transactions)).To(Equal(2))
			Expect(eth.TxModelsContainsCID(cidWrapper.Transactions, mocks.MockCIDWrapper.Transactions[0].CID)).To(BeTrue())
			Expect(eth.TxModelsContainsCID(cidWrapper.Transactions, mocks.MockCIDWrapper.Transactions[1].CID)).To(BeTrue())
			Expect(len(cidWrapper.Receipts)).To(Equal(2))
			Expect(eth.ReceiptModelsContainsCID(cidWrapper.Receipts, mocks.MockCIDWrapper.Receipts[0].CID)).To(BeTrue())
			Expect(eth.ReceiptModelsContainsCID(cidWrapper.Receipts, mocks.MockCIDWrapper.Receipts[1].CID)).To(BeTrue())
			Expect(len(cidWrapper.StateNodes)).To(Equal(2))
			for _, stateNode := range cidWrapper.StateNodes {
				if stateNode.CID == mocks.State1CID.String() {
					Expect(stateNode.StateKey).To(Equal(mocks.ContractLeafKey.Hex()))
					Expect(stateNode.Leaf).To(Equal(true))
				}
				if stateNode.CID == mocks.State2CID.String() {
					Expect(stateNode.StateKey).To(Equal(mocks.AnotherContractLeafKey.Hex()))
					Expect(stateNode.Leaf).To(Equal(true))
				}
			}
			Expect(len(cidWrapper.StorageNodes)).To(Equal(1))
			expectedStorageNodeCIDs := mocks.MockCIDWrapper.StorageNodes
			expectedStorageNodeCIDs[0].ID = cidWrapper.StorageNodes[0].ID
			expectedStorageNodeCIDs[0].StateID = cidWrapper.StorageNodes[0].StateID
			Expect(cidWrapper.StorageNodes).To(Equal(expectedStorageNodeCIDs))
		})

		It("Applies filters from the provided config.Subscription", func() {
			cids1, empty, err := retriever.Retrieve(rctContractFilter, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(empty).ToNot(BeTrue())
			cidWrapper1, ok := cids1.(*eth.CIDWrapper)
			Expect(ok).To(BeTrue())
			Expect(cidWrapper1.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper1.Headers)).To(Equal(0))
			Expect(len(cidWrapper1.Transactions)).To(Equal(0))
			Expect(len(cidWrapper1.StateNodes)).To(Equal(0))
			Expect(len(cidWrapper1.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper1.Receipts)).To(Equal(1))
			expectedReceiptCID := mocks.MockCIDWrapper.Receipts[1]
			expectedReceiptCID.ID = cidWrapper1.Receipts[0].ID
			expectedReceiptCID.TxID = cidWrapper1.Receipts[0].TxID
			Expect(cidWrapper1.Receipts[0]).To(Equal(expectedReceiptCID))

			cids2, empty, err := retriever.Retrieve(rctTopicsFilter, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(empty).ToNot(BeTrue())
			cidWrapper2, ok := cids2.(*eth.CIDWrapper)
			Expect(ok).To(BeTrue())
			Expect(cidWrapper2.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper2.Headers)).To(Equal(0))
			Expect(len(cidWrapper2.Transactions)).To(Equal(0))
			Expect(len(cidWrapper2.StateNodes)).To(Equal(0))
			Expect(len(cidWrapper2.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper2.Receipts)).To(Equal(1))
			expectedReceiptCID = mocks.MockCIDWrapper.Receipts[0]
			expectedReceiptCID.ID = cidWrapper2.Receipts[0].ID
			expectedReceiptCID.TxID = cidWrapper2.Receipts[0].TxID
			Expect(cidWrapper2.Receipts[0]).To(Equal(expectedReceiptCID))

			cids3, empty, err := retriever.Retrieve(rctTopicsAndContractFilter, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(empty).ToNot(BeTrue())
			cidWrapper3, ok := cids3.(*eth.CIDWrapper)
			Expect(ok).To(BeTrue())
			Expect(cidWrapper3.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper3.Headers)).To(Equal(0))
			Expect(len(cidWrapper3.Transactions)).To(Equal(0))
			Expect(len(cidWrapper3.StateNodes)).To(Equal(0))
			Expect(len(cidWrapper3.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper3.Receipts)).To(Equal(1))
			expectedReceiptCID = mocks.MockCIDWrapper.Receipts[0]
			expectedReceiptCID.ID = cidWrapper3.Receipts[0].ID
			expectedReceiptCID.TxID = cidWrapper3.Receipts[0].TxID
			Expect(cidWrapper3.Receipts[0]).To(Equal(expectedReceiptCID))

			cids4, empty, err := retriever.Retrieve(rctContractsAndTopicFilter, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(empty).ToNot(BeTrue())
			cidWrapper4, ok := cids4.(*eth.CIDWrapper)
			Expect(ok).To(BeTrue())
			Expect(cidWrapper4.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper4.Headers)).To(Equal(0))
			Expect(len(cidWrapper4.Transactions)).To(Equal(0))
			Expect(len(cidWrapper4.StateNodes)).To(Equal(0))
			Expect(len(cidWrapper4.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper4.Receipts)).To(Equal(1))
			expectedReceiptCID = mocks.MockCIDWrapper.Receipts[1]
			expectedReceiptCID.ID = cidWrapper4.Receipts[0].ID
			expectedReceiptCID.TxID = cidWrapper4.Receipts[0].TxID
			Expect(cidWrapper4.Receipts[0]).To(Equal(expectedReceiptCID))

			cids5, empty, err := retriever.Retrieve(rctsForAllCollectedTrxs, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(empty).ToNot(BeTrue())
			cidWrapper5, ok := cids5.(*eth.CIDWrapper)
			Expect(ok).To(BeTrue())
			Expect(cidWrapper5.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper5.Headers)).To(Equal(0))
			Expect(len(cidWrapper5.Transactions)).To(Equal(2))
			Expect(eth.TxModelsContainsCID(cidWrapper5.Transactions, mocks.Trx1CID.String())).To(BeTrue())
			Expect(eth.TxModelsContainsCID(cidWrapper5.Transactions, mocks.Trx2CID.String())).To(BeTrue())
			Expect(len(cidWrapper5.StateNodes)).To(Equal(0))
			Expect(len(cidWrapper5.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper5.Receipts)).To(Equal(2))
			Expect(eth.ReceiptModelsContainsCID(cidWrapper5.Receipts, mocks.Rct1CID.String())).To(BeTrue())
			Expect(eth.ReceiptModelsContainsCID(cidWrapper5.Receipts, mocks.Rct2CID.String())).To(BeTrue())

			cids6, empty, err := retriever.Retrieve(rctsForSelectCollectedTrxs, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(empty).ToNot(BeTrue())
			cidWrapper6, ok := cids6.(*eth.CIDWrapper)
			Expect(ok).To(BeTrue())
			Expect(cidWrapper6.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper6.Headers)).To(Equal(0))
			Expect(len(cidWrapper6.Transactions)).To(Equal(1))
			expectedTxCID := mocks.MockCIDWrapper.Transactions[1]
			expectedTxCID.ID = cidWrapper6.Transactions[0].ID
			expectedTxCID.HeaderID = cidWrapper6.Transactions[0].HeaderID
			Expect(cidWrapper6.Transactions[0]).To(Equal(expectedTxCID))
			Expect(len(cidWrapper6.StateNodes)).To(Equal(0))
			Expect(len(cidWrapper6.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper6.Receipts)).To(Equal(1))
			expectedReceiptCID = mocks.MockCIDWrapper.Receipts[1]
			expectedReceiptCID.ID = cidWrapper6.Receipts[0].ID
			expectedReceiptCID.TxID = cidWrapper6.Receipts[0].TxID
			Expect(cidWrapper6.Receipts[0]).To(Equal(expectedReceiptCID))

			cids7, empty, err := retriever.Retrieve(stateFilter, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(empty).ToNot(BeTrue())
			cidWrapper7, ok := cids7.(*eth.CIDWrapper)
			Expect(ok).To(BeTrue())
			Expect(cidWrapper7.BlockNumber).To(Equal(mocks.MockCIDWrapper.BlockNumber))
			Expect(len(cidWrapper7.Headers)).To(Equal(0))
			Expect(len(cidWrapper7.Transactions)).To(Equal(0))
			Expect(len(cidWrapper7.Receipts)).To(Equal(0))
			Expect(len(cidWrapper7.StorageNodes)).To(Equal(0))
			Expect(len(cidWrapper7.StateNodes)).To(Equal(1))
			Expect(cidWrapper7.StateNodes[0]).To(Equal(eth.StateNodeModel{
				ID:       cidWrapper7.StateNodes[0].ID,
				HeaderID: cidWrapper7.StateNodes[0].HeaderID,
				Leaf:     true,
				StateKey: mocks.ContractLeafKey.Hex(),
				CID:      mocks.State1CID.String(),
			}))

			_, empty, err = retriever.Retrieve(rctTopicsAndContractFilterFail, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(empty).To(BeTrue())
		})
	})

	Describe("RetrieveFirstBlockNumber", func() {
		It("Gets the number of the first block that has data in the database", func() {
			err := repo.Index(mocks.MockCIDPayload)
			Expect(err).ToNot(HaveOccurred())
			num, err := retriever.RetrieveFirstBlockNumber()
			Expect(err).ToNot(HaveOccurred())
			Expect(num).To(Equal(int64(1)))
		})

		It("Gets the number of the first block that has data in the database", func() {
			payload := *mocks.MockCIDPayload
			payload.HeaderCID.BlockNumber = "1010101"
			err := repo.Index(&payload)
			Expect(err).ToNot(HaveOccurred())
			num, err := retriever.RetrieveFirstBlockNumber()
			Expect(err).ToNot(HaveOccurred())
			Expect(num).To(Equal(int64(1010101)))
		})

		It("Gets the number of the first block that has data in the database", func() {
			payload1 := *mocks.MockCIDPayload
			payload1.HeaderCID.BlockNumber = "1010101"
			payload2 := payload1
			payload2.HeaderCID.BlockNumber = "5"
			err := repo.Index(&payload1)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Index(&payload2)
			Expect(err).ToNot(HaveOccurred())
			num, err := retriever.RetrieveFirstBlockNumber()
			Expect(err).ToNot(HaveOccurred())
			Expect(num).To(Equal(int64(5)))
		})
	})

	Describe("RetrieveLastBlockNumber", func() {
		It("Gets the number of the latest block that has data in the database", func() {
			err := repo.Index(mocks.MockCIDPayload)
			Expect(err).ToNot(HaveOccurred())
			num, err := retriever.RetrieveLastBlockNumber()
			Expect(err).ToNot(HaveOccurred())
			Expect(num).To(Equal(int64(1)))
		})

		It("Gets the number of the latest block that has data in the database", func() {
			payload := *mocks.MockCIDPayload
			payload.HeaderCID.BlockNumber = "1010101"
			err := repo.Index(&payload)
			Expect(err).ToNot(HaveOccurred())
			num, err := retriever.RetrieveLastBlockNumber()
			Expect(err).ToNot(HaveOccurred())
			Expect(num).To(Equal(int64(1010101)))
		})

		It("Gets the number of the latest block that has data in the database", func() {
			payload1 := *mocks.MockCIDPayload
			payload1.HeaderCID.BlockNumber = "1010101"
			payload2 := payload1
			payload2.HeaderCID.BlockNumber = "5"
			err := repo.Index(&payload1)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Index(&payload2)
			Expect(err).ToNot(HaveOccurred())
			num, err := retriever.RetrieveLastBlockNumber()
			Expect(err).ToNot(HaveOccurred())
			Expect(num).To(Equal(int64(1010101)))
		})
	})

	Describe("RetrieveGapsInData", func() {
		It("Doesn't return gaps if there are none", func() {
			payload1 := *mocks.MockCIDPayload
			payload1.HeaderCID.BlockNumber = "2"
			payload2 := payload1
			payload2.HeaderCID.BlockNumber = "3"
			err := repo.Index(mocks.MockCIDPayload)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Index(&payload1)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Index(&payload2)
			Expect(err).ToNot(HaveOccurred())
			gaps, err := retriever.RetrieveGapsInData()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(gaps)).To(Equal(0))
		})

		It("Doesn't return the gap from 0 to the earliest block", func() {
			payload := *mocks.MockCIDPayload
			payload.HeaderCID.BlockNumber = "5"
			err := repo.Index(&payload)
			Expect(err).ToNot(HaveOccurred())
			gaps, err := retriever.RetrieveGapsInData()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(gaps)).To(Equal(0))
		})

		It("Finds gap between two entries", func() {
			payload1 := *mocks.MockCIDPayload
			payload1.HeaderCID.BlockNumber = "1010101"
			payload2 := payload1
			payload2.HeaderCID.BlockNumber = "5"
			err := repo.Index(&payload1)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Index(&payload2)
			Expect(err).ToNot(HaveOccurred())
			gaps, err := retriever.RetrieveGapsInData()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(gaps)).To(Equal(1))
			Expect(gaps[0].Start).To(Equal(uint64(6)))
			Expect(gaps[0].Stop).To(Equal(uint64(1010100)))
		})

		It("Finds gaps between multiple entries", func() {
			payload1 := *mocks.MockCIDPayload
			payload1.HeaderCID.BlockNumber = "1010101"
			payload2 := payload1
			payload2.HeaderCID.BlockNumber = "5"
			payload3 := payload2
			payload3.HeaderCID.BlockNumber = "100"
			payload4 := payload3
			payload4.HeaderCID.BlockNumber = "101"
			payload5 := payload4
			payload5.HeaderCID.BlockNumber = "102"
			payload6 := payload5
			payload6.HeaderCID.BlockNumber = "1000"
			err := repo.Index(&payload1)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Index(&payload2)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Index(&payload3)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Index(&payload4)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Index(&payload5)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Index(&payload6)
			Expect(err).ToNot(HaveOccurred())
			gaps, err := retriever.RetrieveGapsInData()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(gaps)).To(Equal(3))
			Expect(shared.ListContainsGap(gaps, shared.Gap{Start: 6, Stop: 99})).To(BeTrue())
			Expect(shared.ListContainsGap(gaps, shared.Gap{Start: 103, Stop: 999})).To(BeTrue())
			Expect(shared.ListContainsGap(gaps, shared.Gap{Start: 1001, Stop: 1010100})).To(BeTrue())
		})
	})
})
