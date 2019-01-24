// VulcanizeDB
// Copyright © 2018 Vulcanize

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

var _ = Describe("LogNoteTransformer", func() {
	var (
		repository  mocks.MockRepository
		converter   mocks.MockLogNoteConverter
		headerOne   core.Header
		transformer shared.Transformer
		model       test_data.GenericModel
		config      = test_data.GenericTestConfig
		logs        = test_data.GenericTestLogs
	)

	BeforeEach(func() {
		repository = mocks.MockRepository{}
		converter = mocks.MockLogNoteConverter{}
		transformer = factories.LogNoteTransformer{
			Config:     config,
			Converter:  &converter,
			Repository: &repository,
		}.NewLogNoteTransformer(nil)

		headerOne = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}
	})

	It("sets the database", func() {
		Expect(repository.SetDbCalled).To(BeTrue())
	})

	It("marks header checked if no logs are provided", func() {
		err := transformer.Execute([]types.Log{}, headerOne)

		Expect(err).NotTo(HaveOccurred())
		repository.AssertMarkHeaderCheckedCalledWith(headerOne.Id)
	})

	It("doesn't attempt to convert or persist an empty collection when there are no logs", func() {
		err := transformer.Execute([]types.Log{}, headerOne)

		Expect(err).NotTo(HaveOccurred())
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

	It("converts matching logs to models", func() {
		err := transformer.Execute(logs, headerOne)

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal(logs))
	})

	It("returns error if converter returns error", func() {
		converter.SetConverterError(fakes.FakeError)

		err := transformer.Execute(logs, headerOne)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists the model", func() {
		converter.SetReturnModels([]interface{}{model})
		err := transformer.Execute(logs, headerOne)

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(headerOne.Id))
		Expect(repository.PassedModels).To(Equal([]interface{}{model}))
	})

	It("returns error if repository returns error for create", func() {
		repository.SetCreateError(fakes.FakeError)

		err := transformer.Execute(logs, headerOne)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
