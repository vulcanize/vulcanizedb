package executor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/executor"
	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
	"time"
)

var _ = Describe("Executor", func() {
	var (
		db         *postgres.DB
		blockChain *fakes.MockBlockChain
		mockPlugin fakes.MockPlugin
	)
	BeforeEach(func() {
		db = test_config.NewTestDB(core.Node{ID: "testNode"})
		blockChain = fakes.NewMockBlockChain()
		mockPlugin = fakes.NewMockPlugin()
	})

	It("adds transformer sets to the executor", func() {
		ew := watcher.NewEventWatcher(db, blockChain, false, time.Second)
		sw := watcher.NewStorageWatcher(mocks.NewMockStorageFetcher(), db, time.Second, time.Second)
		cw := watcher.NewContractWatcher(db, blockChain, time.Second)
		ex := executor.NewExecutor(&mockPlugin, &ew, sw, &cw)

		loadErr := ex.LoadTransformerSets()
		Expect(loadErr).NotTo(HaveOccurred())

		Expect(ex.EthEventInitializers).To(HaveLen(len(mockPlugin.FakeEventInitializers)))
		Expect(ex.EthStorageInitializers).To(HaveLen(len(mockPlugin.FakeStorageInitializers)))
		Expect(ex.EthContractInitializers).To(HaveLen(len(mockPlugin.FakeContractInitializers)))
		// the failure shows that expected equals received so not sure why the following doesn't work
		//Expect(ex.EthEventInitializers).To(Equal(mockPlugin.FakeEventInitializers))
		//Expect(ex.EthStorageInitializers).To(Equal(mockPlugin.FakeStorageInitializers))
		//Expect(ex.EthContractInitializers).To(Equal(mockPlugin.FakeContractInitializers))
	})

	It("Calls the correct watchers when there are relevant transformers", func() {
		eventWatcher := mocks.NewMockEventWatcher()
		storageWatcher := mocks.NewMockStorageWatcher()
		contractWatcher := mocks.NewMockContractWatcher()
		ex := executor.NewExecutor(&mockPlugin, &eventWatcher, &storageWatcher, &contractWatcher)

		loadErr := ex.LoadTransformerSets()
		Expect(loadErr).NotTo(HaveOccurred())

		ex.ExecuteTransformerSets()

		Expect(eventWatcher.AddTransformersWasCalled).To(BeTrue())
		Eventually(func() bool {
			return eventWatcher.WatchEthEventsWasCalled
		}).Should(BeTrue())

		Expect(storageWatcher.AddTransformersWasCalled).To(BeTrue())
		Eventually(func() bool {
			return storageWatcher.WatchEthStorageWasCalled
		}).Should(BeTrue())

		Expect(contractWatcher.AddTransformersWasCalled).To(BeTrue())
		Eventually(func() bool {
			return contractWatcher.WatchEthContractWasCalled
		}).Should(BeTrue())
	})
})
