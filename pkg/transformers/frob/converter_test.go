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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Frob converter", func() {
	It("converts a log to an entity", func() {
		converter := frob.FrobConverter{}

		entities, err := converter.ToEntities(shared.PitABI, []types.Log{test_data.EthFrobLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(entities)).To(Equal(1))
		Expect(entities[0]).To(Equal(test_data.FrobEntity))
	})

	It("converts an entity to a model", func() {
		converter := frob.FrobConverter{}

		models, err := converter.ToModels([]frob.FrobEntity{test_data.FrobEntity})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(test_data.FrobModel))
	})
})
