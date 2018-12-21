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
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
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

			methods := mp.GetSelectMethods([]string{"balanceOf"})
			Expect(len(methods)).To(Equal(1))
			balOf := methods[0]
			Expect(balOf.Name).To(Equal("balanceOf"))
			Expect(len(balOf.Args)).To(Equal(1))
			Expect(len(balOf.Return)).To(Equal(1))

			events := mp.GetEvents([]string{"Transfer"})
			_, ok := events["Mint"]
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

	Describe("GetSelectMethods", func() {
		It("Parses and returns only methods specified in passed array", func() {
			contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
			err = p.Parse(contractAddr)
			Expect(err).ToNot(HaveOccurred())

			methods := p.GetSelectMethods([]string{"balanceOf"})
			Expect(len(methods)).To(Equal(1))

			balOf := methods[0]
			Expect(balOf.Name).To(Equal("balanceOf"))
			Expect(len(balOf.Args)).To(Equal(1))
			Expect(len(balOf.Return)).To(Equal(1))

			abiTy := balOf.Args[0].Type.T
			Expect(abiTy).To(Equal(abi.AddressTy))

			pgTy := balOf.Args[0].PgType
			Expect(pgTy).To(Equal("CHARACTER VARYING(66)"))

			abiTy = balOf.Return[0].Type.T
			Expect(abiTy).To(Equal(abi.UintTy))

			pgTy = balOf.Return[0].PgType
			Expect(pgTy).To(Equal("DECIMAL"))

		})

		It("Parses and returns methods in the order they were specified", func() {
			contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
			err = p.Parse(contractAddr)
			Expect(err).ToNot(HaveOccurred())

			selectMethods := p.GetSelectMethods([]string{"balanceOf", "allowance"})
			Expect(len(selectMethods)).To(Equal(2))

			balOf := selectMethods[0]
			allow := selectMethods[1]

			Expect(balOf.Name).To(Equal("balanceOf"))
			Expect(allow.Name).To(Equal("allowance"))
		})

		It("Returns nil if given a nil or empty array", func() {
			contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
			err = p.Parse(contractAddr)
			Expect(err).ToNot(HaveOccurred())

			var nilArr []types.Method
			selectMethods := p.GetSelectMethods([]string{})
			Expect(selectMethods).To(Equal(nilArr))
			selectMethods = p.GetMethods(nil)
			Expect(selectMethods).To(Equal(nilArr))
		})

	})

	Describe("GetMethods", func() {
		It("Parses and returns only methods specified in passed array", func() {
			contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
			err = p.Parse(contractAddr)
			Expect(err).ToNot(HaveOccurred())

			methods := p.GetMethods([]string{"balanceOf"})
			Expect(len(methods)).To(Equal(1))

			balOf := methods[0]
			Expect(balOf.Name).To(Equal("balanceOf"))
			Expect(len(balOf.Args)).To(Equal(1))
			Expect(len(balOf.Return)).To(Equal(1))

			abiTy := balOf.Args[0].Type.T
			Expect(abiTy).To(Equal(abi.AddressTy))

			pgTy := balOf.Args[0].PgType
			Expect(pgTy).To(Equal("CHARACTER VARYING(66)"))

			abiTy = balOf.Return[0].Type.T
			Expect(abiTy).To(Equal(abi.UintTy))

			pgTy = balOf.Return[0].PgType
			Expect(pgTy).To(Equal("DECIMAL"))

		})

		It("Returns nil if given a nil array", func() {
			contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
			err = p.Parse(contractAddr)
			Expect(err).ToNot(HaveOccurred())

			var nilArr []types.Method
			selectMethods := p.GetMethods(nil)
			Expect(selectMethods).To(Equal(nilArr))
		})

		It("Returns every method if given an empty array", func() {
			contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
			err = p.Parse(contractAddr)
			Expect(err).ToNot(HaveOccurred())

			selectMethods := p.GetMethods([]string{})
			Expect(len(selectMethods)).To(Equal(22))
		})
	})
})
