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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Storage Watcher", func() {
	It("adds transformers", func() {
		fakeAddress := common.HexToAddress("0x12345")
		fakeTransformer := &mocks.MockStorageTransformer{Address: fakeAddress}
		w := watcher.NewStorageWatcher(mocks.NewMockStorageFetcher(), test_config.NewTestDB(test_config.NewTestNode()))

		w.AddTransformers([]transformer.StorageTransformerInitializer{fakeTransformer.FakeTransformerInitializer})

		Expect(w.Transformers[fakeAddress]).To(Equal(fakeTransformer))
	})

	Describe("executing watcher", func() {
		var (
			errs            chan error
			mockFetcher     *mocks.MockStorageFetcher
			mockQueue       *mocks.MockStorageQueue
			mockTransformer *mocks.MockStorageTransformer
			row             utils.StorageDiffRow
			rows            chan utils.StorageDiffRow
			storageWatcher  watcher.StorageWatcher
		)

		BeforeEach(func() {
			errs = make(chan error)
			rows = make(chan utils.StorageDiffRow)
			address := common.HexToAddress("0x0123456789abcdef")
			mockFetcher = mocks.NewMockStorageFetcher()
			mockQueue = &mocks.MockStorageQueue{}
			mockTransformer = &mocks.MockStorageTransformer{Address: address}
			row = utils.StorageDiffRow{
				Id:           1337,
				Contract:     address,
				BlockHash:    common.HexToHash("0xfedcba9876543210"),
				BlockHeight:  0,
				StorageKey:   common.HexToHash("0xabcdef1234567890"),
				StorageValue: common.HexToHash("0x9876543210abcdef"),
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

			go storageWatcher.Execute(rows, errs, time.Hour)

			Eventually(func() (string, error) {
				logContent, err := ioutil.ReadFile(tempFile.Name())
				return string(logContent), err
			}).Should(ContainSubstring(fakes.FakeError.Error()))
			close(done)
		})

		Describe("transforming new storage diffs", func() {
			BeforeEach(func() {
				mockFetcher.RowsToReturn = []utils.StorageDiffRow{row}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			It("executes transformer for recognized storage row", func(done Done) {
				go storageWatcher.Execute(rows, errs, time.Hour)

				Eventually(func() utils.StorageDiffRow {
					return mockTransformer.PassedRow
				}).Should(Equal(row))
				close(done)
			})

			It("queues row for later processing if transformer execution fails", func(done Done) {
				mockTransformer.ExecuteErr = fakes.FakeError

				go storageWatcher.Execute(rows, errs, time.Hour)

				Expect(<-errs).To(BeNil())
				Eventually(func() bool {
					return mockQueue.AddCalled
				}).Should(BeTrue())
				Eventually(func() utils.StorageDiffRow {
					return mockQueue.AddPassedRow
				}).Should(Equal(row))
				close(done)
			})

			It("logs error if queueing row fails", func(done Done) {
				mockTransformer.ExecuteErr = utils.ErrStorageKeyNotFound{}
				mockQueue.AddError = fakes.FakeError
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.Execute(rows, errs, time.Hour)

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
				mockQueue.RowsToReturn = []utils.StorageDiffRow{row}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			It("logs error if getting queued storage fails", func(done Done) {
				mockQueue.GetAllErr = fakes.FakeError
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.Execute(rows, errs, time.Nanosecond)

				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring(fakes.FakeError.Error()))
				close(done)
			})

			It("executes transformer for storage row", func(done Done) {
				go storageWatcher.Execute(rows, errs, time.Nanosecond)

				Eventually(func() utils.StorageDiffRow {
					return mockTransformer.PassedRow
				}).Should(Equal(row))
				close(done)
			})

			It("deletes row from queue if transformer execution successful", func(done Done) {
				go storageWatcher.Execute(rows, errs, time.Nanosecond)

				Eventually(func() int {
					return mockQueue.DeletePassedId
				}).Should(Equal(row.Id))
				close(done)
			})

			It("logs error if deleting persisted row fails", func(done Done) {
				mockQueue.DeleteErr = fakes.FakeError
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.Execute(rows, errs, time.Nanosecond)

				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring(fakes.FakeError.Error()))
				close(done)
			})

			It("deletes obsolete row from queue if contract not recognized", func(done Done) {
				obsoleteRow := utils.StorageDiffRow{
					Id:       row.Id + 1,
					Contract: common.HexToAddress("0xfedcba9876543210"),
				}
				mockQueue.RowsToReturn = []utils.StorageDiffRow{obsoleteRow}

				go storageWatcher.Execute(rows, errs, time.Nanosecond)

				Eventually(func() int {
					return mockQueue.DeletePassedId
				}).Should(Equal(obsoleteRow.Id))
				close(done)
			})

			It("logs error if deleting obsolete row fails", func(done Done) {
				obsoleteRow := utils.StorageDiffRow{
					Id:       row.Id + 1,
					Contract: common.HexToAddress("0xfedcba9876543210"),
				}
				mockQueue.RowsToReturn = []utils.StorageDiffRow{obsoleteRow}
				mockQueue.DeleteErr = fakes.FakeError
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.Execute(rows, errs, time.Nanosecond)

				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring(fakes.FakeError.Error()))
				close(done)
			})
		})

	})
})
