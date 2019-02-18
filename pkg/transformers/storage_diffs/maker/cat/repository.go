package cat

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

type CatStorageRepository struct {
	db *postgres.DB
}

func (repository *CatStorageRepository) Create(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case NFlip:
		return repository.insertNFlip(blockNumber, blockHash, value.(string))
	case Live:
		return repository.insertLive(blockNumber, blockHash, value.(string))
	case Vat:
		return repository.insertVat(blockNumber, blockHash, value.(string))
	case Pit:
		return repository.insertPit(blockNumber, blockHash, value.(string))
	case Vow:
		return repository.insertVow(blockNumber, blockHash, value.(string))
	case IlkFlip:
		return repository.insertIlkFlip(blockNumber, blockHash, metadata, value.(string))
	case IlkChop:
		return repository.insertIlkChop(blockNumber, blockHash, metadata, value.(string))
	case IlkLump:
		return repository.insertIlkLump(blockNumber, blockHash, metadata, value.(string))
	case FlipIlk:
		return repository.insertFlipIlk(blockNumber, blockHash, metadata, value.(string))
	case FlipUrn:
		return repository.insertFlipUrn(blockNumber, blockHash, metadata, value.(string))
	case FlipInk:
		return repository.insertFlipInk(blockNumber, blockHash, metadata, value.(string))
	case FlipTab:
		return repository.insertFlipTab(blockNumber, blockHash, metadata, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized cat contract storage name: %s", metadata.Name))
	}
}

func (repository *CatStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *CatStorageRepository) insertNFlip(blockNumber int, blockHash string, nflip string) error {
	_, writeErr := repository.db.Exec(
		`INSERT INTO maker.cat_nflip (block_number, block_hash, nflip) VALUES ($1, $2, $3)`,
		blockNumber, blockHash, nflip)
	return writeErr
}

func (repository *CatStorageRepository) insertLive(blockNumber int, blockHash string, live string) error {
	_, writeErr := repository.db.Exec(
		`INSERT INTO maker.cat_live (block_number, block_hash, live) VALUES ($1, $2, $3 )`,
		blockNumber, blockHash, live)
	return writeErr
}

func (repository *CatStorageRepository) insertVat(blockNumber int, blockHash string, vat string) error {
	_, writeErr := repository.db.Exec(
		`INSERT INTO maker.cat_vat (block_number, block_hash, vat) VALUES ($1, $2, $3 )`,
		blockNumber, blockHash, vat)
	return writeErr
}

func (repository *CatStorageRepository) insertPit(blockNumber int, blockHash string, pit string) error {
	_, writeErr := repository.db.Exec(
		`INSERT INTO maker.cat_pit (block_number, block_hash, pit) VALUES ($1, $2, $3 )`,
		blockNumber, blockHash, pit)
	return writeErr
}

func (repository *CatStorageRepository) insertVow(blockNumber int, blockHash string, vow string) error {
	_, writeErr := repository.db.Exec(
		`INSERT INTO maker.cat_vow (block_number, block_hash, vow) VALUES ($1, $2, $3 )`,
		blockNumber, blockHash, vow)
	return writeErr
}

// Ilks mapping: bytes32 => flip address; chop (ray), lump (wad) uint256
func (repository *CatStorageRepository) insertIlkFlip(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, flip string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(
		`INSERT INTO maker.cat_ilk_flip (block_number, block_hash, ilk, flip) VALUES ($1, $2, $3, $4)`,
		blockNumber, blockHash, ilk, flip)
	return writeErr
}

func (repository *CatStorageRepository) insertIlkChop(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, chop string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(
		`INSERT INTO maker.cat_ilk_chop (block_number, block_hash, ilk, chop) VALUES ($1, $2, $3, $4)`,
		blockNumber, blockHash, ilk, chop)
	return writeErr
}

func (repository *CatStorageRepository) insertIlkLump(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, lump string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(
		`INSERT INTO maker.cat_ilk_lump (block_number, block_hash, ilk, lump) VALUES ($1, $2, $3, $4)`,
		blockNumber, blockHash, ilk, lump)
	return writeErr
}

// Flips mapping: uint256 => ilk, urn bytes32; ink, tab uint256 (both wad)
func (repository *CatStorageRepository) insertFlipIlk(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, ilk string) error {
	nflip, err := getFlip(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(
		`INSERT INTO maker.cat_flip_ilk (block_number, block_hash, nflip, ilk) VALUES ($1, $2, $3, $4)`,
		blockNumber, blockHash, nflip, ilk)
	return writeErr
}

func (repository *CatStorageRepository) insertFlipUrn(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, urn string) error {
	nflip, err := getFlip(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(
		`INSERT INTO maker.cat_flip_urn (block_number, block_hash, nflip, urn) VALUES ($1, $2, $3, $4)`,
		blockNumber, blockHash, nflip, urn)
	return writeErr
}

func (repository *CatStorageRepository) insertFlipInk(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, ink string) error {
	nflip, err := getFlip(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(
		`INSERT INTO maker.cat_flip_ink (block_number, block_hash, nflip, ink) VALUES ($1, $2, $3, $4)`,
		blockNumber, blockHash, nflip, ink)
	return writeErr
}

func (repository *CatStorageRepository) insertFlipTab(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, tab string) error {
	nflip, err := getFlip(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(
		`INSERT INTO maker.cat_flip_tab (block_number, block_hash, nflip, tab) VALUES ($1, $2, $3, $4)`,
		blockNumber, blockHash, nflip, tab)
	return writeErr
}

func getIlk(keys map[shared.Key]string) (string, error) {
	ilk, ok := keys[shared.Ilk]
	if !ok {
		return "", shared.ErrMetadataMalformed{MissingData: shared.Ilk}
	}
	return ilk, nil
}

func getFlip(keys map[shared.Key]string) (string, error) {
	nflip, ok := keys[shared.Flip]
	if !ok {
		return "", shared.ErrMetadataMalformed{MissingData: shared.Flip}
	}
	return nflip, nil
}
