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

package event_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Transformer", func() {
	var (
		repository mocks.MockEventRepository
		converter  mocks.MockConverter
		t          transformer.EventTransformer
		headerOne  core.Header
		config     = test_data.GenericTestConfig
		logs       []core.HeaderSyncLog
	)

	BeforeEach(func() {
		repository = mocks.MockEventRepository{}
		converter = mocks.MockConverter{}

		t = event.Transformer{
			Repository: &repository,
			Converter:  &converter,
			Config:     config,
		}.NewTransformer(nil)

		headerOne = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}

		logs = []core.HeaderSyncLog{{
			ID:          0,
			HeaderID:    headerOne.Id,
			Log:         test_data.GenericTestLog(),
			Transformed: false,
		}}
	})

	It("sets the db", func() {
		Expect(repository.SetDbCalled).To(BeTrue())
	})

	It("doesn't attempt to convert or persist an empty collection when there are no logs", func() {
		err := t.Execute([]core.HeaderSyncLog{})

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.ToEntitiesCalledCounter).To(Equal(0))
		Expect(converter.ToModelsCalledCounter).To(Equal(0))
		Expect(repository.CreateCalledCounter).To(Equal(0))
	})

	It("converts an eth log to an entity", func() {
		err := t.Execute(logs)

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.ContractAbi).To(Equal(config.ContractAbi))
		Expect(converter.LogsToConvert).To(Equal(logs))
	})

	It("returns an error if converter fails", func() {
		converter.ToEntitiesError = fakes.FakeError

		err := t.Execute(logs)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts an entity to a model", func() {
		converter.EntitiesToReturn = []interface{}{test_data.GenericEntity{}}

		err := t.Execute(logs)

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.EntitiesToConvert[0]).To(Equal(test_data.GenericEntity{}))
	})

	It("returns an error if converting to models fails", func() {
		converter.EntitiesToReturn = []interface{}{test_data.GenericEntity{}}
		converter.ToModelsError = fakes.FakeError

		err := t.Execute(logs)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists the record", func() {
		converter.ModelsToReturn = []interface{}{test_data.GenericModel{}}

		err := t.Execute(logs)

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedModels[0]).To(Equal(test_data.GenericModel{}))
	})

	It("returns error if persisting the record fails", func() {
		repository.SetCreateError(fakes.FakeError)
		err := t.Execute(logs)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
