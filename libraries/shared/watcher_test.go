package shared_test

import (
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	shared2 "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Watcher", func() {
	It("initialises correctly", func() {
		db := test_config.NewTestDB(core.Node{ID: "testNode"})
		bc := fakes.NewMockBlockChain()

		watcher := shared.NewWatcher(db, bc)

		Expect(watcher.DB).To(Equal(db))
		Expect(watcher.Fetcher).NotTo(BeNil())
		Expect(watcher.Chunker).NotTo(BeNil())
	})

	It("adds transformers", func() {
		watcher := shared.NewWatcher(nil, nil)
		fakeTransformer := &mocks.MockTransformer{}
		fakeTransformer.SetTransformerConfig(mocks.FakeTransformerConfig)
		watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})

		Expect(len(watcher.Transformers)).To(Equal(1))
		Expect(watcher.Transformers).To(ConsistOf(fakeTransformer))
		Expect(watcher.Topics).To(Equal([]common.Hash{common.HexToHash("FakeTopic")}))
		Expect(watcher.Addresses).To(Equal([]common.Address{common.HexToAddress("FakeAddress")}))
	})

	It("adds transformers from multiple sources", func() {
		watcher := shared.NewWatcher(nil, nil)
		fakeTransformer1 := &mocks.MockTransformer{}
		fakeTransformer1.SetTransformerConfig(mocks.FakeTransformerConfig)

		fakeTransformer2 := &mocks.MockTransformer{}
		fakeTransformer2.SetTransformerConfig(mocks.FakeTransformerConfig)

		watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformer1.FakeTransformerInitializer})
		watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformer2.FakeTransformerInitializer})

		Expect(len(watcher.Transformers)).To(Equal(2))
		Expect(watcher.Topics).To(Equal([]common.Hash{common.HexToHash("FakeTopic"),
			common.HexToHash("FakeTopic")}))
		Expect(watcher.Addresses).To(Equal([]common.Address{common.HexToAddress("FakeAddress"),
			common.HexToAddress("FakeAddress")}))
	})

	It("calculates earliest starting block number", func() {
		fakeTransformer1 := &mocks.MockTransformer{}
		fakeTransformer1.SetTransformerConfig(shared2.TransformerConfig{StartingBlockNumber: 5})

		fakeTransformer2 := &mocks.MockTransformer{}
		fakeTransformer2.SetTransformerConfig(shared2.TransformerConfig{StartingBlockNumber: 3})

		watcher := shared.NewWatcher(nil, nil)
		watcher.AddTransformers([]shared2.TransformerInitializer{
			fakeTransformer1.FakeTransformerInitializer,
			fakeTransformer2.FakeTransformerInitializer,
		})

		Expect(*watcher.StartingBlock).To(Equal(int64(3)))
	})

	It("returns an error when run without transformers", func() {
		watcher := shared.NewWatcher(nil, nil)
		err := watcher.Execute()
		Expect(err).To(MatchError("No transformers added to watcher"))
	})

	Describe("with missing headers", func() {
		var (
			db               *postgres.DB
			watcher          shared.Watcher
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
			watcher = shared.NewWatcher(db, &mockBlockChain)
		})

		It("executes each transformer", func() {
			fakeTransformer := &mocks.MockTransformer{}
			watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})
			repository.SetMissingHeaders([]core.Header{fakes.FakeHeader})

			err := watcher.Execute()

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTransformer.ExecuteWasCalled).To(BeTrue())
		})

		It("returns an error if transformer returns an error", func() {
			fakeTransformer := &mocks.MockTransformer{ExecuteError: errors.New("Something bad happened")}
			watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})
			repository.SetMissingHeaders([]core.Header{fakes.FakeHeader})

			err := watcher.Execute()

			Expect(err).To(HaveOccurred())
			Expect(fakeTransformer.ExecuteWasCalled).To(BeFalse())
		})

		It("passes only relevant logs to each transformer", func() {
			transformerA := &mocks.MockTransformer{}
			transformerB := &mocks.MockTransformer{}

			configA := shared2.TransformerConfig{TransformerName: "transformerA",
				ContractAddresses: []string{"0x000000000000000000000000000000000000000A"},
				Topic:             "0xA"}
			configB := shared2.TransformerConfig{TransformerName: "transformerB",
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
			watcher = shared.NewWatcher(db, &mockBlockChain)
			watcher.AddTransformers([]shared2.TransformerInitializer{
				transformerA.FakeTransformerInitializer, transformerB.FakeTransformerInitializer})

			err := watcher.Execute()
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
				fakeTransformer.SetTransformerConfig(shared2.TransformerConfig{
					Topic: topic, ContractAddresses: addresses})
				watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})

				err := watcher.Execute()
				Expect(err).NotTo(HaveOccurred())

				fakeHash := common.HexToHash(fakes.FakeHeader.Hash)
				mockBlockChain.AssertGetEthLogsWithCustomQueryCalledWith(ethereum.FilterQuery{
					BlockHash: &fakeHash,
					Addresses: shared2.HexStringsToAddresses(addresses),
					Topics:    [][]common.Hash{{common.HexToHash(topic)}},
				})
			})

			It("propagates log fetcher errors", func() {
				fetcherError := errors.New("FetcherError")
				mockBlockChain.SetGetEthLogsWithCustomQueryErr(fetcherError)

				watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})
				err := watcher.Execute()
				Expect(err).To(MatchError(fetcherError))
			})
		})
	})
})
