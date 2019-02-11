/*
 *  Copyright 2018 Vulcanize
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package vow

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

const (
	VowVat  = "vat"
	CowVat  = "cow"
	RowVat  = "row"
	SinVat  = "Sin"
	AshVat  = "ash"
	WaitVat = "wait"
	SumpVat = "sump"
	BumpVat = "bump"
	HumpVat = "hump"
)

var (
	VatKey      = common.HexToHash(storage_diffs.IndexOne)
	VatMetadata = shared.StorageValueMetadata{
		Name: VowVat,
		Keys: nil,
		Type: shared.Address,
	}

	CowKey      = common.HexToHash(storage_diffs.IndexTwo)
	CowMetadata = shared.StorageValueMetadata{
		Name: CowVat,
		Keys: nil,
		Type: shared.Uint256,
	}

	RowKey      = common.HexToHash(storage_diffs.IndexThree)
	RowMetadata = shared.StorageValueMetadata{
		Name: RowVat,
		Keys: nil,
		Type: shared.Uint256,
	}

	SinKey      = common.HexToHash(storage_diffs.IndexFive)
	SinMetadata = shared.StorageValueMetadata{
		Name: SinVat,
		Keys: nil,
		Type: shared.Uint256,
	}

	AshKey      = common.HexToHash(storage_diffs.IndexSeven)
	AshMetadata = shared.StorageValueMetadata{
		Name: AshVat,
		Keys: nil,
		Type: shared.Uint256,
	}

	WaitKey      = common.HexToHash(storage_diffs.IndexEight)
	WaitMetadata = shared.StorageValueMetadata{
		Name: WaitVat,
		Keys: nil,
		Type: shared.Uint256,
	}

	SumpKey      = common.HexToHash(storage_diffs.IndexNine)
	SumpMetadata = shared.StorageValueMetadata{
		Name: SumpVat,
		Keys: nil,
		Type: shared.Uint256,
	}

	BumpKey      = common.HexToHash(storage_diffs.IndexTen)
	BumpMetadata = shared.StorageValueMetadata{
		Name: BumpVat,
		Keys: nil,
		Type: shared.Uint256,
	}

	HumpKey      = common.HexToHash(storage_diffs.IndexEleven)
	HumpMetadata = shared.StorageValueMetadata{
		Name: HumpVat,
		Keys: nil,
		Type: shared.Uint256,
	}
)

type VowMappings struct {
	StorageRepository maker.IMakerStorageRepository
	mappings          map[common.Hash]shared.StorageValueMetadata
}

func (mappings *VowMappings) Lookup(key common.Hash) (shared.StorageValueMetadata, error) {
	metadata, ok := mappings.mappings[key]
	if !ok {
		err := mappings.loadMappings()
		if err != nil {
			return metadata, err
		}
		metadata, ok = mappings.mappings[key]
		if !ok {
			return metadata, shared.ErrStorageKeyNotFound{Key: key.Hex()}
		}
	}
	return metadata, nil
}

func (mappings *VowMappings) loadMappings() error {
	staticMappings := make(map[common.Hash]shared.StorageValueMetadata)
	staticMappings[VatKey] = VatMetadata
	staticMappings[CowKey] = CowMetadata
	staticMappings[RowKey] = RowMetadata
	staticMappings[SinKey] = SinMetadata
	staticMappings[AshKey] = AshMetadata
	staticMappings[WaitKey] = WaitMetadata
	staticMappings[SumpKey] = SumpMetadata
	staticMappings[BumpKey] = BumpMetadata
	staticMappings[HumpKey] = HumpMetadata

	mappings.mappings = staticMappings

	return nil
}

func (mappings *VowMappings) SetDB(db *postgres.DB) {
	mappings.StorageRepository.SetDB(db)
}
