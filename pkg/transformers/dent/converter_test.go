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

package dent_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Dent Converter", func() {
	var converter dent.DentConverter

	BeforeEach(func() {
		converter = dent.NewDentConverter()
	})

	It("converts an eth log to a db model", func() {
		models, err := converter.ToModels([]types.Log{test_data.DentLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0].(dent.DentModel)).To(Equal(test_data.DentModel))
	})

	It("returns an error if the expected amount of topics aren't in the log", func() {
		invalidLog := test_data.DentLog
		invalidLog.Topics = []common.Hash{}
		_, err := converter.ToModels([]types.Log{invalidLog})
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("dent log does not contain expected topics"))
	})

	It("returns an error if the log data is empty", func() {
		emptyDataLog := test_data.DentLog
		emptyDataLog.Data = []byte{}
		_, err := converter.ToModels([]types.Log{emptyDataLog})
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("dent log data is empty"))
	})
})
