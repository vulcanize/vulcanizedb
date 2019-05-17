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

package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/ethereum/level"
)

type Database interface {
	GetBlock(hash []byte, blockNumber int64) *types.Block
	GetBlockHash(blockNumber int64) []byte
	GetBlockReceipts(blockHash []byte, blockNumber int64) types.Receipts
	GetHeadBlockNumber() int64
}

func CreateDatabase(config DatabaseConfig) (Database, error) {
	switch config.Type {
	case Level:
		levelDBConnection, err := rawdb.NewLevelDBDatabase(config.Path, 128, 1024, "vdb")
		if err != nil {
			logrus.Error("CreateDatabase: error connecting to new LDBD: ", err)
			return nil, err
		}
		levelDBReader := level.NewLevelDatabaseReader(levelDBConnection)
		levelDB := level.NewLevelDatabase(levelDBReader)
		return levelDB, nil
	default:
		return nil, fmt.Errorf("Unknown ethereum database: %s", config.Path)
	}
}
