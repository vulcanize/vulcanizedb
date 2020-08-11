// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package watcher_test

import (
	"database/sql"
	"errors"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/mocks"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/watcher"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage Watcher", func() {
	var statusWriter fakes.MockStatusWriter
	Describe("AddTransformer", func() {
		It("adds transformers", func() {
			fakeAddress := fakes.FakeAddress
			fakeTransformer := &mocks.MockStorageTransformer{Address: fakeAddress}
			w := watcher.NewStorageWatcher(test_config.NewTestDB(test_config.NewTestNode()), -1, &statusWriter)

			w.AddTransformers([]storage.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})

			Expect(w.AddressTransformers[fakeAddress]).To(Equal(fakeTransformer))
		})
	})

	Describe("Execute", func() {
		var (
			storageWatcher       watcher.StorageWatcher
			mockDiffsRepository  *mocks.MockStorageDiffRepository
			mockHeaderRepository *fakes.MockHeaderRepository
		)

		BeforeEach(func() {
			mockDiffsRepository = &mocks.MockStorageDiffRepository{}
			mockHeaderRepository = &fakes.MockHeaderRepository{}
			storageWatcher = watcher.NewStorageWatcher(test_config.NewTestDB(test_config.NewTestNode()), -1, &statusWriter)
			storageWatcher.HeaderRepository = mockHeaderRepository
			storageWatcher.StorageDiffRepository = mockDiffsRepository
		})

		It("creates file for health check", func() {
			mockDiffsRepository.GetNewDiffsErrors = []error{fakes.FakeError}

			err := storageWatcher.Execute()

			Expect(err).To(HaveOccurred())
			Expect(statusWriter.WriteCalled).To(BeTrue())
		})

		It("fetches diffs with results limit", func() {
			mockDiffsRepository.GetNewDiffsErrors = []error{fakes.FakeError}

			err := storageWatcher.Execute()

			Expect(err).To(HaveOccurred())
			Expect(mockDiffsRepository.GetNewDiffsPassedLimits).To(ConsistOf(watcher.ResultsLimit))
		})

		It("fetches diffs with min ID from subsequent queries when previous query returns max results", func() {
			var diffs []types.PersistedDiff
			diffID := rand.Int()
			for i := 0; i < watcher.ResultsLimit; i++ {
				diffID = diffID + i
				diff := types.PersistedDiff{
					RawDiff: types.RawDiff{
						Address: test_data.FakeAddress(),
					},
					ID: int64(diffID),
				}
				diffs = append(diffs, diff)
			}
			mockDiffsRepository.GetNewDiffsDiffs = diffs
			mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}

			err := storageWatcher.Execute()

			Expect(err).To(HaveOccurred())
			Expect(mockDiffsRepository.GetNewDiffsPassedMinIDs).To(ConsistOf(0, diffID))
		})

		It("resets min ID to zero when previous query returns fewer than max results", func() {
			var diffs []types.PersistedDiff
			diffID := rand.Int()
			for i := 0; i < watcher.ResultsLimit-1; i++ {
				diffID = diffID + i
				diff := types.PersistedDiff{
					RawDiff: types.RawDiff{
						Address: test_data.FakeAddress(),
					},
					ID: int64(diffID),
				}
				diffs = append(diffs, diff)
			}
			mockDiffsRepository.GetNewDiffsDiffs = diffs
			mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}

			err := storageWatcher.Execute()

			Expect(err).To(HaveOccurred())
			Expect(mockDiffsRepository.GetNewDiffsPassedMinIDs).To(ConsistOf(0, 0))
		})

		It("marks diff as unwatched if no transformer is watching its address", func() {
			unwatchedDiff := types.PersistedDiff{
				RawDiff: types.RawDiff{
					Address: test_data.FakeAddress(),
				},
				ID: rand.Int63(),
			}
			mockDiffsRepository.GetNewDiffsDiffs = []types.PersistedDiff{unwatchedDiff}
			mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}

			err := storageWatcher.Execute()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
			Expect(mockDiffsRepository.MarkUnwatchedPassedID).To(Equal(unwatchedDiff.ID))
		})

		Describe("When the watcher is configured to skip old diffs", func() {
			var diffs []types.PersistedDiff
			var numberOfBlocksFromHeadOfChain = int64(500)

			BeforeEach(func() {
				storageWatcher = watcher.StorageWatcher{
					HeaderRepository:          mockHeaderRepository,
					StorageDiffRepository:     mockDiffsRepository,
					AddressTransformers:       map[common.Address]storage.ITransformer{},
					DiffBlocksFromHeadOfChain: numberOfBlocksFromHeadOfChain,
					StatusWriter:              &statusWriter,
				}
				diffID := rand.Int()
				for i := 0; i < watcher.ResultsLimit; i++ {
					diffID = diffID + i
					diff := types.PersistedDiff{
						RawDiff: types.RawDiff{
							Address: test_data.FakeAddress(),
						},
						ID: int64(diffID),
					}
					diffs = append(diffs, diff)
				}
			})

			It("skips diffs that are from a block more than n from the head of the chain", func() {
				headerBlockNumber := rand.Int63()
				mockHeaderRepository.MostRecentHeaderBlockNumber = headerBlockNumber

				mockDiffsRepository.GetFirstDiffIDToReturn = diffs[0].ID
				mockDiffsRepository.GetNewDiffsDiffs = diffs
				mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}

				expectedFirstMinDiffID := int(diffs[0].ID - 1)
				expectedSecondMinDiffID := int(diffs[len(diffs)-1].ID)

				err := storageWatcher.Execute()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
				Expect(mockDiffsRepository.GetFirstDiffBlockHeightPassed).To(Equal(headerBlockNumber - numberOfBlocksFromHeadOfChain))
				Expect(mockDiffsRepository.GetNewDiffsPassedMinIDs).To(ConsistOf(expectedFirstMinDiffID, expectedSecondMinDiffID))
			})

			It("resets min ID back to new min diff when previous query returns fewer than max results", func() {
				var diffs []types.PersistedDiff
				diffID := rand.Int()
				for i := 0; i < watcher.ResultsLimit-1; i++ {
					diffID = diffID + i
					diff := types.PersistedDiff{
						RawDiff: types.RawDiff{
							Address: test_data.FakeAddress(),
						},
						ID: int64(diffID),
					}
					diffs = append(diffs, diff)
				}

				headerBlockNumber := rand.Int63()
				mockHeaderRepository.MostRecentHeaderBlockNumber = headerBlockNumber

				mockDiffsRepository.GetFirstDiffIDToReturn = diffs[0].ID
				mockDiffsRepository.GetNewDiffsDiffs = diffs
				mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}

				expectedFirstMinDiffID := int(diffs[0].ID - 1)

				err := storageWatcher.Execute()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
				Expect(mockDiffsRepository.GetFirstDiffBlockHeightPassed).To(Equal(headerBlockNumber - numberOfBlocksFromHeadOfChain))
				Expect(mockDiffsRepository.GetNewDiffsPassedMinIDs).To(ConsistOf(expectedFirstMinDiffID, expectedFirstMinDiffID))
			})

			It("sets minID to 0 if there are no headers with the given block height", func() {
				mockHeaderRepository.MostRecentHeaderBlockNumberErr = sql.ErrNoRows
				mockDiffsRepository.GetNewDiffsDiffs = diffs
				mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}
				err := storageWatcher.Execute()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))

				expectedFirstMinDiffID := 0
				expectedSecondMinDiffID := int(diffs[len(diffs)-1].ID)
				Expect(mockDiffsRepository.GetNewDiffsPassedMinIDs).To(ConsistOf(expectedFirstMinDiffID, expectedSecondMinDiffID))
			})

			It("sets minID to 0 if there are no diffs with given block range", func() {
				mockDiffsRepository.GetFirstDiffIDErr = sql.ErrNoRows
				mockDiffsRepository.GetNewDiffsDiffs = diffs
				mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}
				err := storageWatcher.Execute()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))

				expectedFirstMinDiffID := 0
				expectedSecondMinDiffID := int(diffs[len(diffs)-1].ID)
				Expect(mockDiffsRepository.GetNewDiffsPassedMinIDs).To(ConsistOf(expectedFirstMinDiffID, expectedSecondMinDiffID))
			})
		})

		Describe("when diff's address is watched", func() {
			var (
				contractAddress common.Address
				mockTransformer *mocks.MockStorageTransformer
			)

			BeforeEach(func() {
				contractAddress = test_data.FakeAddress()
				mockTransformer = &mocks.MockStorageTransformer{Address: contractAddress}
				storageWatcher.AddTransformers([]storage.TransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			It("does not mark diff checked if no matching header", func() {
				diffWithoutHeader := types.PersistedDiff{
					RawDiff: types.RawDiff{
						Address:     contractAddress,
						BlockHash:   test_data.FakeHash(),
						BlockHeight: rand.Int(),
					},
					ID: rand.Int63(),
				}
				mockDiffsRepository.GetNewDiffsDiffs = []types.PersistedDiff{diffWithoutHeader}
				mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}
				mockHeaderRepository.GetHeaderByBlockNumberError = errors.New("no matching header")

				err := storageWatcher.Execute()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
				Expect(mockDiffsRepository.MarkCheckedPassedID).NotTo(Equal(diffWithoutHeader.ID))
			})

			Describe("when non-matching header found", func() {
				var (
					blockNumber       int
					fakePersistedDiff types.PersistedDiff
				)

				BeforeEach(func() {
					blockNumber = rand.Int()
					fakeRawDiff := types.RawDiff{
						Address:      contractAddress,
						BlockHash:    test_data.FakeHash(),
						BlockHeight:  blockNumber,
						StorageKey:   test_data.FakeHash(),
						StorageValue: test_data.FakeHash(),
					}
					mockHeaderRepository.GetHeaderByBlockNumberReturnID = int64(blockNumber)
					mockHeaderRepository.GetHeaderByBlockNumberReturnHash = test_data.FakeHash().Hex()

					fakePersistedDiff = types.PersistedDiff{
						RawDiff: fakeRawDiff,
						ID:      rand.Int63(),
					}
					mockDiffsRepository.GetNewDiffsDiffs = []types.PersistedDiff{fakePersistedDiff}
				})

				It("does not mark diff checked if getting max known block height fails", func() {
					mockHeaderRepository.MostRecentHeaderBlockNumberErr = errors.New("getting max header failed")
					mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}

					err := storageWatcher.Execute()

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(fakes.FakeError))
					Expect(mockDiffsRepository.MarkCheckedPassedID).NotTo(Equal(fakePersistedDiff.ID))
				})

				It("marks diff noncanonical if block height less than max known block height minus reorg window", func() {
					mockHeaderRepository.MostRecentHeaderBlockNumber = int64(blockNumber + watcher.ReorgWindow + 1)
					mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}

					err := storageWatcher.Execute()

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(fakes.FakeError))
					Expect(mockDiffsRepository.MarkNoncanonicalPassedID).To(Equal(fakePersistedDiff.ID))
				})

				It("does not mark diff checked if block height is within reorg window", func() {
					mockHeaderRepository.MostRecentHeaderBlockNumber = int64(blockNumber + watcher.ReorgWindow)
					mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}

					err := storageWatcher.Execute()

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(fakes.FakeError))
					Expect(mockDiffsRepository.MarkCheckedPassedID).NotTo(Equal(fakePersistedDiff.ID))
				})
			})

			Describe("when matching header exists", func() {
				var fakePersistedDiff types.PersistedDiff

				BeforeEach(func() {
					fakeBlockHash := test_data.FakeHash()
					fakeRawDiff := types.RawDiff{
						Address:   contractAddress,
						BlockHash: fakeBlockHash,
					}

					mockHeaderRepository.GetHeaderByBlockNumberReturnID = rand.Int63()
					mockHeaderRepository.GetHeaderByBlockNumberReturnHash = fakeBlockHash.Hex()

					fakePersistedDiff = types.PersistedDiff{
						RawDiff: fakeRawDiff,
						ID:      rand.Int63(),
					}
					mockDiffsRepository.GetNewDiffsDiffs = []types.PersistedDiff{fakePersistedDiff}
				})

				It("does not mark diff checked if transformer execution fails", func() {
					mockTransformer.ExecuteErr = errors.New("execute failed")
					mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}

					err := storageWatcher.Execute()

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(fakes.FakeError))
					Expect(mockDiffsRepository.MarkCheckedPassedID).NotTo(Equal(fakePersistedDiff.ID))
				})

				It("marks diff as 'unrecognized' when transforming the diff returns a ErrKeyNotFound error", func() {
					mockTransformer.ExecuteErr = types.ErrKeyNotFound
					mockDiffsRepository.GetNewDiffsErrors = []error{nil, types.ErrKeyNotFound}

					err := storageWatcher.Execute()

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(types.ErrKeyNotFound))
					Expect(mockDiffsRepository.MarkUnrecognizedPassedID).To(Equal(fakePersistedDiff.ID))
				})

				It("marks diff checked if transformer execution doesn't fail", func() {
					mockDiffsRepository.GetNewDiffsDiffs = []types.PersistedDiff{fakePersistedDiff}
					mockDiffsRepository.GetNewDiffsErrors = []error{nil, fakes.FakeError}

					err := storageWatcher.Execute()

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(fakes.FakeError))
					Expect(mockDiffsRepository.MarkCheckedPassedID).To(Equal(fakePersistedDiff.ID))
				})
			})
		})
	})
})
