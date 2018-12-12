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

package shared

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Chunker interface {
	AddConfigs(transformerConfigs []TransformerConfig)
	ChunkLogs(logs []types.Log) map[string][]types.Log
}

type LogChunker struct {
	AddressToNames map[string][]string
	NameToTopic0   map[string]common.Hash
}

// Returns a new log chunker with initialised maps.
// Needs to have configs added with `AddConfigs` to consider logs for the respective transformer.
func NewLogChunker() *LogChunker {
	return &LogChunker{
		AddressToNames: map[string][]string{},
		NameToTopic0:   map[string]common.Hash{},
	}
}

// Configures the chunker by adding more addreses and topics to consider.
func (chunker *LogChunker) AddConfigs(transformerConfigs []TransformerConfig) {
	for _, config := range transformerConfigs {
		for _, address := range config.ContractAddresses {
			chunker.AddressToNames[address] = append(chunker.AddressToNames[address], config.TransformerName)
			chunker.NameToTopic0[config.TransformerName] = common.HexToHash(config.Topic)
		}
	}
}

// Goes through an array of logs, associating relevant logs (matching addresses and topic) with transformers
func (chunker *LogChunker) ChunkLogs(logs []types.Log) map[string][]types.Log {
	chunks := map[string][]types.Log{}
	for _, log := range logs {
		// Topic0 is not unique to each transformer, also need to consider the contract address
		relevantTransformers := chunker.AddressToNames[log.Address.String()]

		// TODO What should happen if log can't be assigned?
		for _, transformer := range relevantTransformers {
			if chunker.NameToTopic0[transformer] == log.Topics[0] {
				chunks[transformer] = append(chunks[transformer], log)
			}
		}
	}
	return chunks
}
