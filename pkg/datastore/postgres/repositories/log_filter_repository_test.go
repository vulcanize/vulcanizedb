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

package repositories_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Log Filters Repository", func() {
	var db *postgres.DB
	var filterRepository datastore.FilterRepository
	var node core.Node
	BeforeEach(func() {
		node = core.Node{
			GenesisBlock: "GENESIS",
			NetworkID:    1,
			ID:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
			ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
		}
		db = test_config.NewTestDB(node)
		test_config.CleanTestDB(db)
		filterRepository = repositories.FilterRepository{DB: db}
	})

	Describe("LogFilter", func() {

		It("inserts filter into watched events", func() {

			logFilter := filters.LogFilter{
				Name:      "TestFilter",
				FromBlock: 1,
				ToBlock:   2,
				Address:   "0x8888f1f195afa192cfee860698584c030f4c9db1",
				Topics: core.Topics{
					"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
					"",
					"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
					"",
				},
			}
			err := filterRepository.CreateFilter(logFilter)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if name is not provided", func() {

			logFilter := filters.LogFilter{
				FromBlock: 1,
				ToBlock:   2,
				Address:   "0x8888f1f195afa192cfee860698584c030f4c9db1",
				Topics: core.Topics{
					"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
					"",
					"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
					"",
				},
			}
			err := filterRepository.CreateFilter(logFilter)
			Expect(err).To(HaveOccurred())
		})

		It("gets a log filter", func() {

			expectedLogFilter1 := filters.LogFilter{
				Name:      "TestFilter1",
				FromBlock: 1,
				ToBlock:   2,
				Address:   "0x8888f1f195afa192cfee860698584c030f4c9db1",
				Topics: core.Topics{
					"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
					"",
					"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
					"",
				},
			}
			err := filterRepository.CreateFilter(expectedLogFilter1)
			Expect(err).ToNot(HaveOccurred())
			expectedLogFilter2 := filters.LogFilter{
				Name:      "TestFilter2",
				FromBlock: 10,
				ToBlock:   20,
				Address:   "0x8888f1f195afa192cfee860698584c030f4c9db1",
				Topics: core.Topics{
					"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
					"",
					"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
					"",
				},
			}
			err = filterRepository.CreateFilter(expectedLogFilter2)
			Expect(err).ToNot(HaveOccurred())

			logFilter1, err := filterRepository.GetFilter("TestFilter1")
			Expect(err).ToNot(HaveOccurred())
			Expect(logFilter1).To(Equal(expectedLogFilter1))
			logFilter2, err := filterRepository.GetFilter("TestFilter2")
			Expect(err).ToNot(HaveOccurred())
			Expect(logFilter2).To(Equal(expectedLogFilter2))
		})

		It("returns ErrFilterDoesNotExist error when log does not exist", func() {
			_, err := filterRepository.GetFilter("TestFilter1")
			Expect(err).To(Equal(datastore.ErrFilterDoesNotExist("TestFilter1")))
		})
	})
})
