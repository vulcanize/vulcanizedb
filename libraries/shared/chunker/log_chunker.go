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

package chunker

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type Chunker interface {
	AddConfig(transformerConfig transformer.EventTransformerConfig)
	ChunkLogs(logs []core.HeaderSyncLog) map[string][]core.HeaderSyncLog
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

// Configures the chunker by adding one config with more addresses and topics to consider.
func (chunker *LogChunker) AddConfig(transformerConfig transformer.EventTransformerConfig) {
	for _, address := range transformerConfig.ContractAddresses {
		var lowerCaseAddress = strings.ToLower(address)
		chunker.AddressToNames[lowerCaseAddress] = append(chunker.AddressToNames[lowerCaseAddress], transformerConfig.TransformerName)
		chunker.NameToTopic0[transformerConfig.TransformerName] = common.HexToHash(transformerConfig.Topic)
	}
}

// Goes through a slice of logs, associating relevant logs (matching addresses and topic) with transformers
func (chunker *LogChunker) ChunkLogs(logs []core.HeaderSyncLog) map[string][]core.HeaderSyncLog {
	chunks := map[string][]core.HeaderSyncLog{}
	for _, log := range logs {
		// Topic0 is not unique to each transformer, also need to consider the contract address
		relevantTransformers := chunker.AddressToNames[strings.ToLower(log.Log.Address.Hex())]

		for _, t := range relevantTransformers {
			if chunker.NameToTopic0[t] == log.Log.Topics[0] {
				chunks[t] = append(chunks[t], log)
			}
		}
	}
	return chunks
}
