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

package vow_flog_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vow_flog"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Vow flog converter", func() {
	var converter vow_flog.VowFlogConverter
	BeforeEach(func() {
		converter = vow_flog.VowFlogConverter{}
	})

	It("returns err if log is missing topics", func() {
		badLog := types.Log{
			Data: []byte{1, 1, 1, 1, 1},
		}

		_, err := converter.ToModels([]types.Log{badLog})
		Expect(err).To(HaveOccurred())
	})

	It("returns err if log is missing data", func() {
		badLog := types.Log{
			Topics: []common.Hash{{}, {}, {}, {}},
		}

		_, err := converter.ToModels([]types.Log{badLog})
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a model", func() {
		models, err := converter.ToModels([]types.Log{test_data.EthVowFlogLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0].(vow_flog.VowFlogModel)).To(Equal(test_data.VowFlogModel))
	})
})
