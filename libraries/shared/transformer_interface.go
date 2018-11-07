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
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Transformer interface {
	Execute() error
}

type TransformerInitializer func(db *postgres.DB, blockchain core.BlockChain, con ContractConfig) (Transformer, error)

func HexToInt64(byteString string) int64 {
	value := common.HexToHash(byteString)
	return value.Big().Int64()
}

func HexToString(byteString string) string {
	value := common.HexToHash(byteString)
	return value.Big().String()
}
