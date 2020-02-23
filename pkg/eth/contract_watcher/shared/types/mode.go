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

package types

// Mode is used to explicitly represent the operating mode of the transformer
type Mode int

// Mode enums
const (
	HeaderSync Mode = iota
	FullSync
)

// IsValid returns true is the Mode is valid
func (mode Mode) IsValid() bool {
	return mode >= HeaderSync && mode <= FullSync
}

// String returns the string representation of the mode
func (mode Mode) String() string {
	switch mode {
	case HeaderSync:
		return "header"
	case FullSync:
		return "full"
	default:
		return "unknown"
	}
}
