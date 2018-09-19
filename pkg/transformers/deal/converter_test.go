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

package deal_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/deal"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Flip Deal Converter", func() {
	It("converts logs to models", func() {
		converter := deal.DealConverter{}

		models, err := converter.ToModels([]types.Log{test_data.DealLogNote})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(test_data.DealModel))
	})

	It("returns an error if the expected amount of topics aren't in the log", func() {
		converter := deal.DealConverter{}
		invalidLog := test_data.DealLogNote
		invalidLog.Topics = []common.Hash{}

		_, err := converter.ToModels([]types.Log{invalidLog})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("deal log does not contain expected topics"))
	})
})
