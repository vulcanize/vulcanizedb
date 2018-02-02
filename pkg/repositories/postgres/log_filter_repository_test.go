package postgres

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

var _ = Describe("Logs Repository", func() {
	var repository repositories.FilterRepository
	var node core.Node
	BeforeEach(func() {
		node = core.Node{
			GenesisBlock: "GENESIS",
			NetworkId:    1,
			Id:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
			ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
		}
		repository = BuildRepository(node)
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
			err := repository.AddFilter(logFilter)
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
			err := repository.AddFilter(logFilter)
			Expect(err).To(HaveOccurred())
		})
	})
})
