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
	"bytes"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// ListContainsString used to check if a list of strings contains a particular string
func ListContainsString(sss []string, s string) bool {
	for _, str := range sss {
		if s == str {
			return true
		}
	}
	return false
}

// IPLDsContainBytes used to check if a list of strings contains a particular string
func IPLDsContainBytes(iplds []ipfs.BlockModel, b []byte) bool {
	for _, ipld := range iplds {
		if bytes.Equal(ipld.Data, b) {
			return true
		}
	}
	return false
}

// ListContainsGap used to check if a list of Gaps contains a particular Gap
func ListContainsGap(gapList []Gap, gap Gap) bool {
	for _, listGap := range gapList {
		if listGap == gap {
			return true
		}
	}
	return false
}

// HandleNullAddr converts a nil pointer to an address to a zero-valued hex address string
func HandleNullAddr(to *common.Address) string {
	if to == nil {
		return "0x0000000000000000000000000000000000000000000000000000000000000000"
	}
	return to.Hex()
}
