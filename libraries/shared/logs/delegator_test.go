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
	"github.com/vulcanize/vulcanizedb/libraries/shared/chunker"
	"github.com/vulcanize/vulcanizedb/libraries/shared/logs"
	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"strings"
)

var _ = Describe("Log transformer", func() {
	Describe("AddTransformer", func() {
		It("adds transformers to the delegator", func() {
			fakeTransformer := &mocks.MockEventTransformer{}
			delegator := logs.LogDelegator{Chunker: chunker.NewLogChunker()}

			delegator.AddTransformer(fakeTransformer)

			Expect(delegator.Transformers).To(Equal([]transformer.EventTransformer{fakeTransformer}))
		})

		It("passes transformers' configs to the chunker", func() {
			fakeTransformer := &mocks.MockEventTransformer{}
			fakeConfig := mocks.FakeTransformerConfig
			fakeTransformer.SetTransformerConfig(fakeConfig)
			chunker := chunker.NewLogChunker()
			delegator := logs.LogDelegator{Chunker: chunker}

			delegator.AddTransformer(fakeTransformer)

			expectedName := fakeConfig.TransformerName
			expectedTopic := common.HexToHash(fakeConfig.Topic)
			Expect(chunker.NameToTopic0).To(Equal(map[string]common.Hash{expectedName: expectedTopic}))
			expectedAddress := strings.ToLower(fakeConfig.ContractAddresses[0])
			Expect(chunker.AddressToNames).To(Equal(map[string][]string{expectedAddress: {expectedName}}))
		})
	})

	Describe("DelegateLogs", func() {
		var (
			errsChan  chan error
			logsFound chan bool
		)

		BeforeEach(func() {
			errsChan = make(chan error)
			logsFound = make(chan bool)
		})

		It("returns an error if no transformers configured", func(done Done) {
			delegator := newDelegator(&fakes.MockHeaderSyncLogRepository{})

			go delegator.DelegateLogs(errsChan, logsFound)

			Expect(<-errsChan).To(MatchError(logs.ErrNoTransformers))
			close(done)
		})

		It("gets untransformed logs", func(done Done) {
			mockLogRepository := &fakes.MockHeaderSyncLogRepository{}
			delegator := newDelegator(mockLogRepository)
			delegator.AddTransformer(&mocks.MockEventTransformer{})

			go delegator.DelegateLogs(errsChan, logsFound)

			Eventually(func() bool {
				return mockLogRepository.GetCalled
			}).Should(BeTrue())
			close(done)
		})

		It("emits error if getting untransformed logs fails", func(done Done) {
			mockLogRepository := &fakes.MockHeaderSyncLogRepository{}
			mockLogRepository.GetError = fakes.FakeError
			delegator := newDelegator(mockLogRepository)
			delegator.AddTransformer(&mocks.MockEventTransformer{})

			go delegator.DelegateLogs(errsChan, logsFound)

			Expect(<-errsChan).To(MatchError(fakes.FakeError))
			close(done)
		})

		It("emits that no logs were found if no logs returned", func(done Done) {
			delegator := newDelegator(&fakes.MockHeaderSyncLogRepository{})
			delegator.AddTransformer(&mocks.MockEventTransformer{})

			go delegator.DelegateLogs(errsChan, logsFound)

			Expect(<-logsFound).To(BeFalse())
			close(done)
		})

		It("delegates chunked logs to transformers", func(done Done) {
			fakeTransformer := &mocks.MockEventTransformer{}
			config := mocks.FakeTransformerConfig
			fakeTransformer.SetTransformerConfig(config)
			fakeGethLog := types.Log{
				Address: common.HexToAddress(config.ContractAddresses[0]),
				Topics:  []common.Hash{common.HexToHash(config.Topic)},
			}
			fakeHeaderSyncLogs := []core.HeaderSyncLog{{Log: fakeGethLog}}
			mockLogRepository := &fakes.MockHeaderSyncLogRepository{}
			mockLogRepository.ReturnLogs = fakeHeaderSyncLogs
			delegator := newDelegator(mockLogRepository)
			delegator.AddTransformer(fakeTransformer)

			go delegator.DelegateLogs(errsChan, logsFound)

			Eventually(func() bool {
				return fakeTransformer.ExecuteWasCalled
			}).Should(BeTrue())
			Eventually(func() []core.HeaderSyncLog {
				return fakeTransformer.PassedLogs
			}).Should(Equal(fakeHeaderSyncLogs))
			close(done)
		})

		It("emits error if transformer returns an error", func(done Done) {
			mockLogRepository := &fakes.MockHeaderSyncLogRepository{}
			mockLogRepository.ReturnLogs = []core.HeaderSyncLog{{}}
			delegator := newDelegator(mockLogRepository)
			fakeTransformer := &mocks.MockEventTransformer{ExecuteError: fakes.FakeError}
			delegator.AddTransformer(fakeTransformer)

			go delegator.DelegateLogs(errsChan, logsFound)

			Expect(<-errsChan).To(MatchError(fakes.FakeError))
			close(done)
		})

		It("emits logs found when logs returned and delegated", func(done Done) {
			fakeTransformer := &mocks.MockEventTransformer{}
			config := mocks.FakeTransformerConfig
			fakeTransformer.SetTransformerConfig(config)
			fakeGethLog := types.Log{
				Address: common.HexToAddress(config.ContractAddresses[0]),
				Topics:  []common.Hash{common.HexToHash(config.Topic)},
			}
			fakeHeaderSyncLogs := []core.HeaderSyncLog{{Log: fakeGethLog}}
			mockLogRepository := &fakes.MockHeaderSyncLogRepository{}
			mockLogRepository.ReturnLogs = fakeHeaderSyncLogs
			delegator := newDelegator(mockLogRepository)
			delegator.AddTransformer(fakeTransformer)

			go delegator.DelegateLogs(errsChan, logsFound)

			Expect(<-logsFound).To(BeTrue())
			close(done)
		})
	})

})

func newDelegator(headerSyncLogRepository *fakes.MockHeaderSyncLogRepository) *logs.LogDelegator {
	return &logs.LogDelegator{
		Chunker:       chunker.NewLogChunker(),
		LogRepository: headerSyncLogRepository,
	}
}
