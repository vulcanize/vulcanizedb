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

package drip_drip_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_drip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Drip drip converter", func() {
	It("returns err if log is missing topics", func() {
		converter := drip_drip.DripDripConverter{}
		badLog := types.Log{}

		_, err := converter.ToModels([]types.Log{badLog})

		Expect(err).To(HaveOccurred())
	})

	It("converts a log to an model", func() {
		converter := drip_drip.DripDripConverter{}

		model, err := converter.ToModels([]types.Log{test_data.EthDripDripLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(model).To(Equal([]drip_drip.DripDripModel{test_data.DripDripModel}))
	})
})
