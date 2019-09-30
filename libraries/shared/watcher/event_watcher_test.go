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
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/logs"
	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var errExecuteClosed = errors.New("this error means the mocks were finished executing")

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

			err := eventWatcher.AddTransformers(initializers)
			Expect(err).NotTo(HaveOccurred())
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

		It("extracts watched logs", func(done Done) {
			delegator.DelegateErrors = []error{logs.ErrNoLogs}
			extractor.ExtractLogsErrors = []error{nil, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
			Eventually(func() bool {
				return extractor.ExtractLogsCount > 0
			}).Should(BeTrue())
			close(done)
		})

		It("returns error if extracting logs fails", func(done Done) {
			delegator.DelegateErrors = []error{logs.ErrNoLogs}
			extractor.ExtractLogsErrors = []error{fakes.FakeError}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(fakes.FakeError))
			close(done)
		})

		It("extracts watched logs again if missing headers found", func(done Done) {
			delegator.DelegateErrors = []error{logs.ErrNoLogs}
			extractor.ExtractLogsErrors = []error{nil, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
			Eventually(func() bool {
				return extractor.ExtractLogsCount > 1
			}).Should(BeTrue())
			close(done)
		})

		It("returns error if extracting logs fails on subsequent run", func(done Done) {
			delegator.DelegateErrors = []error{logs.ErrNoLogs}
			extractor.ExtractLogsErrors = []error{nil, fakes.FakeError}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(fakes.FakeError))
			close(done)
		})

		It("delegates untransformed logs", func() {
			delegator.DelegateErrors = []error{nil, errExecuteClosed}
			extractor.ExtractLogsErrors = []error{logs.ErrNoUncheckedHeaders}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
			Eventually(func() bool {
				return delegator.DelegateCallCount > 0
			}).Should(BeTrue())
		})

		It("returns error if delegating logs fails", func(done Done) {
			delegator.DelegateErrors = []error{fakes.FakeError}
			extractor.ExtractLogsErrors = []error{logs.ErrNoUncheckedHeaders}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(fakes.FakeError))
			close(done)
		})

		It("delegates logs again if untransformed logs found", func(done Done) {
			delegator.DelegateErrors = []error{nil, nil, nil, errExecuteClosed}
			extractor.ExtractLogsErrors = []error{logs.ErrNoUncheckedHeaders}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
			Eventually(func() bool {
				return delegator.DelegateCallCount > 1
			}).Should(BeTrue())
			close(done)
		})

		It("returns error if delegating logs fails on subsequent run", func(done Done) {
			delegator.DelegateErrors = []error{nil, fakes.FakeError}
			extractor.ExtractLogsErrors = []error{logs.ErrNoUncheckedHeaders}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(fakes.FakeError))
			close(done)
		})
	})
})
