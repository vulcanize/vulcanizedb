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

package fakes

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var (
	FakeError     = errors.New("failed")
	FakeHash      = common.BytesToHash([]byte{1, 2, 3, 4, 5})
	fakeTimestamp = int64(111111111)
)

var rawFakeHeader, _ = json.Marshal(types.Header{})
var FakeHeader = core.Header{
	Hash:      FakeHash.String(),
	Raw:       rawFakeHeader,
	Timestamp: strconv.FormatInt(fakeTimestamp, 10),
}

func GetFakeHeader(blockNumber int64) core.Header {
	return core.Header{
		Hash:        FakeHash.String(),
		BlockNumber: blockNumber,
		Raw:         rawFakeHeader,
		Timestamp:   strconv.FormatInt(fakeTimestamp, 10),
	}
}
