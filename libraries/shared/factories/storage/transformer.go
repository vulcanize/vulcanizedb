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
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type ITransformer interface {
	Execute(diff types.PersistedDiff) error
	KeccakContractAddress() common.Hash
	GetStorageKeysLookup() KeysLookup
	GetContractAddress() common.Address
}

type TransformerInitializer func(db *postgres.DB) ITransformer

type Transformer struct {
	Address           common.Address
	StorageKeysLookup KeysLookup
	Repository        Repository
	hashedAddress     common.Hash
}

func (transformer Transformer) GetStorageKeysLookup() KeysLookup {
	return transformer.StorageKeysLookup
}

func (transformer Transformer) GetContractAddress() common.Address {
	return transformer.Address
}

func (transformer Transformer) NewTransformer(db *postgres.DB) ITransformer {
	transformer.StorageKeysLookup.SetDB(db)
	transformer.Repository.SetDB(db)
	return &transformer
}

func (transformer *Transformer) KeccakContractAddress() common.Hash {
	emptyHash := common.Hash{}
	if transformer.hashedAddress == emptyHash {
		transformer.hashedAddress = types.HexToKeccak256Hash(transformer.Address.Hex())
		return transformer.hashedAddress
	}

	return transformer.hashedAddress
}

func (transformer Transformer) Execute(diff types.PersistedDiff) error {
	metadata, lookupErr := transformer.StorageKeysLookup.Lookup(diff.StorageKey)
	if lookupErr != nil {
		return lookupErr
	}
	value, decodeErr := storage.Decode(diff, metadata)
	if decodeErr != nil {
		return decodeErr
	}
	return transformer.Repository.Create(diff.ID, diff.HeaderID, metadata, value)
}
