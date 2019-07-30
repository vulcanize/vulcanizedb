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
	"github.com/ethereum/go-ethereum/crypto"
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

var _ = Describe("Geth Storage Watcher", func() {
	It("adds transformers", func() {
		fakeAddress := common.HexToAddress("0x12345")
		fakeTransformer := &mocks.MockStorageTransformer{Address: fakeAddress}
		w := watcher.NewGethStorageWatcher(mocks.NewMockStorageFetcher(), test_config.NewTestDB(test_config.NewTestNode()))

		w.AddTransformers([]transformer.StorageTransformerInitializer{fakeTransformer.FakeTransformerInitializer})

		Expect(w.Transformers[fakeAddress]).To(Equal(fakeTransformer))
	})

	Describe("executing watcher", func() {
		var (
			errs            chan error
			mockFetcher     *mocks.MockStorageFetcher
			mockQueue       *mocks.MockStorageQueue
			mockTransformer *mocks.MockStorageTransformer
			gethDiff        utils.StorageDiff
			diffs           chan utils.StorageDiff
			storageWatcher  watcher.GethStorageWatcher
			address         common.Address
			keccakOfAddress common.Address
		)

		BeforeEach(func() {
			errs = make(chan error)
			diffs = make(chan utils.StorageDiff)
			address = common.HexToAddress("0x0123456789abcdef")
			keccakOfAddress = common.BytesToAddress(crypto.Keccak256(address[:]))
			mockFetcher = mocks.NewMockStorageFetcher()
			mockQueue = &mocks.MockStorageQueue{}
			mockTransformer = &mocks.MockStorageTransformer{Address: address}
			gethDiff = utils.StorageDiff{
				Id:           1338,
				Contract:     keccakOfAddress,
				BlockHash:    common.HexToHash("0xfedcba9876543210"),
				BlockHeight:  0,
				StorageKey:   common.HexToHash("0xabcdef1234567890"),
				StorageValue: common.HexToHash("0x9876543210abcdef"),
			}
		})

		It("logs error if fetching storage diffs fails", func(done Done) {
			mockFetcher.ErrsToReturn = []error{fakes.FakeError}
			storageWatcher = watcher.NewGethStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
			storageWatcher.Queue = mockQueue
			storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			tempFile, fileErr := ioutil.TempFile("", "log")
			Expect(fileErr).NotTo(HaveOccurred())
			defer os.Remove(tempFile.Name())
			logrus.SetOutput(tempFile)

			go storageWatcher.Execute(diffs, errs, time.Hour)

			Eventually(func() (string, error) {
				logContent, err := ioutil.ReadFile(tempFile.Name())
				return string(logContent), err
			}).Should(ContainSubstring(fakes.FakeError.Error()))
			close(done)
		})

		Describe("transforming new storage diffs", func() {
			BeforeEach(func() {
				mockFetcher.DiffsToReturn = []utils.StorageDiff{gethDiff}
				storageWatcher = watcher.NewGethStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			It("executes transformer for recognized storage diff", func(done Done) {
				go storageWatcher.Execute(diffs, errs, time.Hour)

				Eventually(func() utils.StorageDiff {
					return mockTransformer.PassedDiff
				}).Should(Equal(gethDiff))
				close(done)
			})

			It("queues diff for later processing if transformer execution fails", func(done Done) {
				mockTransformer.ExecuteErr = fakes.FakeError

				go storageWatcher.Execute(diffs, errs, time.Hour)

				Expect(<-errs).To(BeNil())
				Eventually(func() bool {
					return mockQueue.AddCalled
				}).Should(BeTrue())
				Eventually(func() utils.StorageDiff {
					return mockQueue.AddPassedDiff
				}).Should(Equal(gethDiff))
				close(done)
			})

			It("logs error if queueing diff fails", func(done Done) {
				mockTransformer.ExecuteErr = utils.ErrStorageKeyNotFound{}
				mockQueue.AddError = fakes.FakeError
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.Execute(diffs, errs, time.Hour)

				Eventually(func() bool {
					return mockQueue.AddCalled
				}).Should(BeTrue())
				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring(fakes.FakeError.Error()))
				close(done)
			})

			It("keeps track transformers by the keccak256 hash of their contract address ", func(done Done) {
				go storageWatcher.Execute(diffs, errs, time.Hour)

				m := make(map[common.Address]transformer.StorageTransformer)
				m[keccakOfAddress] = mockTransformer

				Eventually(func() map[common.Address]transformer.StorageTransformer {
					return storageWatcher.KeccakAddressTransformers
				}).Should(Equal(m))

				close(done)
			})

			It("gets the transformer from the known keccak address map first", func(done Done) {
				anotherAddress := common.HexToAddress("0xafakeaddress")
				anotherTransformer := &mocks.MockStorageTransformer{Address: anotherAddress}
				keccakOfAnotherAddress := common.BytesToAddress(crypto.Keccak256(anotherAddress[:]))

				anotherGethDiff := utils.StorageDiff{
					Id:           1338,
					Contract:     keccakOfAnotherAddress,
					BlockHash:    common.HexToHash("0xfedcba9876543210"),
					BlockHeight:  0,
					StorageKey:   common.HexToHash("0xabcdef1234567890"),
					StorageValue: common.HexToHash("0x9876543210abcdef"),
				}
				mockFetcher.DiffsToReturn = []utils.StorageDiff{anotherGethDiff}
				storageWatcher.KeccakAddressTransformers[keccakOfAnotherAddress] = anotherTransformer

				go storageWatcher.Execute(diffs, errs, time.Hour)

				Eventually(func() utils.StorageDiff {
					return anotherTransformer.PassedDiff
				}).Should(Equal(anotherGethDiff))

				close(done)
			})
		})

		Describe("transforming queued storage diffs", func() {
			BeforeEach(func() {
				mockQueue.DiffsToReturn = []utils.StorageDiff{gethDiff}
				storageWatcher = watcher.NewGethStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			It("executes transformer for storage diff", func(done Done) {
				go storageWatcher.Execute(diffs, errs, time.Nanosecond)

				Eventually(func() utils.StorageDiff {
					return mockTransformer.PassedDiff
				}).Should(Equal(gethDiff))
				close(done)
			})

			It("deletes diff from queue if transformer execution successful", func(done Done) {
				go storageWatcher.Execute(diffs, errs, time.Nanosecond)

				Eventually(func() int {
					return mockQueue.DeletePassedId
				}).Should(Equal(gethDiff.Id))
				close(done)
			})

			It("logs error if deleting persisted diff fails", func(done Done) {
				mockQueue.DeleteErr = fakes.FakeError
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.Execute(diffs, errs, time.Nanosecond)

				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring(fakes.FakeError.Error()))
				close(done)
			})

			It("deletes obsolete diff from queue if contract not recognized", func(done Done) {
				obsoleteDiff := utils.StorageDiff{
					Id:       gethDiff.Id + 1,
					Contract: common.HexToAddress("0xfedcba9876543210"),
				}
				mockQueue.DiffsToReturn = []utils.StorageDiff{obsoleteDiff}

				go storageWatcher.Execute(diffs, errs, time.Nanosecond)

				Eventually(func() int {
					return mockQueue.DeletePassedId
				}).Should(Equal(obsoleteDiff.Id))
				close(done)
			})

			It("logs error if deleting obsolete diff fails", func(done Done) {
				obsoleteDiff := utils.StorageDiff{
					Id:       gethDiff.Id + 1,
					Contract: common.HexToAddress("0xfedcba9876543210"),
				}
				mockQueue.DiffsToReturn = []utils.StorageDiff{obsoleteDiff}
				mockQueue.DeleteErr = fakes.FakeError
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.Execute(diffs, errs, time.Nanosecond)

				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring(fakes.FakeError.Error()))
				close(done)
			})
		})
	})
})
