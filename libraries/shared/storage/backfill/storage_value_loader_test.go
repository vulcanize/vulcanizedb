package backfill_test

import (
	"database/sql"
	"math/big"
	"math/rand"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/mocks"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/backfill"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("getStorageValue Command", func() {
	var (
		bc                             *fakes.MockBlockChain
		keysLookupOne, keysLookupTwo   mocks.MockStorageKeysLookup
		runner                         backfill.StorageValueLoader
		initializerOne, initializerTwo storage.TransformerInitializer
		initializers                   []storage.TransformerInitializer
		keyOne, keyTwo                 common.Hash
		valueOne, valueTwo             common.Hash
		addressOne, addressTwo         common.Address
		blockNumber                    int64
		bigIntBlockNumber              *big.Int
		fakeHeader                     core.Header
		headerRepo                     fakes.MockHeaderRepository
		diffRepo                       mocks.MockStorageDiffRepository
	)

	BeforeEach(func() {
		bc = fakes.NewMockBlockChain()

		keysLookupOne = mocks.MockStorageKeysLookup{}
		keyOne = common.Hash{1, 2, 3}
		addressOne = fakes.FakeAddress
		keysLookupOne.KeysToReturn = []common.Hash{keyOne}
		valueOne = common.BytesToHash([]byte{7, 8, 9})
		bc.SetStorageValuesToReturn(addressOne, valueOne[:])

		keysLookupTwo = mocks.MockStorageKeysLookup{}
		keyTwo = common.Hash{4, 5, 6}
		addressTwo = fakes.AnotherFakeAddress
		keysLookupTwo.KeysToReturn = []common.Hash{keyTwo}
		valueTwo = common.BytesToHash([]byte{10, 11, 12})
		bc.SetStorageValuesToReturn(addressTwo, valueTwo[:])

		initializerOne = storage.Transformer{
			Address:           addressOne,
			StorageKeysLookup: &keysLookupOne,
			Repository:        &mocks.MockStorageRepository{},
		}.NewTransformer

		initializerTwo = storage.Transformer{
			Address:           addressTwo,
			StorageKeysLookup: &keysLookupTwo,
			Repository:        &mocks.MockStorageRepository{},
		}.NewTransformer

		initializers = []storage.TransformerInitializer{initializerOne, initializerTwo}
		blockNumber = rand.Int63()
		bigIntBlockNumber = big.NewInt(blockNumber)

		runner = backfill.NewStorageValueLoader(bc, nil, initializers, blockNumber)

		diffRepo = mocks.MockStorageDiffRepository{}
		runner.StorageDiffRepo = &diffRepo

		headerRepo = fakes.MockHeaderRepository{}
		fakeHeader = fakes.FakeHeader
		fakeHeader.BlockNumber = blockNumber
		headerRepo.GetHeaderReturnHash = fakeHeader.Hash
		runner.HeaderRepo = &headerRepo
	})

	It("gets the storage keys for each transformer", func() {
		runnerErr := runner.Run()
		Expect(runnerErr).NotTo(HaveOccurred())

		Expect(keysLookupOne.GetKeysCalled).To(BeTrue())
		Expect(keysLookupTwo.GetKeysCalled).To(BeTrue())
	})

	It("returns an error if getting the keys from the KeysLookup fails", func() {
		keysLookupTwo.GetKeysError = fakes.FakeError

		runnerErr := runner.Run()
		Expect(keysLookupOne.GetKeysCalled).To(BeTrue())
		Expect(runnerErr).To(HaveOccurred())
		Expect(runnerErr).To(Equal(fakes.FakeError))
	})

	It("fetches the header by the given block number", func() {
		runnerErr := runner.Run()
		Expect(runnerErr).NotTo(HaveOccurred())
		Expect(headerRepo.GetHeaderPassedBlockNumber).To(Equal(blockNumber))
	})

	It("returns an error if a header for the given block cannot be retrieved", func() {
		headerRepo.GetHeaderError = fakes.FakeError
		runnerErr := runner.Run()
		Expect(runnerErr).To(HaveOccurred())
		Expect(runnerErr).To(Equal(fakes.FakeError))
	})

	It("gets the storage values for each of the transformer's keys", func() {
		keysLookupOne.KeysToReturn = []common.Hash{keyOne}
		keysLookupTwo.KeysToReturn = []common.Hash{keyTwo}

		runnerErr := runner.Run()
		Expect(runnerErr).NotTo(HaveOccurred())
		Expect(keysLookupOne.GetKeysCalled).To(BeTrue())
		Expect(keysLookupTwo.GetKeysCalled).To(BeTrue())

		Expect(bc.GetStorageAtPassedBlockNumber).To(Equal(bigIntBlockNumber))
		Expect(bc.GetStorageAtPassedAccounts).To(ConsistOf(addressOne, addressTwo))
		Expect(bc.GetStorageAtPassedKeys).To(ConsistOf(keyOne, keyTwo))
	})

	It("returns an error if blockchain call to GetStorageAt fails", func() {
		keysLookupOne.KeysToReturn = []common.Hash{keyOne}
		bc.SetGetStorageAtError(fakes.FakeError)

		runnerErr := runner.Run()
		Expect(keysLookupOne.GetKeysCalled).To(BeTrue())
		Expect(runnerErr).To(HaveOccurred())
		Expect(runnerErr).To(Equal(fakes.FakeError))
	})

	It("persists the storage values for each transformer", func() {
		runnerErr := runner.Run()
		Expect(runnerErr).NotTo(HaveOccurred())

		trimmedHeaderHash := strings.TrimPrefix(fakeHeader.Hash, "0x")
		headerHashBytes := common.HexToHash(trimmedHeaderHash)
		expectedDiffOne := types.RawDiff{
			BlockHeight:   int(blockNumber),
			BlockHash:     headerHashBytes,
			HashedAddress: crypto.Keccak256Hash(addressOne[:]),
			StorageKey:    keyOne,
			StorageValue:  valueOne,
		}
		expectedDiffTwo := types.RawDiff{
			BlockHeight:   int(blockNumber),
			BlockHash:     headerHashBytes,
			HashedAddress: crypto.Keccak256Hash(addressTwo[:]),
			StorageKey:    keyTwo,
			StorageValue:  valueTwo,
		}

		Expect(diffRepo.CreatePassedRawDiffs).To(ConsistOf(expectedDiffOne, expectedDiffTwo))
	})

	It("ignores sql.ErrNoRows error for duplicate diffs", func() {
		diffRepo.CreateReturnError = sql.ErrNoRows
		runnerErr := runner.Run()
		Expect(runnerErr).NotTo(HaveOccurred())
	})

	It("returns an error if inserting a diff fails", func() {
		diffRepo.CreateReturnError = fakes.FakeError
		runnerErr := runner.Run()
		Expect(runnerErr).To(HaveOccurred())
		Expect(runnerErr).To(Equal(fakes.FakeError))
	})

	It("sets the diff a from_backfill", func() {
		runnerErr := runner.Run()
		Expect(runnerErr).NotTo(HaveOccurred())
		Expect(diffRepo.MarkFromBackfillCalled).To(BeTrue())
	})

	It("returns an error if setting the diff as from_backfill fails", func() {
		diffRepo.MarkFromBackfillError = fakes.FakeError
		runnerErr := runner.Run()
		Expect(runnerErr).To(HaveOccurred())
		Expect(runnerErr).To(Equal(fakes.FakeError))
	})
})
