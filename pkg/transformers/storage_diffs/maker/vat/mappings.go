package vat

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

const (
	Dai     = "dai"
	Gem     = "gem"
	IlkArt  = "Art"
	IlkInk  = "Ink"
	IlkRate = "rate"
	IlkTake = "take"
	Sin     = "sin"
	UrnArt  = "art"
	UrnInk  = "ink"
	VatDebt = "debt"
	VatVice = "vice"
)

var (
	DebtKey      = common.HexToHash(storage_diffs.IndexSix)
	DebtMetadata = shared.StorageValueMetadata{
		Name: VatDebt,
		Keys: nil,
		Type: 0,
	}

	IlksMappingIndex = storage_diffs.IndexOne
	UrnsMappingIndex = storage_diffs.IndexTwo
	GemsMappingIndex = storage_diffs.IndexThree
	DaiMappingIndex  = storage_diffs.IndexFour
	SinMappingIndex  = storage_diffs.IndexFive

	ViceKey      = common.HexToHash(storage_diffs.IndexSeven)
	ViceMetadata = shared.StorageValueMetadata{
		Name: VatVice,
		Keys: nil,
		Type: 0,
	}
)

type VatMappings struct {
	StorageRepository maker.IMakerStorageRepository
	mappings          map[common.Hash]shared.StorageValueMetadata
}

func (mappings VatMappings) Lookup(key common.Hash) (shared.StorageValueMetadata, error) {
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

func (mappings *VatMappings) SetDB(db *postgres.DB) {
	mappings.StorageRepository.SetDB(db)
}

func (mappings *VatMappings) loadMappings() error {
	mappings.mappings = loadStaticMappings()
	daiErr := mappings.loadDaiKeys()
	if daiErr != nil {
		return daiErr
	}
	gemErr := mappings.loadGemKeys()
	if gemErr != nil {
		return gemErr
	}
	ilkErr := mappings.loadIlkKeys()
	if ilkErr != nil {
		return ilkErr
	}
	sinErr := mappings.loadSinKeys()
	if sinErr != nil {
		return sinErr
	}
	urnErr := mappings.loadUrnKeys()
	if urnErr != nil {
		return urnErr
	}
	return nil
}

func loadStaticMappings() map[common.Hash]shared.StorageValueMetadata {
	mappings := make(map[common.Hash]shared.StorageValueMetadata)
	mappings[DebtKey] = DebtMetadata
	mappings[ViceKey] = ViceMetadata
	return mappings
}

func (mappings *VatMappings) loadDaiKeys() error {
	daiKeys, err := mappings.StorageRepository.GetDaiKeys()
	if err != nil {
		return err
	}
	for _, d := range daiKeys {
		mappings.mappings[getDaiKey(d)] = getDaiMetadata(d)
	}
	return nil
}

func (mappings *VatMappings) loadGemKeys() error {
	gemKeys, err := mappings.StorageRepository.GetGemKeys()
	if err != nil {
		return err
	}
	for _, gem := range gemKeys {
		mappings.mappings[getGemKey(gem.Ilk, gem.Guy)] = getGemMetadata(gem.Ilk, gem.Guy)
	}
	return nil
}

func (mappings *VatMappings) loadIlkKeys() error {
	ilks, err := mappings.StorageRepository.GetIlks()
	if err != nil {
		return err
	}
	for _, ilk := range ilks {
		mappings.mappings[getIlkTakeKey(ilk)] = getIlkTakeMetadata(ilk)
		mappings.mappings[getIlkRateKey(ilk)] = getIlkRateMetadata(ilk)
		mappings.mappings[getIlkInkKey(ilk)] = getIlkInkMetadata(ilk)
		mappings.mappings[getIlkArtKey(ilk)] = getIlkArtMetadata(ilk)
	}
	return nil
}

func (mappings *VatMappings) loadSinKeys() error {
	sinKeys, err := mappings.StorageRepository.GetSinKeys()
	if err != nil {
		return err
	}
	for _, s := range sinKeys {
		mappings.mappings[getSinKey(s)] = getSinMetadata(s)
	}
	return nil
}

func (mappings *VatMappings) loadUrnKeys() error {
	urns, err := mappings.StorageRepository.GetUrns()
	if err != nil {
		return err
	}
	for _, urn := range urns {
		mappings.mappings[getUrnInkKey(urn.Ilk, urn.Guy)] = getUrnInkMetadata(urn.Ilk, urn.Guy)
		mappings.mappings[getUrnArtKey(urn.Ilk, urn.Guy)] = getUrnArtMetadata(urn.Ilk, urn.Guy)
	}
	return nil
}

func getIlkTakeKey(ilk string) common.Hash {
	return storage_diffs.GetMapping(IlksMappingIndex, ilk)
}

func getIlkTakeMetadata(ilk string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk}
	return shared.GetStorageValueMetadata(IlkTake, keys, shared.Uint256)
}

func getIlkRateKey(ilk string) common.Hash {
	return storage_diffs.GetIncrementedKey(getIlkTakeKey(ilk), 1)
}

func getIlkRateMetadata(ilk string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk}
	return shared.GetStorageValueMetadata(IlkRate, keys, shared.Uint256)
}

func getIlkInkKey(ilk string) common.Hash {
	return storage_diffs.GetIncrementedKey(getIlkTakeKey(ilk), 2)
}

func getIlkInkMetadata(ilk string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk}
	return shared.GetStorageValueMetadata(IlkInk, keys, shared.Uint256)
}

func getIlkArtKey(ilk string) common.Hash {
	return storage_diffs.GetIncrementedKey(getIlkTakeKey(ilk), 3)
}

func getIlkArtMetadata(ilk string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk}
	return shared.GetStorageValueMetadata(IlkArt, keys, shared.Uint256)
}

func getUrnInkKey(ilk, guy string) common.Hash {
	return storage_diffs.GetNestedMapping(UrnsMappingIndex, ilk, guy)
}

func getUrnInkMetadata(ilk, guy string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk, shared.Guy: guy}
	return shared.GetStorageValueMetadata(UrnInk, keys, shared.Uint256)
}

func getUrnArtKey(ilk, guy string) common.Hash {
	return storage_diffs.GetIncrementedKey(getUrnInkKey(ilk, guy), 1)
}

func getUrnArtMetadata(ilk, guy string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk, shared.Guy: guy}
	return shared.GetStorageValueMetadata(UrnArt, keys, shared.Uint256)
}

func getGemKey(ilk, guy string) common.Hash {
	return storage_diffs.GetNestedMapping(GemsMappingIndex, ilk, guy)
}

func getGemMetadata(ilk, guy string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk, shared.Guy: guy}
	return shared.GetStorageValueMetadata(Gem, keys, shared.Uint256)
}

func getDaiKey(guy string) common.Hash {
	return storage_diffs.GetMapping(DaiMappingIndex, guy)
}

func getDaiMetadata(guy string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Guy: guy}
	return shared.GetStorageValueMetadata(Dai, keys, shared.Uint256)
}

func getSinKey(guy string) common.Hash {
	return storage_diffs.GetMapping(SinMappingIndex, guy)
}

func getSinMetadata(guy string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Guy: guy}
	return shared.GetStorageValueMetadata(Sin, keys, shared.Uint256)
}
