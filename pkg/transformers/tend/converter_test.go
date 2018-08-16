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

package tend_test

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Tend TendConverter", func() {
	var converter tend.TendConverter
	var emptyEntity tend.TendEntity
	var testEntity tend.TendEntity

	BeforeEach(func() {
		converter = tend.TendConverter{}
		emptyEntity = tend.TendEntity{}
		testEntity = test_data.TendEntity
	})

	Describe("ToEntity", func() {
		It("converts a log to an entity", func() {
			entity, err := converter.ToEntity(test_data.FlipAddress, shared.FlipperABI, test_data.TendLog)

			Expect(err).NotTo(HaveOccurred())
			Expect(entity).To(Equal(testEntity))
		})

		It("returns an error if there is a failure in parsing the abi", func() {
			malformedAbi := "bad"
			entity, err := converter.ToEntity(test_data.FlipAddress, malformedAbi, test_data.TendLog)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid abi"))
			Expect(entity).To(Equal(emptyEntity))
		})

		It("returns an error if there is a failure unpacking the log", func() {
			incompleteAbi := "[{}]"
			entity, err := converter.ToEntity(test_data.FlipAddress, incompleteAbi, test_data.TendLog)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("abi: could not locate"))
			Expect(entity).To(Equal(emptyEntity))
		})
	})

	Describe("ToModel", func() {
		It("converts an entity to a model", func() {
			model, err := converter.ToModel(testEntity)

			Expect(err).NotTo(HaveOccurred())
			Expect(model).To(Equal(test_data.TendModel))
		})

		It("handles nil values", func() {
			emptyEntity.Id = big.NewInt(1)
			emptyLog, err := json.Marshal(types.Log{})
			Expect(err).NotTo(HaveOccurred())
			expectedModel := tend.TendModel{
				Id:  "1",
				Lot: "",
				Bid: "",
				Guy: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Tic: "",
				Era: time.Unix(0, 0),
				Raw: string(emptyLog),
			}
			model, err := converter.ToModel(emptyEntity)

			Expect(err).NotTo(HaveOccurred())
			Expect(model).To(Equal(expectedModel))
		})

		It("returns an error if the log Id is nil", func() {
			model, err := converter.ToModel(emptyEntity)

			Expect(err).To(HaveOccurred())
			Expect(model).To(Equal(tend.TendModel{}))
		})
	})
})
