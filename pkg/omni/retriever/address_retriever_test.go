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
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/converter"
	"github.com/vulcanize/vulcanizedb/pkg/omni/parser"
	"github.com/vulcanize/vulcanizedb/pkg/omni/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/retriever"
)

var mockEvent = core.WatchedEvent{
	Name:        constants.TransferEvent.String(),
	BlockNumber: 5488076,
	Address:     constants.TusdContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topic0:      constants.TransferEvent.Signature(),
	Topic1:      "0x000000000000000000000000000000000000000000000000000000000000af21",
	Topic2:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
	Topic3:      "",
	Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var _ = Describe("Address Retriever Test", func() {
	var db *postgres.DB
	var logRepository repositories.LogRepository
	var blockRepository repositories.BlockRepository
	var receiptRepository repositories.ReceiptRepository
	var dataStore repository.DataStore
	var err error
	var info *contract.Contract
	var blockNumber int64
	var blockId int64
	var vulcanizeLogId int64
	var r retriever.AddressRetriever
	var addresses map[common.Address]bool
	rand.Seed(time.Now().UnixNano())

	BeforeEach(func() {
		db, err = postgres.NewDB(config.Database{
			Hostname: "localhost",
			Name:     "vulcanize_private",
			Port:     5432,
		}, core.Node{})
		Expect(err).NotTo(HaveOccurred())

		receiptRepository = repositories.ReceiptRepository{DB: db}
		logRepository = repositories.LogRepository{DB: db}
		blockRepository = *repositories.NewBlockRepository(db)

		blockNumber = rand.Int63()
		blockId = test_helpers.CreateBlock(blockNumber, blockRepository)

		log := core.Log{}
		logs := []core.Log{log}
		receipt := core.Receipt{
			Logs: logs,
		}
		receipts := []core.Receipt{receipt}

		err = receiptRepository.CreateReceiptsAndLogs(blockId, receipts)
		Expect(err).ToNot(HaveOccurred())

		err = logRepository.Get(&vulcanizeLogId, `SELECT id FROM logs`)
		Expect(err).ToNot(HaveOccurred())

		p := parser.NewParser("")
		err = p.Parse(constants.TusdContractAddress)
		Expect(err).ToNot(HaveOccurred())

		info = &contract.Contract{
			Name:          "TrueUSD",
			Address:       constants.TusdContractAddress,
			Abi:           p.Abi(),
			ParsedAbi:     p.ParsedAbi(),
			StartingBlock: 5197514,
			Events:        p.GetEvents(),
			Methods:       p.GetMethods(),
			Addresses:     map[string]bool{},
		}

		event := info.Events["Transfer"]
		err = info.GenerateFilters([]string{"Transfer"})
		Expect(err).ToNot(HaveOccurred())
		c := converter.NewConverter(*info)
		mockEvent.LogID = vulcanizeLogId
		err = c.Convert(mockEvent, event)
		Expect(err).ToNot(HaveOccurred())

		dataStore = repository.NewDataStore(db)
		err = dataStore.PersistEvents(info)
		Expect(err).ToNot(HaveOccurred())

		r = retriever.NewAddressRetriever(db)
	})

	AfterEach(func() {
		db.Query(`DELETE FROM blocks`)
		db.Query(`DELETE FROM logs`)
		db.Query(`DELETE FROM transactions`)
		db.Query(`DELETE FROM receipts`)
		db.Query(`DROP SCHEMA IF EXISTS trueusd CASCADE`)
	})

	It("Retrieves a list of token holder addresses", func() {
		addresses, err = r.RetrieveTokenHolderAddresses(*info)
		Expect(err).ToNot(HaveOccurred())

		_, ok := addresses[common.HexToAddress("0x000000000000000000000000000000000000000000000000000000000000af21")]
		Expect(ok).To(Equal(true))

		_, ok = addresses[common.HexToAddress("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391")]
		Expect(ok).To(Equal(true))

		_, ok = addresses[common.HexToAddress("0x")]
		Expect(ok).To(Equal(false))

	})

	It("Returns empty list when empty contract info is used", func() {
		addresses, err = r.RetrieveTokenHolderAddresses(contract.Contract{})
		Expect(err).ToNot(HaveOccurred())
		Expect(len(addresses)).To(Equal(0))
	})
})
