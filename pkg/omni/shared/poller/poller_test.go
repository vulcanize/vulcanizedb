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
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/poller"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

var _ = Describe("Poller", func() {

	var p poller.Poller
	var con *contract.Contract
	var db *postgres.DB
	var bc core.BlockChain

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	Describe("Full sync mode", func() {
		BeforeEach(func() {
			db, bc = test_helpers.SetupDBandBC()
			p = poller.NewPoller(bc, db, types.FullSync)
		})

		Describe("PollContract", func() {
			It("Polls specified contract methods using contract's token holder address list", func() {
				con = test_helpers.SetupTusdContract(nil, []string{"balanceOf"})
				Expect(con.Abi).To(Equal(constants.TusdAbiString))
				con.StartingBlock = 6707322
				con.LastBlock = 6707323
				con.AddEmittedAddr(common.HexToAddress("0xfE9e8709d3215310075d67E3ed32A380CCf451C8"), common.HexToAddress("0x3f5CE5FBFe3E9af3971dD833D26bA9b5C936f0bE"))

				err := p.PollContract(*con)
				Expect(err).ToNot(HaveOccurred())

				scanStruct := test_helpers.BalanceOf{}

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.balanceof_method WHERE who_ = '0xfE9e8709d3215310075d67E3ed32A380CCf451C8' AND block = '6707322'", constants.TusdContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Balance).To(Equal("66386309548896882859581786"))
				Expect(scanStruct.TokenName).To(Equal("TrueUSD"))

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.balanceof_method WHERE who_ = '0xfE9e8709d3215310075d67E3ed32A380CCf451C8' AND block = '6707323'", constants.TusdContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Balance).To(Equal("66386309548896882859581786"))
				Expect(scanStruct.TokenName).To(Equal("TrueUSD"))

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.balanceof_method WHERE who_ = '0x3f5CE5FBFe3E9af3971dD833D26bA9b5C936f0bE' AND block = '6707322'", constants.TusdContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Balance).To(Equal("17982350181394112023885864"))
				Expect(scanStruct.TokenName).To(Equal("TrueUSD"))

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.balanceof_method WHERE who_ = '0x3f5CE5FBFe3E9af3971dD833D26bA9b5C936f0bE' AND block = '6707323'", constants.TusdContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Balance).To(Equal("17982350181394112023885864"))
				Expect(scanStruct.TokenName).To(Equal("TrueUSD"))
			})

			It("Does not poll and persist any methods if none are specified", func() {
				con = test_helpers.SetupTusdContract(nil, nil)
				Expect(con.Abi).To(Equal(constants.TusdAbiString))
				con.StartingBlock = 6707322
				con.LastBlock = 6707323
				con.AddEmittedAddr(common.HexToAddress("0xfE9e8709d3215310075d67E3ed32A380CCf451C8"), common.HexToAddress("0x3f5CE5FBFe3E9af3971dD833D26bA9b5C936f0bE"))

				err := p.PollContract(*con)
				Expect(err).ToNot(HaveOccurred())

				scanStruct := test_helpers.BalanceOf{}

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.balanceof_method WHERE who_ = '0xfE9e8709d3215310075d67E3ed32A380CCf451C8' AND block = '6707322'", constants.TusdContractAddress)).StructScan(&scanStruct)
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("PollMethod", func() {
			It("Polls a single contract method", func() {
				var name = new(string)
				err := p.FetchContractData(constants.TusdAbiString, constants.TusdContractAddress, "name", nil, &name, 6197514)
				Expect(err).ToNot(HaveOccurred())
				Expect(*name).To(Equal("TrueUSD"))
			})
		})
	})
})
