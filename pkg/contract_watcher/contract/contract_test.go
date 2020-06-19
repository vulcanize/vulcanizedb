// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package contract_test

import (
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/contract"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Contract", func() {
	var info *contract.Contract

	Describe("IsEventAddr", func() {
		BeforeEach(func() {
			info = &contract.Contract{}
			info.FilterArgs = map[string]bool{}
		})

		It("Returns true if address is in event address filter list", func() {
			info.FilterArgs["testAddress1"] = true
			info.FilterArgs["testAddress2"] = true

			is := info.WantedEventArg("testAddress1")
			Expect(is).To(Equal(true))
			is = info.WantedEventArg("testAddress2")
			Expect(is).To(Equal(true))
			is = info.WantedEventArg("testAddress3")
			Expect(is).To(Equal(false))
		})

		It("Returns true if event address filter is empty (no filter)", func() {
			is := info.WantedEventArg("testAddress1")
			Expect(is).To(Equal(true))
			is = info.WantedEventArg("testAddress2")
			Expect(is).To(Equal(true))
		})

		It("Returns false if address is not in event address filter list", func() {
			info.FilterArgs["testAddress1"] = true
			info.FilterArgs["testAddress2"] = true

			is := info.WantedEventArg("testAddress3")
			Expect(is).To(Equal(false))
		})

		It("Returns false if event address filter is nil (block all)", func() {
			info.FilterArgs = nil

			is := info.WantedEventArg("testAddress1")
			Expect(is).To(Equal(false))
			is = info.WantedEventArg("testAddress2")
			Expect(is).To(Equal(false))
		})
	})

	Describe("PassesEventFilter", func() {
		var mapping map[string]string
		BeforeEach(func() {
			info = &contract.Contract{}
			info.FilterArgs = map[string]bool{}
			mapping = map[string]string{}

		})

		It("Return true if event log name-value mapping has filtered for address as a value", func() {
			info.FilterArgs["testAddress1"] = true
			info.FilterArgs["testAddress2"] = true

			mapping["testInputName1"] = "testAddress1"
			mapping["testInputName2"] = "testAddress2"
			mapping["testInputName3"] = "testAddress3"

			pass := info.PassesEventFilter(mapping)
			Expect(pass).To(Equal(true))
		})

		It("Return true if event address filter list is empty (no filter)", func() {
			mapping["testInputName1"] = "testAddress1"
			mapping["testInputName2"] = "testAddress2"
			mapping["testInputName3"] = "testAddress3"

			pass := info.PassesEventFilter(mapping)
			Expect(pass).To(Equal(true))
		})

		It("Return false if event log name-value mapping does not have filtered for address as a value", func() {
			info.FilterArgs["testAddress1"] = true
			info.FilterArgs["testAddress2"] = true

			mapping["testInputName3"] = "testAddress3"

			pass := info.PassesEventFilter(mapping)
			Expect(pass).To(Equal(false))
		})

		It("Return false if event address filter list is nil (block all)", func() {
			info.FilterArgs = nil

			mapping["testInputName1"] = "testAddress1"
			mapping["testInputName2"] = "testAddress2"
			mapping["testInputName3"] = "testAddress3"

			pass := info.PassesEventFilter(mapping)
			Expect(pass).To(Equal(false))
		})
	})
})
