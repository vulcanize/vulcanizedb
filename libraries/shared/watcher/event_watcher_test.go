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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Event Watcher", func() {
	var (
		delegator    *mocks.MockLogDelegator
		extractor    *mocks.MockLogExtractor
		eventWatcher *watcher.EventWatcher
	)

	BeforeEach(func() {
		delegator = &mocks.MockLogDelegator{}
		extractor = &mocks.MockLogExtractor{}
		eventWatcher = &watcher.EventWatcher{
			LogDelegator: delegator,
			LogExtractor: extractor,
		}
	})

	Describe("AddTransformers", func() {
		var (
			fakeTransformerOne, fakeTransformerTwo *mocks.MockEventTransformer
		)

		BeforeEach(func() {
			fakeTransformerOne = &mocks.MockEventTransformer{}
			fakeTransformerOne.SetTransformerConfig(mocks.FakeTransformerConfig)
			fakeTransformerTwo = &mocks.MockEventTransformer{}
			fakeTransformerTwo.SetTransformerConfig(mocks.FakeTransformerConfig)
			initializers := []transformer.EventTransformerInitializer{
				fakeTransformerOne.FakeTransformerInitializer,
				fakeTransformerTwo.FakeTransformerInitializer,
			}

			eventWatcher.AddTransformers(initializers)
		})

		It("adds initialized transformer to log delegator", func() {
			expectedTransformers := []transformer.EventTransformer{
				fakeTransformerOne,
				fakeTransformerTwo,
			}
			Expect(delegator.AddedTransformers).To(Equal(expectedTransformers))
		})

		It("adds transformer config to log extractor", func() {
			expectedConfigs := []transformer.EventTransformerConfig{
				mocks.FakeTransformerConfig,
				mocks.FakeTransformerConfig,
			}
			Expect(extractor.AddedConfigs).To(Equal(expectedConfigs))
		})
	})

	Describe("Execute", func() {
		var errsChan chan error

		BeforeEach(func() {
			errsChan = make(chan error)
		})

		It("extracts watched logs", func(done Done) {
			delegator.DelegateErrors = []error{nil}
			delegator.LogsFound = []bool{false}
			extractor.ExtractLogsErrors = []error{nil}
			extractor.UncheckedHeadersExist = []bool{false}

			go eventWatcher.Execute(constants.HeaderUnchecked, errsChan)

			Eventually(func() int {
				return extractor.ExtractLogsCount
			}).Should(Equal(1))
			close(done)
		})

		It("returns error if extracting logs fails", func(done Done) {
			delegator.DelegateErrors = []error{nil}
			delegator.LogsFound = []bool{false}
			extractor.ExtractLogsErrors = []error{fakes.FakeError}
			extractor.UncheckedHeadersExist = []bool{false}

			go eventWatcher.Execute(constants.HeaderUnchecked, errsChan)

			Expect(<-errsChan).To(MatchError(fakes.FakeError))
			close(done)
		})

		It("extracts watched logs again if missing headers found", func(done Done) {
			delegator.DelegateErrors = []error{nil}
			delegator.LogsFound = []bool{false}
			extractor.ExtractLogsErrors = []error{nil, nil}
			extractor.UncheckedHeadersExist = []bool{true, false}

			go eventWatcher.Execute(constants.HeaderUnchecked, errsChan)

			Eventually(func() int {
				return extractor.ExtractLogsCount
			}).Should(Equal(2))
			close(done)
		})

		It("returns error if extracting logs fails on subsequent run", func(done Done) {
			delegator.DelegateErrors = []error{nil}
			delegator.LogsFound = []bool{false}
			extractor.ExtractLogsErrors = []error{nil, fakes.FakeError}
			extractor.UncheckedHeadersExist = []bool{true, false}

			go eventWatcher.Execute(constants.HeaderUnchecked, errsChan)

			Expect(<-errsChan).To(MatchError(fakes.FakeError))
			close(done)

		})

		It("delegates untransformed logs", func(done Done) {
			delegator.DelegateErrors = []error{nil}
			delegator.LogsFound = []bool{false}
			extractor.ExtractLogsErrors = []error{nil}
			extractor.UncheckedHeadersExist = []bool{false}

			go eventWatcher.Execute(constants.HeaderUnchecked, errsChan)

			Eventually(func() int {
				return delegator.DelegateCallCount
			}).Should(Equal(1))
			close(done)
		})

		It("returns error if delegating logs fails", func(done Done) {
			delegator.LogsFound = []bool{false}
			delegator.DelegateErrors = []error{fakes.FakeError}
			extractor.ExtractLogsErrors = []error{nil}
			extractor.UncheckedHeadersExist = []bool{false}

			go eventWatcher.Execute(constants.HeaderUnchecked, errsChan)

			Expect(<-errsChan).To(MatchError(fakes.FakeError))
			close(done)
		})

		It("delegates logs again if untransformed logs found", func(done Done) {
			delegator.DelegateErrors = []error{nil, nil}
			delegator.LogsFound = []bool{true, false}
			extractor.ExtractLogsErrors = []error{nil}
			extractor.UncheckedHeadersExist = []bool{false}

			go eventWatcher.Execute(constants.HeaderUnchecked, errsChan)

			Eventually(func() int {
				return delegator.DelegateCallCount
			}).Should(Equal(2))
			close(done)
		})

		It("returns error if delegating logs fails on subsequent run", func(done Done) {
			delegator.DelegateErrors = []error{nil, fakes.FakeError}
			delegator.LogsFound = []bool{true, false}
			extractor.ExtractLogsErrors = []error{nil}
			extractor.UncheckedHeadersExist = []bool{false}

			go eventWatcher.Execute(constants.HeaderUnchecked, errsChan)

			Expect(<-errsChan).To(MatchError(fakes.FakeError))
			close(done)
		})
	})
})
