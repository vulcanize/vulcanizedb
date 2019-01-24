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

package filters_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

var _ = Describe("Log filters", func() {
	It("decodes web3 filter to LogFilter", func() {

		var logFilter filters.LogFilter
		jsonFilter := []byte(
			`{
      "name": "TestEvent",
      "fromBlock": "0x1",
      "toBlock": "0x488290",
	  "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
      "topics": ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null, "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null]
             }`)
		err := json.Unmarshal(jsonFilter, &logFilter)

		Expect(err).ToNot(HaveOccurred())
		Expect(logFilter.Name).To(Equal("TestEvent"))
		Expect(logFilter.FromBlock).To(Equal(int64(1)))
		Expect(logFilter.ToBlock).To(Equal(int64(4752016)))
		Expect(logFilter.Address).To(Equal("0x8888f1f195afa192cfee860698584c030f4c9db1"))
		Expect(logFilter.Topics).To(Equal(
			core.Topics{
				"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
				"",
				"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
				""}))
	})

	It("decodes array of web3 filters to  []LogFilter", func() {

		logFilters := make([]filters.LogFilter, 0)
		jsonFilter := []byte(
			`[{
      "name": "TestEvent",
      "fromBlock": "0x1",
      "toBlock": "0x488290",
	  "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
      "topics": ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null, "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null]
             },
      		{
	  "name": "TestEvent2",
      "fromBlock": "0x3",
      "toBlock": "0x4",
	  "address": "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07",
      "topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x0000000000000000000000006b0949d4c6edfe467db78241b7d5566f3c2bb43e", "0x0000000000000000000000005e44c3e467a49c9ca0296a9f130fc433041aaa28"]
             }]`)
		err := json.Unmarshal(jsonFilter, &logFilters)

		Expect(err).ToNot(HaveOccurred())
		Expect(len(logFilters)).To(Equal(2))
		Expect(logFilters[0].Name).To(Equal("TestEvent"))
		Expect(logFilters[1].Name).To(Equal("TestEvent2"))
	})

	It("requires valid ethereum address", func() {

		var logFilter filters.LogFilter
		jsonFilter := []byte(
			`{
      "name": "TestEvent",
      "fromBlock": "0x1",
      "toBlock": "0x2",
	  "address": "0x8888f1f195afa192cf84c030f4c9db1",
      "topics": ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null, "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null]
             }`)
		err := json.Unmarshal(jsonFilter, &logFilter)
		Expect(err).To(HaveOccurred())

	})
	It("requires name", func() {

		var logFilter filters.LogFilter
		jsonFilter := []byte(
			`{
      "fromBlock": "0x1",
      "toBlock": "0x2",
	  "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
      "topics": ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null, "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null]
             }`)
		err := json.Unmarshal(jsonFilter, &logFilter)
		Expect(err).To(HaveOccurred())

	})

	It("maps missing fromBlock to -1", func() {

		var logFilter filters.LogFilter
		jsonFilter := []byte(
			`{
      "name": "TestEvent",
      "toBlock": "0x2",
	  "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
      "topics": ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null, "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null]
             }`)
		err := json.Unmarshal(jsonFilter, &logFilter)
		Expect(err).ToNot(HaveOccurred())
		Expect(logFilter.FromBlock).To(Equal(int64(-1)))

	})

	It("maps missing toBlock to -1", func() {
		var logFilter filters.LogFilter
		jsonFilter := []byte(
			`{
      "name": "TestEvent",
	  "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
      "topics": ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null, "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null]
             }`)
		err := json.Unmarshal(jsonFilter, &logFilter)
		Expect(err).ToNot(HaveOccurred())
		Expect(logFilter.ToBlock).To(Equal(int64(-1)))

	})

})
