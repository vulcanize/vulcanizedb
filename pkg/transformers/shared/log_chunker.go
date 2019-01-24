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

package shared

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"strings"
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
			var lowerCaseAddress = strings.ToLower(address)
			chunker.AddressToNames[lowerCaseAddress] = append(chunker.AddressToNames[lowerCaseAddress], config.TransformerName)
			chunker.NameToTopic0[config.TransformerName] = common.HexToHash(config.Topic)
		}
	}
}

// Goes through an array of logs, associating relevant logs (matching addresses and topic) with transformers
func (chunker *LogChunker) ChunkLogs(logs []types.Log) map[string][]types.Log {
	chunks := map[string][]types.Log{}
	for _, log := range logs {
		// Topic0 is not unique to each transformer, also need to consider the contract address
		relevantTransformers := chunker.AddressToNames[strings.ToLower(log.Address.String())]

		for _, transformer := range relevantTransformers {
			if chunker.NameToTopic0[transformer] == log.Topics[0] {
				chunks[transformer] = append(chunks[transformer], log)
			}
		}
	}
	return chunks
}
