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

package ilk_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Drip file ilk converter", func() {
	It("returns err if log missing topics", func() {
		converter := ilk.DripFileIlkConverter{}
		badLog := types.Log{
			Topics: []common.Hash{{}},
			Data:   []byte{1, 1, 1, 1, 1},
		}

		_, err := converter.ToModel(badLog)

		Expect(err).To(HaveOccurred())
	})

	It("returns err if log missing data", func() {
		converter := ilk.DripFileIlkConverter{}
		badLog := types.Log{
			Topics: []common.Hash{{}, {}, {}, {}},
		}

		_, err := converter.ToModel(badLog)

		Expect(err).To(HaveOccurred())

	})

	It("converts a log to a model", func() {
		converter := ilk.DripFileIlkConverter{}

		model, err := converter.ToModel(test_data.EthDripFileIlkLog)

		Expect(err).NotTo(HaveOccurred())
		Expect(model).To(Equal(test_data.DripFileIlkModel))
	})
})
