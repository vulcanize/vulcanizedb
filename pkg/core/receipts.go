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
	ContractAddress   string
	CumulativeGasUsed uint64
	GasUsed           uint64
	Logs              []Log
	StateRoot         string
	Status            int
	TxHash            string
}

type ReceiptModel struct {
	ContractAddress   string `db:"contract_address"`
	CumulativeGasUsed string `db:"cumulative_gas_used"`
	GasUsed           string `db:"gas_used"`
	StateRoot         string `db:"state_root"`
	Status            int
	TxHash            string `db:"tx_hash"`
}
