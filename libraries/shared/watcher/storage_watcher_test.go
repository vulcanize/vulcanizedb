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
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Storage Watcher", func() {
	var (
		mockFetcher               *mocks.StorageFetcher
		mockQueue                 *mocks.MockStorageQueue
		mockTransformer           *mocks.MockStorageTransformer
		storageWatcher            *watcher.StorageWatcher
		mockStorageDiffRepository *fakes.MockStorageDiffRepository
		fakeDiffId                = rand.Int63()
		hashedAddress             = utils.HexToKeccak256Hash("0x0123456789abcdef")
		csvDiff                   utils.StorageDiffInput
	)
	Describe("AddTransformer", func() {
		It("adds transformers", func() {
			fakeHashedAddress := utils.HexToKeccak256Hash("0x12345")
			fakeTransformer := &mocks.MockStorageTransformer{KeccakOfAddress: fakeHashedAddress}
			w := watcher.NewStorageWatcher(mocks.NewStorageFetcher(), test_config.NewTestDB(test_config.NewTestNode()))

			w.AddTransformers([]transformer.StorageTransformerInitializer{fakeTransformer.FakeTransformerInitializer})

			Expect(w.KeccakAddressTransformers[fakeHashedAddress]).To(Equal(fakeTransformer))
		})
	})
	Describe("Execute", func() {
		BeforeEach(func() {
			mockFetcher = mocks.NewStorageFetcher()
			mockQueue = &mocks.MockStorageQueue{}
			mockStorageDiffRepository = &fakes.MockStorageDiffRepository{}
			mockTransformer = &mocks.MockStorageTransformer{KeccakOfAddress: hashedAddress}
			csvDiff = utils.StorageDiffInput{
				HashedAddress: hashedAddress,
				BlockHash:     common.HexToHash("0xfedcba9876543210"),
				BlockHeight:   0,
				StorageKey:    common.HexToHash("0xabcdef1234567890"),
				StorageValue:  common.HexToHash("0x9876543210abcdef"),
			}
		})

		It("logs error if fetching storage diffs fails", func(done Done) {
			mockFetcher.ErrsToReturn = []error{fakes.FakeError}
			storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
			storageWatcher.Queue = mockQueue
			storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			storageWatcher.StorageDiffRepository = mockStorageDiffRepository
			tempFile, fileErr := ioutil.TempFile("", "log")
			Expect(fileErr).NotTo(HaveOccurred())
			defer os.Remove(tempFile.Name())
			logrus.SetOutput(tempFile)

			go storageWatcher.Execute(time.Hour, false)

			Eventually(func() (string, error) {
				logContent, err := ioutil.ReadFile(tempFile.Name())
				return string(logContent), err
			}).Should(ContainSubstring(fakes.FakeError.Error()))
			close(done)
		})

		Describe("transforming new storage diffs from csv", func() {
			var fakePersistedDiff utils.PersistedStorageDiff
			BeforeEach(func() {
				mockFetcher.DiffsToReturn = []utils.StorageDiffInput{csvDiff}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
				fakePersistedDiff = utils.PersistedStorageDiff{
					ID: fakeDiffId,
					StorageDiffInput: utils.StorageDiffInput{
						HashedAddress: csvDiff.HashedAddress,
						BlockHash:     csvDiff.BlockHash,
						BlockHeight:   csvDiff.BlockHeight,
						StorageValue:  csvDiff.StorageValue,
						StorageKey:    csvDiff.StorageKey,
					},
				}
				mockStorageDiffRepository.CreateReturnID = fakeDiffId
				storageWatcher.StorageDiffRepository = mockStorageDiffRepository
			})

			It("writes raw diff before processing", func(done Done) {
				go storageWatcher.Execute(time.Hour, false)

				Eventually(func() []utils.StorageDiffInput {
					return mockStorageDiffRepository.CreatePassedInputs
				}).Should(ContainElement(csvDiff))
				close(done)
			})

			It("discards raw diff if it's already been persisted", func(done Done) {
				mockStorageDiffRepository.CreateReturnError = repositories.ErrDuplicateDiff

				go storageWatcher.Execute(time.Hour, false)

				Consistently(func() []utils.PersistedStorageDiff {
					return mockTransformer.PassedDiffs
				}).Should(BeZero())
				close(done)
			})

			It("logs error if persisting raw diff fails", func(done Done) {
				mockStorageDiffRepository.CreateReturnError = fakes.FakeError
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.Execute(time.Hour, false)

				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring(fakes.FakeError.Error()))
				close(done)
			})

			It("executes transformer for recognized storage diff", func(done Done) {
				go storageWatcher.Execute(time.Hour, false)

				Eventually(func() []utils.PersistedStorageDiff {
					return mockTransformer.PassedDiffs
				}).Should(Equal([]utils.PersistedStorageDiff{
					fakePersistedDiff,
				}))
				close(done)
			})

			It("queues diff for later processing if transformer execution fails", func(done Done) {
				mockTransformer.ExecuteErr = fakes.FakeError

				go storageWatcher.Execute(time.Hour, false)

				Eventually(func() bool {
					return mockQueue.AddCalled
				}).Should(BeTrue())
				Eventually(func() utils.PersistedStorageDiff {
					if len(mockQueue.AddPassedDiffs) > 0 {
						return mockQueue.AddPassedDiffs[0]
					}
					return utils.PersistedStorageDiff{}
				}).Should(Equal(fakePersistedDiff))
				close(done)
			})

			It("logs error if queueing diff fails", func(done Done) {
				mockTransformer.ExecuteErr = utils.ErrStorageKeyNotFound{}
				mockQueue.AddError = fakes.FakeError
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.Execute(time.Hour, false)

				Eventually(func() bool {
					return mockQueue.AddCalled
				}).Should(BeTrue())
				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring(fakes.FakeError.Error()))
				close(done)
			})
		})

		Describe("transforming queued storage diffs", func() {
			var queuedDiff utils.PersistedStorageDiff
			BeforeEach(func() {
				queuedDiff = utils.PersistedStorageDiff{
					ID: 1337,
					StorageDiffInput: utils.StorageDiffInput{
						HashedAddress: hashedAddress,
						BlockHash:     test_data.FakeHash(),
						BlockHeight:   rand.Int(),
						StorageKey:    test_data.FakeHash(),
						StorageValue:  test_data.FakeHash(),
					},
				}
				mockQueue.DiffsToReturn = []utils.PersistedStorageDiff{queuedDiff}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			It("executes transformer for storage diff", func(done Done) {
				go storageWatcher.Execute(time.Nanosecond, false)

				Eventually(func() utils.PersistedStorageDiff {
					if len(mockTransformer.PassedDiffs) > 0 {
						return mockTransformer.PassedDiffs[0]
					}
					return utils.PersistedStorageDiff{}
				}).Should(Equal(queuedDiff))
				close(done)
			})

			It("deletes diff from queue if transformer execution successful", func(done Done) {
				go storageWatcher.Execute(time.Nanosecond, false)

				Eventually(func() int64 {
					if len(mockQueue.DeletePassedIds) > 0 {
						return mockQueue.DeletePassedIds[0]
					}
					return 0
				}).Should(Equal(queuedDiff.ID))
				close(done)
			})

			It("logs error if deleting persisted diff fails", func(done Done) {
				mockQueue.DeleteErr = fakes.FakeError
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.Execute(time.Nanosecond, false)

				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring(fakes.FakeError.Error()))
				close(done)
			})

			It("deletes obsolete diff from queue if contract not recognized", func(done Done) {
				obsoleteDiff := utils.PersistedStorageDiff{
					ID:               queuedDiff.ID + 1,
					StorageDiffInput: utils.StorageDiffInput{HashedAddress: test_data.FakeHash()},
				}
				mockQueue.DiffsToReturn = []utils.PersistedStorageDiff{obsoleteDiff}

				go storageWatcher.Execute(time.Nanosecond, false)

				Eventually(func() int64 {
					if len(mockQueue.DeletePassedIds) > 0 {
						return mockQueue.DeletePassedIds[0]
					}
					return 0
				}).Should(Equal(obsoleteDiff.ID))
				close(done)
			})

			It("logs error if deleting obsolete diff fails", func(done Done) {
				obsoleteDiff := utils.PersistedStorageDiff{
					ID:               queuedDiff.ID + 1,
					StorageDiffInput: utils.StorageDiffInput{HashedAddress: test_data.FakeHash()},
				}
				mockQueue.DiffsToReturn = []utils.PersistedStorageDiff{obsoleteDiff}
				mockQueue.DeleteErr = fakes.FakeError
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.Execute(time.Nanosecond, false)

				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring(fakes.FakeError.Error()))
				close(done)
			})
		})
	})

	Describe("BackFill", func() {
		var (
			mockBackFiller       *mocks.BackFiller
			mockTransformer2     *mocks.MockStorageTransformer
			mockTransformer3     *mocks.MockStorageTransformer
			createdPersistedDiff = utils.PersistedStorageDiff{
				ID:               fakeDiffId,
				StorageDiffInput: test_data.CreatedExpectedStorageDiff,
			}
			updatedPersistedDiff1 = utils.PersistedStorageDiff{
				ID:               fakeDiffId,
				StorageDiffInput: test_data.UpdatedExpectedStorageDiff,
			}
			deletedPersistedDiff = utils.PersistedStorageDiff{
				ID:               fakeDiffId,
				StorageDiffInput: test_data.DeletedExpectedStorageDiff,
			}
			updatedPersistedDiff2 = utils.PersistedStorageDiff{
				ID:               fakeDiffId,
				StorageDiffInput: test_data.UpdatedExpectedStorageDiff2,
			}
			csvDiff = utils.StorageDiffInput{
				HashedAddress: hashedAddress,
				BlockHash:     common.HexToHash("0xfedcba9876543210"),
				BlockHeight:   int(test_data.BlockNumber2.Int64()) + 1,
				StorageKey:    common.HexToHash("0xabcdef1234567890"),
				StorageValue:  common.HexToHash("0x9876543210abcdef"),
			}
			csvPersistedDiff = utils.PersistedStorageDiff{
				ID:               fakeDiffId,
				StorageDiffInput: csvDiff,
			}
		)

		BeforeEach(func() {
			mockBackFiller = new(mocks.BackFiller)
			hashedAddress = utils.HexToKeccak256Hash("0x0123456789abcdef")
			mockFetcher = mocks.NewStorageFetcher()
			mockQueue = &mocks.MockStorageQueue{}
			mockTransformer = &mocks.MockStorageTransformer{KeccakOfAddress: hashedAddress}
			mockTransformer2 = &mocks.MockStorageTransformer{KeccakOfAddress: common.BytesToHash(test_data.ContractLeafKey[:])}
			mockTransformer3 = &mocks.MockStorageTransformer{KeccakOfAddress: common.BytesToHash(test_data.AnotherContractLeafKey[:])}
			mockStorageDiffRepository = &fakes.MockStorageDiffRepository{}
		})

		Describe("transforming streamed and backfilled storage diffs", func() {
			BeforeEach(func() {
				mockFetcher.DiffsToReturn = []utils.StorageDiffInput{csvDiff}
				mockBackFiller.SetStorageDiffsToReturn([]utils.StorageDiffInput{
					test_data.CreatedExpectedStorageDiff,
					test_data.UpdatedExpectedStorageDiff,
					test_data.DeletedExpectedStorageDiff,
					test_data.UpdatedExpectedStorageDiff2,
				})
				mockQueue.DiffsToReturn = []utils.PersistedStorageDiff{}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{
					mockTransformer.FakeTransformerInitializer,
					mockTransformer2.FakeTransformerInitializer,
					mockTransformer3.FakeTransformerInitializer,
				})
				mockStorageDiffRepository.CreateReturnID = fakeDiffId

				storageWatcher.StorageDiffRepository = mockStorageDiffRepository
			})

			It("executes transformer for storage diffs received from fetcher and backfiller", func(done Done) {
				go storageWatcher.BackFill(test_data.BlockNumber.Uint64(), mockBackFiller)
				go storageWatcher.Execute(time.Hour, true)

				Eventually(func() int {
					return len(mockTransformer.PassedDiffs)
				}).Should(Equal(1))
				Eventually(func() int {
					return len(mockTransformer2.PassedDiffs)
				}).Should(Equal(1))
				Eventually(func() int {
					return len(mockTransformer3.PassedDiffs)
				}).Should(Equal(3))

				Expect(mockBackFiller.PassedEndingBlock).To(Equal(uint64(test_data.BlockNumber2.Int64())))
				Expect(mockTransformer.PassedDiffs[0]).To(Equal(csvPersistedDiff))
				Expect(mockTransformer2.PassedDiffs[0]).To(Equal(createdPersistedDiff))
				Expect(mockTransformer3.PassedDiffs).To(ConsistOf(updatedPersistedDiff1, deletedPersistedDiff, updatedPersistedDiff2))
				close(done)
			})

			It("adds diffs to the queue if transformation fails", func(done Done) {
				mockTransformer3.ExecuteErr = fakes.FakeError
				go storageWatcher.BackFill(test_data.BlockNumber.Uint64(), mockBackFiller)
				go storageWatcher.Execute(time.Hour, true)

				Eventually(func() int {
					return len(mockTransformer.PassedDiffs)
				}).Should(Equal(1))
				Eventually(func() int {
					return len(mockTransformer2.PassedDiffs)
				}).Should(Equal(1))
				Eventually(func() int {
					return len(mockTransformer3.PassedDiffs)
				}).Should(Equal(3))

				Eventually(func() bool {
					return mockQueue.AddCalled
				}).Should(BeTrue())
				Eventually(func() []utils.PersistedStorageDiff {
					if len(mockQueue.AddPassedDiffs) > 2 {
						return mockQueue.AddPassedDiffs
					}
					return []utils.PersistedStorageDiff{}
				}).Should(ConsistOf(updatedPersistedDiff1, deletedPersistedDiff, updatedPersistedDiff2))

				Expect(mockBackFiller.PassedEndingBlock).To(Equal(uint64(test_data.BlockNumber2.Int64())))
				Expect(mockTransformer.PassedDiffs[0]).To(Equal(csvPersistedDiff))
				Expect(mockTransformer2.PassedDiffs[0]).To(Equal(createdPersistedDiff))
				Expect(mockTransformer3.PassedDiffs).To(ConsistOf(updatedPersistedDiff1, deletedPersistedDiff, updatedPersistedDiff2))
				close(done)
			})

			It("logs a backfill error", func(done Done) {
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				mockBackFiller.BackFillErrs = []error{
					nil,
					nil,
					nil,
					errors.New("mock backfiller error"),
				}

				go storageWatcher.BackFill(test_data.BlockNumber.Uint64(), mockBackFiller)
				go storageWatcher.Execute(time.Hour, true)

				Eventually(func() int {
					return len(mockTransformer.PassedDiffs)
				}).Should(Equal(1))
				Eventually(func() int {
					return len(mockTransformer2.PassedDiffs)
				}).Should(Equal(1))
				Eventually(func() int {
					return len(mockTransformer3.PassedDiffs)
				}).Should(Equal(2))
				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring("mock backfiller error"))
				close(done)
			})

			It("logs when backfill finishes", func(done Done) {
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.BackFill(test_data.BlockNumber.Uint64(), mockBackFiller)
				go storageWatcher.Execute(time.Hour, true)

				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring("storage watcher backfill process has finished"))
				close(done)
			})
		})

		Describe("transforms queued storage diffs", func() {
			BeforeEach(func() {
				mockQueue.DiffsToReturn = []utils.PersistedStorageDiff{
					csvPersistedDiff,
					createdPersistedDiff,
					updatedPersistedDiff1,
					deletedPersistedDiff,
					updatedPersistedDiff2,
				}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{
					mockTransformer.FakeTransformerInitializer,
					mockTransformer2.FakeTransformerInitializer,
					mockTransformer3.FakeTransformerInitializer,
				})
			})

			It("executes transformers on queued storage diffs", func(done Done) {
				go storageWatcher.BackFill(test_data.BlockNumber.Uint64(), mockBackFiller)
				go storageWatcher.Execute(time.Nanosecond, true)

				Eventually(func() int {
					return len(mockTransformer.PassedDiffs)
				}).Should(Equal(1))
				Eventually(func() int {
					return len(mockTransformer2.PassedDiffs)
				}).Should(Equal(1))
				Eventually(func() int {
					return len(mockTransformer3.PassedDiffs)
				}).Should(Equal(3))
				Eventually(func() bool {
					return mockQueue.GetAllCalled
				}).Should(BeTrue())
				expectedIDs := []int64{
					fakeDiffId,
					fakeDiffId,
					fakeDiffId,
					fakeDiffId,
					fakeDiffId,
				}
				Eventually(func() []int64 {
					if len(mockQueue.DeletePassedIds) > 4 {
						return mockQueue.DeletePassedIds
					}
					return []int64{}
				}).Should(Equal(expectedIDs))

				Expect(mockQueue.AddCalled).To(Not(BeTrue()))
				Expect(len(mockQueue.DiffsToReturn)).To(Equal(0))
				Expect(mockTransformer.PassedDiffs[0]).To(Equal(csvPersistedDiff))
				Expect(mockTransformer2.PassedDiffs[0]).To(Equal(createdPersistedDiff))
				Expect(mockTransformer3.PassedDiffs).To(ConsistOf(updatedPersistedDiff1, deletedPersistedDiff, updatedPersistedDiff2))
				close(done)
			})
		})
	})
})
