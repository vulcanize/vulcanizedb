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
	"io/ioutil"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
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
			w := watcher.NewStorageWatcher(mocks.NewClosingStorageFetcher(), test_config.NewTestDB(test_config.NewTestNode()))

			w.AddTransformers([]transformer.StorageTransformerInitializer{fakeTransformer.FakeTransformerInitializer})

			Expect(w.KeccakAddressTransformers[fakeHashedAddress]).To(Equal(fakeTransformer))
		})
	})

	Describe("Execute", func() {
		var (
			mockFetcher     *mocks.ClosingStorageFetcher
			mockQueue       *mocks.MockStorageQueue
			mockTransformer *mocks.MockStorageTransformer
			csvDiff         utils.StorageDiff
			storageWatcher  watcher.StorageWatcher
			hashedAddress   common.Hash
		)

		BeforeEach(func() {
			hashedAddress = utils.HexToKeccak256Hash("0x0123456789abcdef")
			mockFetcher = mocks.NewClosingStorageFetcher()
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

				Eventually(func() []utils.StorageDiff {
					return mockTransformer.PassedDiffs
				}).Should(Equal([]utils.StorageDiff{csvDiff}))
				close(done)
			})

			It("queues diff for later processing if transformer execution fails", func(done Done) {
				mockTransformer.ExecuteErr = fakes.FakeError

				go storageWatcher.Execute(time.Hour, false)

				Expect(<-storageWatcher.ErrsChan).To(BeNil())
				Eventually(func() bool {
					return mockQueue.AddCalled
				}).Should(BeTrue())
				Eventually(func() utils.StorageDiff {
					if len(mockQueue.AddPassedDiffs) > 0 {
						return mockQueue.AddPassedDiffs[0]
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
				mockQueue.DiffsToReturn = []utils.StorageDiff{csvDiff}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			It("executes transformer for storage diff", func(done Done) {
				go storageWatcher.Execute(time.Nanosecond, false)

				Eventually(func() utils.StorageDiff {
					if len(mockTransformer.PassedDiffs) > 0 {
						return mockTransformer.PassedDiffs[0]
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
				mockQueue.DiffsToReturn = []utils.StorageDiff{obsoleteDiff}

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
				mockQueue.DiffsToReturn = []utils.StorageDiff{obsoleteDiff}
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
			mockFetcher     *mocks.StorageFetcher
			mockBackFiller  *mocks.BackFiller
			mockQueue       *mocks.MockStorageQueue
			mockTransformer *mocks.MockStorageTransformer
			csvDiff         utils.StorageDiff
			storageWatcher  watcher.StorageWatcher
			hashedAddress   common.Hash
		)

		BeforeEach(func() {
			mockBackFiller = new(mocks.BackFiller)
			mockBackFiller.SetStorageDiffsToReturn([]utils.StorageDiff{
				test_data.CreatedExpectedStorageDiff,
				test_data.UpdatedExpectedStorageDiff,
				test_data.UpdatedExpectedStorageDiff2,
				test_data.DeletedExpectedStorageDiff})
			hashedAddress = utils.HexToKeccak256Hash("0x0123456789abcdef")
			mockFetcher = mocks.NewStorageFetcher()
			mockQueue = &mocks.MockStorageQueue{}
			mockTransformer = &mocks.MockStorageTransformer{KeccakOfAddress: hashedAddress}
			csvDiff = utils.StorageDiff{
				ID:            1337,
				HashedAddress: hashedAddress,
				BlockHash:     common.HexToHash("0xfedcba9876543210"),
				BlockHeight:   int(test_data.BlockNumber2.Int64()) + 1,
				StorageKey:    common.HexToHash("0xabcdef1234567890"),
				StorageValue:  common.HexToHash("0x9876543210abcdef"),
			}
		})

		Describe("transforming streamed and backfilled queued storage diffs", func() {
			BeforeEach(func() {
				mockFetcher.DiffsToReturn = []utils.StorageDiff{csvDiff}
				mockQueue.DiffsToReturn = []utils.StorageDiff{csvDiff,
					test_data.CreatedExpectedStorageDiff,
					test_data.UpdatedExpectedStorageDiff,
					test_data.UpdatedExpectedStorageDiff2,
					test_data.DeletedExpectedStorageDiff}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			It("executes transformer for storage diffs", func(done Done) {
				go storageWatcher.BackFill(mockBackFiller, int(test_data.BlockNumber.Uint64()))
				go storageWatcher.Execute(time.Nanosecond, true)
				expectedDiffsStruct := struct {
					diffs []utils.StorageDiff
				}{
					[]utils.StorageDiff{
						csvDiff,
						test_data.CreatedExpectedStorageDiff,
						test_data.UpdatedExpectedStorageDiff,
						test_data.UpdatedExpectedStorageDiff2,
						test_data.DeletedExpectedStorageDiff,
					},
				}
				expectedDiffsBytes, rlpErr1 := rlp.EncodeToBytes(expectedDiffsStruct)
				Expect(rlpErr1).ToNot(HaveOccurred())
				Eventually(func() []byte {
					diffsStruct := struct {
						diffs []utils.StorageDiff
					}{
						mockTransformer.PassedDiffs,
					}
					diffsBytes, rlpErr2 := rlp.EncodeToBytes(diffsStruct)
					Expect(rlpErr2).ToNot(HaveOccurred())
					return diffsBytes
				}).Should(Equal(expectedDiffsBytes))
				close(done)
			})

			It("deletes diffs from queue if transformer execution is successful", func(done Done) {
				go storageWatcher.Execute(time.Nanosecond, false)
				expectedIdsStruct := struct {
					diffs []int
				}{
					[]int{
						csvDiff.ID,
						test_data.CreatedExpectedStorageDiff.ID,
						test_data.UpdatedExpectedStorageDiff.ID,
						test_data.UpdatedExpectedStorageDiff2.ID,
						test_data.DeletedExpectedStorageDiff.ID,
					},
				}
				expectedIdsBytes, rlpErr1 := rlp.EncodeToBytes(expectedIdsStruct)
				Expect(rlpErr1).ToNot(HaveOccurred())
				Eventually(func() []byte {
					idsStruct := struct {
						diffs []int
					}{
						mockQueue.DeletePassedIds,
					}
					idsBytes, rlpErr2 := rlp.EncodeToBytes(idsStruct)
					Expect(rlpErr2).ToNot(HaveOccurred())
					return idsBytes
				}).Should(Equal(expectedIdsBytes))
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
		})
	})
})
