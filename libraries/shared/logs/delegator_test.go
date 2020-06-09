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
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/chunker"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/logs"
	"github.com/makerdao/vulcanizedb/libraries/shared/mocks"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Log delegator", func() {
	Describe("AddTransformer", func() {
		It("adds transformers to the delegator", func() {
			fakeTransformer := &mocks.MockEventTransformer{}
			delegator := logs.LogDelegator{Chunker: chunker.NewLogChunker()}

			delegator.AddTransformer(fakeTransformer)

			Expect(delegator.Transformers).To(Equal([]event.ITransformer{fakeTransformer}))
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
		It("returns error if no transformers configured", func() {
			delegator := newDelegator(&fakes.MockEventLogRepository{})

			err := delegator.DelegateLogs(0)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(logs.ErrNoTransformers))
		})

		It("returns error if getting untransformed logs fails", func() {
			mockLogRepository := &fakes.MockEventLogRepository{}
			mockLogRepository.GetError = fakes.FakeError
			delegator := newDelegator(mockLogRepository)
			delegator.AddTransformer(&mocks.MockEventTransformer{})

			err := delegator.DelegateLogs(0)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns logs.ErrNoLogs if no logs returned on initial call", func() {
			delegator := newDelegator(&fakes.MockEventLogRepository{})
			delegator.AddTransformer(&mocks.MockEventTransformer{})

			err := delegator.DelegateLogs(0)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(logs.ErrNoLogs))
		})

		It("returns nil for error fewer logs than limit successfully delegated", func() {
			fakeTransformer := &mocks.MockEventTransformer{}
			config := mocks.FakeTransformerConfig
			fakeTransformer.SetTransformerConfig(config)
			fakeGethLog := types.Log{
				Address: common.HexToAddress(config.ContractAddresses[0]),
				Topics:  []common.Hash{common.HexToHash(config.Topic)},
			}
			fakeEventLogs := []core.EventLog{{Log: fakeGethLog}}
			mockLogRepository := &fakes.MockEventLogRepository{}
			mockLogRepository.ReturnLogs = fakeEventLogs
			delegator := newDelegator(mockLogRepository)
			delegator.AddTransformer(fakeTransformer)

			limitGreaterThanUntransformedLogs := 2
			err := delegator.DelegateLogs(limitGreaterThanUntransformedLogs)

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTransformer.ExecuteWasCalled).To(BeTrue())
			Expect(fakeTransformer.PassedLogs).To(Equal(fakeEventLogs))
		})

		It("repeats logs lookup with minID from last result when repository returns maximum number of logs", func() {
			mockLogRepository := &fakes.MockEventLogRepository{}
			returnLogs := []core.EventLog{{ID: 1}, {ID: 2}, {ID: 3}}
			mockLogRepository.ReturnLogs = returnLogs
			delegator := newDelegator(mockLogRepository)
			delegator.AddTransformer(&mocks.MockEventTransformer{})

			limit := len(returnLogs) - 1
			err := delegator.DelegateLogs(limit)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockLogRepository.PassedMinIDs).To(ConsistOf(0, int(returnLogs[1].ID)))
			Expect(mockLogRepository.PassedLimits).To(ConsistOf(limit, limit))
		})

		It("returns logs.ErrNoLogs if no logs returned on subsequent call", func() {
			fakeTransformer := &mocks.MockEventTransformer{}
			config := mocks.FakeTransformerConfig
			fakeTransformer.SetTransformerConfig(config)
			fakeGethLog := types.Log{
				Address: common.HexToAddress(config.ContractAddresses[0]),
				Topics:  []common.Hash{common.HexToHash(config.Topic)},
			}
			fakeEventLogs := []core.EventLog{{Log: fakeGethLog}}
			mockLogRepository := &fakes.MockEventLogRepository{}
			mockLogRepository.ReturnLogs = fakeEventLogs
			delegator := newDelegator(mockLogRepository)
			delegator.AddTransformer(fakeTransformer)

			err := delegator.DelegateLogs(1)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(logs.ErrNoLogs))
			Expect(fakeTransformer.ExecuteWasCalled).To(BeTrue())
			Expect(fakeTransformer.PassedLogs).To(Equal(fakeEventLogs))
		})

		It("returns error if transformer returns an error", func() {
			mockLogRepository := &fakes.MockEventLogRepository{}
			mockLogRepository.ReturnLogs = []core.EventLog{{}}
			delegator := newDelegator(mockLogRepository)
			fakeTransformer := &mocks.MockEventTransformer{ExecuteError: fakes.FakeError}
			delegator.AddTransformer(fakeTransformer)

			err := delegator.DelegateLogs(1)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})

func newDelegator(eventLogRepository *fakes.MockEventLogRepository) *logs.LogDelegator {
	return &logs.LogDelegator{
		Chunker:       chunker.NewLogChunker(),
		LogRepository: eventLogRepository,
	}
}
