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

package parser_test

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/constants"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/helpers/test_helpers/mocks"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/parser"
	"github.com/makerdao/vulcanizedb/pkg/eth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			expectedAbi, err := eth.ParseAbi(constants.DaiAbiString)
			Expect(err).ToNot(HaveOccurred())
			Expect(parsedAbi).To(Equal(expectedAbi))

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

			expectedParsedAbi, err := eth.ParseAbi(expectedAbi)
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
			Expect(pgTy).To(Equal("NUMERIC"))

			_, ok = events["Approval"]
			Expect(ok).To(Equal(false))
		})
	})
})
