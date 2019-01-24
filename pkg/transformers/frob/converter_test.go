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
