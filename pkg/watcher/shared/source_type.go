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

package shared

import (
	"errors"
	"strings"
)

// SourceType enum for specifying source type for raw chain data
type SourceType int

const (
	Unknown SourceType = iota
	VulcanizeDB
	Ethereum
	Bitcoin
)

func (c SourceType) String() string {
	switch c {
	case Ethereum:
		return "Ethereum"
	case Bitcoin:
		return "Bitcoin"
	case VulcanizeDB:
		return "VulcanizeDB"
	default:
		return ""
	}
}

func NewSourceType(name string) (SourceType, error) {
	switch strings.ToLower(name) {
	case "ethereum", "eth":
		return Ethereum, nil
	case "bitcoin", "btc", "xbt":
		return Bitcoin, nil
	case "vulcanizedb", "vdb":
		return VulcanizeDB, nil
	default:
		return Unknown, errors.New("invalid name for data source")
	}
}
