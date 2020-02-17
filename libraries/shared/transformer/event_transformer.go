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

package transformer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type EventTransformer interface {
	Execute(logs []core.EventLog) error
	GetConfig() EventTransformerConfig
}

type EventTransformerInitializer func(db *postgres.DB) EventTransformer

type EventTransformerConfig struct {
	TransformerName     string
	ContractAddresses   []string
	ContractAbi         string
	Topic               string
	StartingBlockNumber int64
	EndingBlockNumber   int64 // Set -1 for indefinite transformer
}

func HexStringsToAddresses(strings []string) (addresses []common.Address) {
	for _, hexString := range strings {
		addresses = append(addresses, common.HexToAddress(hexString))
	}
	return
}
