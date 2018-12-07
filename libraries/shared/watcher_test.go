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
	// TODO Add test for watcher setting the BC
	// TODO Add tests for log chunk delegation
	// TODO Add tests for aggregate fetching
	// TODO Add tests for MissingHeaders

	It("Adds transformers", func() {
		watcher := shared.Watcher{}

		watcher.AddTransformers([]shared2.TransformerInitializer{fakeTransformerInitializer})

		Expect(len(watcher.Transformers)).To(Equal(1))
		Expect(watcher.Transformers).To(ConsistOf(&MockTransformer{}))
	})

	It("Adds transformers from multiple sources", func() {
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
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			watcher = shared.NewWatcher(*db, nil)
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
	})
})
