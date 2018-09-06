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

package shared_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

var _ = Describe("Event signature generator", func() {
	Describe("generating non-anonymous event signatures", func() {
		It("generates bite event signature", func() {
			expected := "0x99b5620489b6ef926d4518936cfec15d305452712b88bd59da2d9c10fb0953e8"
			actual := shared.GetEventSignature("Bite(bytes32,bytes32,uint256,uint256,uint256,uint256,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates frob event signature", func() {
			expected := "0xb2afa28318bcc689926b52835d844de174ef8de97e982a85c0199d584920791b"
			actual := shared.GetEventSignature("Frob(bytes32,bytes32,uint256,uint256,int256,int256,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates flip kick event signature", func() {
			expected := "0xbac86238bdba81d21995024470425ecb370078fa62b7271b90cf28cbd1e3e87e"
			actual := shared.GetEventSignature("Kick(uint256,uint256,uint256,address,uint48,bytes32,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates log value event signature", func() {
			expected := "0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527"
			actual := shared.GetEventSignature("LogValue(bytes32)")

			Expect(expected).To(Equal(actual))
		})
	})

	Describe("generating LogNote event signatures", func() {
		It("generates flip tend event signature", func() {
			expected := "0x4b43ed1200000000000000000000000000000000000000000000000000000000"
			actual := shared.GetLogNoteSignature("tend(uint256,uint256,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates pit file event signature for overloaded function with three arguments", func() {
			expected := "0x1a0b287e00000000000000000000000000000000000000000000000000000000"
			actual := shared.GetLogNoteSignature("file(bytes32,bytes32,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates pit file event signature for overloaded function with two arguments", func() {
			expected := "0x29ae811400000000000000000000000000000000000000000000000000000000"
			actual := shared.GetLogNoteSignature("file(bytes32,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates pit file event signature for overloaded function with two different arguments", func() {
			expected := "0xd4e8be8300000000000000000000000000000000000000000000000000000000"
			actual := shared.GetLogNoteSignature("file(bytes32,address)")

			Expect(expected).To(Equal(actual))
		})
	})
})
