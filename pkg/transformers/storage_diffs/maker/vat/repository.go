package vat

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

type VatStorageRepository struct {
	db *postgres.DB
}

func (repository *VatStorageRepository) Create(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Dai:
		return repository.insertDai(blockNumber, blockHash, metadata, value.(string))
	case Gem:
		return repository.insertGem(blockNumber, blockHash, metadata, value.(string))
	case IlkArt:
		return repository.insertIlkArt(blockNumber, blockHash, metadata, value.(string))
	case IlkInk:
		return repository.insertIlkInk(blockNumber, blockHash, metadata, value.(string))
	case IlkRate:
		return repository.insertIlkRate(blockNumber, blockHash, metadata, value.(string))
	case IlkTake:
		return repository.insertIlkTake(blockNumber, blockHash, metadata, value.(string))
	case Sin:
		return repository.insertSin(blockNumber, blockHash, metadata, value.(string))
	case UrnArt:
		return repository.insertUrnArt(blockNumber, blockHash, metadata, value.(string))
	case UrnInk:
		return repository.insertUrnInk(blockNumber, blockHash, metadata, value.(string))
	case VatDebt:
		return repository.insertVatDebt(blockNumber, blockHash, value.(string))
	case VatVice:
		return repository.insertVatVice(blockNumber, blockHash, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized vat contract storage name: %s", metadata.Name))
	}
}

func (repository *VatStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *VatStorageRepository) insertDai(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, dai string) error {
	guy, err := getGuy(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(`INSERT INTO maker.vat_dai (block_number, block_hash, guy, dai) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, guy, dai)
	return writeErr
}

func (repository *VatStorageRepository) insertGem(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, gem string) error {
	ilk, ilkErr := getIlk(metadata.Keys)
	if ilkErr != nil {
		return ilkErr
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return guyErr
	}
	_, writeErr := repository.db.Exec(`INSERT INTO maker.vat_gem (block_number, block_hash, ilk, guy, gem) VALUES ($1, $2, $3, $4, $5)`, blockNumber, blockHash, ilk, guy, gem)
	return writeErr
}

func (repository *VatStorageRepository) insertIlkArt(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, art string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(`INSERT INTO maker.vat_ilk_art (block_number, block_hash, ilk, art) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, ilk, art)
	return writeErr
}

func (repository *VatStorageRepository) insertIlkInk(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, ink string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(`INSERT INTO maker.vat_ilk_ink (block_number, block_hash, ilk, ink) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, ilk, ink)
	return writeErr
}

func (repository *VatStorageRepository) insertIlkRate(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, rate string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(`INSERT INTO maker.vat_ilk_rate (block_number, block_hash, ilk, rate) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, ilk, rate)
	return writeErr
}

func (repository *VatStorageRepository) insertIlkTake(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, take string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(`INSERT INTO maker.vat_ilk_take (block_number, block_hash, ilk, take) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, ilk, take)
	return writeErr
}

func (repository *VatStorageRepository) insertSin(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, sin string) error {
	guy, err := getGuy(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(`INSERT INTO maker.vat_sin (block_number, block_hash, guy, sin) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, guy, sin)
	return writeErr
}

func (repository *VatStorageRepository) insertUrnArt(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, art string) error {
	ilk, ilkErr := getIlk(metadata.Keys)
	if ilkErr != nil {
		return ilkErr
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return guyErr
	}
	_, writeErr := repository.db.Exec(`INSERT INTO maker.vat_urn_art (block_number, block_hash, ilk, urn, art) VALUES ($1, $2, $3, $4, $5)`, blockNumber, blockHash, ilk, guy, art)
	return writeErr
}

func (repository *VatStorageRepository) insertUrnInk(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, ink string) error {
	ilk, ilkErr := getIlk(metadata.Keys)
	if ilkErr != nil {
		return ilkErr
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return guyErr
	}
	_, writeErr := repository.db.Exec(`INSERT INTO maker.vat_urn_ink (block_number, block_hash, ilk, urn, ink) VALUES ($1, $2, $3, $4, $5)`, blockNumber, blockHash, ilk, guy, ink)
	return writeErr
}

func (repository *VatStorageRepository) insertVatDebt(blockNumber int, blockHash, debt string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vat_debt (block_number, block_hash, debt) VALUES ($1, $2, $3)`, blockNumber, blockHash, debt)
	return err
}

func (repository *VatStorageRepository) insertVatVice(blockNumber int, blockHash, vice string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vat_vice (block_number, block_hash, vice) VALUES ($1, $2, $3)`, blockNumber, blockHash, vice)
	return err
}

func getGuy(keys map[shared.Key]string) (string, error) {
	guy, ok := keys[shared.Guy]
	if !ok {
		return "", shared.ErrMetadataMalformed{MissingData: shared.Guy}
	}
	return guy, nil
}

func getIlk(keys map[shared.Key]string) (string, error) {
	ilk, ok := keys[shared.Ilk]
	if !ok {
		return "", shared.ErrMetadataMalformed{MissingData: shared.Ilk}
	}
	return ilk, nil
}
