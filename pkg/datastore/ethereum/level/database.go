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
	"github.com/ethereum/go-ethereum/core/types"
)

type LevelDatabase struct {
	reader Reader
}

func NewLevelDatabase(ldbReader Reader) *LevelDatabase {
	return &LevelDatabase{
		reader: ldbReader,
	}
}

func (l LevelDatabase) GetBlock(blockHash []byte, blockNumber int64) *types.Block {
	n := uint64(blockNumber)
	h := common.BytesToHash(blockHash)
	return l.reader.GetBlock(h, n)
}

func (l LevelDatabase) GetBlockHash(blockNumber int64) []byte {
	n := uint64(blockNumber)
	h := l.reader.GetCanonicalHash(n)
	return h.Bytes()
}

func (l LevelDatabase) GetBlockReceipts(blockHash []byte, blockNumber int64) types.Receipts {
	n := uint64(blockNumber)
	h := common.BytesToHash(blockHash)
	return l.reader.GetBlockReceipts(h, n)
}

func (l LevelDatabase) GetHeadBlockNumber() int64 {
	h := l.reader.GetHeadBlockHash()
	n := l.reader.GetBlockNumber(h)
	return int64(*n)
}
