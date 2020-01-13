package storage_test

import (
	"math/rand"

	"github.com/makerdao/vulcanizedb/libraries/shared/mocks"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage diff extractor", func() {
	var (
		mockFetcher    *mocks.MockStorageFetcher
		mockRepository *mocks.MockStorageDiffRepository
		extractor      storage.DiffExtractor
	)

	BeforeEach(func() {
		mockFetcher = mocks.NewMockStorageFetcher()
		mockRepository = &mocks.MockStorageDiffRepository{}
		extractor = storage.DiffExtractor{
			StorageDiffRepository: mockRepository,
			StorageFetcher:        mockFetcher,
		}
	})

	Describe("ExtractDiffs", func() {
		It("fetches storage diffs", func() {
			mockFetcher.ErrsToReturn = []error{fakes.FakeError}

			_ = extractor.ExtractDiffs()

			Expect(mockFetcher.FetchStorageDiffsCalled).To(BeTrue())
		})

		It("returns error if fetching storage diffs fails", func() {
			mockFetcher.ErrsToReturn = []error{fakes.FakeError}

			err := extractor.ExtractDiffs()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("persists fetched storage diff", func() {
			fakeDiff := types.RawDiff{
				HashedAddress: test_data.FakeHash(),
				BlockHash:     test_data.FakeHash(),
				BlockHeight:   rand.Int(),
				StorageKey:    test_data.FakeHash(),
				StorageValue:  test_data.FakeHash(),
			}
			mockFetcher.DiffsToReturn = []types.RawDiff{fakeDiff}
			mockFetcher.ErrsToReturn = []error{fakes.FakeError}

			_ = extractor.ExtractDiffs()

			Expect(mockRepository.CreatePassedRawDiffs).To(Equal([]types.RawDiff{fakeDiff}))
		})
	})
})
