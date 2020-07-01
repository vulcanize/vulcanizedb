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

package converter_test

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/contract"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/converter"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/helpers"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/helpers/test_helpers"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/helpers/test_helpers/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Converter", func() {
	var (
		fakeHeaderID = rand.Int63()

		marketPlaceWantedEvents = []string{"OrderCreated"}
		molochWantedEvents      = []string{"SubmitVote"}
		oasisWantedEvents       = []string{"LogMake"}
		tusdWantedEvents        = []string{"Transfer", "Mint"}
	)

	Describe("Update", func() {
		It("Updates contract info held by the converter", func() {
			con := test_helpers.SetupTusdContract(tusdWantedEvents)
			c := converter.NewConverter()
			c.Update(con)
			Expect(c.ContractInfo).To(Equal(con))

			info := test_helpers.SetupTusdContract([]string{})
			c.Update(info)
			Expect(c.ContractInfo).To(Equal(info))
		})
	})

	Describe("Convert", func() {
		It("Converts a watched event log to mapping of event input names to values", func() {
			con := test_helpers.SetupTusdContract(tusdWantedEvents)
			_, ok := con.Events["Approval"]
			Expect(ok).To(Equal(false))

			event, ok := con.Events["Transfer"]
			Expect(ok).To(Equal(true))

			c := converter.NewConverter()
			c.Update(con)
			logs, err := c.Convert([]types.Log{mocks.MockTransferLog1, mocks.MockTransferLog2}, event, fakeHeaderID)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(2))

			sender1 := common.HexToAddress("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391")
			sender2 := common.HexToAddress("0x000000000000000000000000000000000000000000000000000000000000af21")
			value := helpers.BigFromString("1097077688018008265106216665536940668749033598146")

			Expect(logs[0].Values["to"]).To(Equal(sender1.String()))
			Expect(logs[0].Values["from"]).To(Equal(sender2.String()))
			Expect(logs[0].Values["value"]).To(Equal(value.String()))
			Expect(logs[0].HeaderID).To(Equal(fakeHeaderID))
			Expect(logs[1].Values["to"]).To(Equal(sender2.String()))
			Expect(logs[1].Values["from"]).To(Equal(sender1.String()))
			Expect(logs[1].Values["value"]).To(Equal(value.String()))
			Expect(logs[1].HeaderID).To(Equal(fakeHeaderID))
		})

		It("correctly parses bytes32", func() {
			con := test_helpers.SetupMarketPlaceContract(marketPlaceWantedEvents)
			event, ok := con.Events["OrderCreated"]
			Expect(ok).To(BeTrue())

			c := converter.NewConverter()
			c.Update(con)
			result, err := c.Convert([]types.Log{mocks.MockOrderCreatedLog}, event, fakeHeaderID)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(result)).To(Equal(1))
			Expect(result[0].Values["id"]).To(Equal("0x633f94affdcabe07c000231f85c752c97b9cc43966b432ec4d18641e6d178233"))
		})

		It("correctly parses uint8", func() {
			con := test_helpers.SetupMolochContract(molochWantedEvents)
			event, ok := con.Events["SubmitVote"]
			Expect(ok).To(BeTrue())

			c := converter.NewConverter()
			c.Update(con)
			result, err := c.Convert([]types.Log{mocks.MockSubmitVoteLog}, event, fakeHeaderID)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(result)).To(Equal(1))
			Expect(result[0].Values["uintVote"]).To(Equal("1"))
		})

		It("correctly parses uint64", func() {
			con := test_helpers.SetupOasisContract(oasisWantedEvents)
			event, ok := con.Events["LogMake"]
			Expect(ok).To(BeTrue())

			c := converter.NewConverter()
			c.Update(con)
			result, err := c.Convert([]types.Log{mocks.MockLogMakeLog}, event, fakeHeaderID)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(result)).To(Equal(1))
			Expect(result[0].Values["timestamp"]).To(Equal("1580153827"))
		})

		It("Fails with an empty contract", func() {
			con := contract.Contract{}.Init()
			event := con.Events["Transfer"]
			c := converter.NewConverter()
			c.Update(&contract.Contract{})

			_, err := c.Convert([]types.Log{mocks.MockTransferLog1}, event, fakeHeaderID)

			Expect(err).To(HaveOccurred())
		})
	})
})
