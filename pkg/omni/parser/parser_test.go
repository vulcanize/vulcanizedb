// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/omni/parser"
)

var _ = Describe("Parser Test", func() {

	var p parser.Parser
	var err error

	BeforeEach(func() {
		p = parser.NewParser("")
	})

	It("Fetches and parses abi using contract address", func() {
		contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359" // dai contract address
		err = p.Parse(contractAddr)
		Expect(err).ToNot(HaveOccurred())

		expectedAbi := constants.DaiAbiString
		Expect(p.Abi()).To(Equal(expectedAbi))

		expectedParsedAbi, err := geth.ParseAbi(expectedAbi)
		Expect(err).ToNot(HaveOccurred())
		Expect(p.ParsedAbi()).To(Equal(expectedParsedAbi))
	})

	It("Returns parsed methods and events", func() {
		contractAddr := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
		err = p.Parse(contractAddr)
		Expect(err).ToNot(HaveOccurred())

		methods := p.GetMethods()
		events := p.GetEvents()

		_, ok := methods["totalSupply"]
		Expect(ok).To(Equal(true))

		_, ok = events["Transfer"]
		Expect(ok).To(Equal(true))
	})

})
