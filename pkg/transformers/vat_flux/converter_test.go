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

package vat_flux_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_flux"
)

var _ = Describe("VatFlux converter", func() {
	It("Converts logs to models", func() {
		converter := vat_flux.VatFluxConverter{}
		models, err := converter.ToModels([]types.Log{test_data.VatFluxLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0].Ilk).To(Equal(test_data.VatFluxModel.Ilk))
		Expect(models[0].Src).To(Equal(test_data.VatFluxModel.Src))
		Expect(models[0].Dst).To(Equal(test_data.VatFluxModel.Dst))
		Expect(models[0].Rad).To(Equal(test_data.VatFluxModel.Rad))
		Expect(models[0].TransactionIndex).To(Equal(test_data.VatFluxModel.TransactionIndex))
		Expect(models[0].Raw).To(Equal(test_data.VatFluxModel.Raw))
	})

	It("Returns an error there are missing topics", func() {
		converter := vat_flux.VatFluxConverter{}
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
