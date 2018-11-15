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

package converter_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/omni/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/converter"
	"github.com/vulcanize/vulcanizedb/pkg/omni/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/helpers/test_helpers"
)

var mockEvent = core.WatchedEvent{
	LogID:       1,
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

var _ = Describe("Converter", func() {
	var info *contract.Contract
	var wantedEvents = []string{"Transfer"}
	var err error

	BeforeEach(func() {
		info = test_helpers.SetupTusdContract(wantedEvents, []string{})
	})

	Describe("Update", func() {
		It("Updates contract info held by the converter", func() {
			c := converter.NewConverter(info)
			i := c.CheckInfo()
			Expect(i).To(Equal(info))

			info := test_helpers.SetupTusdContract([]string{}, []string{})
			c.Update(info)
			i = c.CheckInfo()
			Expect(i).To(Equal(info))
		})
	})

	Describe("Convert", func() {
		It("Converts a watched event log to mapping of event input names to values", func() {
			_, ok := info.Events["Approval"]
			Expect(ok).To(Equal(false))

			event, ok := info.Events["Transfer"]
			Expect(ok).To(Equal(true))
			err = info.GenerateFilters()
			Expect(err).ToNot(HaveOccurred())

			c := converter.NewConverter(info)
			err = c.Convert(mockEvent, event)
			Expect(err).ToNot(HaveOccurred())

			from := common.HexToAddress("0x000000000000000000000000000000000000000000000000000000000000af21")
			to := common.HexToAddress("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391")
			value := helpers.BigFromString("1097077688018008265106216665536940668749033598146")

			v := event.Logs[1].Values["value"]

			Expect(event.Logs[1].Values["to"]).To(Equal(to.String()))
			Expect(event.Logs[1].Values["from"]).To(Equal(from.String()))
			Expect(v).To(Equal(value.String()))
		})

		It("Fails with an empty contract", func() {
			event := info.Events["Transfer"]
			c := converter.NewConverter(&contract.Contract{})
			err = c.Convert(mockEvent, event)
			Expect(err).To(HaveOccurred())
		})
	})
})
