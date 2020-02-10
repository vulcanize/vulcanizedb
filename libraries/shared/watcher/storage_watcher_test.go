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
	"errors"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/libraries/shared/mocks"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/libraries/shared/watcher"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage Watcher", func() {
	Describe("AddTransformer", func() {
		It("adds transformers", func() {
			fakeHashedAddress := types.HexToKeccak256Hash("0x12345")
			fakeTransformer := &mocks.MockStorageTransformer{KeccakOfAddress: fakeHashedAddress}
			w := watcher.NewStorageWatcher(test_config.NewTestDB(test_config.NewTestNode()), time.Nanosecond)

			w.AddTransformers([]transformer.StorageTransformerInitializer{fakeTransformer.FakeTransformerInitializer})

			Expect(w.KeccakAddressTransformers[fakeHashedAddress]).To(Equal(fakeTransformer))
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
			storageWatcher = watcher.StorageWatcher{
				HeaderRepository:          mockHeaderRepository,
				StorageDiffRepository:     mockDiffsRepository,
				KeccakAddressTransformers: map[common.Hash]transformer.StorageTransformer{},
				RetryInterval:             time.Nanosecond,
			}

		})

		It("marks diff as checked if no transformer is watching its address", func() {
			unwatchedDiff := types.PersistedDiff{
				RawDiff: types.RawDiff{
					HashedAddress: test_data.FakeHash(),
				},
				ID: rand.Int63(),
			}
			mockDiffsRepository.GetNewDiffsDiffs = []types.PersistedDiff{unwatchedDiff}
			mockDiffsRepository.GetNewDiffsErrors = []error{fakes.FakeError}

			err := storageWatcher.Execute()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
			Expect(mockDiffsRepository.MarkCheckedPassedID).To(Equal(unwatchedDiff.ID))
		})

		Describe("when diff's address is watched", func() {
			var (
				hashedAddress   common.Hash
				mockTransformer *mocks.MockStorageTransformer
			)

			BeforeEach(func() {
				hashedAddress = types.HexToKeccak256Hash("0x" + fakes.RandomString(20))
				mockTransformer = &mocks.MockStorageTransformer{KeccakOfAddress: hashedAddress}
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			It("does not mark diff checked if no matching header", func() {
				diffWithoutHeader := types.PersistedDiff{
					RawDiff: types.RawDiff{
						HashedAddress: hashedAddress,
						BlockHash:     test_data.FakeHash(),
						BlockHeight:   rand.Int(),
					},
					ID:       rand.Int63(),
					HeaderID: rand.Int63(),
				}
				mockDiffsRepository.GetNewDiffsDiffs = []types.PersistedDiff{diffWithoutHeader}
				mockHeaderRepository.GetHeaderError = errors.New("no matching header")
				mockDiffsRepository.GetNewDiffsErrors = []error{fakes.FakeError}

				err := storageWatcher.Execute()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
				Expect(mockDiffsRepository.MarkCheckedPassedID).NotTo(Equal(diffWithoutHeader.ID))
			})

			Describe("when matching header exists", func() {
				var (
					fakeBlockHash     common.Hash
					fakePersistedDiff types.PersistedDiff
				)

				BeforeEach(func() {
					fakeBlockHash = test_data.FakeHash()
					fakeRawDiff := types.RawDiff{
						HashedAddress: hashedAddress,
						BlockHash:     fakeBlockHash,
						BlockHeight:   0,
						StorageKey:    common.HexToHash("0xabcdef1234567890"),
						StorageValue:  common.HexToHash("0x9876543210abcdef"),
					}

					fakeHeaderID := rand.Int63()
					mockHeaderRepository.GetHeaderReturnID = fakeHeaderID
					mockHeaderRepository.GetHeaderReturnHash = fakeBlockHash.Hex()

					fakePersistedDiff = types.PersistedDiff{
						RawDiff:  fakeRawDiff,
						ID:       rand.Int63(),
						HeaderID: fakeHeaderID,
					}
					mockDiffsRepository.GetNewDiffsDiffs = []types.PersistedDiff{fakePersistedDiff}
				})

				It("does not mark diff checked if transformer execution fails", func() {
					mockTransformer.ExecuteErr = errors.New("execute failed")
					mockDiffsRepository.GetNewDiffsErrors = []error{fakes.FakeError}

					err := storageWatcher.Execute()

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
					Expect(mockDiffsRepository.MarkCheckedPassedID).NotTo(Equal(fakePersistedDiff.ID))
				})

				Describe("when transformer execution succeeds", func() {
					It("marks diff checked", func() {
						mockDiffsRepository.GetNewDiffsDiffs = []types.PersistedDiff{fakePersistedDiff}
						mockDiffsRepository.GetNewDiffsErrors = []error{fakes.FakeError}

						err := storageWatcher.Execute()

						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
						Eventually(mockDiffsRepository.MarkCheckedPassedID).Should(Equal(fakePersistedDiff.ID))
					})
				})
			})
		})
	})
})
