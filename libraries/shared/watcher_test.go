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
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
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
		db := test_config.NewTestDB(core.Node{ID: "testNode"})
		fetcher := mocks.MockLogFetcher{}
		repository := &mocks.MockWatcherRepository{}
		configA := shared2.TransformerConfig{
			ContractAddresses: []string{"0xA"},
			Topic: "0xA",
		}
		configB := shared2.TransformerConfig{
			ContractAddresses: []string{"0xB"},
			Topic: "0xB",
		}
		configs := []shared2.TransformerConfig{configA, configB}
		watcher := shared.NewWatcher(db, &fetcher, repository, configs)

		Expect(watcher.DB).To(Equal(db))
		Expect(watcher.Fetcher).NotTo(BeNil())
		Expect(watcher.Chunker).NotTo(BeNil())
		Expect(watcher.Repository).To(Equal(repository))
		Expect(watcher.Topics).To(And(
			ContainElement(common.HexToHash("0xA")), ContainElement(common.HexToHash("0xB"))))
		Expect(watcher.Addresses).To(And(
			ContainElement(common.HexToAddress("0xA")), ContainElement(common.HexToAddress("0xB"))))
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
			mockFetcher      mocks.MockLogFetcher
			repository mocks.MockWatcherRepository
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			mockFetcher = mocks.MockLogFetcher{}
			headerRepository = repositories.NewHeaderRepository(db)
			_, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			repository = mocks.MockWatcherRepository{}
			watcher = shared.NewWatcher(db, &mockFetcher, &repository, transformers.TransformerConfigs())
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
			// TODO Test log delegation
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
