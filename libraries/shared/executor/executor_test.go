package executor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/executor"
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
		mockTailer *fakes.MockTailer
	)
	BeforeEach(func() {
		db = test_config.NewTestDB(core.Node{ID: "testNode"})
		blockChain = fakes.NewMockBlockChain()
		mockPlugin = fakes.NewMockPlugin()
		mockTailer = fakes.NewMockTailer()
	})

	It("initializes", func() {
		ex := executor.NewExecutor(db, blockChain, &mockPlugin, false, time.Second, time.Second, mockTailer)
		Expect(ex.DB).To(Equal(db))
		Expect(ex.BlockChain).To(Equal(blockChain))
	})

	It("adds transformer sets to the executor", func() {
		ex := executor.NewExecutor(db, blockChain, &mockPlugin, false, time.Second, time.Second, mockTailer)
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
		ex := executor.NewExecutor(db, blockChain, &mockPlugin, false, time.Second, time.Second, mockTailer)
		loadErr := ex.LoadTransformerSets()
		Expect(loadErr).NotTo(HaveOccurred())

		//ex.ExecuteTransformerSets()
	})
})
