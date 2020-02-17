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
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/mocks"
	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/rand"
)

var _ = Describe("Transformer", func() {
	var (
		converter mocks.MockConverter
		t         transformer.EventTransformer
		headerOne core.Header
		config    = test_data.GenericTestConfig
		logs      []core.EventLog
	)

	BeforeEach(func() {
		converter = mocks.MockConverter{}

		t = event.ConfiguredTransformer{
			Transformer: &converter,
			Config:      config,
		}.NewTransformer(nil)

		headerOne = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}

		logs = []core.EventLog{{
			ID:          0,
			HeaderID:    headerOne.Id,
			Log:         test_data.GenericTestLog(),
			Transformed: false,
		}}
	})

	It("doesn't attempt to convert or persist an empty collection when there are no logs", func() {
		err := t.Execute([]core.EventLog{})

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.ToModelsCalledCounter).To(Equal(0))
	})

	It("converts an eth log to a model", func() {
		err := t.Execute(logs)

		// TODO Mock DB in repo instead?
		Expect(err).To(MatchError(event.ErrEmptyModelSlice))
		Expect(converter.ContractAbi).To(Equal(config.ContractAbi))
		Expect(converter.LogsToConvert).To(Equal(logs))
	})

	It("returns an error if converting to models fails", func() {
		converter.ToModelsError = fakes.FakeError

		err := t.Execute(logs)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
