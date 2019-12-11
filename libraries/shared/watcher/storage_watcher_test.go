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
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/libraries/shared/watcher"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Storage Watcher", func() {
	Describe("AddTransformer", func() {
		It("adds transformers", func() {
			fakeHashedAddress := storage.HexToKeccak256Hash("0x12345")
			fakeTransformer := &mocks.MockStorageTransformer{KeccakOfAddress: fakeHashedAddress}
			w := watcher.NewStorageWatcher(mocks.NewMockStorageFetcher(), test_config.NewTestDB(test_config.NewTestNode()))

			w.AddTransformers([]transformer.StorageTransformerInitializer{fakeTransformer.FakeTransformerInitializer})

			Expect(w.KeccakAddressTransformers[fakeHashedAddress]).To(Equal(fakeTransformer))
		})
	})

	Describe("Execute", func() {
		var (
			hashedAddress        common.Hash
			storageWatcher       watcher.StorageWatcher
			mockFetcher          *mocks.MockStorageFetcher
			mockQueue            *mocks.MockStorageQueue
			mockHeaderRepository *fakes.MockHeaderRepository
			mockTransformer      *mocks.MockStorageTransformer
		)

		BeforeEach(func() {
			mockFetcher = mocks.NewMockStorageFetcher()
			storageWatcher = watcher.NewStorageWatcher(mockFetcher, test_config.NewTestDB(test_config.NewTestNode()))

			hashedAddress = storage.HexToKeccak256Hash("0x0123456789abcdef")
			mockTransformer = &mocks.MockStorageTransformer{KeccakOfAddress: hashedAddress}
			storageWatcher.AddTransformers([]transformer.StorageTransformerInitializer{mockTransformer.FakeTransformerInitializer})

			mockHeaderRepository = &fakes.MockHeaderRepository{}
			storageWatcher.HeaderRepository = mockHeaderRepository

			mockQueue = &mocks.MockStorageQueue{}
			storageWatcher.Queue = mockQueue
		})

		It("logs error if fetching storage diffs fails", func() {
			mockFetcher.ErrsToReturn = []error{fakes.FakeError}
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

		Describe("transforming new raw storage diffs", func() {
			var (
				mockStorageDiffRepository *fakes.MockStorageDiffRepository
				fakeBlockHash             common.Hash
				fakeDiffID                int64
				fakeRawDiff               storage.RawStorageDiff
				fakePersistedDiff         storage.PersistedStorageDiff
			)

			BeforeEach(func() {
				fakeBlockHash = test_data.FakeHash()
				fakeRawDiff = storage.RawStorageDiff{
					HashedAddress: hashedAddress,
					BlockHash:     fakeBlockHash,
					BlockHeight:   0,
					StorageKey:    common.HexToHash("0xabcdef1234567890"),
					StorageValue:  common.HexToHash("0x9876543210abcdef"),
				}
				mockFetcher.DiffsToReturn = []storage.RawStorageDiff{fakeRawDiff}

				fakeHeaderID := rand.Int63()
				mockHeaderRepository.GetHeaderReturnID = fakeHeaderID
				mockHeaderRepository.GetHeaderReturnHash = fakeBlockHash.Hex()

				mockStorageDiffRepository = &fakes.MockStorageDiffRepository{}
				mockStorageDiffRepository.CreateReturnID = fakeDiffID
				storageWatcher.StorageDiffRepository = mockStorageDiffRepository

				fakePersistedDiff = storage.PersistedStorageDiff{
					RawStorageDiff: fakeRawDiff,
					ID:             fakeDiffID,
					HeaderID:       fakeHeaderID,
				}
			})

			It("writes raw diff before processing", func(done Done) {
				go storageWatcher.Execute(time.Hour)

				Eventually(func() []storage.RawStorageDiff {
					return mockStorageDiffRepository.CreatePassedRawDiffs
				}).Should(ContainElement(fakeRawDiff))
				close(done)
			})

			It("discards raw diff if it's already been persisted", func(done Done) {
				mockStorageDiffRepository.CreateReturnError = repositories.ErrDuplicateDiff

				go storageWatcher.Execute(time.Hour)

				Consistently(func() storage.PersistedStorageDiff {
					return mockTransformer.PassedDiff
				}).Should(BeZero())
				close(done)
			})

			It("logs error if persisting raw diff fails", func(done Done) {
				mockStorageDiffRepository.CreateReturnError = fakes.FakeError
				tempFile, fileErr := ioutil.TempFile("", "log")
				Expect(fileErr).NotTo(HaveOccurred())
				defer os.Remove(tempFile.Name())
				logrus.SetOutput(tempFile)

				go storageWatcher.Execute(time.Hour)

				Eventually(func() (string, error) {
					logContent, err := ioutil.ReadFile(tempFile.Name())
					return string(logContent), err
				}).Should(ContainSubstring(fakes.FakeError.Error()))
				close(done)
			})

			Describe("when getting header succeeds", func() {
				It("executes transformer for diff", func(done Done) {
					go func() {
						err := storageWatcher.Execute(time.Hour)
						Expect(err).NotTo(HaveOccurred())
					}()

					Eventually(func() storage.PersistedStorageDiff {
						return mockTransformer.PassedDiff
					}).Should(Equal(fakePersistedDiff))
					close(done)
				})

				It("recognizes header match even if hash hex missing 0x prefix", func(done Done) {
					mockHeaderRepository.GetHeaderReturnHash = fakeBlockHash.Hex()[2:]

					go func() {
						err := storageWatcher.Execute(time.Hour)
						Expect(err).NotTo(HaveOccurred())
					}()

					Eventually(func() storage.PersistedStorageDiff {
						return mockTransformer.PassedDiff
					}).Should(Equal(fakePersistedDiff))
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

						Eventually(func() storage.PersistedStorageDiff {
							return mockQueue.AddPassedDiff
						}).Should(Equal(fakePersistedDiff))
						close(done)
					})

					It("logs error if queueing diff fails", func(done Done) {
						mockTransformer.ExecuteErr = storage.ErrStorageKeyNotFound{}
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
						wrongHash     string
						expectedError error
					)

					BeforeEach(func() {
						wrongHash = fakes.RandomString(64)
						expectedError = watcher.NewErrHeaderMismatch(wrongHash, fakeBlockHash.Hex())
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
			var (
				fakeBlockHash common.Hash
				queuedDiff    storage.PersistedStorageDiff
			)

			BeforeEach(func() {
				fakeHeaderID := rand.Int63()
				fakeBlockHash = test_data.FakeHash()
				mockHeaderRepository.GetHeaderReturnID = fakeHeaderID
				mockHeaderRepository.GetHeaderReturnHash = fakeBlockHash.Hex()

				queuedDiff = storage.PersistedStorageDiff{
					ID:       rand.Int63(),
					HeaderID: fakeHeaderID,
					RawStorageDiff: storage.RawStorageDiff{
						HashedAddress: hashedAddress,
						BlockHash:     fakeBlockHash,
						BlockHeight:   rand.Int(),
						StorageKey:    test_data.FakeHash(),
						StorageValue:  test_data.FakeHash(),
					},
				}
				mockQueue.DiffsToReturn = []storage.PersistedStorageDiff{queuedDiff}
			})

			Describe("when contract recognized", func() {
				Describe("when getting header succeeds", func() {
					It("executes transformer for storage diff", func(done Done) {
						go func() {
							err := storageWatcher.Execute(time.Nanosecond)
							Expect(err).NotTo(HaveOccurred())
						}()

						Eventually(func() storage.PersistedStorageDiff {
							return mockTransformer.PassedDiff
						}).Should(Equal(queuedDiff))
						close(done)
					})

					Describe("when transformer execution successful", func() {
						It("deletes diff from queue", func(done Done) {
							go func() {
								err := storageWatcher.Execute(time.Nanosecond)
								Expect(err).NotTo(HaveOccurred())
							}()

							Eventually(func() int64 {
								return mockQueue.DeletePassedId
							}).Should(Equal(queuedDiff.ID))
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

						expectedError := watcher.NewErrHeaderMismatch(wrongHash, fakeBlockHash.Hex())
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
					obsoleteDiff := storage.PersistedStorageDiff{
						ID: queuedDiff.ID + 1,
						RawStorageDiff: storage.RawStorageDiff{
							HashedAddress: storage.HexToKeccak256Hash("0xfedcba9876543210"),
						},
					}
					mockQueue.DiffsToReturn = []storage.PersistedStorageDiff{obsoleteDiff}

					go func() {
						err := storageWatcher.Execute(time.Nanosecond)
						Expect(err).NotTo(HaveOccurred())
					}()

					Eventually(func() int64 {
						return mockQueue.DeletePassedId
					}).Should(Equal(obsoleteDiff.ID))
					close(done)
				})

				It("logs error if deleting obsolete diff fails", func(done Done) {
					obsoleteDiff := storage.PersistedStorageDiff{
						ID: queuedDiff.ID + 1,
						RawStorageDiff: storage.RawStorageDiff{
							HashedAddress: storage.HexToKeccak256Hash("0xfedcba9876543210"),
						},
					}
					mockQueue.DiffsToReturn = []storage.PersistedStorageDiff{obsoleteDiff}
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
