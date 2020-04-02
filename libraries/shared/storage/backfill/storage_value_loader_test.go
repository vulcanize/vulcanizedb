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

var _ = Describe("StorageValueLoader", func() {
	var (
		bc                             *fakes.MockBlockChain
		keysLookupOne, keysLookupTwo   mocks.MockStorageKeysLookup
		runner                         backfill.StorageValueLoader
		initializerOne, initializerTwo storage.TransformerInitializer
		initializers                   []storage.TransformerInitializer
		keyOne, keyTwo                 common.Hash
		valueOne, valueTwo             common.Hash
		addressOne, addressTwo         common.Address
		blockOne, blockTwo             int64
		bigIntBlockOne, bigIntBlockTwo *big.Int
		fakeHeader                     core.Header
		headerRepo                     fakes.MockHeaderRepository
		diffRepo                       mocks.MockStorageDiffRepository
	)

	BeforeEach(func() {
		bc = fakes.NewMockBlockChain()

		blockOne = rand.Int63()
		bigIntBlockOne = big.NewInt(blockOne)
		blockTwo = blockOne + 1
		bigIntBlockTwo = big.NewInt(blockTwo)

		keysLookupOne = mocks.MockStorageKeysLookup{}
		keyOne = common.Hash{1, 2, 3}
		addressOne = fakes.FakeAddress
		keysLookupOne.KeysToReturn = []common.Hash{keyOne}
		valueOne = common.BytesToHash([]byte{7, 8, 9})
		bc.SetStorageValuesToReturn(blockOne, addressOne, valueOne[:])

		keysLookupTwo = mocks.MockStorageKeysLookup{}
		keyTwo = common.Hash{4, 5, 6}
		addressTwo = fakes.AnotherFakeAddress
		keysLookupTwo.KeysToReturn = []common.Hash{keyTwo}
		valueTwo = common.BytesToHash([]byte{10, 11, 12})
		bc.SetStorageValuesToReturn(blockOne, addressTwo, valueTwo[:])

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
		runner = backfill.NewStorageValueLoader(bc, nil, initializers, blockOne, blockTwo)

		diffRepo = mocks.MockStorageDiffRepository{}
		runner.StorageDiffRepo = &diffRepo

		headerRepo = fakes.MockHeaderRepository{}
		fakeHeader = fakes.FakeHeader
		fakeHeader.BlockNumber = blockOne
		headerRepo.AllHeaders = []core.Header{fakeHeader}
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

	It("fetches headers in the given block range", func() {
		runnerErr := runner.Run()
		Expect(runnerErr).NotTo(HaveOccurred())
		Expect(headerRepo.GetHeadersInRangeStartingBlock).To(Equal(blockOne))
		Expect(headerRepo.GetHeadersInRangeEndingBlock).To(Equal(blockTwo))
	})

	It("returns an error if a header for the given block cannot be retrieved", func() {
		headerRepo.GetHeaderError = fakes.FakeError
		runnerErr := runner.Run()
		Expect(runnerErr).To(HaveOccurred())
		Expect(runnerErr).To(Equal(fakes.FakeError))
	})

	It("gets the storage values for each transformer's keys", func() {
		runnerErr := runner.Run()
		Expect(runnerErr).NotTo(HaveOccurred())
		Expect(keysLookupOne.GetKeysCalled).To(BeTrue())
		Expect(keysLookupTwo.GetKeysCalled).To(BeTrue())

		Expect(bc.BatchGetStorageAtCalls).To(ConsistOf(
			fakes.BatchGetStorageAtCall{BlockNumber: bigIntBlockOne, Account: addressOne, Keys: []common.Hash{keyOne}},
			fakes.BatchGetStorageAtCall{BlockNumber: bigIntBlockOne, Account: addressTwo, Keys: []common.Hash{keyTwo}},
		))
	})

	It("gets storage values from every header in block range", func() {
		headerRepo.AllHeaders = []core.Header{
			{BlockNumber: blockOne},
			{BlockNumber: blockTwo},
		}
		bc.SetStorageValuesToReturn(blockTwo, addressTwo, valueTwo[:])

		runnerErr := runner.Run()
		Expect(runnerErr).NotTo(HaveOccurred())

		Expect(keysLookupOne.GetKeysCalled).To(BeTrue())
		Expect(keysLookupTwo.GetKeysCalled).To(BeTrue())
		Expect(bc.BatchGetStorageAtCalls).To(ContainElement(
			fakes.BatchGetStorageAtCall{BlockNumber: bigIntBlockOne, Account: addressTwo, Keys: []common.Hash{keyTwo}},
		))
		Expect(bc.BatchGetStorageAtCalls).To(ContainElement(
			fakes.BatchGetStorageAtCall{BlockNumber: bigIntBlockTwo, Account: addressTwo, Keys: []common.Hash{keyTwo}},
		))
	})

	It("returns an error if blockchain call to BatchGetStorageAt fails", func() {
		keysLookupOne.KeysToReturn = []common.Hash{keyOne}
		bc.BatchGetStorageAtError = fakes.FakeError

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
			BlockHeight:   int(blockOne),
			BlockHash:     headerHashBytes,
			HashedAddress: crypto.Keccak256Hash(addressOne[:]),
			StorageKey:    keyOne,
			StorageValue:  valueOne,
		}
		expectedDiffTwo := types.RawDiff{
			BlockHeight:   int(blockOne),
			BlockHash:     headerHashBytes,
			HashedAddress: crypto.Keccak256Hash(addressTwo[:]),
			StorageKey:    keyTwo,
			StorageValue:  valueTwo,
		}

		Expect(diffRepo.CreateBackFilledStorageValuePassedRawDiffs).To(ConsistOf(expectedDiffOne, expectedDiffTwo))
	})

	It("ignores sql.ErrNoRows error for duplicate diffs", func() {
		diffRepo.CreateBackFilledStorageValueReturnError = sql.ErrNoRows
		runnerErr := runner.Run()
		Expect(runnerErr).NotTo(HaveOccurred())
	})

	It("returns an error if inserting a diff fails", func() {
		diffRepo.CreateBackFilledStorageValueReturnError = fakes.FakeError
		runnerErr := runner.Run()
		Expect(runnerErr).To(HaveOccurred())
		Expect(runnerErr).To(Equal(fakes.FakeError))
	})
})
