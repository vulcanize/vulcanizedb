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

package storage_diffs

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

type Mappings interface {
	Lookup(key common.Hash) (shared.StorageValueMetadata, error)
	SetDB(db *postgres.DB)
}

const (
	IndexZero  = "0000000000000000000000000000000000000000000000000000000000000000"
	IndexOne   = "0000000000000000000000000000000000000000000000000000000000000001"
	IndexTwo   = "0000000000000000000000000000000000000000000000000000000000000002"
	IndexThree = "0000000000000000000000000000000000000000000000000000000000000003"
	IndexFour  = "0000000000000000000000000000000000000000000000000000000000000004"
	IndexFive  = "0000000000000000000000000000000000000000000000000000000000000005"
	IndexSix   = "0000000000000000000000000000000000000000000000000000000000000006"
	IndexSeven = "0000000000000000000000000000000000000000000000000000000000000007"
)
