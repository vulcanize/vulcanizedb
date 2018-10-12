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

package vat_move_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_move"
)

var _ = Describe("Vat move converter", func() {
	It("returns err if logs are missing topics", func() {
		converter := vat_move.VatMoveConverter{}
		badLog := types.Log{}

		_, err := converter.ToModels([]types.Log{badLog})

		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a model", func() {
		converter := vat_move.VatMoveConverter{}

		models, err := converter.ToModels([]types.Log{test_data.EthVatMoveLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(models[0]).To(Equal(test_data.VatMoveModel))
	})
})
