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
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/libraries/shared/mocks"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/libraries/shared/watcher"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Storage Watcher", func() {
	Describe("AddTransformer", func() {
		It("adds transformers", func() {
			fakeHashedAddress := utils.HexToKeccak256Hash("0x12345")
			fakeTransformer := &mocks.MockStorageTransformer{KeccakOfAddress: fakeHashedAddress}
			w := watcher.NewStorageWatcher(mocks.NewMockStorageFetcher(), test_config.NewTestDB(test_config.NewTestNode()))

			w.AddTransformers([]transformer.StorageTransformerInitializer{fakeTransformer.FakeTransformerInitializer})

			Expect(w.KeccakAddressTransformers[fakeHashedAddress]).To(Equal(fakeTransformer))
		})
	})

	Describe("Execute", func() {
		var (
			mockFetcher          *mocks.MockStorageFetcher
			mockHeaderRepository *fakes.MockHeaderRepository
			mockQueue            *mocks.MockStorageQueue
			mockTransformer      *mocks.MockStorageTransformer
			csvDiff              utils.StorageDiff
			storageWatcher       watcher.StorageWatcher
			hashedAddress        common.Hash
		)

		BeforeEach(func() {
			hashedAddress = utils.HexToKeccak256Hash("0x0123456789abcdef")
			mockFetcher = mocks.NewMockStorageFetcher()
			mockHeaderRepository = fakes.NewMockHeaderRepository()
			fakeHeaderID := rand.Int63()
			mockHeaderRepository.GetHeaderReturnID = fakeHeaderID
			mockHeaderRepository.GetHeaderReturnHash = fakes.FakeHash.Hex()
			mockQueue = &mocks.MockStorageQueue{}
			mockTransformer = &mocks.MockStorageTransformer{KeccakOfAddress: hashedAddress}
			csvDiff = utils.StorageDiff{
				Id:            1337,
				HashedAddress: hashedAddress,
				BlockHash:     fakes.FakeHash,
				BlockHeight:   0,
				StorageKey:    common.HexToHash("0xabcdef1234567890"),
				StorageValue:  common.HexToHash("0x9876543210abcdef"),
				HeaderID:      fakeHeaderID,
			}
		})

		It("logs error if fetching storage diffs fails", func() {
			mockFetcher.ErrsToReturn = []error{fakes.FakeError}
			storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
			storageWatcher.HeaderRepository = mockHeaderRepository
			storageWatcher.Queue = mockQueue
			storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			tempFile, fileErr := ioutil.TempFile("", "log")
			Expect(fileErr).NotTo(HaveOccurred())
			defer os.Remove(tempFile.Name())
			logrus.SetOutput(tempFile)

			err := storageWatcher.Execute(time.Hour)

			Expect(err).To(MatchError(fakes.FakeError))
			logContent, readErr := ioutil.ReadFile(tempFile.Name())
			Expect(readErr).NotTo(HaveOccurred())
			Expect(string(logContent)).To(ContainSubstring(fakes.FakeError.Error()))
		})

		Describe("transforming new storage diffs from csv", func() {
			BeforeEach(func() {
				mockFetcher.DiffsToReturn = []utils.StorageDiff{csvDiff}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.HeaderRepository = mockHeaderRepository
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			Describe("when getting header succeeds", func() {
				It("executes transformer for recognized storage diff with matching header", func(done Done) {
					go func() {
						err := storageWatcher.Execute(time.Hour)
						Expect(err).NotTo(HaveOccurred())
					}()

					Eventually(func() utils.StorageDiff {
						return mockTransformer.PassedDiff
					}).Should(Equal(csvDiff))
					close(done)
				})

				It("recognizes header match even if hash hex missing 0x prefix", func(done Done) {
					mockHeaderRepository.GetHeaderReturnHash = fakes.FakeHash.Hex()[2:]

					go func() {
						err := storageWatcher.Execute(time.Hour)
						Expect(err).NotTo(HaveOccurred())
					}()

					Eventually(func() utils.StorageDiff {
						return mockTransformer.PassedDiff
					}).Should(Equal(csvDiff))
					close(done)
				})

				Describe("when executing transformer fails", func() {
					It("logs error", func(done Done) {
						mockTransformer.ExecuteErr = fakes.FakeError
						tempFile, fileErr := ioutil.TempFile("", "log")
						Expect(fileErr).NotTo(HaveOccurred())
						defer os.Remove(tempFile.Name())
						logrus.SetOutput(tempFile)

						go func() {
							err := storageWatcher.Execute(time.Hour)
							Expect(err).NotTo(HaveOccurred())
						}()

						Eventually(func() (string, error) {
							logContent, readErr := ioutil.ReadFile(tempFile.Name())
							return string(logContent), readErr
						}).Should(ContainSubstring(fakes.FakeError.Error()))
						close(done)
					})

					It("queues diff", func(done Done) {
						mockTransformer.ExecuteErr = fakes.FakeError

						go func() {
							err := storageWatcher.Execute(time.Hour)
							Expect(err).NotTo(HaveOccurred())
						}()

						Eventually(func() utils.StorageDiff {
							return mockQueue.AddPassedDiff
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

						go func() {
							err := storageWatcher.Execute(time.Hour)
							Expect(err).NotTo(HaveOccurred())
						}()

						Eventually(func() (string, error) {
							logContent, readErr := ioutil.ReadFile(tempFile.Name())
							return string(logContent), readErr
						}).Should(ContainSubstring(fakes.FakeError.Error()))
						close(done)
					})
				})
			})

			Describe("when getting header fails", func() {
				Describe("when repository returns error", func() {
					It("queues diff ", func(done Done) {
						mockHeaderRepository.GetHeaderError = fakes.FakeError

						go func() {
							err := storageWatcher.Execute(time.Hour)
							Expect(err).NotTo(HaveOccurred())
						}()

						Eventually(func() bool {
							return mockQueue.AddCalled
						}).Should(BeTrue())
						close(done)
					})

					It("logs error", func(done Done) {
						mockHeaderRepository.GetHeaderError = fakes.FakeError
						tempFile, fileErr := ioutil.TempFile("", "log")
						Expect(fileErr).NotTo(HaveOccurred())
						defer os.Remove(tempFile.Name())
						logrus.SetOutput(tempFile)

						go func() {
							err := storageWatcher.Execute(time.Hour)
							Expect(err).NotTo(HaveOccurred())
						}()

						Eventually(func() (string, error) {
							logContent, readErr := ioutil.ReadFile(tempFile.Name())
							return string(logContent), readErr
						}).Should(ContainSubstring(fakes.FakeError.Error()))
						close(done)
					})
				})

				Describe("when hash doesn't match", func() {
					var (
						wrongHash     = fakes.RandomString(64)
						expectedError = watcher.NewErrHeaderMismatch(wrongHash, fakes.FakeHash.Hex())
					)

					BeforeEach(func() {
						mockHeaderRepository.GetHeaderReturnHash = wrongHash
					})

					It("queues diff", func(done Done) {
						go func() {
							err := storageWatcher.Execute(time.Hour)
							Expect(err).NotTo(HaveOccurred())
						}()

						Eventually(func() bool {
							return mockQueue.AddCalled
						}).Should(BeTrue())
						close(done)
					})

					It("logs error", func(done Done) {
						tempFile, fileErr := ioutil.TempFile("", "log")
						Expect(fileErr).NotTo(HaveOccurred())
						defer os.Remove(tempFile.Name())
						logrus.SetOutput(tempFile)

						go func() {
							err := storageWatcher.Execute(time.Hour)
							Expect(err).NotTo(HaveOccurred())
						}()

						Eventually(func() (string, error) {
							logContent, readErr := ioutil.ReadFile(tempFile.Name())
							return string(logContent), readErr
						}).Should(ContainSubstring(expectedError.Error()))
						close(done)
					})
				})
			})
		})

		Describe("transforming queued storage diffs", func() {
			BeforeEach(func() {
				mockQueue.DiffsToReturn = []utils.StorageDiff{csvDiff}
				storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))
				storageWatcher.HeaderRepository = mockHeaderRepository
				storageWatcher.Queue = mockQueue
				storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})
			})

			Describe("when contract recognized", func() {
				Describe("when getting header succeeds", func() {
					It("executes transformer for storage diff", func(done Done) {
						go func() {
							err := storageWatcher.Execute(time.Nanosecond)
							Expect(err).NotTo(HaveOccurred())
						}()

						Eventually(func() utils.StorageDiff {
							return mockTransformer.PassedDiff
						}).Should(Equal(csvDiff))
						close(done)
					})

					Describe("when transformer execution successful", func() {
						It("deletes diff from queue", func(done Done) {
							go func() {
								err := storageWatcher.Execute(time.Nanosecond)
								Expect(err).NotTo(HaveOccurred())
							}()

							Eventually(func() int {
								return mockQueue.DeletePassedId
							}).Should(Equal(csvDiff.Id))
							close(done)
						})

						It("logs error if deleting queued diff fails", func(done Done) {
							mockQueue.DeleteErr = fakes.FakeError
							tempFile, fileErr := ioutil.TempFile("", "log")
							Expect(fileErr).NotTo(HaveOccurred())
							defer os.Remove(tempFile.Name())
							logrus.SetOutput(tempFile)

							go func() {
								err := storageWatcher.Execute(time.Nanosecond)
								Expect(err).NotTo(HaveOccurred())
							}()

							Eventually(func() (string, error) {
								logContent, readErr := ioutil.ReadFile(tempFile.Name())
								return string(logContent), readErr
							}).Should(ContainSubstring(fakes.FakeError.Error()))
							close(done)
						})
					})

					Describe("when transformer execution fails", func() {
						BeforeEach(func() {
							mockTransformer.ExecuteErr = fakes.FakeError
						})

						It("logs error", func(done Done) {
							tempFile, fileErr := ioutil.TempFile("", "log")
							Expect(fileErr).NotTo(HaveOccurred())
							defer os.Remove(tempFile.Name())
							logrus.SetOutput(tempFile)

							go func() {
								err := storageWatcher.Execute(time.Nanosecond)
								Expect(err).NotTo(HaveOccurred())
							}()

							Eventually(func() (string, error) {
								logContent, readErr := ioutil.ReadFile(tempFile.Name())
								return string(logContent), readErr
							}).Should(ContainSubstring(fakes.FakeError.Error()))
							close(done)
						})

						It("does not delete diff from queue", func(done Done) {
							go func() {
								err := storageWatcher.Execute(time.Nanosecond)
								Expect(err).NotTo(HaveOccurred())
							}()

							Consistently(func() bool {
								return mockQueue.DeleteCalled
							}).Should(BeFalse())
							close(done)
						})
					})
				})

				Describe("when getting header fails", func() {
					It("logs error if repository returns error", func(done Done) {
						mockHeaderRepository.GetHeaderError = fakes.FakeError
						tempFile, fileErr := ioutil.TempFile("", "log")
						Expect(fileErr).NotTo(HaveOccurred())
						defer os.Remove(tempFile.Name())
						logrus.SetOutput(tempFile)

						go func() {
							err := storageWatcher.Execute(time.Nanosecond)
							Expect(err).NotTo(HaveOccurred())
						}()

						Eventually(func() (string, error) {
							logContent, readErr := ioutil.ReadFile(tempFile.Name())
							return string(logContent), readErr
						}).Should(ContainSubstring(fakes.FakeError.Error()))
						close(done)
					})

					It("logs error if header hash doesn't match diff", func(done Done) {
						wrongHash := fakes.RandomString(64)
						mockHeaderRepository.GetHeaderReturnHash = wrongHash
						tempFile, fileErr := ioutil.TempFile("", "log")
						Expect(fileErr).NotTo(HaveOccurred())
						defer os.Remove(tempFile.Name())
						logrus.SetOutput(tempFile)

						go func() {
							err := storageWatcher.Execute(time.Nanosecond)
							Expect(err).NotTo(HaveOccurred())
						}()

						expectedError := watcher.NewErrHeaderMismatch(wrongHash, fakes.FakeHash.Hex())
						Eventually(func() (string, error) {
							logContent, readErr := ioutil.ReadFile(tempFile.Name())
							return string(logContent), readErr
						}).Should(ContainSubstring(expectedError.Error()))
						close(done)
					})

					It("does not delete diff from queue", func(done Done) {
						mockHeaderRepository.GetHeaderError = fakes.FakeError
						go func() {
							err := storageWatcher.Execute(time.Nanosecond)
							Expect(err).NotTo(HaveOccurred())
						}()

						Consistently(func() bool {
							return mockQueue.DeleteCalled
						}).Should(BeFalse())
						close(done)
					})
				})
			})

			Describe("when contract not recognized", func() {
				It("deletes obsolete diff from queue", func(done Done) {
					obsoleteDiff := utils.StorageDiff{
						Id:            csvDiff.Id + 1,
						HashedAddress: utils.HexToKeccak256Hash("0xfedcba9876543210"),
						BlockHash:     fakes.FakeHash,
					}
					mockQueue.DiffsToReturn = []utils.StorageDiff{obsoleteDiff}

					go func() {
						err := storageWatcher.Execute(time.Nanosecond)
						Expect(err).NotTo(HaveOccurred())
					}()

					Eventually(func() int {
						return mockQueue.DeletePassedId
					}).Should(Equal(obsoleteDiff.Id))
					close(done)
				})

				It("logs error if deleting obsolete diff fails", func(done Done) {
					obsoleteDiff := utils.StorageDiff{
						Id:            csvDiff.Id + 1,
						HashedAddress: utils.HexToKeccak256Hash("0xfedcba9876543210"),
						BlockHash:     fakes.FakeHash,
					}
					mockQueue.DiffsToReturn = []utils.StorageDiff{obsoleteDiff}
					mockQueue.DeleteErr = fakes.FakeError
					tempFile, fileErr := ioutil.TempFile("", "log")
					Expect(fileErr).NotTo(HaveOccurred())
					defer os.Remove(tempFile.Name())
					logrus.SetOutput(tempFile)

					go func() {
						err := storageWatcher.Execute(time.Nanosecond)
						Expect(err).NotTo(HaveOccurred())
					}()

					Eventually(func() (string, error) {
						logContent, readErr := ioutil.ReadFile(tempFile.Name())
						return string(logContent), readErr
					}).Should(ContainSubstring(fakes.FakeError.Error()))
					close(done)
				})
			})
		})
	})
})
