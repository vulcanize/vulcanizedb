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
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/omni/light/converter"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers/mocks"
)

var _ = Describe("Converter", func() {
	var con *contract.Contract
	var wantedEvents = []string{"Transfer", "Mint"}
	var err error

	BeforeEach(func() {
		con = test_helpers.SetupTusdContract(wantedEvents, []string{})
	})

	Describe("Update", func() {
		It("Updates contract info held by the converter", func() {
			c := converter.NewConverter(con)
			Expect(c.ContractInfo).To(Equal(con))

			info := test_helpers.SetupTusdContract([]string{}, []string{})
			c.Update(info)
			Expect(c.ContractInfo).To(Equal(info))
		})
	})

	Describe("Convert", func() {
		It("Converts a watched event log to mapping of event input names to values", func() {
			_, ok := con.Events["Approval"]
			Expect(ok).To(Equal(false))

			event, ok := con.Events["Transfer"]
			Expect(ok).To(Equal(true))

			c := converter.NewConverter(con)
			logs, err := c.Convert([]types.Log{mocks.MockTransferLog1, mocks.MockTransferLog2}, event, 232)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(2))

			sender1 := common.HexToAddress("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391")
			sender2 := common.HexToAddress("0x000000000000000000000000000000000000000000000000000000000000af21")
			value := helpers.BigFromString("1097077688018008265106216665536940668749033598146")

			Expect(logs[0].Values["to"]).To(Equal(sender1.String()))
			Expect(logs[0].Values["from"]).To(Equal(sender2.String()))
			Expect(logs[0].Values["value"]).To(Equal(value.String()))
			Expect(logs[0].Id).To(Equal(int64(232)))
			Expect(logs[1].Values["to"]).To(Equal(sender2.String()))
			Expect(logs[1].Values["from"]).To(Equal(sender1.String()))
			Expect(logs[1].Values["value"]).To(Equal(value.String()))
			Expect(logs[1].Id).To(Equal(int64(232)))
		})

		It("Keeps track of addresses it sees to grow a token holder address list for the contract", func() {
			event, ok := con.Events["Transfer"]
			Expect(ok).To(Equal(true))

			c := converter.NewConverter(con)
			_, err := c.Convert([]types.Log{mocks.MockTransferLog1, mocks.MockTransferLog2}, event, 232)
			Expect(err).ToNot(HaveOccurred())

			b, ok := con.TknHolderAddrs["0x000000000000000000000000000000000000Af21"]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = con.TknHolderAddrs["0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			_, ok = con.TknHolderAddrs["0x"]
			Expect(ok).To(Equal(false))

			_, ok = con.TknHolderAddrs[""]
			Expect(ok).To(Equal(false))

			_, ok = con.TknHolderAddrs["0x09THISE21a5IS5cFAKE1D82fAND43bCE06MADEUP"]
			Expect(ok).To(Equal(false))
		})

		It("Fails with an empty contract", func() {
			event := con.Events["Transfer"]
			c := converter.NewConverter(&contract.Contract{})
			_, err = c.Convert([]types.Log{mocks.MockTransferLog1}, event, 232)
			Expect(err).To(HaveOccurred())
		})
	})
})
