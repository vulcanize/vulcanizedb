package backfill_test

import (
	"math/big"
	"math/rand"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/mocks"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/backfill"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StorageValueLoader", func() {
	var (
		bc                                               *fakes.MockBlockChain
		keysLookupOne, keysLookupTwo, keysLookupThree    mocks.MockStorageKeysLookup
		runner                                           backfill.StorageValueLoader
		initializerOne, initializerTwo, initializerThree storage.TransformerInitializer
		initializers                                     []storage.TransformerInitializer
		keyOne, keyTwo, keyThree                         common.Hash
		valueOne, valueTwo, valueThree                   common.Hash
		addressOne, addressTwo, addressThree             common.Address
		blockOne, blockTwo, blockThree                   int64
		bigIntBlockOne, bigIntBlockTwo, bigIntBlockThree *big.Int
		fakeHeader                                       core.Header
		headerRepo                                       fakes.MockHeaderRepository
		diffRepo                                         mocks.MockStorageDiffRepository
	)

	BeforeEach(func() {
		bc = fakes.NewMockBlockChain()

		blockOne = rand.Int63()
		bigIntBlockOne = big.NewInt(blockOne)
		blockTwo = blockOne + 1
		bigIntBlockTwo = big.NewInt(blockTwo)
		blockThree = blockTwo + 1
		bigIntBlockThree = big.NewInt(blockThree)

		keysLookupOne = mocks.MockStorageKeysLookup{}
		keyOne = test_data.FakeHash()
		addressOne = test_data.FakeAddress()
		keysLookupOne.KeysToReturn = []common.Hash{keyOne}
		valueOne = test_data.FakeHash()
		bc.SetStorageValuesToReturn(blockOne, addressOne, valueOne[:])

		keysLookupTwo = mocks.MockStorageKeysLookup{}
		keyTwo = test_data.FakeHash()
		addressTwo = test_data.FakeAddress()
		keysLookupTwo.KeysToReturn = []common.Hash{keyTwo}
		valueTwo = test_data.FakeHash()
		bc.SetStorageValuesToReturn(blockOne, addressTwo, valueTwo[:])

		keysLookupThree = mocks.MockStorageKeysLookup{}
		keyThree = test_data.FakeHash()
		addressThree = test_data.FakeAddress()
		keysLookupThree.KeysToReturn = []common.Hash{keyThree}
		valueThree = common.BytesToHash([]byte{})
		bc.SetStorageValuesToReturn(blockThree, addressThree, valueThree[:])

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

		initializerThree = storage.Transformer{
			Address:           addressThree,
			StorageKeysLookup: &keysLookupThree,
			Repository:        &mocks.MockStorageRepository{},
		}.NewTransformer

		initializers = []storage.TransformerInitializer{initializerOne, initializerTwo, initializerThree}
		runner = backfill.NewStorageValueLoader(bc, nil, initializers, blockOne, blockThree)

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
		Expect(headerRepo.GetHeadersInRangeEndingBlock).To(Equal(blockThree))
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
			fakes.BatchGetStorageAtCall{BlockNumber: bigIntBlockOne, Account: addressThree, Keys: []common.Hash{keyThree}},
		))
	})

	It("chunks requests to avoid 413 Request Entity Too Large reply from server", func() {
		manyKeys := make([]common.Hash, 401)
		for index, _ := range manyKeys {
			manyKeys[index] = test_data.FakeHash()
		}
		keysLookupTwo.KeysToReturn = manyKeys

		runnerErr := runner.Run()
		Expect(runnerErr).NotTo(HaveOccurred())
		Expect(keysLookupOne.GetKeysCalled).To(BeTrue())
		Expect(keysLookupTwo.GetKeysCalled).To(BeTrue())

		Expect(bc.BatchGetStorageAtCalls).To(ContainElement(fakes.BatchGetStorageAtCall{BlockNumber: bigIntBlockOne, Account: addressTwo, Keys: manyKeys[0:400]}))
		Expect(bc.BatchGetStorageAtCalls).To(ContainElement(fakes.BatchGetStorageAtCall{BlockNumber: bigIntBlockOne, Account: addressTwo, Keys: manyKeys[400:]}))

	})

	It("gets storage values from every header in block range", func() {
		headerRepo.AllHeaders = []core.Header{
			{BlockNumber: blockOne},
			{BlockNumber: blockTwo},
			{BlockNumber: blockThree},
		}

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
		Expect(bc.BatchGetStorageAtCalls).To(ContainElement(
			fakes.BatchGetStorageAtCall{BlockNumber: bigIntBlockThree, Account: addressThree, Keys: []common.Hash{keyThree}},
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

	It("persists the non-zero storage values for each transformer", func() {
		runnerErr := runner.Run()
		Expect(runnerErr).NotTo(HaveOccurred())

		trimmedHeaderHash := strings.TrimPrefix(fakeHeader.Hash, "0x")
		headerHashBytes := common.HexToHash(trimmedHeaderHash)
		expectedDiffOne := types.RawDiff{
			BlockHeight:   int(blockOne),
			BlockHash:     headerHashBytes,
			HashedAddress: crypto.Keccak256Hash(addressOne[:]),
			StorageKey:    crypto.Keccak256Hash(keyOne.Bytes()),
			StorageValue:  valueOne,
		}
		expectedDiffTwo := types.RawDiff{
			BlockHeight:   int(blockOne),
			BlockHash:     headerHashBytes,
			HashedAddress: crypto.Keccak256Hash(addressTwo[:]),
			StorageKey:    crypto.Keccak256Hash(keyTwo.Bytes()),
			StorageValue:  valueTwo,
		}

		Expect(diffRepo.CreateBackFilledStorageValuePassedRawDiffs).To(ConsistOf(expectedDiffOne, expectedDiffTwo))
	})

	It("returns an error if inserting a diff fails", func() {
		diffRepo.CreateBackFilledStorageValueReturnError = fakes.FakeError
		runnerErr := runner.Run()
		Expect(runnerErr).To(HaveOccurred())
		Expect(runnerErr).To(Equal(fakes.FakeError))
	})
})
