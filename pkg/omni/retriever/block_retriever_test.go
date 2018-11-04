// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retriever_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/omni/retriever"
)

var _ = Describe("Block Retriever Test", func() {
	var db *postgres.DB
	var r retriever.BlockRetriever
	var blockRepository repositories.BlockRepository

	BeforeEach(func() {
		var err error
		db, err = postgres.NewDB(config.Database{
			Hostname: "localhost",
			Name:     "vulcanize_private",
			Port:     5432,
		}, core.Node{})
		Expect(err).NotTo(HaveOccurred())

		blockRepository = *repositories.NewBlockRepository(db)

		r = retriever.NewBlockRetriever(db)
	})

	AfterEach(func() {
		db.Query(`DELETE FROM blocks`)
		db.Query(`DELETE FROM logs`)
		db.Query(`DELETE FROM transactions`)
		db.Query(`DELETE FROM receipts`)
	})

	It("Retrieve first block number of a contract from receipt if possible", func() {
		log := core.Log{
			BlockNumber: 2,
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
		}

		receipt1 := core.Receipt{
			TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
			ContractAddress: constants.TusdContractAddress,
			Logs:            []core.Log{},
		}

		receipt2 := core.Receipt{
			TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
			ContractAddress: constants.TusdContractAddress,
			Logs:            []core.Log{log},
		}

		transaction1 := core.Transaction{
			Hash:    "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
			Receipt: receipt1,
		}

		transaction2 := core.Transaction{
			Hash:    "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
			Receipt: receipt2,
		}

		block1 := core.Block{
			Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
			Number:       1,
			Transactions: []core.Transaction{transaction1},
		}

		block2 := core.Block{
			Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
			Number:       2,
			Transactions: []core.Transaction{transaction2},
		}

		blockRepository.CreateOrUpdateBlock(block1)
		blockRepository.CreateOrUpdateBlock(block2)

		i, err := r.RetrieveFirstBlock(constants.TusdContractAddress)
		Expect(err).NotTo(HaveOccurred())
		Expect(i).To(Equal(int64(1)))
	})

	It("Retrieves first block number of a contract from event logs", func() {
		log1 := core.Log{
			BlockNumber: 1,
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
		}

		log2 := core.Log{
			BlockNumber: 2,
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
		}

		receipt1 := core.Receipt{
			TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
			ContractAddress: "",
			Logs:            []core.Log{log1},
		}

		receipt2 := core.Receipt{
			TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
			ContractAddress: "",
			Logs:            []core.Log{log2},
		}

		transaction1 := core.Transaction{
			Hash:    "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
			Receipt: receipt1,
		}

		transaction2 := core.Transaction{
			Hash:    "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
			Receipt: receipt2,
		}

		block1 := core.Block{
			Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
			Number:       1,
			Transactions: []core.Transaction{transaction1},
		}

		block2 := core.Block{
			Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
			Number:       2,
			Transactions: []core.Transaction{transaction2},
		}

		blockRepository.CreateOrUpdateBlock(block1)
		blockRepository.CreateOrUpdateBlock(block2)

		i, err := r.RetrieveFirstBlock(constants.DaiContractAddress)
		Expect(err).NotTo(HaveOccurred())
		Expect(i).To(Equal(int64(1)))
	})

	It("Fails if a block cannot be found", func() {

		block1 := core.Block{
			Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
			Number:       1,
			Transactions: []core.Transaction{},
		}

		block2 := core.Block{
			Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
			Number:       2,
			Transactions: []core.Transaction{},
		}

		blockRepository.CreateOrUpdateBlock(block1)
		blockRepository.CreateOrUpdateBlock(block2)

		_, err := r.RetrieveFirstBlock(constants.DaiContractAddress)
		Expect(err).To(HaveOccurred())
	})
})
