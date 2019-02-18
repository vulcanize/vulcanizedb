package cat

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
	"strconv"
)

const (
	NFlip = "nflip"
	Live  = "live"
	Vat   = "vat"
	Pit   = "pit"
	Vow   = "vow"

	IlkFlip = "flip"
	IlkChop = "chop"
	IlkLump = "lump"

	FlipIlk = "ilk"
	FlipUrn = "urn"
	FlipInk = "ink"
	FlipTab = "tab"
)

var (
	// wards takes up index 0
	IlksMappingIndex  = storage_diffs.IndexOne // bytes32 => flip address; chop (ray), lump (wad) uint256
	FlipsMappingIndex = storage_diffs.IndexTwo // uint256 => ilk, urn bytes32; ink, tab uint256 (both wad)

	NFlipKey      = common.HexToHash(storage_diffs.IndexThree)
	NFlipMetadata = shared.GetStorageValueMetadata(NFlip, nil, shared.Uint256)

	LiveKey      = common.HexToHash(storage_diffs.IndexFour)
	LiveMetadata = shared.GetStorageValueMetadata(Live, nil, shared.Uint256)

	VatKey      = common.HexToHash(storage_diffs.IndexFive)
	VatMetadata = shared.GetStorageValueMetadata(Vat, nil, shared.Address)

	PitKey      = common.HexToHash(storage_diffs.IndexSix)
	PitMetadata = shared.GetStorageValueMetadata(Pit, nil, shared.Address)

	VowKey      = common.HexToHash(storage_diffs.IndexSeven)
	VowMetadata = shared.GetStorageValueMetadata(Vow, nil, shared.Address)
)

type CatMappings struct {
	StorageRepository maker.IMakerStorageRepository
	mappings          map[common.Hash]shared.StorageValueMetadata
}

func (mappings CatMappings) Lookup(key common.Hash) (shared.StorageValueMetadata, error) {
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

func (mappings *CatMappings) SetDB(db *postgres.DB) {
	mappings.StorageRepository.SetDB(db)
}

func (mappings *CatMappings) loadMappings() error {
	mappings.mappings = loadStaticMappings()
	ilkErr := mappings.loadIlkKeys()
	if ilkErr != nil {
		return ilkErr
	}

	flipsErr := mappings.loadFlipsKeys()
	if flipsErr != nil {
		return flipsErr
	}

	return nil
}

func loadStaticMappings() map[common.Hash]shared.StorageValueMetadata {
	mappings := make(map[common.Hash]shared.StorageValueMetadata)
	mappings[NFlipKey] = NFlipMetadata
	mappings[LiveKey] = LiveMetadata
	mappings[VatKey] = VatMetadata
	mappings[PitKey] = PitMetadata
	mappings[VowKey] = VowMetadata
	return mappings
}

// Ilks
func (mappings *CatMappings) loadIlkKeys() error {
	ilks, err := mappings.StorageRepository.GetIlks()
	if err != nil {
		return err
	}
	for _, ilk := range ilks {
		mappings.mappings[getIlkFlipKey(ilk)] = getIlkFlipMetadata(ilk)
		mappings.mappings[getIlkChopKey(ilk)] = getIlkChopMetadata(ilk)
		mappings.mappings[getIlkLumpKey(ilk)] = getIlkLumpMetadata(ilk)
	}
	return nil
}

func getIlkFlipKey(ilk string) common.Hash {
	return storage_diffs.GetMapping(IlksMappingIndex, ilk)
}

func getIlkFlipMetadata(ilk string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk}
	return shared.GetStorageValueMetadata(IlkFlip, keys, shared.Address)
}

func getIlkChopKey(ilk string) common.Hash {
	return storage_diffs.GetIncrementedKey(getIlkFlipKey(ilk), 1)
}

func getIlkChopMetadata(ilk string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk}
	return shared.GetStorageValueMetadata(IlkChop, keys, shared.Uint256)
}

func getIlkLumpKey(ilk string) common.Hash {
	return storage_diffs.GetIncrementedKey(getIlkFlipKey(ilk), 2)
}

func getIlkLumpMetadata(ilk string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk}
	return shared.GetStorageValueMetadata(IlkLump, keys, shared.Uint256)
}

// Flip ID increments each time it happens, so we just need the biggest flip ID from the DB
// and we can interpolate the sequence [0..max]. This makes sure we track all earlier flips,
// even if we've missed events
func (mappings CatMappings) loadFlipsKeys() error {
	maxFlip, err := mappings.StorageRepository.GetMaxFlip()
	if err != nil {
		logrus.Error("loadFlipsKeys: error getting max flip: ", err)
		return err
	} else if maxFlip == nil { // No flips occurred yet
		return nil
	}

	last := maxFlip.Int64()
	for flip := 0; int64(flip) <= last; flip++ {
		flipStr := strconv.Itoa(flip)
		mappings.mappings[getFlipIlkKey(flipStr)] = getFlipIlkMetadata(flipStr)
		mappings.mappings[getFlipUrnKey(flipStr)] = getFlipUrnMetadata(flipStr)
		mappings.mappings[getFlipInkKey(flipStr)] = getFlipInkMetadata(flipStr)
		mappings.mappings[getFlipTabKey(flipStr)] = getFlipTabMetadata(flipStr)
	}
	return nil
}

func getFlipIlkKey(flip string) common.Hash {
	return storage_diffs.GetMapping(FlipsMappingIndex, flip)
}

func getFlipIlkMetadata(flip string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Flip: flip}
	return shared.GetStorageValueMetadata(FlipIlk, keys, shared.Bytes32)
}

func getFlipUrnKey(flip string) common.Hash {
	return storage_diffs.GetIncrementedKey(getFlipIlkKey(flip), 1)
}

func getFlipUrnMetadata(flip string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Flip: flip}
	return shared.GetStorageValueMetadata(FlipUrn, keys, shared.Bytes32)
}

func getFlipInkKey(flip string) common.Hash {
	return storage_diffs.GetIncrementedKey(getFlipIlkKey(flip), 2)
}

func getFlipInkMetadata(flip string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Flip: flip}
	return shared.GetStorageValueMetadata(FlipInk, keys, shared.Uint256)
}

func getFlipTabKey(flip string) common.Hash {
	return storage_diffs.GetIncrementedKey(getFlipIlkKey(flip), 3)
}

func getFlipTabMetadata(flip string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Flip: flip}
	return shared.GetStorageValueMetadata(FlipTab, keys, shared.Uint256)
}
