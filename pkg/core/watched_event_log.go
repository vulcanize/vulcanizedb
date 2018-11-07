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

package core

type WatchedEvent struct {
	LogID       int64  `json:"log_id" db:"id"`
	Name        string `json:"name"`
	BlockNumber int64  `json:"block_number" db:"block_number"`
	Address     string `json:"address"`
	TxHash      string `json:"tx_hash" db:"tx_hash"`
	Index       int64  `json:"index"`
	Topic0      string `json:"topic0"`
	Topic1      string `json:"topic1"`
	Topic2      string `json:"topic2"`
	Topic3      string `json:"topic3"`
	Data        string `json:"data"`
}
