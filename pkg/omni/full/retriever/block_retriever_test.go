// VulcanizeDB
// Copyright © 2018 Vulcanize

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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/omni/full/retriever"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers"
)

var _ = Describe("Block Retriever", func() {
	var db *postgres.DB
	var r retriever.BlockRetriever
	var blockRepository repositories.BlockRepository

	// Contains no contract address
	var block1 = core.Block{
		Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
		Number:       1,
		Transactions: []core.Transaction{},
	}

	BeforeEach(func() {
		db, _ = test_helpers.SetupDBandBC()
		blockRepository = *repositories.NewBlockRepository(db)
		r = retriever.NewBlockRetriever(db)
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
				Transactions: []core.Transaction{{
					Hash: "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
					Receipt: core.Receipt{
						TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
						ContractAddress: constants.TusdContractAddress,
						Logs:            []core.Log{},
					},
				}},
			}

			// Contains address in logs
			block3 := core.Block{
				Hash:   "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
				Number: 3,
				Transactions: []core.Transaction{{
					Hash: "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
					Receipt: core.Receipt{
						TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
						ContractAddress: constants.TusdContractAddress,
						Logs: []core.Log{{
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
				}},
			}

			blockRepository.CreateOrUpdateBlock(block1)
			blockRepository.CreateOrUpdateBlock(block2)
			blockRepository.CreateOrUpdateBlock(block3)

			i, err := r.RetrieveFirstBlock(strings.ToLower(constants.TusdContractAddress))
			Expect(err).NotTo(HaveOccurred())
			Expect(i).To(Equal(int64(2)))
		})

		It("Retrieves block number where contract first appears in event logs if it cannot find the address in a receipt", func() {

			block2 := core.Block{
				Hash:   "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
				Number: 2,
				Transactions: []core.Transaction{{
					Hash: "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
					Receipt: core.Receipt{
						TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
						ContractAddress: "",
						Logs: []core.Log{{
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
				}},
			}

			block3 := core.Block{
				Hash:   "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
				Number: 3,
				Transactions: []core.Transaction{{
					Hash: "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
					Receipt: core.Receipt{
						TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
						ContractAddress: "",
						Logs: []core.Log{{
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
				}},
			}

			blockRepository.CreateOrUpdateBlock(block1)
			blockRepository.CreateOrUpdateBlock(block2)
			blockRepository.CreateOrUpdateBlock(block3)

			i, err := r.RetrieveFirstBlock(constants.DaiContractAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(i).To(Equal(int64(2)))
		})

		It("Fails if the contract address cannot be found in any blocks", func() {
			block2 := core.Block{
				Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
				Number:       2,
				Transactions: []core.Transaction{},
			}

			block3 := core.Block{
				Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
				Number:       3,
				Transactions: []core.Transaction{},
			}

			blockRepository.CreateOrUpdateBlock(block1)
			blockRepository.CreateOrUpdateBlock(block2)
			blockRepository.CreateOrUpdateBlock(block3)

			_, err := r.RetrieveFirstBlock(constants.DaiContractAddress)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("RetrieveMostRecentBlock", func() {
		It("Retrieves the latest block", func() {
			block2 := core.Block{
				Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
				Number:       2,
				Transactions: []core.Transaction{},
			}

			block3 := core.Block{
				Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
				Number:       3,
				Transactions: []core.Transaction{},
			}

			blockRepository.CreateOrUpdateBlock(block1)
			blockRepository.CreateOrUpdateBlock(block2)
			blockRepository.CreateOrUpdateBlock(block3)

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
