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

package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
)

type Reader interface {
	GetBlock(hash common.Hash, number uint64) *types.Block
	GetBlockNumber(hash common.Hash) *uint64
	GetBlockReceipts(hash common.Hash, number uint64) types.Receipts
	GetCanonicalHash(number uint64) common.Hash
	GetHeadBlockHash() common.Hash
}

type LevelDatabaseReader struct {
	reader rawdb.DatabaseReader
}

func NewLevelDatabaseReader(reader rawdb.DatabaseReader) *LevelDatabaseReader {
	return &LevelDatabaseReader{reader: reader}
}

func (ldbr *LevelDatabaseReader) GetBlock(hash common.Hash, number uint64) *types.Block {
	return rawdb.ReadBlock(ldbr.reader, hash, number)
}

func (ldbr *LevelDatabaseReader) GetBlockNumber(hash common.Hash) *uint64 {
	return rawdb.ReadHeaderNumber(ldbr.reader, hash)
}

func (ldbr *LevelDatabaseReader) GetBlockReceipts(hash common.Hash, number uint64) types.Receipts {
	return rawdb.ReadReceipts(ldbr.reader, hash, number)
}

func (ldbr *LevelDatabaseReader) GetCanonicalHash(number uint64) common.Hash {
	return rawdb.ReadCanonicalHash(ldbr.reader, number)
}

func (ldbr *LevelDatabaseReader) GetHeadBlockHash() common.Hash {
	return rawdb.ReadHeadBlockHash(ldbr.reader)
}
