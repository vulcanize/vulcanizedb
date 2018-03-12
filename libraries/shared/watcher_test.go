package shared_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockHandler struct {
	executeWasCalled bool
	executeError     error
}

func (mh *MockHandler) Execute() error {
	if mh.executeError != nil {
		return mh.executeError
	}
	mh.executeWasCalled = true
	return nil
}

func fakeHandlerInitializer(db *postgres.DB, blockchain core.Blockchain) shared.Handler {
	return &MockHandler{}
}

var _ = Describe("Watcher", func() {
	It("Adds handlers", func() {
		watcher := shared.Watcher{}

		watcher.AddHandlers([]shared.HandlerInitializer{fakeHandlerInitializer})

		Expect(len(watcher.Handlers)).To(Equal(1))
		Expect(watcher.Handlers).To(ConsistOf(&MockHandler{}))
	})

	It("Adds handlers from multiple sources", func() {
		watcher := shared.Watcher{}

		watcher.AddHandlers([]shared.HandlerInitializer{fakeHandlerInitializer})
		watcher.AddHandlers([]shared.HandlerInitializer{fakeHandlerInitializer})

		Expect(len(watcher.Handlers)).To(Equal(2))
	})

	It("Executes each handler", func() {
		watcher := shared.Watcher{}
		fakeHandler := &MockHandler{}
		watcher.Handlers = []shared.Handler{fakeHandler}

		watcher.Execute()

		Expect(fakeHandler.executeWasCalled).To(BeTrue())
	})

	It("Returns an error if handler returns an error", func() {
		watcher := shared.Watcher{}
		fakeHandler := &MockHandler{executeError: errors.New("Something bad happened")}
		watcher.Handlers = []shared.Handler{fakeHandler}

		err := watcher.Execute()

		Expect(err).To(HaveOccurred())
		Expect(fakeHandler.executeWasCalled).To(BeFalse())
	})
})
