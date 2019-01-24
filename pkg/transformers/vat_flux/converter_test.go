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
		Expect(models[0]).To(Equal(test_data.VatFluxModel))
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
