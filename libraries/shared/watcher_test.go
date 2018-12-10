package shared_test

import (
	"errors"
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
	return "MockTransformer"
}

func fakeTransformerInitializer(db *postgres.DB) shared2.Transformer {
	return &MockTransformer{}
}

var _ = Describe("Watcher", func() {
	It("initialises correctly", func() {
		// TODO Test watcher initialisation
	})

	It("adds transformers", func() {
		watcher := shared.Watcher{}

		watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformerInitializer})

		Expect(len(watcher.Transformers)).To(Equal(1))
		Expect(watcher.Transformers).To(ConsistOf(&MockTransformer{}))
	})

	It("adds transformers from multiple sources", func() {
		watcher := shared.Watcher{}

		watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformerInitializer})
		watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformerInitializer})

		Expect(len(watcher.Transformers)).To(Equal(2))
	})

	Describe("with missing headers", func() {
		var (
			db               *postgres.DB
			watcher          shared.Watcher
			fakeTransformer  *MockTransformer
			headerRepository repositories.HeaderRepository
			mockFetcher      shared2.LogFetcher
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			mockFetcher = &mocks.MockLogFetcher{}
			watcher = shared.NewWatcher(db, mockFetcher)
			headerRepository = repositories.NewHeaderRepository(db)
			_, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("executes each transformer", func() {
			fakeTransformer = &MockTransformer{}
			watcher.Transformers = []shared2.Transformer{fakeTransformer}

			err := watcher.Execute()

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTransformer.executeWasCalled).To(BeTrue())
		})

		It("returns an error if transformer returns an error", func() {
			fakeTransformer = &MockTransformer{executeError: errors.New("Something bad happened")}
			watcher.Transformers = []shared2.Transformer{fakeTransformer}

			err := watcher.Execute()

			Expect(err).To(HaveOccurred())
			Expect(fakeTransformer.executeWasCalled).To(BeFalse())
		})

		It("calls MissingHeaders", func() {
			// TODO Tests for calling MissingHeaders
		})

		It("returns an error if missingHeaders returns an error", func() {
			// TODO Test for propagating missingHeaders error
		})

		It("calls the log fetcher", func() {
			// TODO Test for calling FetchLogs
		})

		It("returns an error if the log fetcher returns an error", func() {
			// TODO Test for propagating log fetcher error
		})

		It("passes only relevant logs to each transformer", func() {
			// TODO Test log delegation
		})
	})
})
