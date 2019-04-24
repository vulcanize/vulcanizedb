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
			mockFetcher     *mocks.MockStorageFetcher
			mockQueue       *mocks.MockStorageQueue
			mockTransformer *mocks.MockStorageTransformer
			row             utils.StorageDiffRow
			storageWatcher  watcher.StorageWatcher
		)

		BeforeEach(func() {
			address := common.HexToAddress("0x0123456789abcdef")
			mockFetcher = mocks.NewMockStorageFetcher()
			mockQueue = &mocks.MockStorageQueue{}
			mockTransformer = &mocks.MockStorageTransformer{Address: address}
			row = utils.StorageDiffRow{
				Contract:     address,
				BlockHash:    common.HexToHash("0xfedcba9876543210"),
				BlockHeight:  0,
				StorageKey:   common.HexToHash("0xabcdef1234567890"),
				StorageValue: common.HexToHash("0x9876543210abcdef"),
			}
			mockFetcher.RowsToReturn = []utils.StorageDiffRow{row}
			storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
			storageWatcher.Queue = mockQueue
			storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
		})

		It("executes transformer for recognized storage row", func() {
			err := storageWatcher.Execute()

			Expect(err).NotTo(HaveOccurred())
			Expect(mockTransformer.PassedRow).To(Equal(row))
		})

		It("queues row for later processing if row's key not recognized", func() {
			mockTransformer.ExecuteErr = utils.ErrStorageKeyNotFound{}

			err := storageWatcher.Execute()

			Expect(err).NotTo(HaveOccurred())
			Expect(mockQueue.AddCalled).To(BeTrue())
			Expect(mockQueue.PassedRow).To(Equal(row))
		})

		It("logs error if queueing row fails", func() {
			mockTransformer.ExecuteErr = utils.ErrStorageKeyNotFound{}
			mockQueue.AddError = fakes.FakeError
			tempFile, fileErr := ioutil.TempFile("", "log")
			Expect(fileErr).NotTo(HaveOccurred())
			defer os.Remove(tempFile.Name())
			logrus.SetOutput(tempFile)

			err := storageWatcher.Execute()

			Expect(err).NotTo(HaveOccurred())
			Expect(mockQueue.AddCalled).To(BeTrue())
			logContent, readErr := ioutil.ReadFile(tempFile.Name())
			Expect(readErr).NotTo(HaveOccurred())
			Expect(string(logContent)).To(ContainSubstring(fakes.FakeError.Error()))
		})

		It("logs error if transformer execution fails for reason other than key not found", func() {
			mockTransformer.ExecuteErr = fakes.FakeError
			tempFile, fileErr := ioutil.TempFile("", "log")
			Expect(fileErr).NotTo(HaveOccurred())
			defer os.Remove(tempFile.Name())
			logrus.SetOutput(tempFile)

			err := storageWatcher.Execute()

			Expect(err).NotTo(HaveOccurred())
			logContent, readErr := ioutil.ReadFile(tempFile.Name())
			Expect(readErr).NotTo(HaveOccurred())
			Expect(string(logContent)).To(ContainSubstring(fakes.FakeError.Error()))
		})
	})
})
