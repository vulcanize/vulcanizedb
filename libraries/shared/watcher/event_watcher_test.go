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
	"io"
	"time"

	"github.com/makerdao/vulcanizedb/libraries/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/logs"
	"github.com/makerdao/vulcanizedb/libraries/shared/mocks"
	"github.com/makerdao/vulcanizedb/libraries/shared/watcher"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var errExecuteClosed = errors.New("this error means the mocks were finished executing")

var _ = Describe("Event Watcher", func() {
	var (
		delegator    *mocks.MockLogDelegator
		extractor    *mocks.MockLogExtractor
		eventWatcher watcher.EventWatcher
	)

	BeforeEach(func() {
		delegator = &mocks.MockLogDelegator{}
		extractor = &mocks.MockLogExtractor{}
		bc := fakes.MockBlockChain{}
		eventWatcher = watcher.NewEventWatcher(nil, &bc, extractor, delegator, 0, time.Nanosecond)
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
			initializers := []event.TransformerInitializer{
				fakeTransformerOne.FakeTransformerInitializer,
				fakeTransformerTwo.FakeTransformerInitializer,
			}

			err := eventWatcher.AddTransformers(initializers)
			Expect(err).NotTo(HaveOccurred())
		})

		It("adds initialized transformer to log delegator", func() {
			expectedTransformers := []event.ITransformer{
				fakeTransformerOne,
				fakeTransformerTwo,
			}
			Expect(delegator.AddedTransformers).To(Equal(expectedTransformers))
		})

		It("adds transformer config to log extractor", func() {
			expectedConfigs := []event.TransformerConfig{
				mocks.FakeTransformerConfig,
				mocks.FakeTransformerConfig,
			}
			Expect(extractor.AddedConfigs).To(Equal(expectedConfigs))
		})
	})

	Describe("Execute", func() {
		It("extracts watched logs", func() {
			extractor.ExtractLogsErrors = []error{nil, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
			Expect(extractor.ExtractLogsCount > 0).To(BeTrue())
		})

		It("returns error if extracting logs fails", func() {
			extractor.ExtractLogsErrors = []error{fakes.FakeError}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("retries on extractor error if watcher configured with greater than zero maximum consecutive errors", func() {
			eventWatcher.MaxConsecutiveUnexpectedErrs = 1
			extractor.ExtractLogsErrors = []error{fakes.FakeError, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
			Expect(extractor.ExtractLogsCount > 1).To(BeTrue())
		})

		It("returns error if maximum consecutive errors exceeded", func() {
			eventWatcher.MaxConsecutiveUnexpectedErrs = 1
			extractor.ExtractLogsErrors = []error{fakes.FakeError, fakes.FakeError}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("does not treat absence of unchecked headers as an unexpected error", func() {
			extractor.ExtractLogsErrors = []error{logs.ErrNoUncheckedHeaders, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
		})

		It("does not treat an io.ErrUnexpectedEOF error from the node as an unexpected error", func() {
			extractor.ExtractLogsErrors = []error{io.ErrUnexpectedEOF, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
		})

		It("extracts watched logs again if missing headers found", func() {
			extractor.ExtractLogsErrors = []error{nil, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
			Expect(extractor.ExtractLogsCount > 1).To(BeTrue())
		})

		It("returns error if extracting logs fails on subsequent run", func() {
			extractor.ExtractLogsErrors = []error{nil, fakes.FakeError}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("delegates untransformed logs", func() {
			delegator.DelegateErrors = []error{nil, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
			Expect(delegator.DelegateCallCount > 0).To(BeTrue())
		})

		It("passes results limit to delegator", func() {
			delegator.DelegateErrors = []error{nil, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
			Expect(delegator.DelegatePassedLimit).To(Equal(watcher.ResultsLimit))
		})

		It("returns error if delegating logs fails", func() {
			delegator.DelegateErrors = []error{fakes.FakeError}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("retries on delegator error if watcher configured with greater than zero maximum consecutive errors", func() {
			eventWatcher.MaxConsecutiveUnexpectedErrs = 1
			delegator.DelegateErrors = []error{fakes.FakeError, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
			Expect(delegator.DelegateCallCount > 1).To(BeTrue())
		})

		It("returns error if maximum consecutive errors exceeded", func() {
			eventWatcher.MaxConsecutiveUnexpectedErrs = 1
			delegator.DelegateErrors = []error{fakes.FakeError, fakes.FakeError}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("does not treat absence of unchecked logs as an unexpected error", func() {
			delegator.DelegateErrors = []error{logs.ErrNoLogs, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
		})

		It("delegates logs again if untransformed logs found", func() {
			delegator.DelegateErrors = []error{nil, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
			Expect(delegator.DelegateCallCount > 1).To(BeTrue())
		})

		It("returns error if delegating logs fails on subsequent run", func() {
			delegator.DelegateErrors = []error{nil, fakes.FakeError}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("doesn't panic if one of the go routines errors and closes the err channel", func() {
			extractor.ExtractLogsErrors = []error{nil, errExecuteClosed, errExecuteClosed}
			delegator.DelegateErrors = []error{nil, errExecuteClosed, errExecuteClosed}

			err := eventWatcher.Execute(constants.HeaderUnchecked)

			Expect(err).To(MatchError(errExecuteClosed))
			Expect(delegator.DelegateCallCount > 0 || extractor.ExtractLogsCount > 0).To(BeTrue())
		})
	})
})
