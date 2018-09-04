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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Tend TendConverter", func() {
	var converter tend.TendConverter

	BeforeEach(func() {
		converter = tend.TendConverter{}
	})

	Describe("Convert", func() {
		It("converts an eth log to a db model", func() {
			model, err := converter.Convert(shared.FlipperContractAddress, shared.FlipperABI, test_data.TendLogNote)
			Expect(err).NotTo(HaveOccurred())
			Expect(model).To(Equal(test_data.TendModel))
		})
	})
})
