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

package retriever_test

import (
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/full/retriever"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

var _ = Describe("Block Retriever", func() {
	var db *postgres.DB
	var r retriever.BlockRetriever
	var rawTransaction []byte
	var blockRepository repositories.BlockRepository

	// Contains no contract address
	var block1 = core.Block{
		Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
		Number:       1,
		Transactions: []core.TransactionModel{},
	}

	BeforeEach(func() {
		db, _ = test_helpers.SetupDBandBC()
		blockRepository = *repositories.NewBlockRepository(db)
		r = retriever.NewBlockRetriever(db)
		gethTransaction := types.Transaction{}
		var err error
		rawTransaction, err = gethTransaction.MarshalJSON()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	Describe("RetrieveFirstBlock", func() {
		It("Retrieves block number where contract first appears in receipt, if available", func() {
			// Contains the address in the receipt
			block2 := core.Block{
				Hash:   "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
				Number: 2,
				Transactions: []core.TransactionModel{{
					Hash:     "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
					GasPrice: 0,
					GasLimit: 0,
					Nonce:    0,
					Raw:      rawTransaction,
					Receipt: core.Receipt{
						TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
						ContractAddress: constants.TusdContractAddress,
						Logs:            []core.FullSyncLog{},
					},
					TxIndex: 0,
					Value:   "0",
				}},
			}

			// Contains address in logs
			block3 := core.Block{
				Hash:   "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
				Number: 3,
				Transactions: []core.TransactionModel{{
					Hash:     "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
					GasPrice: 0,
					GasLimit: 0,
					Nonce:    0,
					Raw:      rawTransaction,
					Receipt: core.Receipt{
						TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
						ContractAddress: constants.TusdContractAddress,
						Logs: []core.FullSyncLog{{
							BlockNumber: 3,
							TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
							Address:     constants.TusdContractAddress,
							Topics: core.Topics{
								constants.TransferEvent.Signature(),
								"0x000000000000000000000000000000000000000000000000000000000000af21",
								"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
								"",
							},
							Index: 1,
							Data:  "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
						}},
					},
					TxIndex: 0,
					Value:   "0",
				}},
			}

			_, insertErrOne := blockRepository.CreateOrUpdateBlock(block1)
			Expect(insertErrOne).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(block2)
			Expect(insertErrTwo).NotTo(HaveOccurred())
			_, insertErrThree := blockRepository.CreateOrUpdateBlock(block3)
			Expect(insertErrThree).NotTo(HaveOccurred())

			i, err := r.RetrieveFirstBlock(strings.ToLower(constants.TusdContractAddress))
			Expect(err).NotTo(HaveOccurred())
			Expect(i).To(Equal(int64(2)))
		})

		It("Retrieves block number where contract first appears in event logs if it cannot find the address in a receipt", func() {
			block2 := core.Block{
				Hash:   "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
				Number: 2,
				Transactions: []core.TransactionModel{{
					Hash:     "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
					GasPrice: 0,
					GasLimit: 0,
					Nonce:    0,
					Raw:      rawTransaction,
					Receipt: core.Receipt{
						TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
						ContractAddress: "",
						Logs: []core.FullSyncLog{{
							BlockNumber: 2,
							TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
							Address:     constants.DaiContractAddress,
							Topics: core.Topics{
								constants.TransferEvent.Signature(),
								"0x000000000000000000000000000000000000000000000000000000000000af21",
								"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
								"",
							},
							Index: 1,
							Data:  "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
						}},
					},
					TxIndex: 0,
					Value:   "0",
				}},
			}

			block3 := core.Block{
				Hash:   "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
				Number: 3,
				Transactions: []core.TransactionModel{{
					Hash:     "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
					GasPrice: 0,
					GasLimit: 0,
					Nonce:    0,
					Raw:      rawTransaction,
					Receipt: core.Receipt{
						TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
						ContractAddress: "",
						Logs: []core.FullSyncLog{{
							BlockNumber: 3,
							TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
							Address:     constants.DaiContractAddress,
							Topics: core.Topics{
								constants.TransferEvent.Signature(),
								"0x000000000000000000000000000000000000000000000000000000000000af21",
								"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
								"",
							},
							Index: 1,
							Data:  "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
						}},
					},
					TxIndex: 0,
					Value:   "0",
				}},
			}

			_, insertErrOne := blockRepository.CreateOrUpdateBlock(block1)
			Expect(insertErrOne).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(block2)
			Expect(insertErrTwo).NotTo(HaveOccurred())
			_, insertErrThree := blockRepository.CreateOrUpdateBlock(block3)
			Expect(insertErrThree).NotTo(HaveOccurred())

			i, err := r.RetrieveFirstBlock(constants.DaiContractAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(i).To(Equal(int64(2)))
		})

		It("Fails if the contract address cannot be found in any blocks", func() {
			block2 := core.Block{
				Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
				Number:       2,
				Transactions: []core.TransactionModel{},
			}

			block3 := core.Block{
				Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
				Number:       3,
				Transactions: []core.TransactionModel{},
			}

			_, insertErrOne := blockRepository.CreateOrUpdateBlock(block1)
			Expect(insertErrOne).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(block2)
			Expect(insertErrTwo).NotTo(HaveOccurred())
			_, insertErrThree := blockRepository.CreateOrUpdateBlock(block3)
			Expect(insertErrThree).NotTo(HaveOccurred())

			_, err := r.RetrieveFirstBlock(constants.DaiContractAddress)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("RetrieveMostRecentBlock", func() {
		It("Retrieves the latest block", func() {
			block2 := core.Block{
				Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
				Number:       2,
				Transactions: []core.TransactionModel{},
			}

			block3 := core.Block{
				Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
				Number:       3,
				Transactions: []core.TransactionModel{},
			}

			_, insertErrOne := blockRepository.CreateOrUpdateBlock(block1)
			Expect(insertErrOne).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(block2)
			Expect(insertErrTwo).NotTo(HaveOccurred())
			_, insertErrThree := blockRepository.CreateOrUpdateBlock(block3)
			Expect(insertErrThree).NotTo(HaveOccurred())

			i, err := r.RetrieveMostRecentBlock()
			Expect(err).ToNot(HaveOccurred())
			Expect(i).To(Equal(int64(3)))
		})

		It("Fails if it cannot retrieve the latest block", func() {
			i, err := r.RetrieveMostRecentBlock()
			Expect(err).To(HaveOccurred())
			Expect(i).To(Equal(int64(0)))
		})
	})
})
