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
	"os"
	"sort"
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
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Storage Watcher", func() {
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
		var (
			mockFetcher     *mocks.StorageFetcher
			mockQueue       *mocks.MockStorageQueue
			mockTransformer *mocks.MockStorageTransformer
			csvDiff         utils.StorageDiff
			storageWatcher  *watcher.StorageWatcher
			hashedAddress   common.Hash
		)

		BeforeEach(func() {
			hashedAddress = utils.HexToKeccak256Hash("0x0123456789abcdef")
			mockFetcher = mocks.NewStorageFetcher()
			mockQueue = &mocks.MockStorageQueue{}
			mockTransformer = &mocks.MockStorageTransformer{KeccakOfAddress: hashedAddress}
			csvDiff = utils.StorageDiff{
				ID:            1337,
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
			BeforeEach(func() {
				mockFetcher.DiffsToReturn = []utils.StorageDiff{csvDiff}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			It("executes transformer for recognized storage diff", func(done Done) {
				go storageWatcher.Execute(time.Hour, false)

				Eventually(func() map[int]utils.StorageDiff {
					return mockTransformer.PassedDiffs
				}).Should(Equal(map[int]utils.StorageDiff{
					csvDiff.ID: csvDiff,
				}))
				close(done)
			})

			It("queues diff for later processing if transformer execution fails", func(done Done) {
				mockTransformer.ExecuteErr = fakes.FakeError

				go storageWatcher.Execute(time.Hour, false)

				Eventually(func() bool {
					return mockQueue.AddCalled
				}).Should(BeTrue())
				Eventually(func() utils.StorageDiff {
					if len(mockQueue.AddPassedDiffs) > 0 {
						return mockQueue.AddPassedDiffs[csvDiff.ID]
					}
					return utils.StorageDiff{}
				}).Should(Equal(csvDiff))
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
			BeforeEach(func() {
				mockQueue.DiffsToReturn = map[int]utils.StorageDiff{
					csvDiff.ID: csvDiff,
				}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			It("executes transformer for storage diff", func(done Done) {
				go storageWatcher.Execute(time.Nanosecond, false)

				Eventually(func() utils.StorageDiff {
					if len(mockTransformer.PassedDiffs) > 0 {
						return mockTransformer.PassedDiffs[csvDiff.ID]
					}
					return utils.StorageDiff{}
				}).Should(Equal(csvDiff))
				close(done)
			})

			It("deletes diff from queue if transformer execution successful", func(done Done) {
				go storageWatcher.Execute(time.Nanosecond, false)

				Eventually(func() int {
					if len(mockQueue.DeletePassedIds) > 0 {
						return mockQueue.DeletePassedIds[0]
					}
					return 0
				}).Should(Equal(csvDiff.ID))
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
				obsoleteDiff := utils.StorageDiff{
					ID:            csvDiff.ID + 1,
					HashedAddress: utils.HexToKeccak256Hash("0xfedcba9876543210"),
				}
				mockQueue.DiffsToReturn = map[int]utils.StorageDiff{
					obsoleteDiff.ID: obsoleteDiff,
				}

				go storageWatcher.Execute(time.Nanosecond, false)

				Eventually(func() int {
					if len(mockQueue.DeletePassedIds) > 0 {
						return mockQueue.DeletePassedIds[0]
					}
					return 0
				}).Should(Equal(obsoleteDiff.ID))
				close(done)
			})

			It("logs error if deleting obsolete diff fails", func(done Done) {
				obsoleteDiff := utils.StorageDiff{
					ID:            csvDiff.ID + 1,
					HashedAddress: utils.HexToKeccak256Hash("0xfedcba9876543210"),
				}
				mockQueue.DiffsToReturn = map[int]utils.StorageDiff{
					obsoleteDiff.ID: obsoleteDiff,
				}
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
			mockFetcher                                          *mocks.StorageFetcher
			mockBackFiller                                       *mocks.BackFiller
			mockQueue                                            *mocks.MockStorageQueue
			mockTransformer                                      *mocks.MockStorageTransformer
			mockTransformer2                                     *mocks.MockStorageTransformer
			mockTransformer3                                     *mocks.MockStorageTransformer
			csvDiff                                              utils.StorageDiff
			storageWatcher                                       *watcher.StorageWatcher
			hashedAddress                                        common.Hash
			createdDiff, updatedDiff1, deletedDiff, updatedDiff2 utils.StorageDiff
		)

		BeforeEach(func() {
			createdDiff = test_data.CreatedExpectedStorageDiff
			createdDiff.ID = 1333
			updatedDiff1 = test_data.UpdatedExpectedStorageDiff
			updatedDiff1.ID = 1334
			deletedDiff = test_data.DeletedExpectedStorageDiff
			deletedDiff.ID = 1335
			updatedDiff2 = test_data.UpdatedExpectedStorageDiff2
			updatedDiff2.ID = 1336
			mockBackFiller = new(mocks.BackFiller)
			hashedAddress = utils.HexToKeccak256Hash("0x0123456789abcdef")
			mockFetcher = mocks.NewStorageFetcher()
			mockQueue = &mocks.MockStorageQueue{}
			mockTransformer = &mocks.MockStorageTransformer{KeccakOfAddress: hashedAddress}
			mockTransformer2 = &mocks.MockStorageTransformer{KeccakOfAddress: common.BytesToHash(test_data.ContractLeafKey[:])}
			mockTransformer3 = &mocks.MockStorageTransformer{KeccakOfAddress: common.BytesToHash(test_data.AnotherContractLeafKey[:])}
			csvDiff = utils.StorageDiff{
				ID:            1337,
				HashedAddress: hashedAddress,
				BlockHash:     common.HexToHash("0xfedcba9876543210"),
				BlockHeight:   int(test_data.BlockNumber2.Int64()) + 1,
				StorageKey:    common.HexToHash("0xabcdef1234567890"),
				StorageValue:  common.HexToHash("0x9876543210abcdef"),
			}
		})

		Describe("transforming streamed and backfilled storage diffs", func() {
			BeforeEach(func() {
				mockFetcher.DiffsToReturn = []utils.StorageDiff{csvDiff}
				mockBackFiller.SetStorageDiffsToReturn([]utils.StorageDiff{
					createdDiff,
					updatedDiff1,
					deletedDiff,
					updatedDiff2,
				})
				mockQueue.DiffsToReturn = map[int]utils.StorageDiff{}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{
					mockTransformer.FakeTransformerInitializer,
					mockTransformer2.FakeTransformerInitializer,
					mockTransformer3.FakeTransformerInitializer,
				})
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
				Expect(mockTransformer.PassedDiffs[csvDiff.ID]).To(Equal(csvDiff))
				Expect(mockTransformer2.PassedDiffs[createdDiff.ID]).To(Equal(createdDiff))
				Expect(mockTransformer3.PassedDiffs[updatedDiff1.ID]).To(Equal(updatedDiff1))
				Expect(mockTransformer3.PassedDiffs[deletedDiff.ID]).To(Equal(deletedDiff))
				Expect(mockTransformer3.PassedDiffs[updatedDiff2.ID]).To(Equal(updatedDiff2))
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
				Eventually(func() map[int]utils.StorageDiff {
					if len(mockQueue.AddPassedDiffs) > 2 {
						return mockQueue.AddPassedDiffs
					}
					return map[int]utils.StorageDiff{}
				}).Should(Equal(map[int]utils.StorageDiff{
					updatedDiff1.ID: updatedDiff1,
					deletedDiff.ID:  deletedDiff,
					updatedDiff2.ID: updatedDiff2,
				}))

				Expect(mockBackFiller.PassedEndingBlock).To(Equal(uint64(test_data.BlockNumber2.Int64())))
				Expect(mockTransformer.PassedDiffs[csvDiff.ID]).To(Equal(csvDiff))
				Expect(mockTransformer2.PassedDiffs[createdDiff.ID]).To(Equal(createdDiff))
				Expect(mockTransformer3.PassedDiffs[updatedDiff1.ID]).To(Equal(updatedDiff1))
				Expect(mockTransformer3.PassedDiffs[deletedDiff.ID]).To(Equal(deletedDiff))
				Expect(mockTransformer3.PassedDiffs[updatedDiff2.ID]).To(Equal(updatedDiff2))
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
				mockQueue.DiffsToReturn = map[int]utils.StorageDiff{
					csvDiff.ID:      csvDiff,
					createdDiff.ID:  createdDiff,
					updatedDiff1.ID: updatedDiff1,
					deletedDiff.ID:  deletedDiff,
					updatedDiff2.ID: updatedDiff2,
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
				sortedExpectedIDs := []int{
					csvDiff.ID,
					createdDiff.ID,
					updatedDiff1.ID,
					deletedDiff.ID,
					updatedDiff2.ID,
				}
				sort.Ints(sortedExpectedIDs)
				Eventually(func() []int {
					if len(mockQueue.DeletePassedIds) > 4 {
						sort.Ints(mockQueue.DeletePassedIds)
						return mockQueue.DeletePassedIds
					}
					return []int{}
				}).Should(Equal(sortedExpectedIDs))

				Expect(mockQueue.AddCalled).To(Not(BeTrue()))
				Expect(len(mockQueue.DiffsToReturn)).To(Equal(0))
				Expect(mockTransformer.PassedDiffs[csvDiff.ID]).To(Equal(csvDiff))
				Expect(mockTransformer2.PassedDiffs[createdDiff.ID]).To(Equal(createdDiff))
				Expect(mockTransformer3.PassedDiffs[updatedDiff1.ID]).To(Equal(updatedDiff1))
				Expect(mockTransformer3.PassedDiffs[deletedDiff.ID]).To(Equal(deletedDiff))
				Expect(mockTransformer3.PassedDiffs[updatedDiff2.ID]).To(Equal(updatedDiff2))
				close(done)
			})
		})
	})
})
