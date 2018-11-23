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

package tusd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/generic/event_triggered/tusd"
	"github.com/vulcanize/vulcanizedb/examples/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
)

var blockID1 = int64(5428074)
var logID1 = int64(113)
var blockID2 = int64(5428405)
var logID2 = int64(100)

var fakeWatchedEvents = []*core.WatchedEvent{
	{
		LogID:       logID1,
		Name:        constants.BurnEvent.String(),
		BlockNumber: blockID1,
		Address:     constants.TusdContractAddress,
		TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
		Index:       110,
		Topic0:      constants.BurnEvent.Signature(),
		Topic1:      "0x000000000000000000000000000000000000000000000000000000000000af21",
		Topic2:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
		Topic3:      "",
		Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
	},
	{
		LogID:       logID2,
		Name:        constants.MintEvent.String(),
		BlockNumber: blockID2,
		Address:     constants.TusdContractAddress,
		TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
		Index:       110,
		Topic0:      constants.MintEvent.Signature(),
		Topic1:      "0x000000000000000000000000000000000000000000000000000000000000af21",
		Topic2:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
		Topic3:      "",
		Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
	},
}

var _ = Describe("Mock ERC20 transformer", func() {
	var mockERC20Converter mocks.MockERC20Converter
	var watchedEventsRepo mocks.MockWatchedEventsRepository
	var mockEventRepo mocks.MockEventRepo
	var filterRepo mocks.MockFilterRepository
	var transformer tusd.GenericTransformer

	BeforeEach(func() {
		mockERC20Converter = mocks.MockERC20Converter{}
		watchedEventsRepo = mocks.MockWatchedEventsRepository{}
		watchedEventsRepo.SetWatchedEvents(fakeWatchedEvents)
		mockEventRepo = mocks.MockEventRepo{}
		filterRepo = mocks.MockFilterRepository{}
		filters := constants.TusdGenericFilters

		transformer = tusd.GenericTransformer{
			Converter:              &mockERC20Converter,
			WatchedEventRepository: &watchedEventsRepo,
			FilterRepository:       filterRepo,
			Repository:             &mockEventRepo,
			Filters:                filters,
		}
	})

	It("calls the watched events repo with correct filter", func() {
		transformer.Execute()
		Expect(len(watchedEventsRepo.Names)).To(Equal(2))
		Expect(watchedEventsRepo.Names).To(ConsistOf([]string{constants.BurnEvent.String(), constants.MintEvent.String()}))
	})

	It("calls the mock ERC20 converter with the watched events", func() {
		transformer.Execute()
		Expect(len(mockERC20Converter.WatchedEvents)).To(Equal(2))
		Expect(mockERC20Converter.WatchedEvents).To(ConsistOf(fakeWatchedEvents))
	})

	It("converts a Burn and Mint entity to their models", func() {
		transformer.Execute()
		Expect(len(mockERC20Converter.BurnsToConvert)).To(Equal(1))
		Expect(mockERC20Converter.BurnsToConvert[0].Block).To(Equal(blockID1))

		Expect(len(mockERC20Converter.MintsToConvert)).To(Equal(1))
		Expect(mockERC20Converter.MintsToConvert[0].Block).To(Equal(blockID2))
	})

	It("persists Burn and Mint data for each watched Burn or Mint event", func() {
		transformer.Execute()
		Expect(len(mockEventRepo.BurnLogs)).To(Equal(1))
		Expect(len(mockEventRepo.MintLogs)).To(Equal(1))
		Expect(mockEventRepo.VulcanizeLogIDs).To(ConsistOf(logID1, logID2))
	})

})
