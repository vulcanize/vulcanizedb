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

package vat_heal_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_heal"
)

var _ = Describe("VatHeal converter", func() {
	It("Converts logs to models", func() {
		converter := vat_heal.VatHealConverter{}
		models, err := converter.ToModels([]types.Log{test_data.VatHealLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(models[0].Urn).To(Equal(test_data.VatHealModel.Urn))
		Expect(models[0].V).To(Equal(test_data.VatHealModel.V))
		Expect(models[0].Rad).To(Equal(test_data.VatHealModel.Rad))
		Expect(models[0].TransactionIndex).To(Equal(test_data.VatHealModel.TransactionIndex))
		Expect(models[0].Raw).To(Equal(test_data.VatHealModel.Raw))
	})

	It("Returns an error there are missing topics", func() {
		converter := vat_heal.VatHealConverter{}
		badLog := types.Log{
			Topics: []common.Hash{
				common.HexToHash("0x"),
				common.HexToHash("0x"),
				common.HexToHash("0x"),
			},
		}
		_, err := converter.ToModels([]types.Log{badLog})

		Expect(err).To(HaveOccurred())
	})
})
