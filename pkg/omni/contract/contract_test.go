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

package contract_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/parser"
)

var expectedLogFilter = filters.LogFilter{
	Name:      "Transfer",
	Address:   constants.TusdContractAddress,
	ToBlock:   -1,
	FromBlock: 5197514,
	Topics:    core.Topics{constants.TransferEvent.Signature()},
}

var _ = Describe("Contract test", func() {
	var p parser.Parser
	var err error

	BeforeEach(func() {
		p = parser.NewParser("")
		err = p.Parse(constants.TusdContractAddress)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Creates filters from stored data", func() {
		info := contract.Contract{
			Name:          "TrueUSD",
			Address:       constants.TusdContractAddress,
			Abi:           p.Abi(),
			ParsedAbi:     p.ParsedAbi(),
			StartingBlock: 5197514,
			Events:        p.GetEvents(),
			Methods:       p.GetMethods(),
			Addresses:     map[string]bool{},
		}

		err = info.GenerateFilters([]string{"Transfer"})
		Expect(err).ToNot(HaveOccurred())
		val, ok := info.Filters["Transfer"]
		Expect(ok).To(Equal(true))
		Expect(val).To(Equal(expectedLogFilter))
	})

	It("Fails with an empty contract", func() {
		info := contract.Contract{}
		err = info.GenerateFilters([]string{"Transfer"})
		Expect(err).To(HaveOccurred())
	})
})
