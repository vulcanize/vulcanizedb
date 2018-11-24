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
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/full/converter"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/retriever"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
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
	var dataStore repository.EventRepository
	var err error
	var info *contract.Contract
	var vulcanizeLogId int64
	var log *types.Log
	var r retriever.AddressRetriever
	var addresses map[common.Address]bool
	var wantedEvents = []string{"Transfer"}

	BeforeEach(func() {
		db, info = test_helpers.SetupTusdRepo(&vulcanizeLogId, wantedEvents, []string{})
		mockEvent.LogID = vulcanizeLogId

		event := info.Events["Transfer"]
		err = info.GenerateFilters()
		Expect(err).ToNot(HaveOccurred())

		c := converter.NewConverter(info)
		log, err = c.Convert(mockEvent, event)
		Expect(err).ToNot(HaveOccurred())

		dataStore = repository.NewEventRepository(db, types.FullSync)
		dataStore.PersistLogs([]types.Log{*log}, event, info.Address, info.Name)
		Expect(err).ToNot(HaveOccurred())

		r = retriever.NewAddressRetriever(db, types.FullSync)
	})

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	Describe("RetrieveTokenHolderAddresses", func() {
		It("Retrieves a list of token holder addresses from persisted event logs", func() {
			addresses, err = r.RetrieveTokenHolderAddresses(*info)
			Expect(err).ToNot(HaveOccurred())

			_, ok := addresses[common.HexToAddress("0x000000000000000000000000000000000000000000000000000000000000af21")]
			Expect(ok).To(Equal(true))

			_, ok = addresses[common.HexToAddress("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391")]
			Expect(ok).To(Equal(true))

			_, ok = addresses[common.HexToAddress("0x")]
			Expect(ok).To(Equal(false))

			_, ok = addresses[common.HexToAddress(constants.TusdContractAddress)]
			Expect(ok).To(Equal(false))

		})

		It("Returns empty list when empty contract info is used", func() {
			addresses, err = r.RetrieveTokenHolderAddresses(contract.Contract{})
			Expect(err).ToNot(HaveOccurred())
			Expect(len(addresses)).To(Equal(0))
		})
	})
})
