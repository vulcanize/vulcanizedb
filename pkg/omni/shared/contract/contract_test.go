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

package contract_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

var _ = Describe("Contract", func() {
	var err error
	var info *contract.Contract
	var wantedEvents = []string{"Transfer", "Approval"}

	Describe("GenerateFilters", func() {

		It("Generates filters from contract data", func() {
			info = test_helpers.SetupTusdContract(wantedEvents, nil)
			err = info.GenerateFilters()
			Expect(err).ToNot(HaveOccurred())

			val, ok := info.Filters["Transfer"]
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal(mocks.ExpectedTransferFilter))

			val, ok = info.Filters["Approval"]
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal(mocks.ExpectedApprovalFilter))

			val, ok = info.Filters["Mint"]
			Expect(ok).To(Equal(false))

		})

		It("Fails with an empty contract", func() {
			info = &contract.Contract{}
			err = info.GenerateFilters()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("IsEventAddr", func() {

		BeforeEach(func() {
			info = &contract.Contract{}
			info.MethodArgs = map[string]bool{}
			info.FilterArgs = map[string]bool{}
		})

		It("Returns true if address is in event address filter list", func() {
			info.FilterArgs["testAddress1"] = true
			info.FilterArgs["testAddress2"] = true

			is := info.WantedEventArg("testAddress1")
			Expect(is).To(Equal(true))
			is = info.WantedEventArg("testAddress2")
			Expect(is).To(Equal(true))

			info.MethodArgs["testAddress3"] = true
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

	Describe("IsMethodAddr", func() {
		BeforeEach(func() {
			info = &contract.Contract{}
			info.MethodArgs = map[string]bool{}
			info.FilterArgs = map[string]bool{}
		})

		It("Returns true if address is in method address filter list", func() {
			info.MethodArgs["testAddress1"] = true
			info.MethodArgs["testAddress2"] = true

			is := info.WantedMethodArg("testAddress1")
			Expect(is).To(Equal(true))
			is = info.WantedMethodArg("testAddress2")
			Expect(is).To(Equal(true))

			info.FilterArgs["testAddress3"] = true
			is = info.WantedMethodArg("testAddress3")
			Expect(is).To(Equal(false))
		})

		It("Returns true if method address filter list is empty (no filter)", func() {
			is := info.WantedMethodArg("testAddress1")
			Expect(is).To(Equal(true))
			is = info.WantedMethodArg("testAddress2")
			Expect(is).To(Equal(true))
		})

		It("Returns false if address is not in method address filter list", func() {
			info.MethodArgs["testAddress1"] = true
			info.MethodArgs["testAddress2"] = true

			is := info.WantedMethodArg("testAddress3")
			Expect(is).To(Equal(false))
		})

		It("Returns false if method address filter list is nil (block all)", func() {
			info.MethodArgs = nil

			is := info.WantedMethodArg("testAddress1")
			Expect(is).To(Equal(false))
			is = info.WantedMethodArg("testAddress2")
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

	Describe("AddEmittedAddr", func() {
		BeforeEach(func() {
			info = &contract.Contract{}
			info.FilterArgs = map[string]bool{}
			info.MethodArgs = map[string]bool{}
			info.Methods = []types.Method{}
			info.EmittedAddrs = map[interface{}]bool{}
		})

		It("Adds address to list if it is on the method filter address list", func() {
			info.MethodArgs["testAddress2"] = true
			info.AddEmittedAddr("testAddress2")
			b := info.EmittedAddrs["testAddress2"]
			Expect(b).To(Equal(true))
		})

		It("Adds address to list if method filter is empty", func() {
			info.AddEmittedAddr("testAddress2")
			b := info.EmittedAddrs["testAddress2"]
			Expect(b).To(Equal(true))
		})

		It("Does not add address to list if both filters are closed (nil)", func() {
			info.FilterArgs = nil // close both
			info.MethodArgs = nil
			info.AddEmittedAddr("testAddress1")
			b := info.EmittedAddrs["testAddress1"]
			Expect(b).To(Equal(false))
		})
	})
})
