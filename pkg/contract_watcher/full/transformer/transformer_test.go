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
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/full/retriever"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/full/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/helpers/test_helpers/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/parser"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/poller"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Transformer", func() {
	var fakeAddress = "0x1234567890abcdef"
	rand.Seed(time.Now().UnixNano())

	Describe("Init", func() {
		It("Initializes transformer's contract objects", func() {
			blockRetriever := &fakes.MockFullSyncBlockRetriever{}
			firstBlock := int64(1)
			mostRecentBlock := int64(2)
			blockRetriever.FirstBlock = firstBlock
			blockRetriever.MostRecentBlock = mostRecentBlock

			parsr := &fakes.MockParser{}
			fakeAbi := "fake_abi"
			eventName := "Transfer"
			event := types.Event{}
			parsr.AbiToReturn = fakeAbi
			parsr.EventName = eventName
			parsr.Event = event

			pollr := &fakes.MockPoller{}
			fakeContractName := "fake_contract_name"
			pollr.ContractName = fakeContractName

			t := getTransformer(blockRetriever, parsr, pollr)

			err := t.Init()

			Expect(err).ToNot(HaveOccurred())

			c, ok := t.Contracts[fakeAddress]
			Expect(ok).To(Equal(true))

			Expect(c.StartingBlock).To(Equal(firstBlock))
			Expect(t.LastBlock).To(Equal(mostRecentBlock))
			Expect(c.Abi).To(Equal(fakeAbi))
			Expect(c.Name).To(Equal(fakeContractName))
			Expect(c.Address).To(Equal(fakeAddress))
		})

		It("Fails to initialize if first and most recent blocks cannot be fetched from vDB", func() {
			blockRetriever := &fakes.MockFullSyncBlockRetriever{}
			blockRetriever.FirstBlockErr = fakes.FakeError
			t := getTransformer(blockRetriever, &fakes.MockParser{}, &fakes.MockPoller{})

			err := t.Init()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})

func getTransformer(blockRetriever retriever.BlockRetriever, parsr parser.Parser, pollr poller.Poller) transformer.Transformer {
	return transformer.Transformer{
		FilterRepository: &fakes.MockFilterRepository{},
		Parser:           parsr,
		Retriever:        blockRetriever,
		Poller:           pollr,
		Contracts:        map[string]*contract.Contract{},
		Config:           mocks.MockConfig,
	}
}
