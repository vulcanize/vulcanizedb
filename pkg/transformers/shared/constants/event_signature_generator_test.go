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

package constants_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Event signature generator", func() {
	Describe("generating non-anonymous event signatures", func() {
		It("generates bite event signature", func() {
			expected := "0x99b5620489b6ef926d4518936cfec15d305452712b88bd59da2d9c10fb0953e8"
			actual := constants.GetEventSignature("Bite(bytes32,bytes32,uint256,uint256,uint256,uint256,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates frob event signature", func() {
			expected := "0xb2afa28318bcc689926b52835d844de174ef8de97e982a85c0199d584920791b"
			actual := constants.GetEventSignature("Frob(bytes32,bytes32,uint256,uint256,int256,int256,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates the flap kick event signature", func() {
			expected := "0xefa52d9342a199cb30efd2692463f2c2bef63cd7186b50382d4fb94ad207880e"
			actual := constants.GetEventSignature("Kick(uint256,uint256,uint256,address,uint48)")

			Expect(expected).To(Equal(actual))
		})

		It("generates flip kick event signature", func() {
			expected := "0xbac86238bdba81d21995024470425ecb370078fa62b7271b90cf28cbd1e3e87e"
			actual := constants.GetEventSignature("Kick(uint256,uint256,uint256,address,uint48,bytes32,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates log value event signature", func() {
			expected := "0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527"
			actual := constants.GetEventSignature("LogValue(bytes32)")

			Expect(expected).To(Equal(actual))
		})
	})

	Describe("generating LogNote event signatures", func() {
		It("generates flip tend event signature", func() {
			expected := "0x4b43ed1200000000000000000000000000000000000000000000000000000000"
			actual := constants.GetLogNoteSignature("tend(uint256,uint256,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates pit file event signature for overloaded function with three arguments", func() {
			expected := "0x1a0b287e00000000000000000000000000000000000000000000000000000000"
			actual := constants.GetLogNoteSignature("file(bytes32,bytes32,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates pit file event signature for overloaded function with two arguments", func() {
			expected := "0x29ae811400000000000000000000000000000000000000000000000000000000"
			actual := constants.GetLogNoteSignature("file(bytes32,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates pit file event signature for overloaded function with two different arguments", func() {
			expected := "0xd4e8be8300000000000000000000000000000000000000000000000000000000"
			actual := constants.GetLogNoteSignature("file(bytes32,address)")

			Expect(expected).To(Equal(actual))
		})
	})

	Describe("getting the solidity method/event signature from the abi", func() {
		Describe("it handles methods", func() {
			Describe("from the cat contract", func() {
				It("gets the file method signature", func() {
					expected := "file(bytes32,bytes32,address)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanCatABI, "file")

					Expect(expected).To(Equal(actual))
				})
			})

			Describe("from the drip contract", func() {
				It("gets the drip method signature", func() {
					expected := "drip(bytes32)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanDripABI, "drip")

					Expect(expected).To(Equal(actual))
				})

				It("gets the file method signature", func() {
					expected := "file(bytes32,uint256)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanDripABI, "file")

					Expect(expected).To(Equal(actual))
				})
			})

			Describe("from the flipper contract", func() {
				It("gets the deal method signature", func() {
					expected := "deal(uint256)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanFlipperABI, "deal")

					Expect(expected).To(Equal(actual))
				})

				It("gets the dent method signature", func() {
					expected := "dent(uint256,uint256,uint256)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanFlipperABI, "dent")

					Expect(expected).To(Equal(actual))
				})

				It("gets the tend method signature", func() {
					expected := "tend(uint256,uint256,uint256)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanFlipperABI, "tend")

					Expect(expected).To(Equal(actual))
				})
			})

			Describe("from the pit contract", func() {
				It("gets the file method signature", func() {
					expected := "file(bytes32,address)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanPitABI, "file")

					Expect(expected).To(Equal(actual))
				})
			})

			Describe("from the vat contract", func() {
				It("gets the init method signature", func() {
					expected := "init(bytes32)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanVatABI, "init")

					Expect(expected).To(Equal(actual))
				})

				It("gets the heal method signature", func() {
					expected := "heal(bytes32,bytes32,int256)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanVatABI, "heal")

					Expect(expected).To(Equal(actual))
				})

				It("gets the fold method signature", func() {
					expected := "fold(bytes32,bytes32,int256)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanVatABI, "fold")

					Expect(expected).To(Equal(actual))
				})

				It("gets the flux method signature", func() {
					expected := "flux(bytes32,bytes32,bytes32,int256)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanVatABI, "flux")

					Expect(expected).To(Equal(actual))
				})

				It("gets the grab method signature", func() {
					expected := "grab(bytes32,bytes32,bytes32,bytes32,int256,int256)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanVatABI, "grab")

					Expect(expected).To(Equal(actual))
				})

				It("gets the move method signature", func() {
					expected := "move(bytes32,bytes32,int256)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanVatABI, "move")

					Expect(expected).To(Equal(actual))
				})

				It("gets the slip method signature", func() {
					expected := "slip(bytes32,bytes32,int256)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanVatABI, "slip")

					Expect(expected).To(Equal(actual))
				})

				It("gets the toll method signature", func() {
					expected := "toll(bytes32,bytes32,int256)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanVatABI, "toll")

					Expect(expected).To(Equal(actual))
				})

				It("gets the tune method signature", func() {
					expected := "tune(bytes32,bytes32,bytes32,bytes32,int256,int256)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanVatABI, "tune")

					Expect(expected).To(Equal(actual))
				})
			})

			Describe("from the vow contract", func() {
				It("gets the flog method signature", func() {
					expected := "flog(uint48)"
					actual := constants.GetSolidityMethodSignature(test_data.KovanVowABI, "flog")

					Expect(expected).To(Equal(actual))
				})
			})
		})

		Describe("it handles events", func() {
			It("gets the Bite event signature", func() {
				expected := "Bite(bytes32,bytes32,uint256,uint256,uint256,uint256,uint256)"
				actual := constants.GetSolidityMethodSignature(test_data.KovanCatABI, "Bite")

				Expect(expected).To(Equal(actual))
			})

			It("gets the flap Kick event signature", func() {
				expected := "Kick(uint256,uint256,uint256,address,uint48)"
				actual := constants.GetSolidityMethodSignature(test_data.KovanFlapperABI, "Kick")

				Expect(expected).To(Equal(actual))
			})

			It("gets the flip Kick event signature", func() {
				expected := "Kick(uint256,uint256,uint256,address,uint48,bytes32,uint256)"
				actual := constants.GetSolidityMethodSignature(test_data.KovanFlipperABI, "Kick")

				Expect(expected).To(Equal(actual))
			})

			It("gets the flop Kick event signature", func() {
				expected := "Kick(uint256,uint256,uint256,address,uint48)"
				actual := constants.GetSolidityMethodSignature(test_data.KovanFlopperABI, "Kick")

				Expect(expected).To(Equal(actual))
			})

			It("gets the pit frob event signature", func() {
				expected := "Frob(bytes32,bytes32,uint256,uint256,int256,int256,uint256)"
				actual := constants.GetSolidityMethodSignature(test_data.KovanPitABI, "Frob")

				Expect(expected).To(Equal(actual))
			})

			It("gets the log value method signature", func() {
				expected := "LogValue(bytes32)"
				actual := constants.GetSolidityMethodSignature(test_data.KovanMedianizerABI, "LogValue")

				Expect(expected).To(Equal(actual))
			})
		})
	})
})
