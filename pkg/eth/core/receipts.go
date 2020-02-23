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

type Receipt struct {
	Bloom             string
	ContractAddress   string `db:"contract_address"`
	CumulativeGasUsed uint64 `db:"cumulative_gas_used"`
	GasUsed           uint64 `db:"gas_used"`
	Logs              []FullSyncLog
	StateRoot         string `db:"state_root"`
	Status            int
	TxHash            string `db:"tx_hash"`
	Rlp               []byte `db:"rlp"`
}
