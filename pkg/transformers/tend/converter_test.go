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

package tend_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Tend TendConverter", func() {
	var converter tend.TendConverter

	BeforeEach(func() {
		converter = tend.NewTendConverter()
	})

	Describe("Convert", func() {
		It("converts an eth log to a db model", func() {
			model, err := converter.Convert(shared.FlipperContractAddress, shared.FlipperABI, test_data.TendLogNote)

			Expect(err).NotTo(HaveOccurred())
			Expect(model).To(Equal(test_data.TendModel))
		})

		It("returns an error if the log data is empty", func() {
			emptyDataLog := test_data.TendLogNote
			emptyDataLog.Data = []byte{}
			model, err := converter.Convert(shared.FlipperContractAddress, shared.FlipperABI, emptyDataLog)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("tend log note data is empty"))
			Expect(model).To(Equal(tend.TendModel{}))
		})

		It("returns an error if the expected amount of topics aren't in the log", func() {
			invalidLog := test_data.TendLogNote
			invalidLog.Topics = []common.Hash{}
			model, err := converter.Convert(shared.FlipperContractAddress, shared.FlipperABI, invalidLog)

			Expect(err).To(MatchError("tend log does not contain expected topics"))
			Expect(model).To(Equal(tend.TendModel{}))
		})
	})
})
