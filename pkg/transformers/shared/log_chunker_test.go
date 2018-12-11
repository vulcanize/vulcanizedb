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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

var _ = Describe("Log chunker", func() {
	var (
		configs []shared.TransformerConfig
		chunker shared.LogChunker
	)

	BeforeEach(func() {
		configA := shared.TransformerConfig{
			TransformerName:   "TransformerA",
			ContractAddresses: []string{"0x00000000000000000000000000000000000000A1", "0x00000000000000000000000000000000000000A2"},
			Topic:             "0xA",
		}
		configB := shared.TransformerConfig{
			TransformerName:   "TransformerB",
			ContractAddresses: []string{"0x00000000000000000000000000000000000000B1"},
			Topic:             "0xB",
		}

		configC := shared.TransformerConfig{
			TransformerName:   "TransformerC",
			ContractAddresses: []string{"0x00000000000000000000000000000000000000A2"},
			Topic:             "0xC",
		}

		configs = []shared.TransformerConfig{configA, configB, configC}
		chunker = shared.NewLogChunker(configs)
	})

	Describe("initialisation", func() {
		It("creates lookup maps correctly", func() {
			Expect(chunker.AddressToNames).To(Equal(map[string][]string{
				"0x00000000000000000000000000000000000000A1": []string{"TransformerA"},
				"0x00000000000000000000000000000000000000A2": []string{"TransformerA", "TransformerC"},
				"0x00000000000000000000000000000000000000B1": []string{"TransformerB"},
			}))

			Expect(chunker.NameToTopic0).To(Equal(map[string]common.Hash{
				"TransformerA": common.HexToHash("0xA"),
				"TransformerB": common.HexToHash("0xB"),
				"TransformerC": common.HexToHash("0xC"),
			}))
		})
	})

	Describe("ChunkLogs", func() {
		It("only associates logs with relevant topic0 and address to transformers", func() {
			logs := []types.Log{log1, log2, log3, log4, log5}
			chunks := chunker.ChunkLogs(logs)

			Expect(chunks["TransformerA"]).To(And(ContainElement(log1), ContainElement(log4)))
			Expect(chunks["TransformerB"]).To(BeEmpty())
			Expect(chunks["TransformerC"]).To(ContainElement(log5))
		})
	})
})

var (
	// Match TransformerA
	log1 = types.Log{
		Address: common.HexToAddress("0xA1"),
		Topics: []common.Hash{
			common.HexToHash("0xA"),
			common.HexToHash("0xLogTopic1"),
		},
	}
	// Match TransformerA address, but not topic0
	log2 = types.Log{
		Address: common.HexToAddress("0xA1"),
		Topics: []common.Hash{
			common.HexToHash("0xB"),
			common.HexToHash("0xLogTopic2"),
		},
	}
	// Match TransformerA topic, but TransformerB address
	log3 = types.Log{
		Address: common.HexToAddress("0xB1"),
		Topics: []common.Hash{
			common.HexToHash("0xA"),
			common.HexToHash("0xLogTopic3"),
		},
	}
	// Match TransformerA, with the other address
	log4 = types.Log{
		Address: common.HexToAddress("0xA2"),
		Topics: []common.Hash{
			common.HexToHash("0xA"),
			common.HexToHash("0xLogTopic4"),
		},
	}
	// Match TransformerC, which shares address with TransformerA
	log5 = types.Log{
		Address: common.HexToAddress("0xA2"),
		Topics: []common.Hash{
			common.HexToHash("0xC"),
			common.HexToHash("0xLogTopic5"),
		},
	}
)
