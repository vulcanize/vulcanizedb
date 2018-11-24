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

package parser_test

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/parser"
)

var _ = Describe("Parser", func() {

	var p parser.Parser
	var err error

	BeforeEach(func() {
		p = parser.NewParser("")
	})

	Describe("Mock Parse", func() {
		It("Uses parses given abi string", func() {
			mp := mocks.NewParser(constants.DaiAbiString)
			err = mp.Parse()
			Expect(err).ToNot(HaveOccurred())

			parsedAbi := mp.ParsedAbi()
			expectedAbi, err := geth.ParseAbi(constants.DaiAbiString)
			Expect(err).ToNot(HaveOccurred())
			Expect(parsedAbi).To(Equal(expectedAbi))

			methods := mp.GetMethods([]string{"balanceOf"})
			_, ok := methods["totalSupply"]
			Expect(ok).To(Equal(false))
			m, ok := methods["balanceOf"]
			Expect(ok).To(Equal(true))
			Expect(len(m.Args)).To(Equal(1))
			Expect(len(m.Return)).To(Equal(1))

			events := mp.GetEvents([]string{"Transfer"})
			_, ok = events["Mint"]
			Expect(ok).To(Equal(false))
			e, ok := events["Transfer"]
			Expect(ok).To(Equal(true))
			Expect(len(e.Fields)).To(Equal(3))
		})
	})

	Describe("Parse", func() {
		It("Fetches and parses abi from etherscan using contract address", func() {
			contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359" // dai contract address
			err = p.Parse(contractAddr)
			Expect(err).ToNot(HaveOccurred())

			expectedAbi := constants.DaiAbiString
			Expect(p.Abi()).To(Equal(expectedAbi))

			expectedParsedAbi, err := geth.ParseAbi(expectedAbi)
			Expect(err).ToNot(HaveOccurred())
			Expect(p.ParsedAbi()).To(Equal(expectedParsedAbi))
		})

		It("Fails with a normal, non-contract, account address", func() {
			addr := "0xAb2A8F7cB56D9EC65573BA1bE0f92Fa2Ff7dd165"
			err = p.Parse(addr)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetEvents", func() {
		It("Returns parsed events", func() {
			contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
			err = p.Parse(contractAddr)
			Expect(err).ToNot(HaveOccurred())

			events := p.GetEvents([]string{"Transfer"})

			e, ok := events["Transfer"]
			Expect(ok).To(Equal(true))

			abiTy := e.Fields[0].Type.T
			Expect(abiTy).To(Equal(abi.AddressTy))

			pgTy := e.Fields[0].PgType
			Expect(pgTy).To(Equal("CHARACTER VARYING(66)"))

			abiTy = e.Fields[1].Type.T
			Expect(abiTy).To(Equal(abi.AddressTy))

			pgTy = e.Fields[1].PgType
			Expect(pgTy).To(Equal("CHARACTER VARYING(66)"))

			abiTy = e.Fields[2].Type.T
			Expect(abiTy).To(Equal(abi.UintTy))

			pgTy = e.Fields[2].PgType
			Expect(pgTy).To(Equal("DECIMAL"))

			_, ok = events["Approval"]
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetMethods", func() {
		It("Parses and returns only methods specified in passed array", func() {
			contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
			err = p.Parse(contractAddr)
			Expect(err).ToNot(HaveOccurred())

			selectMethods := p.GetMethods([]string{"balanceOf"})

			m, ok := selectMethods["balanceOf"]
			Expect(ok).To(Equal(true))

			abiTy := m.Args[0].Type.T
			Expect(abiTy).To(Equal(abi.AddressTy))

			pgTy := m.Args[0].PgType
			Expect(pgTy).To(Equal("CHARACTER VARYING(66)"))

			abiTy = m.Return[0].Type.T
			Expect(abiTy).To(Equal(abi.UintTy))

			pgTy = m.Return[0].PgType
			Expect(pgTy).To(Equal("DECIMAL"))

			_, ok = selectMethods["totalSupply"]
			Expect(ok).To(Equal(false))
		})

		It("Parses and returns all methods if passed an empty array", func() {
			contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
			err = p.Parse(contractAddr)
			Expect(err).ToNot(HaveOccurred())

			selectMethods := p.GetMethods([]string{})

			_, ok := selectMethods["balanceOf"]
			Expect(ok).To(Equal(true))

			_, ok = selectMethods["totalSupply"]
			Expect(ok).To(Equal(true))

			_, ok = selectMethods["allowance"]
			Expect(ok).To(Equal(true))
		})

		It("Parses and returns no methods if pass a nil array", func() {
			contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
			err = p.Parse(contractAddr)
			Expect(err).ToNot(HaveOccurred())

			selectMethods := p.GetMethods(nil)
			Expect(len(selectMethods)).To(Equal(0))
		})
	})

	Describe("GetAddrMethods", func() {
		It("Parses and returns only methods whose inputs, if any, are all addresses", func() {
			contractAddr := "0xDdE2D979e8d39BB8416eAfcFC1758f3CaB2C9C72"
			err = p.Parse(contractAddr)
			Expect(err).ToNot(HaveOccurred())
			wanted := []string{"isApprovedForAll", "supportsInterface", "getApproved", "totalSupply", "balanceOf"}

			methods := p.GetMethods(wanted)
			selectMethods := p.GetAddrMethods(wanted)

			_, ok := selectMethods["totalSupply"]
			Expect(ok).To(Equal(true))
			_, ok = methods["totalSupply"]
			Expect(ok).To(Equal(true))

			_, ok = selectMethods["balanceOf"]
			Expect(ok).To(Equal(true))
			_, ok = methods["balanceOf"]
			Expect(ok).To(Equal(true))

			_, ok = selectMethods["isApprovedForAll"]
			Expect(ok).To(Equal(true))
			_, ok = methods["isApprovedForAll"]
			Expect(ok).To(Equal(true))

			_, ok = selectMethods["supportsInterface"]
			Expect(ok).To(Equal(false))
			_, ok = methods["supportsInterface"]
			Expect(ok).To(Equal(true))

			_, ok = selectMethods["getApproved"]
			Expect(ok).To(Equal(false))
			_, ok = methods["getApproved"]
			Expect(ok).To(Equal(true))

			_, ok = selectMethods["name"]
			Expect(ok).To(Equal(false))
			_, ok = methods["name"]
			Expect(ok).To(Equal(false))
		})
	})
})
