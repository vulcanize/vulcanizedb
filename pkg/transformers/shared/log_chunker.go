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

type LogChunker struct {
	addressToNames map[string][]string
	nameToTopic0   map[string]common.Hash
}

// Initialises a chunker by creating efficient lookup maps
func NewLogChunker(transformerConfigs []TransformerConfig) LogChunker {
	addressToNames := map[string][]string{}
	nameToTopic0 := map[string]common.Hash{}

	for _, config := range transformerConfigs {
		for _, address := range config.ContractAddresses {
			addressToNames[address] = append(addressToNames[address], config.TransformerName)
			nameToTopic0[config.TransformerName] = common.HexToHash(config.Topic)
		}
	}

	return LogChunker{
		addressToNames,
		nameToTopic0,
	}
}

// Goes through an array of logs, associating relevant logs with transformers
func (chunker LogChunker) ChunkLogs(logs []types.Log) (chunks map[string][]types.Log) {
	for _, log := range logs {
		// Topic0 is not unique to each transformer, also need to consider the contract address
		relevantTransformers := chunker.addressToNames[log.Address.String()]
		for _, transformer := range relevantTransformers {
			if chunker.nameToTopic0[transformer] == log.Topics[0] {
				chunks[transformer] = append(chunks[transformer], log)
			}
		}
	}
	return
}
