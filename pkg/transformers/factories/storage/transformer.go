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

package storage

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/storage"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

type Transformer struct {
	Address    common.Address
	Mappings   storage_diffs.Mappings
	Repository storage_diffs.Repository
}

func (transformer Transformer) NewTransformer(db *postgres.DB) storage.Transformer {
	transformer.Mappings.SetDB(db)
	transformer.Repository.SetDB(db)
	return transformer
}

func (transformer Transformer) ContractAddress() common.Address {
	return transformer.Address
}

func (transformer Transformer) Execute(row shared.StorageDiffRow) error {
	metadata, lookupErr := transformer.Mappings.Lookup(row.StorageKey)
	if lookupErr != nil {
		return lookupErr
	}
	value, decodeErr := shared.Decode(row, metadata)
	if decodeErr != nil {
		return decodeErr
	}
	return transformer.Repository.Create(row.BlockHeight, row.BlockHash.Hex(), metadata, value)
}
