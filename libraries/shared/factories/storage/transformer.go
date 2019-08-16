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

package storage

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct {
	HashedAddress common.Hash
	Mappings      storage.Mappings
	Repository    repository.StorageRepository
}

func (transformer Transformer) NewTransformer(db *postgres.DB) transformer.StorageTransformer {
	transformer.Mappings.SetDB(db)
	transformer.Repository.SetDB(db)
	return transformer
}

func (transformer Transformer) KeccakContractAddress() common.Hash {
	return transformer.HashedAddress
}

func (transformer Transformer) Execute(diff utils.StorageDiff) error {
	metadata, lookupErr := transformer.Mappings.Lookup(diff.StorageKey)
	if lookupErr != nil {
		return lookupErr
	}
	value, decodeErr := utils.Decode(diff, metadata)
	if decodeErr != nil {
		return decodeErr
	}
	return transformer.Repository.Create(diff.BlockHeight, diff.BlockHash.Hex(), metadata, value)
}
