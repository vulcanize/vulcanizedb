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
		It("returns an error if no transformers configured", func() {
			delegator := newDelegator(&fakes.MockHeaderSyncLogRepository{})

			err := delegator.DelegateLogs()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(logs.ErrNoTransformers))
		})

		It("gets untransformed logs", func() {
			mockLogRepository := &fakes.MockHeaderSyncLogRepository{}
			delegator := newDelegator(mockLogRepository)
			delegator.AddTransformer(&mocks.MockEventTransformer{})

			err := delegator.DelegateLogs()

			Expect(err).NotTo(HaveOccurred())
			Expect(mockLogRepository.GetCalled).To(BeTrue())
		})

		It("returns error if getting untransformed logs fails", func() {
			mockLogRepository := &fakes.MockHeaderSyncLogRepository{}
			mockLogRepository.GetError = fakes.FakeError
			delegator := newDelegator(mockLogRepository)
			delegator.AddTransformer(&mocks.MockEventTransformer{})

			err := delegator.DelegateLogs()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("delegates chunked logs to transformers", func() {
			fakeTransformer := &mocks.MockEventTransformer{}
			fakeTransformer.SetTransformerConfig(mocks.FakeTransformerConfig)
			fakeGethLog := types.Log{
				Address: common.HexToAddress(fakeTransformer.GetConfig().ContractAddresses[0]),
				Topics:  []common.Hash{common.HexToHash(fakeTransformer.GetConfig().Topic)},
			}
			fakeHeaderSyncLog := core.HeaderSyncLog{Log: fakeGethLog}
			fakeHeaderSyncLogs := []core.HeaderSyncLog{fakeHeaderSyncLog}
			mockLogRepository := &fakes.MockHeaderSyncLogRepository{}
			mockLogRepository.ReturnLogs = fakeHeaderSyncLogs
			delegator := newDelegator(mockLogRepository)
			delegator.AddTransformer(fakeTransformer)

			err := delegator.DelegateLogs()

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTransformer.ExecuteWasCalled).To(BeTrue())
			Expect(fakeTransformer.PassedLogs).To(Equal(fakeHeaderSyncLogs))
		})

		It("returns an error if transformer returns an error", func() {
			delegator := newDelegator(&fakes.MockHeaderSyncLogRepository{})
			fakeTransformer := &mocks.MockEventTransformer{ExecuteError: fakes.FakeError}
			delegator.AddTransformer(fakeTransformer)

			err := delegator.DelegateLogs()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

})

func newDelegator(headerSyncLogRepository *fakes.MockHeaderSyncLogRepository) logs.LogDelegator {
	return logs.LogDelegator{
		Chunker:       chunker.NewLogChunker(),
		LogRepository: headerSyncLogRepository,
	}
}
