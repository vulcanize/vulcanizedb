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

package frob_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Frob converter", func() {
	var converter = frob.FrobConverter{}
	It("converts a log to an entity", func() {
		entities, err := converter.ToEntities(test_data.KovanPitABI, []types.Log{test_data.EthFrobLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(entities)).To(Equal(1))
		Expect(entities[0]).To(Equal(test_data.FrobEntity))
	})

	It("returns an error if converting to an entity fails", func() {
		_, err := converter.ToEntities("bad abi", []types.Log{test_data.EthFrobLog})

		Expect(err).To(HaveOccurred())
	})

	It("converts an entity to a model", func() {
		models, err := converter.ToModels([]interface{}{test_data.FrobEntity})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(test_data.FrobModel))
	})

	It("returns an error if the entity type is wrong", func() {
		_, err := converter.ToModels([]interface{}{test_data.WrongEntity{}})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("entity of type test_data.WrongEntity, not frob.FrobEntity"))
	})
})
