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

package chunker_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	chunk "github.com/vulcanize/vulcanizedb/libraries/shared/chunker"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var _ = Describe("Log chunker", func() {
	var (
		chunker *chunk.LogChunker
	)

	BeforeEach(func() {
		chunker = chunk.NewLogChunker()

		configA := transformer.EventTransformerConfig{
			TransformerName:   "TransformerA",
			ContractAddresses: []string{"0x00000000000000000000000000000000000000A1", "0x00000000000000000000000000000000000000A2"},
			Topic:             "0xA",
		}
		chunker.AddConfig(configA)

		configB := transformer.EventTransformerConfig{
			TransformerName:   "TransformerB",
			ContractAddresses: []string{"0x00000000000000000000000000000000000000B1"},
			Topic:             "0xB",
		}
		chunker.AddConfig(configB)

		configC := transformer.EventTransformerConfig{
			TransformerName:   "TransformerC",
			ContractAddresses: []string{"0x00000000000000000000000000000000000000A2"},
			Topic:             "0xC",
		}
		chunker.AddConfig(configC)
	})

	Describe("initialisation", func() {
		It("creates lookup maps correctly", func() {
			Expect(chunker.AddressToNames).To(Equal(map[string][]string{
				"0x00000000000000000000000000000000000000a1": {"TransformerA"},
				"0x00000000000000000000000000000000000000a2": {"TransformerA", "TransformerC"},
				"0x00000000000000000000000000000000000000b1": {"TransformerB"},
			}))

			Expect(chunker.NameToTopic0).To(Equal(map[string]common.Hash{
				"TransformerA": common.HexToHash("0xA"),
				"TransformerB": common.HexToHash("0xB"),
				"TransformerC": common.HexToHash("0xC"),
			}))
		})
	})

	Describe("AddConfig", func() {
		It("can add more configs later", func() {
			configD := transformer.EventTransformerConfig{
				TransformerName:   "TransformerD",
				ContractAddresses: []string{"0x000000000000000000000000000000000000000D"},
				Topic:             "0xD",
			}
			chunker.AddConfig(configD)

			Expect(chunker.AddressToNames).To(ContainElement([]string{"TransformerD"}))
			Expect(chunker.NameToTopic0).To(ContainElement(common.HexToHash("0xD")))
		})

		It("lower cases address", func() {
			configD := transformer.EventTransformerConfig{
				TransformerName:   "TransformerD",
				ContractAddresses: []string{"0x000000000000000000000000000000000000000D"},
				Topic:             "0xD",
			}
			chunker.AddConfig(configD)

			Expect(chunker.AddressToNames["0x000000000000000000000000000000000000000d"]).To(Equal([]string{"TransformerD"}))
		})
	})

	Describe("ChunkLogs", func() {
		It("only associates logs with relevant topic0 and address to transformers", func() {
			logs := []core.HeaderSyncLog{log1, log2, log3, log4, log5}
			chunks := chunker.ChunkLogs(logs)

			Expect(chunks["TransformerA"]).To(And(ContainElement(log1), ContainElement(log4)))
			Expect(chunks["TransformerB"]).To(BeEmpty())
			Expect(chunks["TransformerC"]).To(ContainElement(log5))
		})
	})
})

var (
	// Match TransformerA
	log1 = core.HeaderSyncLog{
		Log: types.Log{
			Address: common.HexToAddress("0xA1"),
			Topics: []common.Hash{
				common.HexToHash("0xA"),
				common.HexToHash("0xLogTopic1"),
			},
		},
	}
	// Match TransformerA address, but not topic0
	log2 = core.HeaderSyncLog{
		Log: types.Log{
			Address: common.HexToAddress("0xA1"),
			Topics: []common.Hash{
				common.HexToHash("0xB"),
				common.HexToHash("0xLogTopic2"),
			},
		},
	}
	// Match TransformerA topic, but TransformerB address
	log3 = core.HeaderSyncLog{
		Log: types.Log{
			Address: common.HexToAddress("0xB1"),
			Topics: []common.Hash{
				common.HexToHash("0xA"),
				common.HexToHash("0xLogTopic3"),
			},
		},
	}
	// Match TransformerA, with the other address
	log4 = core.HeaderSyncLog{
		Log: types.Log{
			Address: common.HexToAddress("0xA2"),
			Topics: []common.Hash{
				common.HexToHash("0xA"),
				common.HexToHash("0xLogTopic4"),
			},
		},
	}
	// Match TransformerC, which shares address with TransformerA
	log5 = core.HeaderSyncLog{
		Log: types.Log{
			Address: common.HexToAddress("0xA2"),
			Topics: []common.Hash{
				common.HexToHash("0xC"),
				common.HexToHash("0xLogTopic5"),
			},
		},
	}
)
