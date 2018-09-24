package shared_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockTransformer struct {
	executeWasCalled bool
	executeError     error
}

func (mh *MockTransformer) Execute() error {
	if mh.executeError != nil {
		return mh.executeError
	}
	mh.executeWasCalled = true
	return nil
}

func fakeTransformerInitializer(db *postgres.DB, blockchain core.BlockChain, con shared.ContractConfig) (shared.Transformer, error) {
	return &MockTransformer{}, nil
}

var _ = Describe("Watcher", func() {
	It("Adds transformers", func() {
		watcher := shared.Watcher{}
		con := shared.ContractConfig{}

		watcher.AddTransformers([]shared.TransformerInitializer{fakeTransformerInitializer}, con)

		Expect(len(watcher.Transformers)).To(Equal(1))
		Expect(watcher.Transformers).To(ConsistOf(&MockTransformer{}))
	})

	It("Adds transformers from multiple sources", func() {
		watcher := shared.Watcher{}
		con := shared.ContractConfig{}

		watcher.AddTransformers([]shared.TransformerInitializer{fakeTransformerInitializer}, con)
		watcher.AddTransformers([]shared.TransformerInitializer{fakeTransformerInitializer}, con)

		Expect(len(watcher.Transformers)).To(Equal(2))
	})

	It("Executes each transformer", func() {
		watcher := shared.Watcher{}
		fakeTransformer := &MockTransformer{}
		watcher.Transformers = []shared.Transformer{fakeTransformer}

		watcher.Execute()

		Expect(fakeTransformer.executeWasCalled).To(BeTrue())
	})

	It("Returns an error if transformer returns an error", func() {
		watcher := shared.Watcher{}
		fakeTransformer := &MockTransformer{executeError: errors.New("Something bad happened")}
		watcher.Transformers = []shared.Transformer{fakeTransformer}

		err := watcher.Execute()

		Expect(err).To(HaveOccurred())
		Expect(fakeTransformer.executeWasCalled).To(BeFalse())
	})
})
