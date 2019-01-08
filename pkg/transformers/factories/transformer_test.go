// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package factories_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	"math/rand"
)

var _ = Describe("Transformer", func() {
	var (
		repository  mocks.MockRepository
		converter   mocks.MockConverter
		transformer shared.Transformer
		headerOne   core.Header
		config      = test_data.GenericTestConfig
		logs        = test_data.GenericTestLogs
	)

	BeforeEach(func() {
		repository = mocks.MockRepository{}
		converter = mocks.MockConverter{}

		transformer = factories.Transformer{
			Repository: &repository,
			Converter:  &converter,
			Config:     config,
		}.NewTransformer(nil)

		headerOne = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}
	})

	It("sets the db", func() {
		Expect(repository.SetDbCalled).To(BeTrue())
	})

	It("marks header checked if no logs returned", func() {
		err := transformer.Execute([]types.Log{}, headerOne)

		Expect(err).NotTo(HaveOccurred())
		repository.AssertMarkHeaderCheckedCalledWith(headerOne.Id)
	})

	It("doesn't attempt to convert or persist an empty collection when there are no logs", func() {
		err := transformer.Execute([]types.Log{}, headerOne)

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.ToEntitiesCalledCounter).To(Equal(0))
		Expect(converter.ToModelsCalledCounter).To(Equal(0))
		Expect(repository.CreateCalledCounter).To(Equal(0))
	})

	It("does not call repository.MarkCheckedHeader when there are logs", func() {
		err := transformer.Execute(logs, headerOne)

		Expect(err).NotTo(HaveOccurred())
		repository.AssertMarkHeaderCheckedNotCalled()
	})

	It("returns error if marking header checked returns err", func() {
		repository.SetMarkHeaderCheckedError(fakes.FakeError)

		err := transformer.Execute([]types.Log{}, headerOne)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts an eth log to an entity", func() {
		err := transformer.Execute(logs, headerOne)

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.ContractAbi).To(Equal(config.ContractAbi))
		Expect(converter.LogsToConvert).To(Equal(logs))
	})

	It("returns an error if converter fails", func() {
		converter.ToEntitiesError = fakes.FakeError

		err := transformer.Execute(logs, headerOne)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts an entity to a model", func() {
		converter.EntitiesToReturn = []interface{}{test_data.GenericEntity{}}

		err := transformer.Execute(logs, headerOne)

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.EntitiesToConvert[0]).To(Equal(test_data.GenericEntity{}))
	})

	It("returns an error if converting to models fails", func() {
		converter.EntitiesToReturn = []interface{}{test_data.GenericEntity{}}
		converter.ToModelsError = fakes.FakeError

		err := transformer.Execute(logs, headerOne)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists the record", func() {
		converter.ModelsToReturn = []interface{}{test_data.GenericModel{}}

		err := transformer.Execute(logs, headerOne)

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(headerOne.Id))
		Expect(repository.PassedModels[0]).To(Equal(test_data.GenericModel{}))
	})

	It("returns error if persisting the record fails", func() {
		repository.SetCreateError(fakes.FakeError)
		err := transformer.Execute(logs, headerOne)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
