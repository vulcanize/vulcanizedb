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

package watcher_test

import (
	"errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Watcher", func() {
	It("initialises correctly", func() {
		db := test_config.NewTestDB(core.Node{ID: "testNode"})
		bc := fakes.NewMockBlockChain()

		w := watcher.NewEventWatcher(db, bc)

		Expect(w.DB).To(Equal(db))
		Expect(w.Fetcher).NotTo(BeNil())
		Expect(w.Chunker).NotTo(BeNil())
	})

	It("adds transformers", func() {
		w := watcher.NewEventWatcher(nil, nil)
		fakeTransformer := &mocks.MockTransformer{}
		fakeTransformer.SetTransformerConfig(mocks.FakeTransformerConfig)
		w.AddTransformers([]transformer.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})

		Expect(len(w.Transformers)).To(Equal(1))
		Expect(w.Transformers).To(ConsistOf(fakeTransformer))
		Expect(w.Topics).To(Equal([]common.Hash{common.HexToHash("FakeTopic")}))
		Expect(w.Addresses).To(Equal([]common.Address{common.HexToAddress("FakeAddress")}))
	})

	It("adds transformers from multiple sources", func() {
		w := watcher.NewEventWatcher(nil, nil)
		fakeTransformer1 := &mocks.MockTransformer{}
		fakeTransformer1.SetTransformerConfig(mocks.FakeTransformerConfig)

		fakeTransformer2 := &mocks.MockTransformer{}
		fakeTransformer2.SetTransformerConfig(mocks.FakeTransformerConfig)

		w.AddTransformers([]transformer.TransformerInitializer{fakeTransformer1.FakeTransformerInitializer})
		w.AddTransformers([]transformer.TransformerInitializer{fakeTransformer2.FakeTransformerInitializer})

		Expect(len(w.Transformers)).To(Equal(2))
		Expect(w.Topics).To(Equal([]common.Hash{common.HexToHash("FakeTopic"),
			common.HexToHash("FakeTopic")}))
		Expect(w.Addresses).To(Equal([]common.Address{common.HexToAddress("FakeAddress"),
			common.HexToAddress("FakeAddress")}))
	})

	It("calculates earliest starting block number", func() {
		fakeTransformer1 := &mocks.MockTransformer{}
		fakeTransformer1.SetTransformerConfig(transformer.TransformerConfig{StartingBlockNumber: 5})

		fakeTransformer2 := &mocks.MockTransformer{}
		fakeTransformer2.SetTransformerConfig(transformer.TransformerConfig{StartingBlockNumber: 3})

		w := watcher.NewEventWatcher(nil, nil)
		w.AddTransformers([]transformer.TransformerInitializer{
			fakeTransformer1.FakeTransformerInitializer,
			fakeTransformer2.FakeTransformerInitializer,
		})

		Expect(*w.StartingBlock).To(Equal(int64(3)))
	})

	It("returns an error when run without transformers", func() {
		w := watcher.NewEventWatcher(nil, nil)
		err := w.Execute(constants.HeaderMissing)
		Expect(err).To(MatchError("No transformers added to watcher"))
	})

	Describe("with missing headers", func() {
		var (
			db               *postgres.DB
			w                watcher.EventWatcher
			mockBlockChain   fakes.MockBlockChain
			headerRepository repositories.HeaderRepository
			repository       mocks.MockWatcherRepository
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			mockBlockChain = fakes.MockBlockChain{}
			headerRepository = repositories.NewHeaderRepository(db)
			_, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			repository = mocks.MockWatcherRepository{}
			w = watcher.NewEventWatcher(db, &mockBlockChain)
		})

		It("executes each transformer", func() {
			fakeTransformer := &mocks.MockTransformer{}
			w.AddTransformers([]transformer.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})
			repository.SetMissingHeaders([]core.Header{fakes.FakeHeader})

			err := w.Execute(constants.HeaderMissing)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTransformer.ExecuteWasCalled).To(BeTrue())
		})

		It("returns an error if transformer returns an error", func() {
			fakeTransformer := &mocks.MockTransformer{ExecuteError: errors.New("Something bad happened")}
			w.AddTransformers([]transformer.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})
			repository.SetMissingHeaders([]core.Header{fakes.FakeHeader})

			err := w.Execute(constants.HeaderMissing)
			Expect(err).To(HaveOccurred())
			Expect(fakeTransformer.ExecuteWasCalled).To(BeFalse())
		})

		It("passes only relevant logs to each transformer", func() {
			transformerA := &mocks.MockTransformer{}
			transformerB := &mocks.MockTransformer{}

			configA := transformer.TransformerConfig{TransformerName: "transformerA",
				ContractAddresses: []string{"0x000000000000000000000000000000000000000A"},
				Topic:             "0xA"}
			configB := transformer.TransformerConfig{TransformerName: "transformerB",
				ContractAddresses: []string{"0x000000000000000000000000000000000000000b"},
				Topic:             "0xB"}

			transformerA.SetTransformerConfig(configA)
			transformerB.SetTransformerConfig(configB)

			logA := types.Log{Address: common.HexToAddress("0xA"),
				Topics: []common.Hash{common.HexToHash("0xA")}}
			logB := types.Log{Address: common.HexToAddress("0xB"),
				Topics: []common.Hash{common.HexToHash("0xB")}}
			mockBlockChain.SetGetEthLogsWithCustomQueryReturnLogs([]types.Log{logA, logB})

			repository.SetMissingHeaders([]core.Header{fakes.FakeHeader})
			w = watcher.NewEventWatcher(db, &mockBlockChain)
			w.AddTransformers([]transformer.TransformerInitializer{
				transformerA.FakeTransformerInitializer, transformerB.FakeTransformerInitializer})

			err := w.Execute(constants.HeaderMissing)
			Expect(err).NotTo(HaveOccurred())
			Expect(transformerA.PassedLogs).To(Equal([]types.Log{logA}))
			Expect(transformerB.PassedLogs).To(Equal([]types.Log{logB}))
		})

		Describe("uses the LogFetcher correctly:", func() {
			var fakeTransformer mocks.MockTransformer
			BeforeEach(func() {
				repository.SetMissingHeaders([]core.Header{fakes.FakeHeader})
				fakeTransformer = mocks.MockTransformer{}
			})

			It("fetches logs for added transformers", func() {
				addresses := []string{"0xA", "0xB"}
				topic := "0x1"
				fakeTransformer.SetTransformerConfig(transformer.TransformerConfig{
					Topic: topic, ContractAddresses: addresses})
				w.AddTransformers([]transformer.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})

				err := w.Execute(constants.HeaderMissing)
				Expect(err).NotTo(HaveOccurred())

				fakeHash := common.HexToHash(fakes.FakeHeader.Hash)
				mockBlockChain.AssertGetEthLogsWithCustomQueryCalledWith(ethereum.FilterQuery{
					BlockHash: &fakeHash,
					Addresses: transformer.HexStringsToAddresses(addresses),
					Topics:    [][]common.Hash{{common.HexToHash(topic)}},
				})
			})

			It("propagates log fetcher errors", func() {
				fetcherError := errors.New("FetcherError")
				mockBlockChain.SetGetEthLogsWithCustomQueryErr(fetcherError)

				w.AddTransformers([]transformer.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})
				err := w.Execute(constants.HeaderMissing)
				Expect(err).To(MatchError(fetcherError))
			})
		})
	})
})
