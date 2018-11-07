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

package poller_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/omni/helpers"
	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/omni/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/poller"
)

var _ = Describe("Poller", func() {

	var p poller.Poller
	var con *contract.Contract

	BeforeEach(func() {
		p = poller.NewPoller(test_helpers.SetupBC())
	})

	Describe("PollContract", func() {
		It("Polls contract methods using token holder address list", func() {
			con = test_helpers.SetupTusdContract(nil, []string{"balanceOf"})
			Expect(con.Abi).To(Equal(constants.TusdAbiString))
			con.StartingBlock = 6707322
			con.LastBlock = 6707323
			con.TknHolderAddrs = map[string]bool{
				"0xfE9e8709d3215310075d67E3ed32A380CCf451C8": true,
				"0x3f5CE5FBFe3E9af3971dD833D26bA9b5C936f0bE": true,
			}

			err := p.PollContract(con)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(con.Methods["balanceOf"].Results)).To(Equal(2))
			res1 := con.Methods["balanceOf"].Results[0]
			res2 := con.Methods["balanceOf"].Results[1]
			expectedRes11 := helpers.BigFromString("66386309548896882859581786")
			expectedRes12 := helpers.BigFromString("66386309548896882859581786")
			expectedRes21 := helpers.BigFromString("17982350181394112023885864")
			expectedRes22 := helpers.BigFromString("17982350181394112023885864")
			if res1.Inputs[0].(common.Address).String() == "0xfE9e8709d3215310075d67E3ed32A380CCf451C8" {
				Expect(res2.Inputs[0].(common.Address).String()).To(Equal("0x3f5CE5FBFe3E9af3971dD833D26bA9b5C936f0bE"))
				Expect(res1.Outputs[6707322].(*big.Int).String()).To(Equal(expectedRes11.String()))
				Expect(res1.Outputs[6707323].(*big.Int).String()).To(Equal(expectedRes12.String()))
				Expect(res2.Outputs[6707322].(*big.Int).String()).To(Equal(expectedRes21.String()))
				Expect(res2.Outputs[6707323].(*big.Int).String()).To(Equal(expectedRes22.String()))
			} else {
				Expect(res2.Inputs[0].(common.Address).String()).To(Equal("0xfE9e8709d3215310075d67E3ed32A380CCf451C8"))
				Expect(res2.Outputs[6707322].(*big.Int).String()).To(Equal(expectedRes11.String()))
				Expect(res2.Outputs[6707323].(*big.Int).String()).To(Equal(expectedRes12.String()))
				Expect(res1.Outputs[6707322].(*big.Int).String()).To(Equal(expectedRes21.String()))
				Expect(res1.Outputs[6707323].(*big.Int).String()).To(Equal(expectedRes22.String()))
			}
		})
	})

	Describe("PollMethod", func() {
		It("Polls a single contract method", func() {
			var name = new(string)
			err := p.PollMethod(constants.TusdAbiString, constants.TusdContractAddress, "name", nil, &name, 6197514)
			Expect(err).ToNot(HaveOccurred())
			Expect(*name).To(Equal("TrueUSD"))
		})
	})
})
