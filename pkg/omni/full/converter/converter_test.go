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

	"github.com/vulcanize/vulcanizedb/pkg/omni/full/converter"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers"
)

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
			Expect(c.ContractInfo).To(Equal(info))

			info := test_helpers.SetupTusdContract([]string{}, []string{})
			c.Update(info)
			Expect(c.ContractInfo).To(Equal(info))
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
			log, err := c.Convert(test_helpers.MockTranferEvent, event)
			Expect(err).ToNot(HaveOccurred())

			from := common.HexToAddress("0x000000000000000000000000000000000000000000000000000000000000af21")
			to := common.HexToAddress("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391")
			value := helpers.BigFromString("1097077688018008265106216665536940668749033598146")

			v := log.Values["value"]

			Expect(log.Values["to"]).To(Equal(to.String()))
			Expect(log.Values["from"]).To(Equal(from.String()))
			Expect(v).To(Equal(value.String()))
		})

		It("Fails with an empty contract", func() {
			event := info.Events["Transfer"]
			c := converter.NewConverter(&contract.Contract{})
			_, err = c.Convert(test_helpers.MockTranferEvent, event)
			Expect(err).To(HaveOccurred())
		})
	})
})
