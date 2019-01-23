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
			It("Polls specified contract methods using contract's argument list", func() {
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

			It("Polls specified contract methods using contract's hash list", func() {
				con = test_helpers.SetupENSContract(nil, []string{"owner"})
				Expect(con.Abi).To(Equal(constants.ENSAbiString))
				Expect(len(con.Methods)).To(Equal(1))
				con.AddEmittedHash(common.HexToHash("0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae"), common.HexToHash("0x7e74a86b6e146964fb965db04dc2590516da77f720bb6759337bf5632415fd86"))

				err := p.PollContractAt(*con, 6885877)
				Expect(err).ToNot(HaveOccurred())

				scanStruct := test_helpers.Owner{}

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.owner_method WHERE node_ = '0x7e74a86b6e146964fb965db04dc2590516da77f720bb6759337bf5632415fd86' AND block = '6885877'", constants.EnsContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Address).To(Equal("0x546aA2EaE2514494EeaDb7bbb35243348983C59d"))
				Expect(scanStruct.TokenName).To(Equal("ENS-Registry"))

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.owner_method WHERE node_ = '0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae' AND block = '6885877'", constants.EnsContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Address).To(Equal("0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef"))
				Expect(scanStruct.TokenName).To(Equal("ENS-Registry"))
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

		Describe("FetchContractData", func() {
			It("Calls a single contract method", func() {
				var name = new(string)
				err := p.FetchContractData(constants.TusdAbiString, constants.TusdContractAddress, "name", nil, &name, 6197514)
				Expect(err).ToNot(HaveOccurred())
				Expect(*name).To(Equal("TrueUSD"))
			})
		})
	})

	Describe("Light sync mode", func() {
		BeforeEach(func() {
			db, bc = test_helpers.SetupDBandBC()
			p = poller.NewPoller(bc, db, types.LightSync)
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

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.balanceof_method WHERE who_ = '0xfE9e8709d3215310075d67E3ed32A380CCf451C8' AND block = '6707322'", constants.TusdContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Balance).To(Equal("66386309548896882859581786"))
				Expect(scanStruct.TokenName).To(Equal("TrueUSD"))

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.balanceof_method WHERE who_ = '0xfE9e8709d3215310075d67E3ed32A380CCf451C8' AND block = '6707323'", constants.TusdContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Balance).To(Equal("66386309548896882859581786"))
				Expect(scanStruct.TokenName).To(Equal("TrueUSD"))

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.balanceof_method WHERE who_ = '0x3f5CE5FBFe3E9af3971dD833D26bA9b5C936f0bE' AND block = '6707322'", constants.TusdContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Balance).To(Equal("17982350181394112023885864"))
				Expect(scanStruct.TokenName).To(Equal("TrueUSD"))

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.balanceof_method WHERE who_ = '0x3f5CE5FBFe3E9af3971dD833D26bA9b5C936f0bE' AND block = '6707323'", constants.TusdContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Balance).To(Equal("17982350181394112023885864"))
				Expect(scanStruct.TokenName).To(Equal("TrueUSD"))
			})

			It("Polls specified contract methods using contract's hash list", func() {
				con = test_helpers.SetupENSContract(nil, []string{"owner"})
				Expect(con.Abi).To(Equal(constants.ENSAbiString))
				Expect(len(con.Methods)).To(Equal(1))
				con.AddEmittedHash(common.HexToHash("0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae"), common.HexToHash("0x7e74a86b6e146964fb965db04dc2590516da77f720bb6759337bf5632415fd86"))

				err := p.PollContractAt(*con, 6885877)
				Expect(err).ToNot(HaveOccurred())

				scanStruct := test_helpers.Owner{}

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.owner_method WHERE node_ = '0x7e74a86b6e146964fb965db04dc2590516da77f720bb6759337bf5632415fd86' AND block = '6885877'", constants.EnsContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Address).To(Equal("0x546aA2EaE2514494EeaDb7bbb35243348983C59d"))
				Expect(scanStruct.TokenName).To(Equal("ENS-Registry"))

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.owner_method WHERE node_ = '0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae' AND block = '6885877'", constants.EnsContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Address).To(Equal("0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef"))
				Expect(scanStruct.TokenName).To(Equal("ENS-Registry"))
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

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.balanceof_method WHERE who_ = '0xfE9e8709d3215310075d67E3ed32A380CCf451C8' AND block = '6707322'", constants.TusdContractAddress)).StructScan(&scanStruct)
				Expect(err).To(HaveOccurred())
			})

			It("Caches returned values of the appropriate types for downstream method polling if method piping is turned on", func() {
				con = test_helpers.SetupENSContract(nil, []string{"resolver"})
				Expect(con.Abi).To(Equal(constants.ENSAbiString))
				con.StartingBlock = 6921967
				con.LastBlock = 6921968
				con.EmittedAddrs = map[interface{}]bool{}
				con.Piping = false
				con.AddEmittedHash(common.HexToHash("0x495b6e6efdedb750aa519919b5cf282bdaa86067b82a2293a3ff5723527141e8"))
				err := p.PollContract(*con)
				Expect(err).ToNot(HaveOccurred())

				scanStruct := test_helpers.Resolver{}
				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.resolver_method WHERE node_ = '0x495b6e6efdedb750aa519919b5cf282bdaa86067b82a2293a3ff5723527141e8' AND block = '6921967'", constants.EnsContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Address).To(Equal("0x5FfC014343cd971B7eb70732021E26C35B744cc4"))
				Expect(scanStruct.TokenName).To(Equal("ENS-Registry"))
				Expect(len(con.EmittedAddrs)).To(Equal(0)) // With piping off the address is not saved

				test_helpers.TearDown(db)
				db, bc = test_helpers.SetupDBandBC()
				p = poller.NewPoller(bc, db, types.LightSync)

				con.Piping = true
				err = p.PollContract(*con)
				Expect(err).ToNot(HaveOccurred())

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.resolver_method WHERE node_ = '0x495b6e6efdedb750aa519919b5cf282bdaa86067b82a2293a3ff5723527141e8' AND block = '6921967'", constants.EnsContractAddress)).StructScan(&scanStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(scanStruct.Address).To(Equal("0x5FfC014343cd971B7eb70732021E26C35B744cc4"))
				Expect(scanStruct.TokenName).To(Equal("ENS-Registry"))
				Expect(len(con.EmittedAddrs)).To(Equal(1)) // With piping on it is saved
				Expect(con.EmittedAddrs[common.HexToAddress("0x5FfC014343cd971B7eb70732021E26C35B744cc4")]).To(Equal(true))
			})
		})

		Describe("FetchContractData", func() {
			It("Calls a single contract method", func() {
				var name = new(string)
				err := p.FetchContractData(constants.TusdAbiString, constants.TusdContractAddress, "name", nil, &name, 6197514)
				Expect(err).ToNot(HaveOccurred())
				Expect(*name).To(Equal("TrueUSD"))
			})
		})
	})
})
