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
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
)

type KeysLookup interface {
	Lookup(key common.Hash) (utils.StorageValueMetadata, error)
	SetDB(db *postgres.DB)
}

type keysLookup struct {
	loader   KeysLoader
	mappings map[common.Hash]utils.StorageValueMetadata
}

func NewKeysLookup(loader KeysLoader) KeysLookup {
	return &keysLookup{loader: loader, mappings: make(map[common.Hash]utils.StorageValueMetadata)}
}

func (lookup *keysLookup) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
	metadata, ok := lookup.mappings[key]
	if !ok {
		refreshErr := lookup.refreshMappings()
		if refreshErr != nil {
			return metadata, refreshErr
		}
		metadata, ok = lookup.mappings[key]
		if !ok {
			return metadata, utils.ErrStorageKeyNotFound{Key: key.Hex()}
		}
	}
	return metadata, nil
}

func (lookup *keysLookup) refreshMappings() error {
	var err error
	lookup.mappings, err = lookup.loader.LoadMappings()
	if err != nil {
		return err
	}
	lookup.mappings = utils.AddHashedKeys(lookup.mappings)
	return nil
}

func (lookup *keysLookup) SetDB(db *postgres.DB) {
	lookup.loader.SetDB(db)
}
