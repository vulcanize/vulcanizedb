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

package transformer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/retriever"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/helpers/test_helpers/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/parser"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/poller"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Transformer", func() {
	var fakeAddress = "0x1234567890abcdef"
	Describe("Init", func() {
		It("Initializes transformer's contract objects", func() {
			blockRetriever := &fakes.MockLightBlockRetriever{}
			firstBlock := int64(1)
			blockRetriever.FirstBlock = firstBlock

			parsr := &fakes.MockParser{}
			fakeAbi := "fake_abi"
			parsr.AbiToReturn = fakeAbi

			pollr := &fakes.MockPoller{}
			fakeContractName := "fake_contract_name"
			pollr.ContractName = fakeContractName

			t := getFakeTransformer(blockRetriever, parsr, pollr)

			err := t.Init()

			Expect(err).ToNot(HaveOccurred())

			c, ok := t.Contracts[fakeAddress]
			Expect(ok).To(Equal(true))

			Expect(c.StartingBlock).To(Equal(firstBlock))
			Expect(c.LastBlock).To(Equal(int64(-1)))
			Expect(c.Abi).To(Equal(fakeAbi))
			Expect(c.Name).To(Equal(fakeContractName))
			Expect(c.Address).To(Equal(fakeAddress))
		})

		It("Fails to initialize if first block cannot be fetched from vDB headers table", func() {
			blockRetriever := &fakes.MockLightBlockRetriever{}
			blockRetriever.FirstBlockErr = fakes.FakeError
			t := getFakeTransformer(blockRetriever, &fakes.MockParser{}, &fakes.MockPoller{})

			err := t.Init()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})

func getFakeTransformer(blockRetriever retriever.BlockRetriever, parsr parser.Parser, pollr poller.Poller) transformer.Transformer {
	return transformer.Transformer{
		Parser:           parsr,
		BlockRetriever:   blockRetriever,
		Poller:           pollr,
		HeaderRepository: &fakes.MockLightHeaderRepository{},
		Contracts:        map[string]*contract.Contract{},
		Config:           mocks.MockConfig,
	}
}
