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

package logs_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/logs"
	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
)

var _ = Describe("Log extractor", func() {
	var (
		checkedHeadersRepository *fakes.MockCheckedHeadersRepository
		checkedLogsRepository    *fakes.MockCheckedLogsRepository
		extractor                *logs.LogExtractor
	)

	BeforeEach(func() {
		checkedHeadersRepository = &fakes.MockCheckedHeadersRepository{}
		checkedLogsRepository = &fakes.MockCheckedLogsRepository{}
		extractor = &logs.LogExtractor{
			CheckedHeadersRepository: checkedHeadersRepository,
			CheckedLogsRepository:    checkedLogsRepository,
			Fetcher:                  &mocks.MockLogFetcher{},
			LogRepository:            &fakes.MockHeaderSyncLogRepository{},
			Syncer:                   &fakes.MockTransactionSyncer{},
		}
	})

	Describe("AddTransformerConfig", func() {
		It("updates extractor's starting block number to earliest available", func() {
			earlierStartingBlockNumber := rand.Int63()
			laterStartingBlockNumber := earlierStartingBlockNumber + 1

			errOne := extractor.AddTransformerConfig(getTransformerConfig(laterStartingBlockNumber))
			Expect(errOne).NotTo(HaveOccurred())
			errTwo := extractor.AddTransformerConfig(getTransformerConfig(earlierStartingBlockNumber))
			Expect(errTwo).NotTo(HaveOccurred())

			Expect(*extractor.StartingBlock).To(Equal(earlierStartingBlockNumber))
		})

		It("adds transformer's addresses to extractor's watched addresses", func() {
			addresses := []string{"0xA", "0xB"}
			configWithAddresses := transformer.EventTransformerConfig{
				ContractAddresses:   addresses,
				StartingBlockNumber: rand.Int63(),
			}

			err := extractor.AddTransformerConfig(configWithAddresses)

			Expect(err).NotTo(HaveOccurred())
			expectedAddresses := transformer.HexStringsToAddresses(addresses)
			Expect(extractor.Addresses).To(Equal(expectedAddresses))
		})

		It("adds transformer's topic to extractor's watched topics", func() {
			topic := "0x1"
			configWithTopic := transformer.EventTransformerConfig{
				ContractAddresses:   []string{fakes.FakeAddress.Hex()},
				Topic:               topic,
				StartingBlockNumber: rand.Int63(),
			}

			err := extractor.AddTransformerConfig(configWithTopic)

			Expect(err).NotTo(HaveOccurred())
			Expect(extractor.Topics).To(Equal([]common.Hash{common.HexToHash(topic)}))
		})

		It("returns error if checking whether log has been checked returns error", func() {
			checkedLogsRepository.HasLogBeenCheckedError = fakes.FakeError

			err := extractor.AddTransformerConfig(getTransformerConfig(rand.Int63()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		Describe("when log has previously been checked", func() {
			It("does not mark any headers unchecked", func() {
				checkedLogsRepository.HasLogBeenCheckedReturn = true

				err := extractor.AddTransformerConfig(getTransformerConfig(rand.Int63()))

				Expect(err).NotTo(HaveOccurred())
				Expect(checkedHeadersRepository.MarkHeadersUncheckedCalled).To(BeFalse())
			})
		})

		Describe("when log has not previously been checked", func() {
			BeforeEach(func() {
				checkedLogsRepository.HasLogBeenCheckedReturn = false
			})

			It("marks headers since transformer's starting block number as unchecked", func() {
				blockNumber := rand.Int63()

				err := extractor.AddTransformerConfig(getTransformerConfig(blockNumber))

				Expect(err).NotTo(HaveOccurred())
				Expect(checkedHeadersRepository.MarkHeadersUncheckedCalled).To(BeTrue())
				Expect(checkedHeadersRepository.MarkHeadersUncheckedStartingBlockNumber).To(Equal(blockNumber))
			})

			It("returns error if marking headers unchecked returns error", func() {
				checkedHeadersRepository.MarkHeadersUncheckedReturnError = fakes.FakeError

				err := extractor.AddTransformerConfig(getTransformerConfig(rand.Int63()))

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})

			It("persists that tranformer's log has been checked", func() {
				config := getTransformerConfig(rand.Int63())

				err := extractor.AddTransformerConfig(config)

				Expect(err).NotTo(HaveOccurred())
				Expect(checkedLogsRepository.MarkLogCheckedAddresses).To(Equal(config.ContractAddresses))
				Expect(checkedLogsRepository.MarkLogCheckedTopicZero).To(Equal(config.Topic))
			})

			It("returns error if marking logs checked returns error", func() {
				checkedLogsRepository.MarkLogCheckedError = fakes.FakeError

				err := extractor.AddTransformerConfig(getTransformerConfig(rand.Int63()))

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})
		})
	})

	Describe("ExtractLogs", func() {
		It("returns error if no watched addresses configured", func() {
			err, _ := extractor.ExtractLogs(constants.HeaderMissing)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(logs.ErrNoWatchedAddresses))
		})

		Describe("when checking missing headers", func() {
			It("gets missing headers since configured starting block with check_count < 1", func() {
				mockCheckedHeadersRepository := &fakes.MockCheckedHeadersRepository{}
				mockCheckedHeadersRepository.MissingHeadersReturnHeaders = []core.Header{{}}
				extractor.CheckedHeadersRepository = mockCheckedHeadersRepository
				startingBlockNumber := rand.Int63()
				extractor.AddTransformerConfig(getTransformerConfig(startingBlockNumber))

				err, _ := extractor.ExtractLogs(constants.HeaderMissing)

				Expect(err).NotTo(HaveOccurred())
				Expect(mockCheckedHeadersRepository.MissingHeadersStartingBlockNumber).To(Equal(startingBlockNumber))
				Expect(mockCheckedHeadersRepository.MissingHeadersEndingBlockNumber).To(Equal(int64(-1)))
				Expect(mockCheckedHeadersRepository.MissingHeadersCheckCount).To(Equal(int64(1)))
			})
		})

		Describe("when rechecking headers", func() {
			It("gets missing headers since configured starting block with check_count < RecheckHeaderCap", func() {
				mockCheckedHeadersRepository := &fakes.MockCheckedHeadersRepository{}
				mockCheckedHeadersRepository.MissingHeadersReturnHeaders = []core.Header{{}}
				extractor.CheckedHeadersRepository = mockCheckedHeadersRepository
				startingBlockNumber := rand.Int63()
				extractor.AddTransformerConfig(getTransformerConfig(startingBlockNumber))

				err, _ := extractor.ExtractLogs(constants.HeaderRecheck)

				Expect(err).NotTo(HaveOccurred())
				Expect(mockCheckedHeadersRepository.MissingHeadersStartingBlockNumber).To(Equal(startingBlockNumber))
				Expect(mockCheckedHeadersRepository.MissingHeadersEndingBlockNumber).To(Equal(int64(-1)))
				Expect(mockCheckedHeadersRepository.MissingHeadersCheckCount).To(Equal(constants.RecheckHeaderCap))
			})
		})

		It("emits error if getting missing headers fails", func() {
			addTransformerConfig(extractor)
			mockCheckedHeadersRepository := &fakes.MockCheckedHeadersRepository{}
			mockCheckedHeadersRepository.MissingHeadersReturnError = fakes.FakeError
			extractor.CheckedHeadersRepository = mockCheckedHeadersRepository

			err, _ := extractor.ExtractLogs(constants.HeaderMissing)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		Describe("when no missing headers", func() {
			It("does not fetch logs", func() {
				addTransformerConfig(extractor)
				mockLogFetcher := &mocks.MockLogFetcher{}
				extractor.Fetcher = mockLogFetcher

				err, _ := extractor.ExtractLogs(constants.HeaderMissing)

				Expect(err).NotTo(HaveOccurred())
				Expect(mockLogFetcher.FetchCalled).To(BeFalse())
			})

			It("emits that no missing headers were found", func() {
				addTransformerConfig(extractor)
				mockLogFetcher := &mocks.MockLogFetcher{}
				extractor.Fetcher = mockLogFetcher

				_, missingHeadersFound := extractor.ExtractLogs(constants.HeaderMissing)

				Expect(missingHeadersFound).To(BeFalse())
			})
		})

		Describe("when there are missing headers", func() {
			It("fetches logs for missing headers", func() {
				addMissingHeader(extractor)
				config := transformer.EventTransformerConfig{
					ContractAddresses:   []string{fakes.FakeAddress.Hex()},
					Topic:               fakes.FakeHash.Hex(),
					StartingBlockNumber: rand.Int63(),
				}
				extractor.AddTransformerConfig(config)
				mockLogFetcher := &mocks.MockLogFetcher{}
				extractor.Fetcher = mockLogFetcher

				err, _ := extractor.ExtractLogs(constants.HeaderMissing)

				Expect(err).NotTo(HaveOccurred())
				Expect(mockLogFetcher.FetchCalled).To(BeTrue())
				expectedTopics := []common.Hash{common.HexToHash(config.Topic)}
				Expect(mockLogFetcher.Topics).To(Equal(expectedTopics))
				expectedAddresses := transformer.HexStringsToAddresses(config.ContractAddresses)
				Expect(mockLogFetcher.ContractAddresses).To(Equal(expectedAddresses))
			})

			It("returns error if fetching logs fails", func() {
				addMissingHeader(extractor)
				addTransformerConfig(extractor)
				mockLogFetcher := &mocks.MockLogFetcher{}
				mockLogFetcher.ReturnError = fakes.FakeError
				extractor.Fetcher = mockLogFetcher

				err, _ := extractor.ExtractLogs(constants.HeaderMissing)

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})

			Describe("when no fetched logs", func() {
				It("does not sync transactions", func() {
					addMissingHeader(extractor)
					addTransformerConfig(extractor)
					mockTransactionSyncer := &fakes.MockTransactionSyncer{}
					extractor.Syncer = mockTransactionSyncer

					err, _ := extractor.ExtractLogs(constants.HeaderMissing)

					Expect(err).NotTo(HaveOccurred())
					Expect(mockTransactionSyncer.SyncTransactionsCalled).To(BeFalse())
				})
			})

			Describe("when there are fetched logs", func() {
				It("syncs transactions", func() {
					addMissingHeader(extractor)
					addFetchedLog(extractor)
					addTransformerConfig(extractor)
					mockTransactionSyncer := &fakes.MockTransactionSyncer{}
					extractor.Syncer = mockTransactionSyncer

					err, _ := extractor.ExtractLogs(constants.HeaderMissing)

					Expect(err).NotTo(HaveOccurred())
					Expect(mockTransactionSyncer.SyncTransactionsCalled).To(BeTrue())
				})

				It("returns error if syncing transactions fails", func() {
					addMissingHeader(extractor)
					addFetchedLog(extractor)
					addTransformerConfig(extractor)
					mockTransactionSyncer := &fakes.MockTransactionSyncer{}
					mockTransactionSyncer.SyncTransactionsError = fakes.FakeError
					extractor.Syncer = mockTransactionSyncer

					err, _ := extractor.ExtractLogs(constants.HeaderMissing)

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(fakes.FakeError))
				})

				It("persists fetched logs", func() {
					addMissingHeader(extractor)
					addTransformerConfig(extractor)
					fakeLogs := []types.Log{{
						Address: common.HexToAddress("0xA"),
						Topics:  []common.Hash{common.HexToHash("0xA")},
						Data:    []byte{},
						Index:   0,
					}}
					mockLogFetcher := &mocks.MockLogFetcher{ReturnLogs: fakeLogs}
					extractor.Fetcher = mockLogFetcher
					mockLogRepository := &fakes.MockHeaderSyncLogRepository{}
					extractor.LogRepository = mockLogRepository

					err, _ := extractor.ExtractLogs(constants.HeaderMissing)

					Expect(err).NotTo(HaveOccurred())
					Expect(mockLogRepository.PassedLogs).To(Equal(fakeLogs))
				})

				It("returns error if persisting logs fails", func() {
					addMissingHeader(extractor)
					addFetchedLog(extractor)
					addTransformerConfig(extractor)
					mockLogRepository := &fakes.MockHeaderSyncLogRepository{}
					mockLogRepository.CreateError = fakes.FakeError
					extractor.LogRepository = mockLogRepository

					err, _ := extractor.ExtractLogs(constants.HeaderMissing)

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(fakes.FakeError))
				})
			})

			It("marks header checked", func() {
				addFetchedLog(extractor)
				addTransformerConfig(extractor)
				mockCheckedHeadersRepository := &fakes.MockCheckedHeadersRepository{}
				headerID := rand.Int63()
				mockCheckedHeadersRepository.MissingHeadersReturnHeaders = []core.Header{{Id: headerID}}
				extractor.CheckedHeadersRepository = mockCheckedHeadersRepository

				err, _ := extractor.ExtractLogs(constants.HeaderMissing)

				Expect(err).NotTo(HaveOccurred())
				Expect(mockCheckedHeadersRepository.MarkHeaderCheckedHeaderID).To(Equal(headerID))
			})

			It("returns error if marking header checked fails", func() {
				addFetchedLog(extractor)
				addTransformerConfig(extractor)
				mockCheckedHeadersRepository := &fakes.MockCheckedHeadersRepository{}
				mockCheckedHeadersRepository.MissingHeadersReturnHeaders = []core.Header{{Id: rand.Int63()}}
				mockCheckedHeadersRepository.MarkHeaderCheckedReturnError = fakes.FakeError
				extractor.CheckedHeadersRepository = mockCheckedHeadersRepository

				err, _ := extractor.ExtractLogs(constants.HeaderMissing)

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})

			It("emits that missing headers were found", func() {
				addMissingHeader(extractor)
				addTransformerConfig(extractor)

				err, missingHeadersFound := extractor.ExtractLogs(constants.HeaderMissing)

				Expect(err).NotTo(HaveOccurred())
				Expect(missingHeadersFound).To(BeTrue())
			})
		})
	})
})

func addTransformerConfig(extractor *logs.LogExtractor) {
	fakeConfig := transformer.EventTransformerConfig{
		ContractAddresses:   []string{fakes.FakeAddress.Hex()},
		Topic:               fakes.FakeHash.Hex(),
		StartingBlockNumber: rand.Int63(),
	}
	extractor.AddTransformerConfig(fakeConfig)
}

func addMissingHeader(extractor *logs.LogExtractor) {
	mockCheckedHeadersRepository := &fakes.MockCheckedHeadersRepository{}
	mockCheckedHeadersRepository.MissingHeadersReturnHeaders = []core.Header{{}}
	extractor.CheckedHeadersRepository = mockCheckedHeadersRepository
}

func addFetchedLog(extractor *logs.LogExtractor) {
	mockLogFetcher := &mocks.MockLogFetcher{}
	mockLogFetcher.ReturnLogs = []types.Log{{}}
	extractor.Fetcher = mockLogFetcher
}

func getTransformerConfig(startingBlockNumber int64) transformer.EventTransformerConfig {
	return transformer.EventTransformerConfig{
		ContractAddresses:   []string{fakes.FakeAddress.Hex()},
		Topic:               fakes.FakeHash.Hex(),
		StartingBlockNumber: startingBlockNumber,
	}
}
