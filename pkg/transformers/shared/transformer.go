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

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Transformer interface {
	Execute() error
}

type TransformerInitializer func(db *postgres.DB, blockChain core.BlockChain) Transformer

type TransformerConfig struct {
	TransformerName     string
	ContractAddresses   []string
	ContractAbi         string
	Topic               string
	StartingBlockNumber int64
	EndingBlockNumber   int64 // Set -1 for indefinite transformer
}

func HexToInt64(byteString string) int64 {
	value := common.HexToHash(byteString)
	return value.Big().Int64()
}

func HexToString(byteString string) string {
	value := common.HexToHash(byteString)
	return value.Big().String()
}
