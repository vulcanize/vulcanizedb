package shared_test

import (
	"errors"
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

type MockTransformer struct {
	executeWasCalled bool
	executeError     error
	passedLogs       []types.Log
	passedHeader     core.Header
	transformerName  string
}

func (mh *MockTransformer) Execute(logs []types.Log, header core.Header) error {
	if mh.executeError != nil {
		return mh.executeError
	}
	mh.executeWasCalled = true
	mh.passedLogs = logs
	mh.passedHeader = header
	return nil
}

func (mh *MockTransformer) Name() string {
	return mh.transformerName
}

func (mh *MockTransformer) SetTransformerName(name string) {
	mh.transformerName = name
}

func fakeTransformerInitializer(db *postgres.DB) shared2.Transformer {
	return &MockTransformer{}
}

var fakeTransformerConfig = []shared2.TransformerConfig{{
	TransformerName:   "FakeTransformer",
	ContractAddresses: []string{"FakeAddress"},
	Topic:             "FakeTopic",
}}

var _ = Describe("Watcher", func() {
	It("initialises correctly", func() {
		db := test_config.NewTestDB(core.Node{ID: "testNode"})
		fetcher := &mocks.MockLogFetcher{}
		repository := &mocks.MockWatcherRepository{}
		chunker := shared2.NewLogChunker()

		watcher := shared.NewWatcher(db, fetcher, repository, chunker)

		Expect(watcher.DB).To(Equal(db))
		Expect(watcher.Fetcher).To(Equal(fetcher))
		Expect(watcher.Chunker).To(Equal(chunker))
		Expect(watcher.Repository).To(Equal(repository))
	})

	It("adds transformers", func() {
		chunker := shared2.NewLogChunker()
		watcher := shared.NewWatcher(nil, nil, nil, chunker)

		watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformerInitializer}, fakeTransformerConfig)

		Expect(len(watcher.Transformers)).To(Equal(1))
		Expect(watcher.Transformers).To(ConsistOf(&MockTransformer{}))
		Expect(watcher.Topics).To(Equal([]common.Hash{common.HexToHash("FakeTopic")}))
		Expect(watcher.Addresses).To(Equal([]common.Address{common.HexToAddress("FakeAddress")}))
	})

	It("adds transformers from multiple sources", func() {
		chunker := shared2.NewLogChunker()
		watcher := shared.NewWatcher(nil, nil, nil, chunker)

		watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformerInitializer}, fakeTransformerConfig)
		watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformerInitializer}, fakeTransformerConfig)

		Expect(len(watcher.Transformers)).To(Equal(2))
		Expect(watcher.Topics).To(Equal([]common.Hash{common.HexToHash("FakeTopic"),
			common.HexToHash("FakeTopic")}))
		Expect(watcher.Addresses).To(Equal([]common.Address{common.HexToAddress("FakeAddress"),
			common.HexToAddress("FakeAddress")}))
	})

	Describe("with missing headers", func() {
		var (
			db               *postgres.DB
			watcher          shared.Watcher
			fakeTransformer  *MockTransformer
			headerRepository repositories.HeaderRepository
			mockFetcher      mocks.MockLogFetcher
			repository       mocks.MockWatcherRepository
			chunker          *shared2.LogChunker
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			mockFetcher = mocks.MockLogFetcher{}
			headerRepository = repositories.NewHeaderRepository(db)
			_, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			chunker = shared2.NewLogChunker()

			repository = mocks.MockWatcherRepository{}
			watcher = shared.NewWatcher(db, &mockFetcher, &repository, chunker)
		})

		It("executes each transformer", func() {
			fakeTransformer = &MockTransformer{}
			watcher.Transformers = []shared2.Transformer{fakeTransformer}
			repository.SetMissingHeaders([]core.Header{fakes.FakeHeader})

			err := watcher.Execute()

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTransformer.executeWasCalled).To(BeTrue())
		})

		It("returns an error if transformer returns an error", func() {
			fakeTransformer = &MockTransformer{executeError: errors.New("Something bad happened")}
			watcher.Transformers = []shared2.Transformer{fakeTransformer}
			repository.SetMissingHeaders([]core.Header{fakes.FakeHeader})

			err := watcher.Execute()

			Expect(err).To(HaveOccurred())
			Expect(fakeTransformer.executeWasCalled).To(BeFalse())
		})

		It("passes only relevant logs to each transformer", func() {
			transformerA := &MockTransformer{}
			transformerA.SetTransformerName("transformerA")
			transformerB := &MockTransformer{}
			transformerB.SetTransformerName("transformerB")

			configA := shared2.TransformerConfig{TransformerName: "transformerA",
				ContractAddresses: []string{"0x000000000000000000000000000000000000000A"},
				Topic:             "0xA"}
			configB := shared2.TransformerConfig{TransformerName: "transformerB",
				ContractAddresses: []string{"0x000000000000000000000000000000000000000b"},
				Topic:             "0xB"}
			configs := []shared2.TransformerConfig{configA, configB}

			logA := types.Log{Address: common.HexToAddress("0xA"),
				Topics: []common.Hash{common.HexToHash("0xA")}}
			logB := types.Log{Address: common.HexToAddress("0xB"),
				Topics: []common.Hash{common.HexToHash("0xB")}}
			mockFetcher.SetFetchedLogs([]types.Log{logA, logB})

			chunker.AddConfigs(configs)

			repository.SetMissingHeaders([]core.Header{fakes.FakeHeader})
			watcher = shared.NewWatcher(db, &mockFetcher, &repository, chunker)
			watcher.Transformers = []shared2.Transformer{transformerA, transformerB}

			err := watcher.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(transformerA.passedLogs).To(Equal([]types.Log{logA}))
			Expect(transformerB.passedLogs).To(Equal([]types.Log{logB}))
		})

		Describe("uses the repository correctly:", func() {

			It("calls MissingHeaders", func() {
				err := watcher.Execute()
				Expect(err).To(Not(HaveOccurred()))
				Expect(repository.MissingHeadersCalled).To(BeTrue())
			})

			It("propagates MissingHeaders errors", func() {
				missingHeadersError := errors.New("MissingHeadersError")
				repository.MissingHeadersError = missingHeadersError

				err := watcher.Execute()
				Expect(err).To(MatchError(missingHeadersError))
			})

			It("calls CreateNotCheckedSQL", func() {
				err := watcher.Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(repository.CreateNotCheckedSQLCalled).To(BeTrue())
			})

			It("calls GetCheckedColumnNames", func() {
				err := watcher.Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(repository.GetCheckedColumnNamesCalled).To(BeTrue())
			})

			It("propagates GetCheckedColumnNames errors", func() {
				getCheckedColumnNamesError := errors.New("GetCheckedColumnNamesError")
				repository.GetCheckedColumnNamesError = getCheckedColumnNamesError

				err := watcher.Execute()
				Expect(err).To(MatchError(getCheckedColumnNamesError))
			})
		})

		Describe("uses the LogFetcher correctly:", func() {
			BeforeEach(func() {
				repository.SetMissingHeaders([]core.Header{fakes.FakeHeader})
			})

			It("fetches logs", func() {
				err := watcher.Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(mockFetcher.FetchLogsCalled).To(BeTrue())
			})

			It("propagates log fetcher errors", func() {
				fetcherError := errors.New("FetcherError")
				mockFetcher.SetFetcherError(fetcherError)

				err := watcher.Execute()
				Expect(err).To(MatchError(fetcherError))
			})
		})
	})
})
