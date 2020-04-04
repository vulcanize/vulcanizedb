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

package eth

import "github.com/ethereum/go-ethereum/statediff"

func ResolveFromNodeType(nodeType statediff.NodeType) int {
	switch nodeType {
	case statediff.Branch:
		return 0
	case statediff.Extension:
		return 1
	case statediff.Leaf:
		return 2
	default:
		return -1
	}
}

func ResolveToNodeType(nodeType int) statediff.NodeType {
	switch nodeType {
	case 0:
		return statediff.Branch
	case 1:
		return statediff.Extension
	case 2:
		return statediff.Leaf
	default:
		return statediff.Unknown
	}
}
