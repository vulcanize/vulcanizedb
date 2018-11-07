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

package common

import (
	"bytes"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type HeaderConverter struct{}

func (converter HeaderConverter) Convert(gethHeader *types.Header) (core.Header, error) {
	writer := new(bytes.Buffer)
	err := rlp.Encode(writer, &gethHeader)
	if err != nil {
		panic(err)
	}
	coreHeader := core.Header{
		Hash:        gethHeader.Hash().Hex(),
		BlockNumber: gethHeader.Number.Int64(),
		Raw:         writer.Bytes(),
	}
	return coreHeader, nil
}
