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
	var extractor *logs.LogExtractor

	BeforeEach(func() {
		extractor = &logs.LogExtractor{
			Fetcher:                  &mocks.MockLogFetcher{},
			CheckedHeadersRepository: &fakes.MockCheckedHeadersRepository{},
			LogRepository:            &fakes.MockHeaderSyncLogRepository{},
			Syncer:                   &fakes.MockTransactionSyncer{},
		}
	})

	Describe("AddTransformerConfig", func() {
		It("it includes earliest starting block number in fetch logs query", func() {
			earlierStartingBlockNumber := rand.Int63()
			laterStartingBlockNumber := earlierStartingBlockNumber + 1

			extractor.AddTransformerConfig(getTransformerConfig(laterStartingBlockNumber))
			extractor.AddTransformerConfig(getTransformerConfig(earlierStartingBlockNumber))

			Expect(*extractor.StartingBlock).To(Equal(earlierStartingBlockNumber))
		})

		It("includes added addresses in fetch logs query", func() {
			addresses := []string{"0xA", "0xB"}
			configWithAddresses := transformer.EventTransformerConfig{
				ContractAddresses:   addresses,
				StartingBlockNumber: rand.Int63(),
			}

			extractor.AddTransformerConfig(configWithAddresses)

			expectedAddresses := transformer.HexStringsToAddresses(addresses)
			Expect(extractor.Addresses).To(Equal(expectedAddresses))
		})

		It("includes added topics in fetch logs query", func() {
			topic := "0x1"
			configWithTopic := transformer.EventTransformerConfig{
				ContractAddresses:   []string{fakes.FakeAddress.Hex()},
				Topic:               topic,
				StartingBlockNumber: rand.Int63(),
			}

			extractor.AddTransformerConfig(configWithTopic)

			Expect(extractor.Topics).To(Equal([]common.Hash{common.HexToHash(topic)}))
		})
	})

	Describe("ExtractLogs", func() {
		var (
			errsChan            chan error
			missingHeadersFound chan bool
		)

		BeforeEach(func() {
			errsChan = make(chan error)
			missingHeadersFound = make(chan bool)
		})

		It("returns error if no watched addresses configured", func(done Done) {
			go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

			Expect(<-errsChan).To(MatchError(logs.ErrNoWatchedAddresses))
			close(done)
		})

		Describe("when checking missing headers", func() {
			It("gets missing headers since configured starting block with check_count < 1", func(done Done) {
				mockCheckedHeadersRepository := &fakes.MockCheckedHeadersRepository{}
				mockCheckedHeadersRepository.ReturnHeaders = []core.Header{{}}
				extractor.CheckedHeadersRepository = mockCheckedHeadersRepository
				startingBlockNumber := rand.Int63()
				extractor.AddTransformerConfig(getTransformerConfig(startingBlockNumber))

				go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

				Eventually(func() int64 {
					return mockCheckedHeadersRepository.StartingBlockNumber
				}).Should(Equal(startingBlockNumber))
				Eventually(func() int64 {
					return mockCheckedHeadersRepository.EndingBlockNumber
				}).Should(Equal(int64(-1)))
				Eventually(func() int64 {
					return mockCheckedHeadersRepository.CheckCount
				}).Should(Equal(int64(1)))
				close(done)
			})
		})

		Describe("when rechecking headers", func() {
			It("gets missing headers since configured starting block with check_count < RecheckHeaderCap", func(done Done) {
				mockCheckedHeadersRepository := &fakes.MockCheckedHeadersRepository{}
				mockCheckedHeadersRepository.ReturnHeaders = []core.Header{{}}
				extractor.CheckedHeadersRepository = mockCheckedHeadersRepository
				startingBlockNumber := rand.Int63()
				extractor.AddTransformerConfig(getTransformerConfig(startingBlockNumber))

				go extractor.ExtractLogs(constants.HeaderRecheck, errsChan, missingHeadersFound)

				Eventually(func() int64 {
					return mockCheckedHeadersRepository.StartingBlockNumber
				}).Should(Equal(startingBlockNumber))
				Eventually(func() int64 {
					return mockCheckedHeadersRepository.EndingBlockNumber
				}).Should(Equal(int64(-1)))
				Eventually(func() int64 {
					return mockCheckedHeadersRepository.CheckCount
				}).Should(Equal(constants.RecheckHeaderCap))
				close(done)
			})
		})

		It("emits error if getting missing headers fails", func(done Done) {
			addTransformerConfig(extractor)
			mockCheckedHeadersRepository := &fakes.MockCheckedHeadersRepository{}
			mockCheckedHeadersRepository.MissingHeadersReturnError = fakes.FakeError
			extractor.CheckedHeadersRepository = mockCheckedHeadersRepository

			go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

			Expect(<-errsChan).To(MatchError(fakes.FakeError))
			close(done)
		})

		Describe("when no missing headers", func() {
			It("does not fetch logs", func(done Done) {
				addTransformerConfig(extractor)
				mockLogFetcher := &mocks.MockLogFetcher{}
				extractor.Fetcher = mockLogFetcher

				go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

				Consistently(func() bool {
					return mockLogFetcher.FetchCalled
				}).Should(BeFalse())
				close(done)
			})

			It("emits that no missing headers were found", func(done Done) {
				addTransformerConfig(extractor)
				mockLogFetcher := &mocks.MockLogFetcher{}
				extractor.Fetcher = mockLogFetcher

				go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

				Expect(<-missingHeadersFound).To(BeFalse())
				close(done)
			})
		})

		Describe("when there are missing headers", func() {
			It("fetches logs for missing headers", func(done Done) {
				addMissingHeader(extractor)
				config := transformer.EventTransformerConfig{
					ContractAddresses:   []string{fakes.FakeAddress.Hex()},
					Topic:               fakes.FakeHash.Hex(),
					StartingBlockNumber: rand.Int63(),
				}
				extractor.AddTransformerConfig(config)
				mockLogFetcher := &mocks.MockLogFetcher{}
				extractor.Fetcher = mockLogFetcher

				go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

				Eventually(func() bool {
					return mockLogFetcher.FetchCalled
				}).Should(BeTrue())
				expectedTopics := []common.Hash{common.HexToHash(config.Topic)}
				Eventually(func() []common.Hash {
					return mockLogFetcher.Topics
				}).Should(Equal(expectedTopics))
				expectedAddresses := transformer.HexStringsToAddresses(config.ContractAddresses)
				Eventually(func() []common.Address {
					return mockLogFetcher.ContractAddresses
				}).Should(Equal(expectedAddresses))
				close(done)
			})

			It("returns error if fetching logs fails", func(done Done) {
				addMissingHeader(extractor)
				addTransformerConfig(extractor)
				mockLogFetcher := &mocks.MockLogFetcher{}
				mockLogFetcher.ReturnError = fakes.FakeError
				extractor.Fetcher = mockLogFetcher

				go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

				Expect(<-errsChan).To(MatchError(fakes.FakeError))
				close(done)
			})

			Describe("when no fetched logs", func() {
				It("does not sync transactions", func(done Done) {
					addMissingHeader(extractor)
					addTransformerConfig(extractor)
					mockTransactionSyncer := &fakes.MockTransactionSyncer{}
					extractor.Syncer = mockTransactionSyncer

					go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

					Consistently(func() bool {
						return mockTransactionSyncer.SyncTransactionsCalled
					}).Should(BeFalse())
					close(done)
				})
			})

			Describe("when there are fetched logs", func() {
				It("syncs transactions", func(done Done) {
					addMissingHeader(extractor)
					addFetchedLog(extractor)
					addTransformerConfig(extractor)
					mockTransactionSyncer := &fakes.MockTransactionSyncer{}
					extractor.Syncer = mockTransactionSyncer

					go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

					Eventually(func() bool {
						return mockTransactionSyncer.SyncTransactionsCalled
					}).Should(BeTrue())
					close(done)
				})

				It("returns error if syncing transactions fails", func(done Done) {
					addMissingHeader(extractor)
					addFetchedLog(extractor)
					addTransformerConfig(extractor)
					mockTransactionSyncer := &fakes.MockTransactionSyncer{}
					mockTransactionSyncer.SyncTransactionsError = fakes.FakeError
					extractor.Syncer = mockTransactionSyncer

					go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

					Expect(<-errsChan).To(MatchError(fakes.FakeError))
					close(done)
				})

				It("persists fetched logs", func(done Done) {
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

					go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

					Eventually(func() []types.Log {
						return mockLogRepository.PassedLogs
					}).Should(Equal(fakeLogs))
					close(done)
				})

				It("returns error if persisting logs fails", func(done Done) {
					addMissingHeader(extractor)
					addFetchedLog(extractor)
					addTransformerConfig(extractor)
					mockLogRepository := &fakes.MockHeaderSyncLogRepository{}
					mockLogRepository.CreateError = fakes.FakeError
					extractor.LogRepository = mockLogRepository

					go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

					Expect(<-errsChan).To(MatchError(fakes.FakeError))
					close(done)
				})
			})

			It("marks header checked", func(done Done) {
				addFetchedLog(extractor)
				addTransformerConfig(extractor)
				mockCheckedHeadersRepository := &fakes.MockCheckedHeadersRepository{}
				headerID := rand.Int63()
				mockCheckedHeadersRepository.ReturnHeaders = []core.Header{{Id: headerID}}
				extractor.CheckedHeadersRepository = mockCheckedHeadersRepository

				go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

				Eventually(func() int64 {
					return mockCheckedHeadersRepository.HeaderID
				}).Should(Equal(headerID))
				close(done)
			})

			It("returns error if marking header checked fails", func(done Done) {
				addFetchedLog(extractor)
				addTransformerConfig(extractor)
				mockCheckedHeadersRepository := &fakes.MockCheckedHeadersRepository{}
				mockCheckedHeadersRepository.ReturnHeaders = []core.Header{{Id: rand.Int63()}}
				mockCheckedHeadersRepository.MarkHeaderCheckedReturnError = fakes.FakeError
				extractor.CheckedHeadersRepository = mockCheckedHeadersRepository

				go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

				Expect(<-errsChan).To(MatchError(fakes.FakeError))
				close(done)
			})

			It("emits that missing headers were found", func(done Done) {
				addMissingHeader(extractor)
				addTransformerConfig(extractor)

				go extractor.ExtractLogs(constants.HeaderMissing, errsChan, missingHeadersFound)

				Expect(<-missingHeadersFound).To(BeTrue())
				close(done)
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
	mockCheckedHeadersRepository.ReturnHeaders = []core.Header{{}}
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
