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

package core

type Block struct {
	Reward       string `db:"reward"`
	Difficulty   int64  `db:"difficulty"`
	ExtraData    string `db:"extra_data"`
	GasLimit     uint64 `db:"gas_limit"`
	GasUsed      uint64 `db:"gas_used"`
	Hash         string `db:"hash"`
	IsFinal      bool   `db:"is_final"`
	Miner        string `db:"miner"`
	Nonce        string `db:"nonce"`
	Number       int64  `db:"number"`
	ParentHash   string `db:"parent_hash"`
	Size         string `db:"size"`
	Time         uint64 `db:"time"`
	Transactions []TransactionModel
	UncleHash    string `db:"uncle_hash"`
	UnclesReward string `db:"uncles_reward"`
	Uncles       []Uncle
}
