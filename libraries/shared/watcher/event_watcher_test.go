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
		eventWatcher watcher.EventWatcher
	)

	BeforeEach(func() {
		delegator = &mocks.MockLogDelegator{}
		extractor = &mocks.MockLogExtractor{}
		eventWatcher = watcher.EventWatcher{
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
		It("extracts watched logs", func() {
			err := eventWatcher.Execute(constants.HeaderMissing)

			Expect(err).NotTo(HaveOccurred())
			Expect(extractor.ExtractLogsCalled).To(BeTrue())
		})

		It("returns error if extracting logs fails", func() {
			extractor.ExtractLogsError = fakes.FakeError

			err := eventWatcher.Execute(constants.HeaderMissing)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("delegates untransformed logs", func() {
			err := eventWatcher.Execute(constants.HeaderMissing)

			Expect(err).NotTo(HaveOccurred())
			Expect(delegator.DelegateCalled).To(BeTrue())
		})

		It("returns error if delegating logs fails", func() {
			delegator.DelegateError = fakes.FakeError

			err := eventWatcher.Execute(constants.HeaderMissing)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})
